package server

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

type MixInterceptor interface {
	UnaryInterceptor() grpc.UnaryServerInterceptor
	StreamInterceptor() grpc.StreamServerInterceptor
}

type MixAuthInterceptor struct {
	Token string
}

func (i *MixAuthInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return i.unaryInterceptor
}

func (i *MixAuthInterceptor) StreamInterceptor() grpc.StreamServerInterceptor {
	return i.streamInterceptor
}

// valid validates the authorization.
func (i *MixAuthInterceptor) validToken(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	// Perform the token validation here. For the sake of this example, the code
	// here forgoes any of the usual OAuth2 token validation and instead checks
	// for a token matching an arbitrary string.
	if token != i.Token {
		return false
	}
	return true
}

// unaryInterceptor ensures a valid token exists within a request's metadata. If
// the token is missing or invalid, the interceptor blocks execution of the
// handler and returns an error. Otherwise, the interceptor invokes the unary
// handler.
func (i *MixAuthInterceptor) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	// The keys within metadata.MD are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	if !i.validToken(md["authorization"]) {
		return nil, errInvalidToken
	}
	// Continue execution of handler after ensuring a valid token.
	return handler(ctx, req)
}

// streamInterceptor ensure a valid token exists within a stream's metadata.
func (i *MixAuthInterceptor) streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return errMissingMetadata
	}
	if !i.validToken(md["authorization"]) {
		return errInvalidToken
	}
	return handler(srv, ss)
}

var (
	// to check EnsureValidToken whether implement grpc.UnaryServerInterceptor
	_ MixInterceptor = &MixAuthInterceptor{Token: ""}
)
