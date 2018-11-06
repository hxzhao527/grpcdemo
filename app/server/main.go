package main

import (
	helloworld_impl "grpcdemo/pkg/service/helloword"
	routeguide_impl "grpcdemo/pkg/service/routeguide"
	"grpcdemo/proto/helloworld"
	"grpcdemo/proto/routeguide"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {
	done := make(chan error)
	sigs := make(chan os.Signal)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	rpcServer := grpc.NewServer()

	helloworld.RegisterHelloServer(rpcServer, helloworld_impl.NewServer())
	routeguide.RegisterRouteGuideServer(rpcServer, routeguide_impl.NewServer())
	go func() {
		done <- rpcServer.Serve(lis)
	}()

	select {
	case err := <-done:
		{
			log.Fatalf("failed to serve: %v", err)
		}
	case <-sigs:
		{
			log.Println("Signal received: terminated by user")
			rpcServer.GracefulStop()
			log.Println("rpc server graceful stopped successfully...")
		}
	}
}
