// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: zkp_auth.proto

package pb

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

const (
	Auth_Register_FullMethodName                      = "/zkp_auth.Auth/Register"
	Auth_CreateAuthenticationChallenge_FullMethodName = "/zkp_auth.Auth/CreateAuthenticationChallenge"
	Auth_VerifyAuthentication_FullMethodName          = "/zkp_auth.Auth/VerifyAuthentication"
)

// AuthClient is the client API for Auth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	CreateAuthenticationChallenge(ctx context.Context, in *AuthenticationChallengeRequest, opts ...grpc.CallOption) (*AuthenticationChallengeResponse, error)
	VerifyAuthentication(ctx context.Context, in *AuthenticationAnswerRequest, opts ...grpc.CallOption) (*AuthenticationAnswerResponse, error)
}

type authClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthClient(cc grpc.ClientConnInterface) AuthClient {
	return &authClient{cc}
}

func (c *authClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, Auth_Register_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) CreateAuthenticationChallenge(ctx context.Context, in *AuthenticationChallengeRequest, opts ...grpc.CallOption) (*AuthenticationChallengeResponse, error) {
	out := new(AuthenticationChallengeResponse)
	err := c.cc.Invoke(ctx, Auth_CreateAuthenticationChallenge_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) VerifyAuthentication(ctx context.Context, in *AuthenticationAnswerRequest, opts ...grpc.CallOption) (*AuthenticationAnswerResponse, error) {
	out := new(AuthenticationAnswerResponse)
	err := c.cc.Invoke(ctx, Auth_VerifyAuthentication_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServer is the server API for Auth service.
// All implementations must embed UnimplementedAuthServer
// for forward compatibility
type AuthServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	CreateAuthenticationChallenge(context.Context, *AuthenticationChallengeRequest) (*AuthenticationChallengeResponse, error)
	VerifyAuthentication(context.Context, *AuthenticationAnswerRequest) (*AuthenticationAnswerResponse, error)
	mustEmbedUnimplementedAuthServer()
}

// UnimplementedAuthServer must be embedded to have forward compatible implementations.
type UnimplementedAuthServer struct {
}

func (UnimplementedAuthServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedAuthServer) CreateAuthenticationChallenge(context.Context, *AuthenticationChallengeRequest) (*AuthenticationChallengeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAuthenticationChallenge not implemented")
}
func (UnimplementedAuthServer) VerifyAuthentication(context.Context, *AuthenticationAnswerRequest) (*AuthenticationAnswerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyAuthentication not implemented")
}
func (UnimplementedAuthServer) mustEmbedUnimplementedAuthServer() {}

// UnsafeAuthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServer will
// result in compilation errors.
type UnsafeAuthServer interface {
	mustEmbedUnimplementedAuthServer()
}

func RegisterAuthServer(s grpc.ServiceRegistrar, srv AuthServer) {
	s.RegisterService(&Auth_ServiceDesc, srv)
}

func _Auth_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auth_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_CreateAuthenticationChallenge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthenticationChallengeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).CreateAuthenticationChallenge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auth_CreateAuthenticationChallenge_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).CreateAuthenticationChallenge(ctx, req.(*AuthenticationChallengeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_VerifyAuthentication_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthenticationAnswerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).VerifyAuthentication(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auth_VerifyAuthentication_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).VerifyAuthentication(ctx, req.(*AuthenticationAnswerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Auth_ServiceDesc is the grpc.ServiceDesc for Auth service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Auth_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "zkp_auth.Auth",
	HandlerType: (*AuthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _Auth_Register_Handler,
		},
		{
			MethodName: "CreateAuthenticationChallenge",
			Handler:    _Auth_CreateAuthenticationChallenge_Handler,
		},
		{
			MethodName: "VerifyAuthentication",
			Handler:    _Auth_VerifyAuthentication_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "zkp_auth.proto",
}
