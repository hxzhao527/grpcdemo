//go:generate  openssl req -x509 -nodes -newkey rsa:2048 -keyout assets/private.key -out assets/public.pem -days 3650 -subj "/CN=localhost"

package main

import (
	"flag"
	helloworld_impl "grpcdemo/pkg/service/helloworld"
	routeguide_impl "grpcdemo/pkg/service/routeguide"
	"grpcdemo/proto/helloworld"
	"grpcdemo/proto/routeguide"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

const (
	port         = ":50051"
	certFilePath = "assets/public.pem"
	keyFilePath  = "assets/private.key"
)

var (
	ssl = flag.Bool("ssl", false, "whether TLS enabled")
)

func main() {
	flag.Parse()

	var opts []grpc.ServerOption
	done := make(chan error)
	sigs := make(chan os.Signal)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	lis, err := net.Listen("tcp", port)
	defer lis.Close()
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if *ssl {
		log.Println("server in secure mode")
		creds, err := credentials.NewServerTLSFromFile(certFilePath, keyFilePath)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.Creds(creds))
	}

	rpcServer := grpc.NewServer(opts...)
	log.Println("server started")

	helloworld.RegisterHelloServer(rpcServer, helloworld_impl.NewServer())
	routeguide.RegisterRouteGuideServer(rpcServer, routeguide_impl.NewServer())
	log.Println("service Registered")

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
