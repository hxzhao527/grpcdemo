package server

import (
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
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

func WithAuthInterceptor(token string) RPCServerOption {
	authInterceptor := &MixAuthInterceptor{Token: token}
	return WithMixInterceptor(authInterceptor)
}

type MixAuthInterceptor struct {
	Token string
}

func (i *MixAuthInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return grpc_auth.UnaryServerInterceptor(i.validToken)
}

func (i *MixAuthInterceptor) StreamInterceptor() grpc.StreamServerInterceptor {
	return grpc_auth.StreamServerInterceptor(i.validToken)
}

// valid validates the authorization.
func (i *MixAuthInterceptor) validToken(ctx context.Context) (context.Context, error) {
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
	if token != i.Token {
		return ctx, errInvalidToken
	}
	return ctx, nil
}

var (
	// to check EnsureValidToken whether implement grpc.UnaryServerInterceptor
	_ MixInterceptor = &MixAuthInterceptor{Token: ""}
)
