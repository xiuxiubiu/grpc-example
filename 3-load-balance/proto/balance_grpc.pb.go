// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// BalanceClient is the client API for Balance service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BalanceClient interface {
	Call(ctx context.Context, in *CallRequest, opts ...grpc.CallOption) (*CallResponse, error)
}

type balanceClient struct {
	cc grpc.ClientConnInterface
}

func NewBalanceClient(cc grpc.ClientConnInterface) BalanceClient {
	return &balanceClient{cc}
}

func (c *balanceClient) Call(ctx context.Context, in *CallRequest, opts ...grpc.CallOption) (*CallResponse, error) {
	out := new(CallResponse)
	err := c.cc.Invoke(ctx, "/proto.Balance/Call", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BalanceServer is the server API for Balance service.
// All implementations must embed UnimplementedBalanceServer
// for forward compatibility
type BalanceServer interface {
	Call(context.Context, *CallRequest) (*CallResponse, error)
	mustEmbedUnimplementedBalanceServer()
}

// UnimplementedBalanceServer must be embedded to have forward compatible implementations.
type UnimplementedBalanceServer struct {
}

func (UnimplementedBalanceServer) Call(context.Context, *CallRequest) (*CallResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Call not implemented")
}
func (UnimplementedBalanceServer) mustEmbedUnimplementedBalanceServer() {}

// UnsafeBalanceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BalanceServer will
// result in compilation errors.
type UnsafeBalanceServer interface {
	mustEmbedUnimplementedBalanceServer()
}

func RegisterBalanceServer(s grpc.ServiceRegistrar, srv BalanceServer) {
	s.RegisterService(&Balance_ServiceDesc, srv)
}

func _Balance_Call_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CallRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BalanceServer).Call(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Balance/Call",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BalanceServer).Call(ctx, req.(*CallRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Balance_ServiceDesc is the grpc.ServiceDesc for Balance service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Balance_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Balance",
	HandlerType: (*BalanceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Call",
			Handler:    _Balance_Call_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "balance.proto",
}
