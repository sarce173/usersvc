package grpc

import (
	"context"
	"net"
)

type ServerOption interface{}

type Server struct {
	methods []*ServiceDesc
}

func NewServer(opts ...ServerOption) *Server {
	return &Server{}
}

func (s *Server) RegisterService(sd *ServiceDesc, impl interface{}) {
	s.methods = append(s.methods, sd)
}

func (s *Server) Serve(lis net.Listener) error {
	// No-op stub to keep binaries buildable in environments without the real gRPC runtime.
	return nil
}

func (s *Server) GracefulStop() {}

type UnaryServerInterceptor func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (interface{}, error)

type UnaryHandler func(ctx context.Context, req interface{}) (interface{}, error)

type UnaryServerInfo struct {
	Server     interface{}
	FullMethod string
}

type MethodDesc struct {
	MethodName string
	Handler    func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor UnaryServerInterceptor) (interface{}, error)
}

type StreamDesc struct{}

type ServiceDesc struct {
	ServiceName string
	HandlerType interface{}
	Methods     []MethodDesc
	Streams     []StreamDesc
	Metadata    interface{}
}

type ServiceRegistrar interface {
	RegisterService(*ServiceDesc, interface{})
}

type CallOption interface{}

type ClientConnInterface interface {
	Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...CallOption) error
}
