package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	consulApi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
)

type RPCServerOption func(*RPCServer)

type RPCServer struct {
	host     string
	port     int
	listener net.Listener

	grpcsrv *grpc.Server

	grpcUnaryInterceptors  []grpc.UnaryServerInterceptor
	grpcStreamInterceptors []grpc.StreamServerInterceptor
	grpcopts               []grpc.ServerOption

	grpcsvc   map[string]Service
	healthSvc *health.Server

	consulClient     *consulApi.Client
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
func NewRPCServer(host string, port int, opts ...RPCServerOption) *RPCServer {
	srv := &RPCServer{grpcopts: make([]grpc.ServerOption, 0),
		grpcUnaryInterceptors:  make([]grpc.UnaryServerInterceptor, 0),
		grpcStreamInterceptors: make([]grpc.StreamServerInterceptor, 0),
		host:                   host,
		port:                   port,
	}
	for _, opt := range opts {
		opt(srv)
	}

	srv.grpcopts = append(srv.grpcopts, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(srv.grpcStreamInterceptors...)))
	srv.grpcopts = append(srv.grpcopts, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(srv.grpcUnaryInterceptors...)))

	srv.grpcsrv = grpc.NewServer(srv.grpcopts...)
	srv.grpcsvc = make(map[string]Service)
	srv.done = make(chan struct{})
	return srv
}

func (srv *RPCServer) Run() error {
	if srv.running {
		return errors.New("server has started, no need to started again")
	}
	var err error
	srv.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", srv.host, srv.port))
	if err != nil {
		return err
	}
	if srv.healthSvc != nil {
		srv.initServiceStatus()
		go srv.checkServiceStatusInterval(30 * time.Second)
	}
	if srv.consulClient != nil {
		srv.registerWithConsul()
	}
	srv.running = true
	return srv.grpcsrv.Serve(srv.listener)
}

func (srv *RPCServer) Stop() {
	if srv.healthSvc != nil {
		srv.healthCheckTimer.Stop()
		close(srv.done)
	}
	if srv.consulClient != nil {
		srv.deRegisterWithConsul()
	}

	srv.grpcsrv.GracefulStop()
	_ = srv.listener.Close()
	log.Println("rpc server graceful stopped successfully...")
}

// if want to use health and consul, use AttachService instead
func (srv *RPCServer) RegisterService(srs ...ServiceRegister) {
	if srv.running {
		return
	}
	for _, sr := range srs {
		sr.Register(srv.grpcsrv)
	}
}

func (srv *RPCServer) AttachService(svcs ...Service) {
	if srv.running {
		return
	}
	for _, svc := range svcs {
		srv.grpcsvc[svc.Name()] = svc
		svc.Register(srv.grpcsrv)
	}
}
