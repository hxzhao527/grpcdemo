package server

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"grpcdemo/pkg/util"

	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func WithAuthInterceptor(authFunc grpc_auth.AuthFunc, excludeSvc ...string) RPCServerOption {
	authInterceptor := &MixAuthInterceptor{authFunc: authFunc}
	authInterceptor.exclude = make(map[string]struct{})
	for _, svc := range excludeSvc {
		authInterceptor.exclude[svc] = struct{}{}
	}
	return WithMixInterceptor(authInterceptor)
}

type MixAuthInterceptor struct {
	authFunc grpc_auth.AuthFunc
	exclude  map[string]struct{}
}

func (i *MixAuthInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if _, ok := i.exclude[util.GetServiceNameFromFullMethod(info.FullMethod)]; ok {
			return handler(ctx, req)
		}
		newCtx, err := i.authFunc(ctx)

		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

func (i *MixAuthInterceptor) StreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if _, ok := i.exclude[util.GetServiceNameFromFullMethod(info.FullMethod)]; ok {
			return handler(srv, stream)
		}
		newCtx, err := i.authFunc(stream.Context())

		if err != nil {
			return err
		}
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(srv, stream)
	}
}

var (
	// to check EnsureValidToken whether implement grpc.UnaryServerInterceptor
	_ MixInterceptor = &MixAuthInterceptor{}
)
