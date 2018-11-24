package server

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"grpcdemo/pkg/service"
	"log"
	"net"
	"time"
)

type RPCServerOption func(*RPCServer)

type RPCServer struct {
	grpcsrv *grpc.Server

	grpcUnaryInterceptors  []grpc.UnaryServerInterceptor
	grpcStreamInterceptors []grpc.StreamServerInterceptor
	grpcopts               []grpc.ServerOption

	grpcsvc   map[string]service.Service
	healthSvc *health.Server

	healthCheckTimer *time.Ticker

	running bool
	done    chan struct{}
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
	srv.grpcsvc = make(map[string]service.Service)
	srv.done = make(chan struct{})
	return srv
}

func (srv *RPCServer) Run(lis net.Listener) error {
	if srv.healthSvc != nil {
		srv.initServiceStatus()
		go srv.checkServiceStatusInterval(30 * time.Second)
	}
	srv.running = true
	return srv.grpcsrv.Serve(lis)
}

func (srv *RPCServer) Stop() {
	srv.healthCheckTimer.Stop()
	close(srv.done)
	srv.grpcsrv.GracefulStop()
	log.Println("rpc server graceful stopped successfully...")
}

// if want to use health, use AttachService instead
func (srv *RPCServer) RegisterService(srs ...service.ServiceRegister) {
	if srv.running {
		return
	}
	for _, sr := range srs {
		sr.Register(srv.grpcsrv)
	}
}

func (srv *RPCServer) AttachService(name string, svc service.Service) {
	if srv.running {
		return
	}
	srv.grpcsvc[name] = svc
	svc.Register(srv.grpcsrv)
}
