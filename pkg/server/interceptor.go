package server

import (
	"google.golang.org/grpc"
)

type MixInterceptor interface {
	UnaryInterceptor() grpc.UnaryServerInterceptor
	StreamInterceptor() grpc.StreamServerInterceptor
}

func WithMixInterceptor(mi MixInterceptor) RPCServerOption {
	return func(server *RPCServer) {
		server.grpcUnaryInterceptors = append(server.grpcUnaryInterceptors, mi.UnaryInterceptor())
		server.grpcStreamInterceptors = append(server.grpcStreamInterceptors, mi.StreamInterceptor())
	}
}

func WithUnaryInterceptor(interceptor ...grpc.UnaryServerInterceptor) RPCServerOption {
	return func(srv *RPCServer) {
		srv.grpcUnaryInterceptors = append(srv.grpcUnaryInterceptors, interceptor...)
	}
}

func WithStreamInterceptor(interceptor ...grpc.StreamServerInterceptor) RPCServerOption {
	return func(srv *RPCServer) {
		srv.grpcStreamInterceptors = append(srv.grpcStreamInterceptors, interceptor...)
	}
}
