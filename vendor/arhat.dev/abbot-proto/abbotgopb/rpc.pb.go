// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: rpc.proto

package abbotgopb

import (
	context "context"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

func init() { proto.RegisterFile("rpc.proto", fileDescriptor_77a6da22d6a3feb1) }

var fileDescriptor_77a6da22d6a3feb1 = []byte{
	// 185 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2c, 0x2a, 0x48, 0xd6,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4d, 0x4c, 0x4a, 0xca, 0x2f, 0x91, 0xe2, 0x06, 0xf3,
	0x20, 0x62, 0x46, 0x36, 0x5c, 0x7c, 0x7e, 0xa9, 0x25, 0xe5, 0xf9, 0x45, 0xd9, 0xbe, 0x89, 0x79,
	0x89, 0xe9, 0xa9, 0x45, 0x42, 0x5a, 0x5c, 0xec, 0x01, 0x45, 0xf9, 0xc9, 0xa9, 0xc5, 0xc5, 0x42,
	0x7c, 0x7a, 0x60, 0x1d, 0x7a, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x52, 0xfc, 0x70, 0x7e,
	0x71, 0x41, 0x7e, 0x5e, 0x71, 0xaa, 0x53, 0xe8, 0x85, 0x87, 0x72, 0x0c, 0x37, 0x1e, 0xca, 0x31,
	0x7c, 0x78, 0x28, 0xc7, 0xd8, 0xf0, 0x48, 0x8e, 0x71, 0xc5, 0x23, 0x39, 0xc6, 0x13, 0x8f, 0xe4,
	0x18, 0x2f, 0x3c, 0x92, 0x63, 0x7c, 0xf0, 0x48, 0x8e, 0xf1, 0xc5, 0x23, 0x39, 0x86, 0x0f, 0x8f,
	0xe4, 0x18, 0x27, 0x3c, 0x96, 0x63, 0xb8, 0xf0, 0x58, 0x8e, 0xe1, 0xc6, 0x63, 0x39, 0x86, 0x28,
	0xf9, 0xc4, 0xa2, 0x8c, 0xc4, 0x12, 0xbd, 0x94, 0xd4, 0x32, 0x7d, 0xb0, 0x79, 0xba, 0x60, 0xa7,
	0x40, 0xd8, 0xe9, 0xf9, 0x05, 0x49, 0x49, 0x6c, 0x60, 0x01, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x4b, 0x9f, 0xe4, 0x36, 0xbc, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// NetworkManagerClient is the client API for NetworkManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type NetworkManagerClient interface {
	Process(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type networkManagerClient struct {
	cc *grpc.ClientConn
}

func NewNetworkManagerClient(cc *grpc.ClientConn) NetworkManagerClient {
	return &networkManagerClient{cc}
}

func (c *networkManagerClient) Process(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/abbot.NetworkManager/Process", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NetworkManagerServer is the server API for NetworkManager service.
type NetworkManagerServer interface {
	Process(context.Context, *Request) (*Response, error)
}

// UnimplementedNetworkManagerServer can be embedded to have forward compatible implementations.
type UnimplementedNetworkManagerServer struct {
}

func (*UnimplementedNetworkManagerServer) Process(ctx context.Context, req *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Process not implemented")
}

func RegisterNetworkManagerServer(s *grpc.Server, srv NetworkManagerServer) {
	s.RegisterService(&_NetworkManager_serviceDesc, srv)
}

func _NetworkManager_Process_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkManagerServer).Process(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/abbot.NetworkManager/Process",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkManagerServer).Process(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _NetworkManager_serviceDesc = grpc.ServiceDesc{
	ServiceName: "abbot.NetworkManager",
	HandlerType: (*NetworkManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Process",
			Handler:    _NetworkManager_Process_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpc.proto",
}
