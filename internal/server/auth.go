package server

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

// AuthFunc generate grpc_auth.AuthFunc with authToken.
func AuthFunc(authToken string) grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
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
}
