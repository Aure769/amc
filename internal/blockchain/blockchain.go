// Copyright 2022 The AmazeChain Authors
// This file is part of the AmazeChain library.
//
// The AmazeChain library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The AmazeChain library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the AmazeChain library. If not, see <http://www.gnu.org/licenses/>.

package blockchain

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/amazechain/amc/api/protocol/msg_proto"
	"github.com/amazechain/amc/api/protocol/types_pb"
	"github.com/amazechain/amc/common"
	block2 "github.com/amazechain/amc/common/block"
	"github.com/amazechain/amc/common/db"
	"github.com/amazechain/amc/common/message"
	"github.com/amazechain/amc/common/types"
	"github.com/amazechain/amc/internal/avm"
	"github.com/amazechain/amc/internal/consensus"
	"github.com/amazechain/amc/internal/kv"
	"github.com/amazechain/amc/log"
	event "github.com/amazechain/amc/modules/event/v2"
	"github.com/amazechain/amc/modules/rawdb"
	"github.com/amazechain/amc/modules/statedb"
	"github.com/gogo/protobuf/proto"
	lru "github.com/hashicorp/golang-lru"
	"github.com/libp2p/go-libp2p-core/peer"
)

var (
	ErrKnownBlock           = errors.New("block already known")
	ErrUnknownAncestor      = errors.New("unknown ancestor")
	ErrPrunedAncestor       = errors.New("pruned ancestor")
	ErrFutureBlock          = errors.New("block in the future")
	ErrInvalidNumber        = errors.New("invalid block number")
	ErrInvalidTerminalBlock = errors.New("insertion is interrupted")
	errChainStopped         = errors.New("blockchain is stopped")
	errInsertionInterrupted = errors.New("insertion is interrupted")
)

type WriteStatus byte

const (
	NonStatTy   WriteStatus = iota //
	CanonStatTy                    //
	SideStatTy                     //
)

const (
	//maxTimeFutureBlocks
	receiptsCacheLimit  = 32
	maxFutureBlocks     = 256
	maxTimeFutureBlocks = 10
)

type BlockChain struct {
	ctx          context.Context
	cancel       context.CancelFunc
	genesisBlock block2.IBlock
	blocks       []block2.IBlock
	headers      []block2.IHeader
	currentBlock block2.IBlock
	//state        *statedb.StateDB
	chainDB  db.IDatabase
	changeDB kv.RwDB
	engine   consensus.Engine

	insertLock    chan struct{}
	latestBlockCh chan block2.IBlock
	lock          sync.Mutex

	peers map[peer.ID]bool

	downloader common.IDownloader

	chBlocks chan block2.IBlock

	pubsub common.IPubSub

	errorCh chan error

	process *avm.VMProcessor

	wg sync.WaitGroup //

	procInterrupt int32                       // insert chain
	tdCache       map[types.Hash]types.Int256 // td cache
	futureBlocks  *lru.Cache                  //
	receiptCache  *lru.Cache

	forker    *ForkChoice
	validator Validator
}

type insertStats struct {
	queued, processed, ignored int
	usedGas                    uint64
	lastIndex                  int
	startTime                  time.Time
}

//func (bc *BlockChain) GetState() *statedb.StateDB {
//	return bc.state
//}

func (bc *BlockChain) Engine() consensus.Engine {
	return bc.engine
}

func NewBlockChain(ctx context.Context, genesisBlock block2.IBlock, engine consensus.Engine, downloader common.IDownloader, database db.IDatabase, changeDB kv.RwDB, pubsub common.IPubSub) (common.IBlockChain, error) {
	c, cancel := context.WithCancel(ctx)
	current, _ := rawdb.GetLatestBlock(database)
	if current == nil {
		current = genesisBlock
	}

	futureBlocks, _ := lru.New(maxFutureBlocks)
	receiptsCache, _ := lru.New(receiptsCacheLimit)
	bc := &BlockChain{
		genesisBlock:  genesisBlock,
		blocks:        []block2.IBlock{},
		currentBlock:  current,
		chainDB:       database,
		changeDB:      changeDB,
		ctx:           c,
		cancel:        cancel,
		insertLock:    make(chan struct{}, 1),
		peers:         make(map[peer.ID]bool),
		chBlocks:      make(chan block2.IBlock, 100),
		errorCh:       make(chan error),
		pubsub:        pubsub,
		downloader:    downloader,
		latestBlockCh: make(chan block2.IBlock, 50),
		engine:        engine,
		tdCache:       make(map[types.Hash]types.Int256),
		futureBlocks:  futureBlocks,
		receiptCache:  receiptsCache,
	}

	bc.forker = NewForkChoice(bc, nil)
	bc.process = avm.NewVMProcessor(ctx, bc, engine)
	bc.validator = NewBlockValidator(bc, engine)

	return bc, nil
}

func (bc *BlockChain) StateAt(root types.Hash) common.IStateDB {
	return statedb.NewStateDB(root, bc.chainDB, bc.changeDB)
}

func (bc *BlockChain) CurrentBlock() block2.IBlock {
	return bc.currentBlock
}

func (bc *BlockChain) Blocks() []block2.IBlock {
	return bc.blocks
}

func (bc *BlockChain) InsertHeader(headers []block2.IHeader) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (bc *BlockChain) GenesisBlock() block2.IBlock {
	return bc.genesisBlock
}

func (bc *BlockChain) Start() error {
	if bc.pubsub == nil {
		return ErrInvalidPubSub
	}

	bc.wg.Add(3)
	go bc.runLoop()
	go bc.newBlockLoop()
	go bc.updateFutureBlocksLoop()

	return nil
}

// verifyBody
// Deprecated:
func (bc *BlockChain) verifyBody(block block2.IBlock) error {
	return nil
}

// verifyState
// Deprecated:
func (bc *BlockChain) verifyState(block block2.IBlock, state *statedb.StateDB, receipts block2.Receipts, usedGas uint64) error {
	return nil
}

func (bc *BlockChain) AddPeer(hash string, remoteBlock uint64, peerID peer.ID) error {
	if bc.genesisBlock.Hash().String() != hash {
		return fmt.Errorf("failed to addPeer, err: genesis block different")
	}
	if _, ok := bc.peers[peerID]; ok {
		return fmt.Errorf("failed to addPeer, err: the peer already exists")
	}

	log.Debugf("local heigth:%d --> remote height: %d", bc.currentBlock.Number64(), remoteBlock)

	bc.peers[peerID] = true
	//if remoteBlock > bc.currentBlock.Number64().Uint64() {
	//	bc.syncChain(remoteBlock, peerID)
	//}

	return nil
}

func (bc *BlockChain) GetReceipts(blockHash types.Hash) (block2.Receipts, error) {
	return rawdb.GetReceipts(bc.chainDB, blockHash)
}

func (bc *BlockChain) GetLogs(blockHash types.Hash) ([][]*block2.Log, error) {
	receipts, err := bc.GetReceipts(blockHash)
	if err != nil {
		return nil, err
	}

	logs := make([][]*block2.Log, len(receipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	return logs, nil
}

// InsertBlock
// Deprecated:
func (bc *BlockChain) InsertBlock(blocks []block2.IBlock, isSync bool) (int, error) {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	runBlock := func(b block2.IBlock) error {
		//todo copy?
		stateDB := statedb.NewStateDB(b.ParentHash(), bc.chainDB, bc.changeDB)

		receipts, logs, usedGas, err := bc.process.Processor(b, stateDB)
		if err != nil {
			return err
		}
		//verify state
		if err = bc.verifyState(b, stateDB, receipts, usedGas); err != nil {
			return err
		}
		_, err = stateDB.Commit(b.Number64())
		if err != nil {
			return err
		}

		rawdb.StoreReceipts(bc.chainDB, b.Hash(), receipts)

		if len(logs) > 0 {
			event.GlobalEvent.Send(&common.NewLogsEvent{Logs: logs})
		}
		if len(receipts) > 0 {
			log.Infof("Receipt len(%d), receipts: [%v]", len(receipts), receipts)
		}
		return nil
	}

	current := bc.CurrentBlock()
	var insertBlocks []block2.IBlock
	for i, block := range blocks {
		if block.Number64() == current.Number64() && block.Difficulty().Compare(current.Difficulty()) == 1 {
			if err := bc.engine.VerifyHeader(bc, block.Header(), false); err != nil {
				log.Errorf("failed verify block err: %v", err)
				continue
			}
			//verify body
			if err := bc.verifyBody(block); err != nil {
				log.Errorf("failed verify block err: %v", err)
				continue
			}
			if err := runBlock(block); err != nil {
				log.Errorf("failed runblock, err:%v", err)
			}
			insertBlocks = append(insertBlocks, block)
			current = blocks[i]

		} else if block.Number64().Equal(current.Number64().Add(types.NewInt64(1))) && block.ParentHash().String() == current.Hash().String() {
			if err := bc.engine.VerifyHeader(bc, block.Header(), false); err != nil {
				log.Errorf("failed verify block err: %v", err)
				continue
			}

			if err := runBlock(block); err != nil {
				log.Errorf("failed runblock, err:%v", err)
			} else {
				insertBlocks = append(insertBlocks, block)
				current = blocks[i]
			}
		} else {
			author, _ := bc.engine.Author(block.Header())
			log.Errorf("failed instert mew block, hash: %s, number: %s, diff: %s, miner: %s, txs: %d", block.Hash(), block.Number64().String(), block.Difficulty().String(), author.String(), len(block.Transactions()))
			log.Errorf("failed instert cur block, hash: %s, number: %s, diff: %s, miner: %s, txs: %d", current.Hash(), current.Number64().String(), current.Difficulty().String(), author.String(), len(current.Transactions()))
		}

	}

	if len(insertBlocks) > 0 {
		if _, err := rawdb.SaveBlocks(bc.chainDB, insertBlocks); err != nil {
			log.Errorf("failed to save blocks, err: %v", err)
			return 0, err
		}

		if err := rawdb.SaveLatestBlock(bc.chainDB, current); err != nil {
			log.Errorf("failed to save lates blocks, err: %v", err)
			return 0, err
		}
		bc.currentBlock = current

		if !isSync {
			i := event.GlobalEvent.Send(&current)
			author, _ := bc.engine.Author(current.Header())
			log.Debugf("current number:%d, miner: %s, feed send count: %d", current.Number64().Uint64(), author, i)
		}

		return len(insertBlocks), nil
	}

	return 0, fmt.Errorf("invalid block len(%d)", len(blocks))
}

func (bc *BlockChain) LatestBlockCh() (block2.IBlock, error) {
	select {
	case <-bc.ctx.Done():
		return nil, fmt.Errorf("the main chain is closed")
	case block, ok := <-bc.latestBlockCh:
		if !ok {
			return nil, fmt.Errorf("the main chain is closed")
		}

		return block, nil
	}
}

func (bc *BlockChain) newBlockLoop() {
	bc.wg.Done()
	if bc.pubsub == nil {
		bc.errorCh <- ErrInvalidPubSub
		return
	}

	topic, err := bc.pubsub.JoinTopic(message.GossipBlockMessage)
	if err != nil {
		bc.errorCh <- ErrInvalidPubSub
		return
	}

	sub, err := topic.Subscribe()
	if err != nil {
		bc.errorCh <- ErrInvalidPubSub
		return
	}

	for {
		select {
		case <-bc.ctx.Done():
			log.Infof("block chain quit...")
			return
		default:
			msg, err := sub.Next(bc.ctx)
			if err != nil {
				bc.errorCh <- err
				return
			}

			var newBlock types_pb.PBlock
			if err := proto.Unmarshal(msg.Data, &newBlock); err == nil {
				var block block2.Block
				if err := block.FromProtoMessage(&newBlock); err == nil {
					var inserted bool
					log.Info("Subscribe new block", "hash", block.Hash())
					if _, err := bc.InsertChain([]block2.IBlock{&block}); err != nil {
						inserted = false
						log.Errorf("failed to inster new block in blockchain, err:%v", err)
					} else {
						inserted = true
						log.Debugf("successfully inserted a new block number is %d", block.Number64().Uint64())
					}
					event.GlobalEvent.Send(&common.ChainHighestBlock{Block: block, Inserted: inserted})

				} else {
					log.Errorf("unmarshal err: %v", err)
				}

			} else {
				log.Warnf("failed to unmarshal pubsub message, err:%v", err)
			}
		}
	}

}

func (bc *BlockChain) runLoop() {
	defer func() {
		bc.wg.Done()
		bc.cancel()
		bc.StopInsert()
		close(bc.errorCh)
		bc.wg.Wait()
	}()

	for {
		select {
		case <-bc.ctx.Done():
			return
		case err, ok := <-bc.errorCh:
			if ok {
				log.Errorf("receive error from action, err:%v", err)
				return
			}
		}
	}
}

// updateFutureBlocksLoop
func (bc *BlockChain) updateFutureBlocksLoop() {
	futureTimer := time.NewTicker(2 * time.Second)
	defer futureTimer.Stop()
	defer bc.wg.Done()
	for {
		select {
		case <-futureTimer.C:
			//
			blocks := make([]block2.IBlock, 0, bc.futureBlocks.Len())
			for _, key := range bc.futureBlocks.Keys() {
				if value, ok := bc.futureBlocks.Get(key); ok {
					blocks = append(blocks, value.(block2.IBlock))
				}
			}
			sort.Slice(blocks, func(i, j int) bool {
				return blocks[i].Number64().Compare(blocks[j].Number64()) < 0
			})
			for i := range blocks {
				bc.InsertChain(blocks[i : i+1])
			}
		case <-bc.ctx.Done():
			return
		}
	}
}

func (bc *BlockChain) runNewBlockMessage() {
	newBlockCh := make(chan msg_proto.NewBlockMessageData, 10)
	sub := event.GlobalEvent.Subscribe(newBlockCh)
	defer sub.Unsubscribe()

	for {
		select {
		case <-bc.ctx.Done():
			return
		case err := <-sub.Err():
			log.Errorf("failed subscribe new block at blockchain err :%v", err)
			return
		case block, ok := <-bc.chBlocks:
			if ok {
				rawdb.SaveBlocks(bc.chainDB, []block2.IBlock{block})
				rawdb.SaveLatestBlock(bc.chainDB, block)
				if block, err := rawdb.GetLatestBlock(bc.chainDB); err == nil {
					log.Debugf("latest block: %v", block.Header())
				}
			}
		case msg, ok := <-newBlockCh:
			if ok {
				block := block2.Block{}
				if err := block.FromProtoMessage(msg.Block); err == nil {
					rawdb.SaveBlocks(bc.chainDB, []block2.IBlock{&block})
					rawdb.SaveLatestBlock(bc.chainDB, &block)
					if block, err := rawdb.GetLatestBlock(bc.chainDB); err == nil {
						log.Debugf("latest block: %v", block.Header())
					}
				}
			}
		}
	}
}

func (bc *BlockChain) syncChain(remoteBlock uint64, peerID peer.ID) {
	/*sync chain
	 */
	//if remoteBlock < bc.currentBlock.Header.Number {
	//	return
	//}
	//var startNumber uint64
	//if bc.currentBlock.Header.Number == 0 {
	//	startNumber = bc.currentBlock.Header.Number
	//}
	log.Debugf("syncChain.......")
}

func (bc *BlockChain) GetHeader(h types.Hash, number types.Int256) block2.IHeader {
	header, err := bc.GetHeaderByHash(h)
	if err != nil {
		return nil
	}
	return header
}

func (bc *BlockChain) GetHeaderByNumber(number types.Int256) (block2.IHeader, error) {
	hash, err := rawdb.ReadCanonicalHash(bc.chainDB, number.Uint64())
	if nil != err {
		return nil, err
	}

	header, _, err := rawdb.GetHeader(bc.chainDB, hash)
	if err != nil {
		return nil, err
	}

	return header, nil
}

func (bc *BlockChain) GetHeaderByHash(h types.Hash) (block2.IHeader, error) {
	header, _, err := rawdb.GetHeaderByHash(bc.chainDB, h)
	return header, err
}

// GetCanonicalHash returns the canonical hash for a given block number
func (bc *BlockChain) GetCanonicalHash(number types.Int256) types.Hash {
	block, err := bc.GetBlockByNumber(number)
	if nil != err {
		return types.BytesToHash([]byte{})
	}

	return block.Hash()
}

func (bc *BlockChain) GetBlockByHash(h types.Hash) (block2.IBlock, error) {
	number, err := rawdb.GetHashNumber(bc.chainDB, h)
	if err != nil {
		return nil, err
	}
	return bc.GetBlockByNumber(number)
}

func (bc *BlockChain) GetBlockByNumber(number types.Int256) (block2.IBlock, error) {
	hash, err := rawdb.ReadCanonicalHash(bc.chainDB, number.Uint64())
	if hash == (types.Hash{}) {
		return nil, err
	}
	return bc.GetBlock(hash), nil
}

func (bc *BlockChain) NewBlockHandler(payload []byte, peer peer.ID) error {

	var nweBlock msg_proto.NewBlockMessageData
	if err := proto.Unmarshal(payload, &nweBlock); err != nil {
		log.Errorf("failed unmarshal to msg, from peer:%s", peer)
		return err
	} else {
		var block block2.Block
		if err := block.FromProtoMessage(nweBlock.GetBlock()); err == nil {
			var block block2.Block
			if err := block.FromProtoMessage(nweBlock.GetBlock()); err == nil {
				bc.chBlocks <- &block
			}
		}
	}
	return nil
}

func (bc *BlockChain) SetEngine(engine consensus.Engine) {
	bc.engine = engine
	bc.process.SetEngine(engine)
}

func (bc *BlockChain) GetBlocksFromHash(hash types.Hash, n int) (blocks []block2.IBlock) {
	h := hash
	for i := 0; i < n; i++ {
		block := bc.GetBlock(h)
		if block == nil {
			break
		}

		blocks = append(blocks, block)
		h = block.ParentHash()
	}
	return blocks
}

func (bc *BlockChain) GetBlock(hash types.Hash) block2.IBlock {
	if hash == (types.Hash{}) {
		return nil
	}
	header, h, err := rawdb.GetHeaderByHash(bc.chainDB, hash)
	if err != nil {
		return nil
	}

	if hash.String() != h.String() {
		log.Error("Failed to get block, the hash is differ", "hash", hash.String(), "headerHash", h.String())
		return nil
	}

	body, err := rawdb.GetBody(bc.chainDB, header.Hash())
	if err != nil {
		log.Error("Failed to get block body", "err", err)
		return nil
	}

	return block2.NewBlock(header, body.Transactions())
}

func (bc *BlockChain) SealedBlock(b block2.IBlock) {
	pbBlock := b.ToProtoMessage()

	_ = bc.pubsub.Publish(message.GossipBlockMessage, pbBlock)
}

// StopInsert stop insert
func (bc *BlockChain) StopInsert() {
	atomic.StoreInt32(&bc.procInterrupt, 1)
}

// insertStopped returns true after StopInsert has been called.
func (bc *BlockChain) insertStopped() bool {
	return atomic.LoadInt32(&bc.procInterrupt) == 1
}

// HasBlockAndState
func (bc *BlockChain) HasBlockAndState(hash types.Hash) bool {
	block := bc.GetBlock(hash)
	if block == nil {
		return false
	}
	return bc.HasState(block.Hash())
}

// HasState
func (bc *BlockChain) HasState(hash types.Hash) bool {
	_, err := rawdb.GetHashNumber(bc.chainDB, hash)
	if nil != err {
		return false
	}
	return true
}

func (bc *BlockChain) HasBlock(hash types.Hash) bool {
	block := bc.GetBlock(hash)
	if block == nil {
		return false
	}
	return true
}

// GetTd
func (bc *BlockChain) GetTd(hash types.Hash) types.Int256 {
	if td, ok := bc.tdCache[hash]; ok {
		return td
	}
	td, err := rawdb.ReadTd(bc.chainDB, hash)
	if nil != err {
		return types.Int256{}
	}
	bc.tdCache[hash] = td
	return bc.tdCache[hash]
}

func (bc *BlockChain) skipBlock(err error) bool {
	if !errors.Is(err, ErrKnownBlock) {
		return false
	}
	return true
}

// InsertChain
func (bc *BlockChain) InsertChain(chain []block2.IBlock) (int, error) {
	if len(chain) == 0 {
		return 0, nil
	}
	//
	for i := 1; i < len(chain); i++ {
		block, prev := chain[i], chain[i-1]
		if !block.Number64().Equal(prev.Number64().Add(types.NewInt64(1))) || block.ParentHash() != prev.Hash() {
			log.Error("Non contiguous block insert",
				"number", block.Number64().String(),
				"hash", block.Hash(),
				"parent", block.ParentHash(),
				"prev number", prev.Number64(),
				"prev hash", prev.Hash(),
			)
			return 0, fmt.Errorf("non contiguous insert: item %s is #%d [%x..], item %d is #%d [%x..] (parent [%x..])", i-1, prev.Number64().String(),
				prev.Hash().Bytes()[:4], i, block.Number64().String(), block.Hash().Bytes()[:4], block.ParentHash().Bytes()[:4])
		}
	}
	bc.lock.Lock()
	defer bc.lock.Unlock()
	return bc.insertChain(chain)
}

func (bc *BlockChain) insertChain(chain []block2.IBlock) (int, error) {
	if bc.insertStopped() {
		return 0, nil
	}

	var (
		stats     = insertStats{startTime: time.Now()}
		lastCanon block2.IBlock
	)

	defer func() {
		if lastCanon != nil && bc.CurrentBlock().Hash() == lastCanon.Hash() {
			// todo
			// event.GlobalEvent.Send(&common.ChainHighestBlock{Block: lastCanon, Inserted: true})
		}
	}()

	// Start the parallel header verifier
	headers := make([]block2.IHeader, len(chain))
	seals := make([]bool, len(chain))

	for i, block := range chain {
		headers[i] = block.Header()
		seals[i] = true
	}
	abort, results := bc.engine.VerifyHeaders(bc, headers, seals)
	defer close(abort)

	// Peek the error for the first block to decide the directing import logic
	it := newInsertIterator(chain, results, bc.validator)
	block, err := it.next()

	if bc.skipBlock(err) {
		var (
			reorg   bool
			current = bc.CurrentBlock()
		)
		for block != nil && bc.skipBlock(err) {
			reorg, err = bc.forker.ReorgNeeded(current.Header(), block.Header())
			if err != nil {
				return it.index, err
			}
			if reorg {
				// Switch to import mode if the forker says the reorg is necessary
				// and also the block is not on the canonical chain.
				// In eth2 the forker always returns true for reorg decision (blindly trusting
				// the external consensus engine), but in order to prevent the unnecessary
				// reorgs when importing known blocks, the special case is handled here.
				if block.Number64().Uint64() > current.Number64().Uint64() || bc.GetCanonicalHash(block.Number64()) != block.Hash() {
					break
				}
			}
			log.Debug("Ignoring already known block", "number", block.Number64(), "hash", block.Hash())
			stats.ignored++
			block, err = it.next()
		}
		// The remaining blocks are still known blocks, the only scenario here is:
		// During the fast sync, the pivot point is already submitted but rollback
		// happens. Then node resets the head full block to a lower height via `rollback`
		// and leaves a few known blocks in the database.
		//
		// When node runs a fast sync again, it can re-import a batch of known blocks via
		// `insertChain` while a part of them have higher total difficulty than current
		// head full block(new pivot point).
		for block != nil && bc.skipBlock(err) {
			log.Debug("Writing previously known block", "number", block.Number64(), "hash", block.Hash())
			if err := bc.writeKnownBlock(block); err != nil {
				return it.index, err
			}
			lastCanon = block

			block, err = it.next()
		}
	}

	switch {
	// First block is pruned
	case errors.Is(err, ErrPrunedAncestor):
		// First block is pruned, insert as sidechain and reorg only if TD grows enough
		log.Debug("Pruned ancestor, inserting as sidechain", "number", block.Number64(), "hash", block.Hash())
		return bc.insertSideChain(block, it)

	// First block is future, shove it (and all children) to the future queue (unknown ancestor)
	case errors.Is(err, ErrFutureBlock) || (errors.Is(err, ErrUnknownAncestor) && bc.futureBlocks.Contains(it.first().ParentHash())):
		for block != nil && (it.index == 0 || errors.Is(err, ErrUnknownAncestor)) {
			log.Debug("Future block, postponing import", "number", block.Number64(), "hash", block.Hash())
			if err := bc.addFutureBlock(block); err != nil {
				return it.index, err
			}
			block, err = it.next()
		}
		stats.queued += it.processed()
		stats.ignored += it.remaining()

		// If there are any still remaining, mark as ignored
		return it.index, err

	// Some other error(except ErrKnownBlock) occurred, abort.
	// ErrKnownBlock is allowed here since some known blocks
	// still need re-execution to generate snapshots that are missing
	case err != nil && !errors.Is(err, ErrKnownBlock):
		bc.futureBlocks.Remove(block.Hash())
		stats.ignored += len(it.chain)
		bc.reportBlock(block, nil, err)
		return it.index, err
	}

	for ; block != nil && err == nil || errors.Is(err, ErrKnownBlock); block, err = it.next() {
		// If the chain is terminating, stop processing blocks
		if bc.insertStopped() {
			log.Debug("Abort during block processing")
			break
		}

		log.Debugf("Current block: number=%v, hash=%v, difficult=%v | Insert block block: number=%v, hash=%v, difficult= %v",
			bc.CurrentBlock().Number64(), bc.CurrentBlock().Hash(), bc.CurrentBlock().Difficulty(), block.Number64(), block.Hash(), block.Difficulty())
		// Retrieve the parent block and it's state to execute on top
		start := time.Now()
		stateDB := statedb.NewStateDB(block.ParentHash(), bc.chainDB, bc.changeDB)
		var followupInterrupt uint32
		receipts, logs, usedGas, err := bc.process.Processor(block, stateDB)
		if err != nil {
			bc.reportBlock(block, receipts, err)
			atomic.StoreUint32(&followupInterrupt, 1)
			return it.index, err
		}

		if len(logs) > 0 {
			event.GlobalEvent.Send(&common.NewLogsEvent{Logs: logs})
		}

		if err := bc.validator.ValidateState(block, stateDB, receipts, usedGas); err != nil {
			bc.reportBlock(block, receipts, err)
			atomic.StoreUint32(&followupInterrupt, 1)
			return it.index, err
		}

		stateDB.Commit(block.Number64())

		var status WriteStatus
		status, err = bc.writeBlockWithState(block, receipts)
		atomic.StoreUint32(&followupInterrupt, 1)
		if err != nil {
			return it.index, err
		}

		// Report the import stats before returning the various results
		stats.processed++
		stats.usedGas += usedGas

		switch status {
		case CanonStatTy:
			log.Debug("Inserted new block ", "number ", block.Number64(), "hash", block.Hash(),
				"txs", len(block.Transactions()), "gas", block.GasUsed(),
				"elapsed", time.Since(start).Seconds(),
				"root", block.StateRoot())

			lastCanon = block

		case SideStatTy:
			log.Debug("Inserted forked block", "number", block.Number64(), "hash", block.Hash(),
				"diff", block.Difficulty(), "elapsed", time.Since(start).Seconds(),
				"txs", len(block.Transactions()), "gas", block.GasUsed(),
				"root", block.StateRoot())

		default:
			// This in theory is impossible, but lets be nice to our future selves and leave
			// a log, instead of trying to track down blocks imports that don't emit logs.
			log.Warn("Inserted block with unknown status", "number", block.Number64(), "hash", block.Hash(),
				"diff", block.Difficulty(), "elapsed", time.Since(start).Seconds(),
				"txs", len(block.Transactions()), "gas", block.GasUsed(),
				"root", block.StateRoot())
		}
	}

	// Any blocks remaining here? The only ones we care about are the future ones
	if block != nil && errors.Is(err, ErrFutureBlock) {
		if err := bc.addFutureBlock(block); err != nil {
			return it.index, err
		}
		block, err = it.next()

		for ; block != nil && errors.Is(err, ErrUnknownAncestor); block, err = it.next() {
			if err := bc.addFutureBlock(block); err != nil {
				return it.index, err
			}
			stats.queued++
		}
	}
	stats.ignored += it.remaining()

	return it.index, err
}

// insertSideChain
func (bc *BlockChain) insertSideChain(block block2.IBlock, it *insertIterator) (int, error) {
	var (
		externTd  types.Int256
		lastBlock = block
		current   = bc.CurrentBlock()
	)
	err := ErrPrunedAncestor
	for ; block != nil && errors.Is(err, ErrPrunedAncestor); block, err = it.next() {
		// Check the canonical state root for that number
		if number := block.Number64(); current.Number64().Compare(number) >= 0 {
			canonical, err := bc.GetBlockByNumber(number)
			if nil != err {
				return 0, err
			}

			if canonical != nil && canonical.Hash() == block.Hash() {
				// Not a sidechain block, this is a re-import of a canon block which has it's state pruned

				// Collect the TD of the block. Since we know it's a canon one,
				// we can get it directly, and not (like further below) use
				// the parent and then add the block on top
				externTd = bc.GetTd(block.Hash())
				continue
			}
			if canonical != nil && canonical.StateRoot() == block.StateRoot() {
				// This is most likely a shadow-state attack. When a fork is imported into the
				// database, and it eventually reaches a block height which is not pruned, we
				// just found that the state already exist! This means that the sidechain block
				// refers to a state which already exists in our canon chain.
				//
				// If left unchecked, we would now proceed importing the blocks, without actually
				// having verified the state of the previous blocks.
				log.Warn("Sidechain ghost-state attack detected", "number", block.Number64(), "sideroot", block.StateRoot(), "canonroot", canonical.StateRoot())

				// If someone legitimately side-mines blocks, they would still be imported as usual. However,
				// we cannot risk writing unverified blocks to disk when they obviously target the pruning
				// mechanism.
				return it.index, errors.New("sidechain ghost-state attack")
			}
		}
		if externTd == types.NewInt64(0) {
			externTd = bc.GetTd(block.ParentHash())
		}
		externTd = externTd.Add(block.Difficulty())

		if !bc.HasBlock(block.Hash()) {
			start := time.Now()
			if err := bc.WriteBlockWithoutState(block); err != nil {
				return it.index, err
			}
			log.Debug("Injected sidechain block", "number", block.Number64(), "hash", block.Hash(),
				"diff", block.Difficulty(), "elapsed", time.Since(start).Seconds(),
				"txs", len(block.Transactions()), "gas", block.GasUsed(),
				"root", block.StateRoot())
		}
		lastBlock = block
	}

	reorg, err := bc.forker.ReorgNeeded(current.Header(), lastBlock.Header())
	if err != nil {
		return it.index, err
	}

	if !reorg {
		localTd := bc.GetTd(current.Hash())
		log.Info("Sidechain written to disk", "start", it.first().Number64(), "end", it.previous().Number64(), "sidetd", externTd, "localtd", localTd)
		return it.index, err
	}
	var (
		hashes  []types.Hash
		numbers []uint64
	)
	parent := it.previous()
	for parent != nil && !bc.HasState(parent.StateRoot()) {
		hashes = append(hashes, parent.Hash())
		numbers = append(numbers, parent.Number64().Uint64())

		parent = bc.GetHeader(parent.(*block2.Header).ParentHash, parent.Number64().Sub(types.NewInt64(1)))
	}
	if parent == nil {
		return it.index, errors.New("missing parent")
	}
	// Import all the pruned blocks to make the state available
	var (
		blocks []block2.IBlock
		memory types.StorageSize
	)
	for i := len(hashes) - 1; i >= 0; i-- {
		// Append the next block to our batch
		block := bc.GetBlock(hashes[i])

		blocks = append(blocks, block)

		// If memory use grew too large, import and continue. Sadly we need to discard
		// all raised events and logs from notifications since we're too heavy on the
		// memory here.
		if len(blocks) >= 2048 || memory > 64*1024*1024 {
			log.Info("Importing heavy sidechain segment", "blocks", len(blocks), "start", blocks[0].Number64(), "end", block.Number64())
			if _, err := bc.insertChain(blocks); err != nil {
				return 0, err
			}
			blocks, memory = blocks[:0], 0
			memory += block.Size()
			// If the chain is terminating, stop processing blocks
			if bc.insertStopped() {
				log.Debug("Abort during blocks processing")
				return 0, nil
			}
		}
	}
	if len(blocks) > 0 {
		log.Info("Importing sidechain segment", "start", blocks[0].Number64(), "end", blocks[len(blocks)-1].Number64())
		return bc.insertChain(blocks)
	}
	return 0, nil
}

// recoverAncestors
func (bc *BlockChain) recoverAncestors(block block2.IBlock) (types.Hash, error) {
	var (
		hashes  []types.Hash
		numbers []types.Int256
		parent  = block
	)
	for parent != nil && !bc.HasState(parent.Hash()) {
		hashes = append(hashes, parent.Hash())
		numbers = append(numbers, parent.Number64())
		parent = bc.GetBlock(parent.ParentHash())

	}
	if parent == nil {
		return types.Hash{}, errors.New("missing parent")
	}
	for i := len(hashes) - 1; i >= 0; i-- {

		var b block2.IBlock
		if i == 0 {
			b = block
		} else {
			b = bc.GetBlock(hashes[i])
		}
		if _, err := bc.insertChain([]block2.IBlock{b}); err != nil {
			return b.ParentHash(), err
		}
	}
	return block.Hash(), nil
}

// ValidateBody verify body
func (bc *BlockChain) ValidateBody(block block2.IBlock) error {
	if bc.HasBlockAndState(block.Hash()) {
		return ErrKnownBlock
	}

	// todo verify Transactions root

	if !bc.HasBlockAndState(block.ParentHash()) {
		if !bc.HasBlock(block.ParentHash()) {
			return ErrUnknownAncestor
		}
		return ErrPrunedAncestor
	}
	return nil
}

// ValidateState verify state
func (v *BlockChain) ValidateState(block block2.IBlock, receipts []*block2.Receipt, usedGas uint64) error {
	return nil
}

// WriteBlockWithoutState without state
func (bc *BlockChain) WriteBlockWithoutState(block block2.IBlock) (err error) {
	if bc.insertStopped() {
		return errInsertionInterrupted
	}
	//if err := bc.state.WriteTD(block.Hash(), td); err != nil {
	//	return err
	//}
	if _, err := rawdb.SaveBlocks(bc.chainDB, []block2.IBlock{block}); err != nil {
		return err
	}
	return nil
}

// WriteBlockWithState
func (bc *BlockChain) writeBlockWithState(block block2.IBlock, receipts []*block2.Receipt) (status WriteStatus, err error) {
	ptd := bc.GetTd(block.ParentHash())
	externTd := ptd.Add(block.Difficulty())
	if err := rawdb.WriteTd(bc.chainDB, block.Hash(), externTd); nil != err {
		return NonStatTy, err
	}
	log.Debugf("writeTd: hash=%s, td=%d", block.Hash(), externTd.Uint64())

	if _, err := rawdb.SaveBlocks(bc.chainDB, []block2.IBlock{block}); err != nil {
		return NonStatTy, err
	}

	reorg, err := bc.forker.ReorgNeeded(bc.currentBlock.Header(), block.Header())
	if nil != err {
		return NonStatTy, err
	}
	if reorg {
		// Reorganise the chain if the parent is not the head block
		if block.ParentHash() != bc.currentBlock.Hash() {
			if err := bc.reorg(bc.currentBlock, block); err != nil {
				return NonStatTy, err
			}
		}
		status = CanonStatTy
	} else {
		status = SideStatTy
	}
	// Set new head.
	if status == CanonStatTy {
		if err := bc.writeHeadBlock(block); nil != err {
			log.Errorf("failed to save lates blocks, err: %v", err)
			return NonStatTy, err
		}
	}
	if len(receipts) > 0 {
		if err := rawdb.StoreReceipts(bc.chainDB, block.Hash(), receipts); nil != err {
			log.Errorf("rawdb.StoreReceipts failed err= %v", err)
			return NonStatTy, err
		}
	}
	//
	if _, ok := bc.futureBlocks.Get(block.Hash()); ok {
		bc.futureBlocks.Remove(block.Hash())
	}

	return status, nil
}

// writeHeadBlock head
func (bc *BlockChain) writeHeadBlock(block block2.IBlock) error {

	if err := rawdb.SaveLatestBlock(bc.chainDB, block); err != nil {
		log.Errorf("failed to save last block, err: %v", err)
		return err
	}
	rawdb.WriteCanonicalHash(bc.chainDB, block.Hash(), block.Number64().Uint64())
	bc.currentBlock = block
	return nil
}

// reportBlock logs a bad block error.
func (bc *BlockChain) reportBlock(block block2.IBlock, receipts []*block2.Receipt, err error) {

	var receiptString string
	for i, receipt := range receipts {
		receiptString += fmt.Sprintf("\t %d: cumulative: %v gas: %v contract: %v status: %v tx: %v logs: %v bloom: %x state: %x\n",
			i, receipt.CumulativeGasUsed, receipt.GasUsed, receipt.ContractAddress.String(),
			receipt.Status, receipt.TxHash.String(), "Logs", receipt.Bloom, receipt.PostState)
	}
	log.Error(fmt.Sprintf(`
########## BAD BLOCK #########

Number: %v
Hash: %#x
%v

Error: %v
##############################
`, block.Number64().String(), block.Hash(), receiptString, err))
}

// ReorgNeeded
func (bc *BlockChain) ReorgNeeded(current block2.IBlock, header block2.IBlock) bool {
	switch current.Number64().Compare(header.Number64()) {
	case 1:
		return false
	case 0:
		return current.Difficulty().Compare(types.NewInt64(2)) != 0
	}
	return true
}

// SetHead set new head
func (bc *BlockChain) SetHead(head uint64) error {
	// todo state rawdb 。。。。。
	newHeadBlock, err := bc.GetBlockByNumber(types.NewInt64(head))
	if err != nil {
		return nil
	}
	err = rawdb.SaveLatestBlock(bc.chainDB, newHeadBlock)
	return err
}

// addFutureBlock checks if the block is within the max allowed window to get
// accepted for future processing, and returns an error if the block is too far
// ahead and was not added.
//
// TODO after the transition, the future block shouldn't be kept. Because
// it's not checked in the Geth side anymore.
func (bc *BlockChain) addFutureBlock(block block2.IBlock) error {
	max := uint64(time.Now().Unix() + maxTimeFutureBlocks)
	if block.Time() > max {
		return fmt.Errorf("future block timestamp %v > allowed %v", block.Time(), max)
	}
	if block.Difficulty().Uint64() == 0 {
		// Never add PoS blocks into the future queue
		return nil
	}
	bc.futureBlocks.Add(block.Hash(), block)
	return nil
}

// writeKnownBlock updates the head block flag with a known block
// and introduces chain reorg if necessary.
func (bc *BlockChain) writeKnownBlock(block block2.IBlock) error {
	current := bc.CurrentBlock()
	if block.ParentHash() != current.Hash() {
		if err := bc.reorg(current, block); err != nil {
			return err
		}
	}
	bc.writeHeadBlock(block)
	return nil
}

// reorg takes two blocks, an old chain and a new chain and will reconstruct the
// blocks and inserts them to be part of the new canonical chain and accumulates
// potential missing transactions and post an event about them.
// Note the new head block won't be processed here, callers need to handle it
// externally.
func (bc *BlockChain) reorg(oldBlock, newBlock block2.IBlock) error {
	var (
		newChain    block2.Blocks
		oldChain    block2.Blocks
		commonBlock block2.IBlock

		deletedTxs []types.Hash
		addedTxs   []types.Hash
	)
	// Reduce the longer chain to the same number as the shorter one
	if oldBlock.Number64().Uint64() > newBlock.Number64().Uint64() {
		// Old chain is longer, gather all transactions and logs as deleted ones
		for ; oldBlock != nil && oldBlock.Number64().Uint64() != newBlock.Number64().Uint64(); oldBlock = bc.GetBlock(oldBlock.ParentHash()) {
			oldChain = append(oldChain, oldBlock)
			for _, tx := range oldBlock.Transactions() {
				hash, _ := tx.Hash()
				deletedTxs = append(deletedTxs, hash)
			}
		}
	} else {
		// New chain is longer, stash all blocks away for subsequent insertion
		for ; newBlock != nil && newBlock.Number64() != oldBlock.Number64(); newBlock = bc.GetBlock(newBlock.ParentHash()) {
			newChain = append(newChain, newBlock)
		}
	}
	if oldBlock == nil {
		return fmt.Errorf("invalid old chain")
	}
	if newBlock == nil {
		return fmt.Errorf("invalid new chain")
	}
	// Both sides of the reorg are at the same number, reduce both until the common
	// ancestor is found
	for {
		// If the common ancestor was found, bail out
		if oldBlock.Hash() == newBlock.Hash() {
			commonBlock = oldBlock
			break
		}
		// Remove an old block as well as stash away a new block
		oldChain = append(oldChain, oldBlock)
		for _, tx := range oldBlock.Transactions() {
			h, _ := tx.Hash()
			deletedTxs = append(deletedTxs, h)
		}
		newChain = append(newChain, newBlock)

		// Step back with both chains
		oldBlock = bc.GetBlock(oldBlock.ParentHash())
		if oldBlock == nil {
			return fmt.Errorf("invalid old chain")
		}
		newBlock = bc.GetBlock(newBlock.ParentHash())
		if newBlock == nil {
			return fmt.Errorf("invalid new chain")
		}
	}

	// Ensure the user sees large reorgs
	if len(oldChain) > 0 && len(newChain) > 0 {
		logFn := log.Info
		msg := "Chain reorg detected"
		if len(oldChain) > 63 {
			msg = "Large chain reorg detected"
			logFn = log.Warn
		}
		logFn(msg, "number", commonBlock.Number64(), "hash", commonBlock.Hash(),
			"drop", len(oldChain), "dropfrom", oldChain[0].Hash(), "add", len(newChain), "addfrom", newChain[0].Hash())
	} else if len(newChain) > 0 {
		// Special case happens in the post merge stage that current head is
		// the ancestor of new head while these two blocks are not consecutive
		log.Info("Extend chain", "add", len(newChain), "number", newChain[0].Number64(), "hash", newChain[0].Hash())
	} else {
		// len(newChain) == 0 && len(oldChain) > 0
		// rewind the canonical chain to a lower point.
		log.Error("Impossible reorg, please file an issue", "oldnum", oldBlock.Number64(), "oldhash", oldBlock.Hash(), "oldblocks", len(oldChain), "newnum", newBlock.Number64(), "newhash", newBlock.Hash(), "newblocks", len(newChain))
	}
	// Insert the new chain(except the head block(reverse order)),
	// taking care of the proper incremental order.
	for i := len(newChain) - 1; i >= 1; i-- {
		// Insert the block in the canonical way, re-writing history
		bc.writeHeadBlock(newChain[i])

		// Collect the new added transactions.
		for _, tx := range newChain[i].Transactions() {
			h, _ := tx.Hash()
			addedTxs = append(addedTxs, h)
		}
	}

	// Delete useless indexes right now which includes the non-canonical
	// transaction indexes, canonical chain indexes which above the head.
	for _, tx := range types.HashDifference(deletedTxs, addedTxs) {
		rawdb.DeleteTransactionIndex(bc.chainDB, tx)
	}

	// Delete all hash markers that are not part of the new canonical chain.
	// Because the reorg function does not handle new chain head, all hash
	// markers greater than or equal to new chain head should be deleted.
	number := commonBlock.Number64().Uint64()
	if len(newChain) > 1 {
		number = newChain[1].Number64().Uint64()
	}
	for i := number + 1; ; i++ {
		hash, _ := rawdb.ReadCanonicalHash(bc.chainDB, i)
		if hash == (types.Hash{}) {
			break
		}
		rawdb.DeleteCanonicalHash(bc.chainDB, i)
	}
	return nil
}