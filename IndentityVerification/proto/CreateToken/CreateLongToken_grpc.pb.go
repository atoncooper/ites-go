// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.20.3
// source: CreateLongToken.proto

package __

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	CreateLongToken_CreateLongToken_FullMethodName = "/CreateLongToken.CreateLongToken/createLongToken"
)

// CreateLongTokenClient is the client API for CreateLongToken service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CreateLongTokenClient interface {
	CreateLongToken(ctx context.Context, in *CreateLongTokenRequest, opts ...grpc.CallOption) (*CreateLongTokenResponse, error)
}

type createLongTokenClient struct {
	cc grpc.ClientConnInterface
}

func NewCreateLongTokenClient(cc grpc.ClientConnInterface) CreateLongTokenClient {
	return &createLongTokenClient{cc}
}

func (c *createLongTokenClient) CreateLongToken(ctx context.Context, in *CreateLongTokenRequest, opts ...grpc.CallOption) (*CreateLongTokenResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateLongTokenResponse)
	err := c.cc.Invoke(ctx, CreateLongToken_CreateLongToken_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CreateLongTokenServer is the server API for CreateLongToken service.
// All implementations must embed UnimplementedCreateLongTokenServer
// for forward compatibility.
type CreateLongTokenServer interface {
	CreateLongToken(context.Context, *CreateLongTokenRequest) (*CreateLongTokenResponse, error)
	mustEmbedUnimplementedCreateLongTokenServer()
}

// UnimplementedCreateLongTokenServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCreateLongTokenServer struct{}

func (UnimplementedCreateLongTokenServer) CreateLongToken(context.Context, *CreateLongTokenRequest) (*CreateLongTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateLongToken not implemented")
}
func (UnimplementedCreateLongTokenServer) mustEmbedUnimplementedCreateLongTokenServer() {}
func (UnimplementedCreateLongTokenServer) testEmbeddedByValue()                         {}

// UnsafeCreateLongTokenServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CreateLongTokenServer will
// result in compilation errors.
type UnsafeCreateLongTokenServer interface {
	mustEmbedUnimplementedCreateLongTokenServer()
}

func RegisterCreateLongTokenServer(s grpc.ServiceRegistrar, srv CreateLongTokenServer) {
	// If the following call pancis, it indicates UnimplementedCreateLongTokenServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CreateLongToken_ServiceDesc, srv)
}

func _CreateLongToken_CreateLongToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateLongTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CreateLongTokenServer).CreateLongToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CreateLongToken_CreateLongToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CreateLongTokenServer).CreateLongToken(ctx, req.(*CreateLongTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CreateLongToken_ServiceDesc is the grpc.ServiceDesc for CreateLongToken service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CreateLongToken_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "CreateLongToken.CreateLongToken",
	HandlerType: (*CreateLongTokenServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "createLongToken",
			Handler:    _CreateLongToken_CreateLongToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "CreateLongToken.proto",
}
