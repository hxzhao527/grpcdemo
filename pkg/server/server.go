package server

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"grpcdemo/pkg/service"
	"log"
	"net"
)

type RPCServerOption func(*RPCServer)

type RPCServer struct {
	grpcsrv *grpc.Server

	grpcUnaryInterceptors  []grpc.UnaryServerInterceptor
	grpcStreamInterceptors []grpc.StreamServerInterceptor
	grpcopts               []grpc.ServerOption
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

func WithGrpcServerOption(grpcopts ...grpc.ServerOption) RPCServerOption {
	return func(srv *RPCServer) {
		srv.grpcopts = append(srv.grpcopts, grpcopts...)
	}
}

// https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
func NewRPCServer(opts ...RPCServerOption) *RPCServer {
	srv := &RPCServer{grpcopts: make([]grpc.ServerOption, 0), grpcUnaryInterceptors: make([]grpc.UnaryServerInterceptor, 0), grpcStreamInterceptors: make([]grpc.StreamServerInterceptor, 0)}
	for _, opt := range opts {
		opt(srv)
	}

	srv.grpcopts = append(srv.grpcopts, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(srv.grpcStreamInterceptors...)))
	srv.grpcopts = append(srv.grpcopts, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(srv.grpcUnaryInterceptors...)))

	srv.grpcsrv = grpc.NewServer(srv.grpcopts...)
	return srv
}

func (srv *RPCServer) Run(lis net.Listener) error {
	return srv.grpcsrv.Serve(lis)
}

func (srv *RPCServer) Stop() {
	srv.grpcsrv.GracefulStop()
	log.Println("rpc server graceful stopped successfully...")
}

func (srv *RPCServer) RegisterService(srs ...service.ServiceRegister) {
	for _, sr := range srs {
		sr.Register(srv.grpcsrv)
	}
}
