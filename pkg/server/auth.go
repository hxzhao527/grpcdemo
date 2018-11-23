package server

import (
	"grpcdemo/pkg/util"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

func WithAuthInterceptor(token string, excludeSvc ...string) RPCServerOption {
	authInterceptor := &MixAuthInterceptor{token: token}
	authInterceptor.exclude = make(map[string]struct{})
	for _, svc := range excludeSvc {
		authInterceptor.exclude[svc] = struct{}{}
	}
	return WithMixInterceptor(authInterceptor)
}

type MixAuthInterceptor struct {
	token   string
	exclude map[string]struct{}
}

func (i *MixAuthInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if _, ok := i.exclude[util.GetServiceNameFromFullMethod(info.FullMethod)]; ok {
			return handler(ctx, req)
		}
		err := i.validToken(ctx)

		if err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func (i *MixAuthInterceptor) StreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if _, ok := i.exclude[util.GetServiceNameFromFullMethod(info.FullMethod)]; ok {
			return handler(srv, stream)
		}
		err := i.validToken(stream.Context())

		if err != nil {
			return err
		}
		return handler(srv, stream)
	}
}

// valid validates the authorization.
func (i *MixAuthInterceptor) validToken(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errMissingMetadata
	}

	// The keys within metadata.MD are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	authorization := md["authorization"]
	if len(authorization) < 1 {
		return errInvalidToken
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	// Perform the token validation here.
	if token != i.token {
		return errInvalidToken
	}
	return nil
}

var (
	// to check EnsureValidToken whether implement grpc.UnaryServerInterceptor
	_ MixInterceptor = &MixAuthInterceptor{token: ""}
)
