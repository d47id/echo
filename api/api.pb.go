// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

package api

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ShoutRequest struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShoutRequest) Reset()         { *m = ShoutRequest{} }
func (m *ShoutRequest) String() string { return proto.CompactTextString(m) }
func (*ShoutRequest) ProtoMessage()    {}
func (*ShoutRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{0}
}

func (m *ShoutRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShoutRequest.Unmarshal(m, b)
}
func (m *ShoutRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShoutRequest.Marshal(b, m, deterministic)
}
func (m *ShoutRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShoutRequest.Merge(m, src)
}
func (m *ShoutRequest) XXX_Size() int {
	return xxx_messageInfo_ShoutRequest.Size(m)
}
func (m *ShoutRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ShoutRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ShoutRequest proto.InternalMessageInfo

func (m *ShoutRequest) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type ShoutReply struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShoutReply) Reset()         { *m = ShoutReply{} }
func (m *ShoutReply) String() string { return proto.CompactTextString(m) }
func (*ShoutReply) ProtoMessage()    {}
func (*ShoutReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{1}
}

func (m *ShoutReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShoutReply.Unmarshal(m, b)
}
func (m *ShoutReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShoutReply.Marshal(b, m, deterministic)
}
func (m *ShoutReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShoutReply.Merge(m, src)
}
func (m *ShoutReply) XXX_Size() int {
	return xxx_messageInfo_ShoutReply.Size(m)
}
func (m *ShoutReply) XXX_DiscardUnknown() {
	xxx_messageInfo_ShoutReply.DiscardUnknown(m)
}

var xxx_messageInfo_ShoutReply proto.InternalMessageInfo

func (m *ShoutReply) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*ShoutRequest)(nil), "api.ShoutRequest")
	proto.RegisterType((*ShoutReply)(nil), "api.ShoutReply")
}

func init() { proto.RegisterFile("api.proto", fileDescriptor_00212fb1f9d3bf1c) }

var fileDescriptor_00212fb1f9d3bf1c = []byte{
	// 123 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4c, 0x2c, 0xc8, 0xd4,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4e, 0x2c, 0xc8, 0x54, 0xd2, 0xe0, 0xe2, 0x09, 0xce,
	0xc8, 0x2f, 0x2d, 0x09, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x11, 0x92, 0xe0, 0x62, 0xcf, 0x4d,
	0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x82, 0x71, 0x95, 0xd4,
	0xb8, 0xb8, 0xa0, 0x2a, 0x0b, 0x72, 0x2a, 0x71, 0xab, 0x33, 0x32, 0xe5, 0x62, 0x71, 0x4d, 0xce,
	0xc8, 0x17, 0xd2, 0xe5, 0x62, 0x05, 0xab, 0x17, 0x12, 0xd4, 0x03, 0xd9, 0x89, 0x6c, 0x8b, 0x14,
	0x3f, 0xb2, 0x50, 0x41, 0x4e, 0xa5, 0x12, 0x43, 0x12, 0x1b, 0xd8, 0x51, 0xc6, 0x80, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x36, 0xc0, 0xab, 0x76, 0xa1, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// EchoClient is the client API for Echo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EchoClient interface {
	Shout(ctx context.Context, in *ShoutRequest, opts ...grpc.CallOption) (*ShoutReply, error)
}

type echoClient struct {
	cc *grpc.ClientConn
}

func NewEchoClient(cc *grpc.ClientConn) EchoClient {
	return &echoClient{cc}
}

func (c *echoClient) Shout(ctx context.Context, in *ShoutRequest, opts ...grpc.CallOption) (*ShoutReply, error) {
	out := new(ShoutReply)
	err := c.cc.Invoke(ctx, "/api.Echo/Shout", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EchoServer is the server API for Echo service.
type EchoServer interface {
	Shout(context.Context, *ShoutRequest) (*ShoutReply, error)
}

func RegisterEchoServer(s *grpc.Server, srv EchoServer) {
	s.RegisterService(&_Echo_serviceDesc, srv)
}

func _Echo_Shout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShoutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EchoServer).Shout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Echo/Shout",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoServer).Shout(ctx, req.(*ShoutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Echo_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Echo",
	HandlerType: (*EchoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Shout",
			Handler:    _Echo_Shout_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}
