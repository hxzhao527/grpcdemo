//go:generate  openssl req -x509 -nodes -newkey rsa:2048 -keyout assets/private.key -out assets/public.pem -days 3650 -subj "/CN=localhost"

package main

import (
	"flag"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"grpcdemo/pkg/server"
	helloworld_impl "grpcdemo/pkg/service/helloworld"
	routeguide_impl "grpcdemo/pkg/service/routeguide"
	"grpcdemo/pkg/util"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	port         = ":50051"
	certFilePath = "assets/public.pem"
	keyFilePath  = "assets/private.key"
	authToken    = "grpcdemo"
)

var (
	ssl  = flag.Bool("ssl", false, "whether TLS enabled")
	auth = flag.Bool("auth", false, "whether oauth enabled")
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

var (
	rpcServer *server.RPCServer
)

func validToken(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, errMissingMetadata
	}

	// The keys within metadata.MD are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	authorization := md["authorization"]
	if len(authorization) < 1 {
		return ctx, errInvalidToken
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	// Perform the token validation here.
	if token != authToken {
		return ctx, errInvalidToken
	}
	return ctx, nil
}

func recoveryHandle(ctx context.Context, method string, p interface{}) (err error) {
	rpcServer.UpdateSericeStatus(util.GetServiceNameFromFullMethod(method), grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	return server.DefaultRecoveryHandler(ctx, method, p)
}

func main() {
	flag.Parse()

	var opts []server.RPCServerOption
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
		opts = append(opts, server.WithGrpcServerOption(grpc.Creds(creds)))
	}
	if *auth {
		log.Println("server enable auth")
		opts = append(opts, server.WithAuthInterceptor(validToken, "grpc.health.v1.Health"))
	}

	//opts = append(opts, server.WithUnaryInterceptor(temp.Interceptor))
	opts = append(opts, server.WithRecovery(recoveryHandle))

	rpcServer = server.NewRPCServer(opts...)
	rpcServer.EnableHealth()

	rpcServer.RegisterService(routeguide_impl.NewServer())
	rpcServer.AttachService("helloworld.Hello", helloworld_impl.NewServer())
	log.Println("service Registered")

	go func() {
		done <- rpcServer.Run(lis)
	}()

	select {
	case err := <-done:
		{
			log.Fatalf("failed to serve: %v", err)
		}
	case <-sigs:
		{
			log.Println("Signal received: terminated by user")
			rpcServer.Stop()
		}
	}
}
