//go:generate  openssl req -x509 -nodes -newkey rsa:2048 -keyout ../../assets/private.key -out ../../assets/public.pem -days 3650 -subj "/CN=*"

package main

import (
	"flag"
	"grpcdemo/internal/server"
	helloworldImpl "grpcdemo/internal/service/helloworld"
	routeguideImpl "grpcdemo/internal/service/routeguide"
	rpcServer "grpcdemo/pkg/server"
	"log"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	port         = 50051
	certFilePath = "assets/public.pem"
	keyFilePath  = "assets/private.key"
	authToken    = "grpcdemo"
)

var (
	ssl           = flag.Bool("ssl", false, "whether TLS enabled")
	auth          = flag.Bool("auth", false, "whether oauth enabled")
	health        = flag.Bool("health", false, "whether enable health")
	consulAddress = flag.String("consul", "", "consul address to register svc")
)

var (
	myServer = new(rpcServer.RPCServer)
)

func main() {
	flag.Parse()

	var opts []rpcServer.RPCServerOption
	done := make(chan error)
	sigs := make(chan os.Signal)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	if *ssl {
		log.Println("server in secure mode")
		creds, err := credentials.NewServerTLSFromFile(certFilePath, keyFilePath)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, rpcServer.WithGrpcServerOption(grpc.Creds(creds)))
	}
	if *auth {
		log.Println("server enable auth")
		opts = append(opts, rpcServer.WithAuthInterceptor(server.AuthFunc(authToken), "grpc.health.v1.Health"))
	}

	opts = append(opts, rpcServer.WithRecovery(server.RecoveryFunc(myServer)))

	if len(*consulAddress) > 0 {
		opts = append(opts, rpcServer.WithConsulIntegration(*consulAddress))
	}

	*myServer = *rpcServer.NewRPCServer("", port, opts...) // ??
	if *health {
		myServer.EnableHealth()
	}
	log.Printf("2: %p", myServer)
	myServer.RegisterService(routeguideImpl.NewServer())
	myServer.AttachService(helloworldImpl.NewServer())
	log.Println("service Registered")

	go func() {
		done <- myServer.Run()
	}()

	select {
	case err := <-done:
		{
			log.Fatalf("failed to serve: %v", err)
		}
	case <-sigs:
		{
			log.Println("Signal received: terminated by user")
			myServer.Stop()
		}
	}
}
