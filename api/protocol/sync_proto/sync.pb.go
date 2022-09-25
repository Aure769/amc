// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sync.proto

package sync_proto

import (
	fmt "fmt"
	types_pb "github.com/amazechain/amc/api/protocol/types_pb"
	github_com_amazechain_amc_common_types "github.com/amazechain/amc/common/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type SyncType int32

const (
	SyncType_FINDReq           SyncType = 0
	SyncType_FindRes           SyncType = 1
	SyncType_HeaderReq         SyncType = 2
	SyncType_HeaderRes         SyncType = 3
	SyncType_BodyReq           SyncType = 4
	SyncType_BodyRes           SyncType = 5
	SyncType_StateReq          SyncType = 6
	SyncType_StateRes          SyncType = 7
	SyncType_TransactionReq    SyncType = 8
	SyncType_TransactionRes    SyncType = 9
	SyncType_PeerInfoBroadcast SyncType = 10
)

var SyncType_name = map[int32]string{
	0:  "FINDReq",
	1:  "FindRes",
	2:  "HeaderReq",
	3:  "HeaderRes",
	4:  "BodyReq",
	5:  "BodyRes",
	6:  "StateReq",
	7:  "StateRes",
	8:  "TransactionReq",
	9:  "TransactionRes",
	10: "PeerInfoBroadcast",
}

var SyncType_value = map[string]int32{
	"FINDReq":           0,
	"FindRes":           1,
	"HeaderReq":         2,
	"HeaderRes":         3,
	"BodyReq":           4,
	"BodyRes":           5,
	"StateReq":          6,
	"StateRes":          7,
	"TransactionReq":    8,
	"TransactionRes":    9,
	"PeerInfoBroadcast": 10,
}

func (x SyncType) String() string {
	return proto.EnumName(SyncType_name, int32(x))
}

func (SyncType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_5273b98214de8075, []int{0}
}

type SyncProtocol struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SyncProtocol) Reset()         { *m = SyncProtocol{} }
func (m *SyncProtocol) String() string { return proto.CompactTextString(m) }
func (*SyncProtocol) ProtoMessage()    {}
func (*SyncProtocol) Descriptor() ([]byte, []int) {
	return fileDescriptor_5273b98214de8075, []int{0}
}
func (m *SyncProtocol) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SyncProtocol.Unmarshal(m, b)
}
func (m *SyncProtocol) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SyncProtocol.Marshal(b, m, deterministic)
}
func (m *SyncProtocol) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncProtocol.Merge(m, src)
}
func (m *SyncProtocol) XXX_Size() int {
	return xxx_messageInfo_SyncProtocol.Size(m)
}
func (m *SyncProtocol) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncProtocol.DiscardUnknown(m)
}

var xxx_messageInfo_SyncProtocol proto.InternalMessageInfo

type Value struct {
	Hash                 string   `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Height               uint64   `protobuf:"varint,2,opt,name=height,proto3" json:"height,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Value) Reset()         { *m = Value{} }
func (m *Value) String() string { return proto.CompactTextString(m) }
func (*Value) ProtoMessage()    {}
func (*Value) Descriptor() ([]byte, []int) {
	return fileDescriptor_5273b98214de8075, []int{1}
}
func (m *Value) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Value.Unmarshal(m, b)
}
func (m *Value) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Value.Marshal(b, m, deterministic)
}
func (m *Value) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Value.Merge(m, src)
}
func (m *Value) XXX_Size() int {
	return xxx_messageInfo_Value.Size(m)
}
func (m *Value) XXX_DiscardUnknown() {
	xxx_messageInfo_Value.DiscardUnknown(m)
}

var xxx_messageInfo_Value proto.InternalMessageInfo

func (m *Value) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *Value) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

type SyncBlockRequest struct {
	Number               []github_com_amazechain_amc_common_types.Int256 `protobuf:"bytes,1,rep,name=number,proto3,customtype=github.com/amazechain/amc/common/types.Int256" json:"number"`
	XXX_NoUnkeyedLiteral struct{}                                        `json:"-"`
	XXX_unrecognized     []byte                                          `json:"-"`
	XXX_sizecache        int32                                           `json:"-"`
}

func (m *SyncBlockRequest) Reset()         { *m = SyncBlockRequest{} }
func (m *SyncBlockRequest) String() string { return proto.CompactTextString(m) }
func (*SyncBlockRequest) ProtoMessage()    {}
func (*SyncBlockRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_5273b98214de8075, []int{2}
}
func (m *SyncBlockRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SyncBlockRequest.Unmarshal(m, b)
}
func (m *SyncBlockRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SyncBlockRequest.Marshal(b, m, deterministic)
}
func (m *SyncBlockRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncBlockRequest.Merge(m, src)
}
func (m *SyncBlockRequest) XXX_Size() int {
	return xxx_messageInfo_SyncBlockRequest.Size(m)
}
func (m *SyncBlockRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncBlockRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SyncBlockRequest proto.InternalMessageInfo

type SyncBlockResponse struct {
	Blocks               []*types_pb.PBlock `protobuf:"bytes,1,rep,name=blocks,proto3" json:"blocks,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *SyncBlockResponse) Reset()         { *m = SyncBlockResponse{} }
func (m *SyncBlockResponse) String() string { return proto.CompactTextString(m) }
func (*SyncBlockResponse) ProtoMessage()    {}
func (*SyncBlockResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_5273b98214de8075, []int{3}
}
func (m *SyncBlockResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SyncBlockResponse.Unmarshal(m, b)
}
func (m *SyncBlockResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SyncBlockResponse.Marshal(b, m, deterministic)
}
func (m *SyncBlockResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncBlockResponse.Merge(m, src)
}
func (m *SyncBlockResponse) XXX_Size() int {
	return xxx_messageInfo_SyncBlockResponse.Size(m)
}
func (m *SyncBlockResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncBlockResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SyncBlockResponse proto.InternalMessageInfo

func (m *SyncBlockResponse) GetBlocks() []*types_pb.PBlock {
	if m != nil {
		return m.Blocks
	}
	return nil
}

type SyncHeaderRequest struct {
	Number               github_com_amazechain_amc_common_types.Int256 `protobuf:"bytes,1,opt,name=number,proto3,customtype=github.com/amazechain/amc/common/types.Int256" json:"number"`
	Amount               github_com_amazechain_amc_common_types.Int256 `protobuf:"bytes,3,opt,name=amount,proto3,customtype=github.com/amazechain/amc/common/types.Int256" json:"amount"`
	XXX_NoUnkeyedLiteral struct{}                                      `json:"-"`
	XXX_unrecognized     []byte                                        `json:"-"`
	XXX_sizecache        int32                                         `json:"-"`
}

func (m *SyncHeaderRequest) Reset()         { *m = SyncHeaderRequest{} }
func (m *SyncHeaderRequest) String() string { return proto.CompactTextString(m) }
func (*SyncHeaderRequest) ProtoMessage()    {}
func (*SyncHeaderRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_5273b98214de8075, []int{4}
}
func (m *SyncHeaderRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SyncHeaderRequest.Unmarshal(m, b)
}
func (m *SyncHeaderRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SyncHeaderRequest.Marshal(b, m, deterministic)
}
func (m *SyncHeaderRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncHeaderRequest.Merge(m, src)
}
func (m *SyncHeaderRequest) XXX_Size() int {
	return xxx_messageInfo_SyncHeaderRequest.Size(m)
}
func (m *SyncHeaderRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncHeaderRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SyncHeaderRequest proto.InternalMessageInfo

type SyncHeaderResponse struct {
	Headers              []*types_pb.PBHeader `protobuf:"bytes,1,rep,name=headers,proto3" json:"headers,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *SyncHeaderResponse) Reset()         { *m = SyncHeaderResponse{} }
func (m *SyncHeaderResponse) String() string { return proto.CompactTextString(m) }
func (*SyncHeaderResponse) ProtoMessage()    {}
func (*SyncHeaderResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_5273b98214de8075, []int{5}
}
func (m *SyncHeaderResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SyncHeaderResponse.Unmarshal(m, b)
}
func (m *SyncHeaderResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SyncHeaderResponse.Marshal(b, m, deterministic)
}
func (m *SyncHeaderResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncHeaderResponse.Merge(m, src)
}
func (m *SyncHeaderResponse) XXX_Size() int {
	return xxx_messageInfo_SyncHeaderResponse.Size(m)
}
func (m *SyncHeaderResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncHeaderResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SyncHeaderResponse proto.InternalMessageInfo

func (m *SyncHeaderResponse) GetHeaders() []*types_pb.PBHeader {
	if m != nil {
		return m.Headers
	}
	return nil
}

type SyncTransactionRequest struct {
	Bloom                []byte   `protobuf:"bytes,1,opt,name=bloom,proto3" json:"bloom,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SyncTransactionRequest) Reset()         { *m = SyncTransactionRequest{} }
func (m *SyncTransactionRequest) String() string { return proto.CompactTextString(m) }
func (*SyncTransactionRequest) ProtoMessage()    {}
func (*SyncTransactionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_5273b98214de8075, []int{6}
}
func (m *SyncTransactionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SyncTransactionRequest.Unmarshal(m, b)
}
func (m *SyncTransactionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SyncTransactionRequest.Marshal(b, m, deterministic)
}
func (m *SyncTransactionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncTransactionRequest.Merge(m, src)
}
func (m *SyncTransactionRequest) XXX_Size() int {
	return xxx_messageInfo_SyncTransactionRequest.Size(m)
}
func (m *SyncTransactionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncTransactionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SyncTransactionRequest proto.InternalMessageInfo

func (m *SyncTransactionRequest) GetBloom() []byte {
	if m != nil {
		return m.Bloom
	}
	return nil
}

type SyncTransactionResponse struct {
	Transactions         []*types_pb.Transaction `protobuf:"bytes,1,rep,name=transactions,proto3" json:"transactions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *SyncTransactionResponse) Reset()         { *m = SyncTransactionResponse{} }
func (m *SyncTransactionResponse) String() string { return proto.CompactTextString(m) }
func (*SyncTransactionResponse) ProtoMessage()    {}
func (*SyncTransactionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_5273b98214de8075, []int{7}
}
func (m *SyncTransactionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SyncTransactionResponse.Unmarshal(m, b)
}
func (m *SyncTransactionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SyncTransactionResponse.Marshal(b, m, deterministic)
}
func (m *SyncTransactionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncTransactionResponse.Merge(m, src)
}
func (m *SyncTransactionResponse) XXX_Size() int {
	return xxx_messageInfo_SyncTransactionResponse.Size(m)
}
func (m *SyncTransactionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncTransactionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SyncTransactionResponse proto.InternalMessageInfo

func (m *SyncTransactionResponse) GetTransactions() []*types_pb.Transaction {
	if m != nil {
		return m.Transactions
	}
	return nil
}

type SyncPeerInfoBroadcast struct {
	Difficulty           github_com_amazechain_amc_common_types.Int256 `protobuf:"bytes,1,opt,name=Difficulty,proto3,customtype=github.com/amazechain/amc/common/types.Int256" json:"Difficulty"`
	Number               github_com_amazechain_amc_common_types.Int256 `protobuf:"bytes,2,opt,name=Number,proto3,customtype=github.com/amazechain/amc/common/types.Int256" json:"Number"`
	XXX_NoUnkeyedLiteral struct{}                                      `json:"-"`
	XXX_unrecognized     []byte                                        `json:"-"`
	XXX_sizecache        int32                                         `json:"-"`
}

func (m *SyncPeerInfoBroadcast) Reset()         { *m = SyncPeerInfoBroadcast{} }
func (m *SyncPeerInfoBroadcast) String() string { return proto.CompactTextString(m) }
func (*SyncPeerInfoBroadcast) ProtoMessage()    {}
func (*SyncPeerInfoBroadcast) Descriptor() ([]byte, []int) {
	return fileDescriptor_5273b98214de8075, []int{8}
}
func (m *SyncPeerInfoBroadcast) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SyncPeerInfoBroadcast.Unmarshal(m, b)
}
func (m *SyncPeerInfoBroadcast) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SyncPeerInfoBroadcast.Marshal(b, m, deterministic)
}
func (m *SyncPeerInfoBroadcast) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncPeerInfoBroadcast.Merge(m, src)
}
func (m *SyncPeerInfoBroadcast) XXX_Size() int {
	return xxx_messageInfo_SyncPeerInfoBroadcast.Size(m)
}
func (m *SyncPeerInfoBroadcast) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncPeerInfoBroadcast.DiscardUnknown(m)
}

var xxx_messageInfo_SyncPeerInfoBroadcast proto.InternalMessageInfo

type SyncTask struct {
	Id       uint64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Ok       bool     `protobuf:"varint,2,opt,name=ok,proto3" json:"ok,omitempty"`
	SyncType SyncType `protobuf:"varint,3,opt,name=syncType,proto3,enum=sync_proto.SyncType" json:"syncType,omitempty"`
	// Types that are valid to be assigned to Payload:
	//	*SyncTask_SyncHeaderRequest
	//	*SyncTask_SyncHeaderResponse
	//	*SyncTask_SyncBlockRequest
	//	*SyncTask_SyncBlockResponse
	//	*SyncTask_SyncTransactionRequest
	//	*SyncTask_SyncTransactionResponse
	//	*SyncTask_SyncPeerInfoBroadcast
	Payload              isSyncTask_Payload `protobuf_oneof:"payload"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *SyncTask) Reset()         { *m = SyncTask{} }
func (m *SyncTask) String() string { return proto.CompactTextString(m) }
func (*SyncTask) ProtoMessage()    {}
func (*SyncTask) Descriptor() ([]byte, []int) {
	return fileDescriptor_5273b98214de8075, []int{9}
}
func (m *SyncTask) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SyncTask.Unmarshal(m, b)
}
func (m *SyncTask) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SyncTask.Marshal(b, m, deterministic)
}
func (m *SyncTask) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncTask.Merge(m, src)
}
func (m *SyncTask) XXX_Size() int {
	return xxx_messageInfo_SyncTask.Size(m)
}
func (m *SyncTask) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncTask.DiscardUnknown(m)
}

var xxx_messageInfo_SyncTask proto.InternalMessageInfo

type isSyncTask_Payload interface {
	isSyncTask_Payload()
}

type SyncTask_SyncHeaderRequest struct {
	SyncHeaderRequest *SyncHeaderRequest `protobuf:"bytes,4,opt,name=syncHeaderRequest,proto3,oneof" json:"syncHeaderRequest,omitempty"`
}
type SyncTask_SyncHeaderResponse struct {
	SyncHeaderResponse *SyncHeaderResponse `protobuf:"bytes,5,opt,name=syncHeaderResponse,proto3,oneof" json:"syncHeaderResponse,omitempty"`
}
type SyncTask_SyncBlockRequest struct {
	SyncBlockRequest *SyncBlockRequest `protobuf:"bytes,6,opt,name=syncBlockRequest,proto3,oneof" json:"syncBlockRequest,omitempty"`
}
type SyncTask_SyncBlockResponse struct {
	SyncBlockResponse *SyncBlockResponse `protobuf:"bytes,7,opt,name=syncBlockResponse,proto3,oneof" json:"syncBlockResponse,omitempty"`
}
type SyncTask_SyncTransactionRequest struct {
	SyncTransactionRequest *SyncTransactionRequest `protobuf:"bytes,8,opt,name=syncTransactionRequest,proto3,oneof" json:"syncTransactionRequest,omitempty"`
}
type SyncTask_SyncTransactionResponse struct {
	SyncTransactionResponse *SyncTransactionResponse `protobuf:"bytes,9,opt,name=syncTransactionResponse,proto3,oneof" json:"syncTransactionResponse,omitempty"`
}
type SyncTask_SyncPeerInfoBroadcast struct {
	SyncPeerInfoBroadcast *SyncPeerInfoBroadcast `protobuf:"bytes,10,opt,name=syncPeerInfoBroadcast,proto3,oneof" json:"syncPeerInfoBroadcast,omitempty"`
}

func (*SyncTask_SyncHeaderRequest) isSyncTask_Payload()       {}
func (*SyncTask_SyncHeaderResponse) isSyncTask_Payload()      {}
func (*SyncTask_SyncBlockRequest) isSyncTask_Payload()        {}
func (*SyncTask_SyncBlockResponse) isSyncTask_Payload()       {}
func (*SyncTask_SyncTransactionRequest) isSyncTask_Payload()  {}
func (*SyncTask_SyncTransactionResponse) isSyncTask_Payload() {}
func (*SyncTask_SyncPeerInfoBroadcast) isSyncTask_Payload()   {}

func (m *SyncTask) GetPayload() isSyncTask_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *SyncTask) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *SyncTask) GetOk() bool {
	if m != nil {
		return m.Ok
	}
	return false
}

func (m *SyncTask) GetSyncType() SyncType {
	if m != nil {
		return m.SyncType
	}
	return SyncType_FINDReq
}

func (m *SyncTask) GetSyncHeaderRequest() *SyncHeaderRequest {
	if x, ok := m.GetPayload().(*SyncTask_SyncHeaderRequest); ok {
		return x.SyncHeaderRequest
	}
	return nil
}

func (m *SyncTask) GetSyncHeaderResponse() *SyncHeaderResponse {
	if x, ok := m.GetPayload().(*SyncTask_SyncHeaderResponse); ok {
		return x.SyncHeaderResponse
	}
	return nil
}

func (m *SyncTask) GetSyncBlockRequest() *SyncBlockRequest {
	if x, ok := m.GetPayload().(*SyncTask_SyncBlockRequest); ok {
		return x.SyncBlockRequest
	}
	return nil
}

func (m *SyncTask) GetSyncBlockResponse() *SyncBlockResponse {
	if x, ok := m.GetPayload().(*SyncTask_SyncBlockResponse); ok {
		return x.SyncBlockResponse
	}
	return nil
}

func (m *SyncTask) GetSyncTransactionRequest() *SyncTransactionRequest {
	if x, ok := m.GetPayload().(*SyncTask_SyncTransactionRequest); ok {
		return x.SyncTransactionRequest
	}
	return nil
}

func (m *SyncTask) GetSyncTransactionResponse() *SyncTransactionResponse {
	if x, ok := m.GetPayload().(*SyncTask_SyncTransactionResponse); ok {
		return x.SyncTransactionResponse
	}
	return nil
}

func (m *SyncTask) GetSyncPeerInfoBroadcast() *SyncPeerInfoBroadcast {
	if x, ok := m.GetPayload().(*SyncTask_SyncPeerInfoBroadcast); ok {
		return x.SyncPeerInfoBroadcast
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*SyncTask) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*SyncTask_SyncHeaderRequest)(nil),
		(*SyncTask_SyncHeaderResponse)(nil),
		(*SyncTask_SyncBlockRequest)(nil),
		(*SyncTask_SyncBlockResponse)(nil),
		(*SyncTask_SyncTransactionRequest)(nil),
		(*SyncTask_SyncTransactionResponse)(nil),
		(*SyncTask_SyncPeerInfoBroadcast)(nil),
	}
}

func init() {
	proto.RegisterEnum("sync_proto.SyncType", SyncType_name, SyncType_value)
	proto.RegisterType((*SyncProtocol)(nil), "sync_proto.SyncProtocol")
	proto.RegisterType((*Value)(nil), "sync_proto.Value")
	proto.RegisterType((*SyncBlockRequest)(nil), "sync_proto.SyncBlockRequest")
	proto.RegisterType((*SyncBlockResponse)(nil), "sync_proto.SyncBlockResponse")
	proto.RegisterType((*SyncHeaderRequest)(nil), "sync_proto.SyncHeaderRequest")
	proto.RegisterType((*SyncHeaderResponse)(nil), "sync_proto.SyncHeaderResponse")
	proto.RegisterType((*SyncTransactionRequest)(nil), "sync_proto.SyncTransactionRequest")
	proto.RegisterType((*SyncTransactionResponse)(nil), "sync_proto.SyncTransactionResponse")
	proto.RegisterType((*SyncPeerInfoBroadcast)(nil), "sync_proto.SyncPeerInfoBroadcast")
	proto.RegisterType((*SyncTask)(nil), "sync_proto.SyncTask")
}

func init() { proto.RegisterFile("sync.proto", fileDescriptor_5273b98214de8075) }

var fileDescriptor_5273b98214de8075 = []byte{
	// 697 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x95, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0xc7, 0xe3, 0x34, 0x9f, 0xd3, 0x10, 0xb9, 0xab, 0xa6, 0xb5, 0x10, 0xd0, 0x60, 0x2e, 0x11,
	0x02, 0x1b, 0x52, 0x15, 0x89, 0x03, 0x17, 0xab, 0x42, 0x29, 0x52, 0xab, 0x6a, 0x5b, 0x90, 0x40,
	0x48, 0xd5, 0xda, 0xde, 0xc6, 0x56, 0x62, 0xaf, 0x9b, 0x75, 0x0e, 0xe6, 0xad, 0x38, 0x71, 0xea,
	0x9d, 0x67, 0xe0, 0xd0, 0x67, 0x41, 0x5e, 0x6f, 0x1a, 0x27, 0x76, 0x38, 0xd0, 0x53, 0x77, 0x66,
	0xff, 0xf3, 0x9b, 0x99, 0xf5, 0x4c, 0x03, 0xc0, 0x93, 0xd0, 0x31, 0xa2, 0x19, 0x8b, 0x19, 0x12,
	0xe7, 0x2b, 0x71, 0x7e, 0xbc, 0x1b, 0x27, 0x11, 0xe5, 0x57, 0x91, 0x6d, 0x8a, 0x83, 0x21, 0xbd,
	0x63, 0x36, 0x66, 0xe2, 0x68, 0xa6, 0xa7, 0xcc, 0xab, 0x77, 0xa1, 0x73, 0x91, 0x84, 0xce, 0x79,
	0x6a, 0x38, 0x6c, 0xaa, 0x1f, 0x42, 0xfd, 0x0b, 0x99, 0xce, 0x29, 0x42, 0x50, 0xf3, 0x08, 0xf7,
	0x34, 0xa5, 0xaf, 0x0c, 0xda, 0x58, 0x9c, 0xd1, 0x1e, 0x34, 0x3c, 0xea, 0x8f, 0xbd, 0x58, 0xab,
	0xf6, 0x95, 0x41, 0x0d, 0x4b, 0x4b, 0x27, 0xa0, 0xa6, 0x10, 0x6b, 0xca, 0x9c, 0x09, 0xa6, 0x37,
	0x73, 0xca, 0x63, 0x74, 0x0a, 0x8d, 0x70, 0x1e, 0xd8, 0x74, 0xa6, 0x29, 0xfd, 0xad, 0x41, 0xdb,
	0x3a, 0xfa, 0x7d, 0x77, 0x50, 0xf9, 0x73, 0x77, 0xf0, 0x7a, 0xec, 0xc7, 0xde, 0xdc, 0x36, 0x1c,
	0x16, 0x98, 0x24, 0x20, 0x3f, 0xa8, 0xe3, 0x11, 0x3f, 0x34, 0x49, 0xe0, 0x98, 0x0e, 0x0b, 0x02,
	0x16, 0xca, 0xa2, 0x4f, 0xc2, 0x78, 0x78, 0xf4, 0x0e, 0x4b, 0x88, 0xfe, 0x01, 0x76, 0x72, 0x29,
	0x78, 0xc4, 0x42, 0x4e, 0xd1, 0x00, 0x1a, 0x76, 0xea, 0xe0, 0x22, 0xc7, 0xf6, 0x50, 0x35, 0x64,
	0xc3, 0xb6, 0x71, 0x9e, 0x29, 0xe5, 0xbd, 0xfe, 0x53, 0xc9, 0xe2, 0x47, 0x94, 0xb8, 0x74, 0x56,
	0x56, 0xa3, 0x32, 0xe8, 0x3c, 0xb0, 0xc6, 0x14, 0x47, 0x02, 0x36, 0x0f, 0x63, 0x6d, 0xeb, 0x41,
	0xb8, 0x0c, 0xa2, 0x5b, 0x80, 0xf2, 0x25, 0xcb, 0x9e, 0x5f, 0x41, 0xd3, 0x13, 0x9e, 0x45, 0xd3,
	0x28, 0xdf, 0xb4, 0x14, 0x2f, 0x24, 0xba, 0x01, 0x7b, 0x29, 0xe3, 0x72, 0x46, 0x42, 0x4e, 0x9c,
	0xd8, 0x67, 0xe1, 0xa2, 0xf7, 0x5d, 0xa8, 0xdb, 0x53, 0xc6, 0x82, 0xac, 0x75, 0x9c, 0x19, 0xfa,
	0x25, 0xec, 0x17, 0xf4, 0x32, 0xf1, 0x7b, 0xe8, 0xc4, 0x4b, 0xf7, 0x22, 0x7b, 0x6f, 0x99, 0x3d,
	0x1f, 0xb4, 0x22, 0xd5, 0x6f, 0x15, 0xe8, 0x89, 0x29, 0xa3, 0x74, 0x76, 0x12, 0x5e, 0x33, 0x6b,
	0xc6, 0x88, 0xeb, 0x10, 0x1e, 0xa3, 0xcf, 0x00, 0xc7, 0xfe, 0xf5, 0xb5, 0xef, 0xcc, 0xa7, 0x71,
	0x92, 0xcd, 0xda, 0xff, 0x3e, 0x5b, 0x0e, 0x94, 0x7e, 0x89, 0xb3, 0xec, 0xc3, 0x56, 0x1f, 0x82,
	0x94, 0x10, 0xfd, 0xb6, 0x0e, 0x2d, 0xf1, 0x2c, 0x84, 0x4f, 0x50, 0x17, 0xaa, 0xbe, 0x2b, 0x4a,
	0xad, 0xe1, 0xaa, 0xef, 0xa6, 0x36, 0x9b, 0x88, 0x3c, 0x2d, 0x5c, 0x65, 0x13, 0xf4, 0x06, 0x5a,
	0xe9, 0x2e, 0x5e, 0x26, 0x11, 0x15, 0x73, 0xd0, 0x1d, 0xee, 0x1a, 0xcb, 0xe5, 0x34, 0x2e, 0xe4,
	0x1d, 0xbe, 0x57, 0xa1, 0x53, 0xd8, 0xe1, 0xeb, 0xb3, 0xa9, 0xd5, 0xfa, 0xca, 0x60, 0x7b, 0xf8,
	0x74, 0x3d, 0x74, 0x45, 0x34, 0xaa, 0xe0, 0x62, 0x24, 0x3a, 0x07, 0xc4, 0x0b, 0x73, 0xa3, 0xd5,
	0x05, 0xef, 0xd9, 0x26, 0x5e, 0xa6, 0x1a, 0x55, 0x70, 0x49, 0x2c, 0xfa, 0x04, 0x2a, 0x5f, 0xdb,
	0x6f, 0xad, 0x21, 0x78, 0x4f, 0xd6, 0x79, 0x79, 0xcd, 0xa8, 0x82, 0x0b, 0x71, 0x8b, 0x66, 0x57,
	0x16, 0x59, 0x6b, 0x96, 0x37, 0xbb, 0x22, 0x5a, 0x34, 0xbb, 0xfa, 0x2f, 0xe0, 0x3b, 0xec, 0xf1,
	0xd2, 0x01, 0xd7, 0x5a, 0x82, 0xa9, 0x17, 0xde, 0xbe, 0xa0, 0x1c, 0x55, 0xf0, 0x06, 0x06, 0xba,
	0x82, 0x7d, 0x5e, 0xbe, 0x0e, 0x5a, 0x5b, 0xe0, 0x5f, 0xfc, 0x13, 0x7f, 0x5f, 0xf8, 0x26, 0x0a,
	0xfa, 0x0a, 0x3d, 0x5e, 0xb6, 0x18, 0x1a, 0x08, 0xfc, 0xf3, 0x75, 0x7c, 0x41, 0x38, 0xaa, 0xe0,
	0x72, 0x82, 0xd5, 0x86, 0x66, 0x44, 0x92, 0x29, 0x23, 0xee, 0xcb, 0x5f, 0x8a, 0x9c, 0xdf, 0x74,
	0xda, 0xb6, 0xa1, 0xf9, 0xf1, 0xe4, 0xec, 0x18, 0xd3, 0x1b, 0xb5, 0x22, 0x0c, 0x3f, 0x74, 0x31,
	0xe5, 0xaa, 0x82, 0x1e, 0x41, 0xfb, 0x7e, 0x92, 0xd4, 0x6a, 0xde, 0xe4, 0xea, 0x56, 0x2a, 0xb5,
	0x98, 0x9b, 0xa4, 0x77, 0xb5, 0xa5, 0xc1, 0xd5, 0x3a, 0xea, 0x40, 0xeb, 0x22, 0x26, 0x31, 0x4d,
	0xaf, 0x1a, 0x39, 0x8b, 0xab, 0x4d, 0x84, 0xa0, 0xbb, 0xfa, 0xae, 0x6a, 0xab, 0xe0, 0xe3, 0x6a,
	0x1b, 0xf5, 0x60, 0xa7, 0xd0, 0x82, 0x0a, 0xd6, 0xe1, 0xb7, 0xb7, 0x9b, 0x57, 0x96, 0x44, 0xbe,
	0x19, 0xc9, 0x1f, 0x2e, 0x73, 0xf9, 0x50, 0x76, 0x43, 0xfc, 0x39, 0xfc, 0x1b, 0x00, 0x00, 0xff,
	0xff, 0xd7, 0x0d, 0x36, 0x68, 0x20, 0x07, 0x00, 0x00,
}