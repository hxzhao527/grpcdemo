package server

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// RecoveryHandlerFunc is a function that recovers from the panic `p` by returning an `error`.
// It refers to grpc_recovery.RecoveryHandlerFunc. The change is adding an argument `ctx` according to https://github.com/grpc-ecosystem/go-grpc-middleware/issues/168
// If you need caller-info, `package google.golang.org/grpc/peer` maybe helpful.
// It add a new argument `serverInfo` to pass server-info to handler.
// If you also need request-info, just change RPCCallInfo by yourself.
// If you want attach stacks of panic to error, `github.com/go-errors/errors` maybe helpful.
type RecoveryHandler func(ctx context.Context, si *RPCCallInfo, p interface{}) (err error)

func WithRecovery(handle RecoveryHandler) RPCServerOption {
	mi := &MixRecoveryInterceptor{rcFunc: handle}
	return WithMixInterceptor(mi)
}

type MixRecoveryInterceptor struct {
	rcFunc RecoveryHandler
}

func (i *MixRecoveryInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if p := recover(); p != nil {
				err = recoverFrom(ctx, ServerInfoFromGrpc(*info), p, i.rcFunc)
			}
		}()
		return handler(ctx, req)
	}
}
func (i *MixRecoveryInterceptor) StreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if p := recover(); p != nil {
				err = recoverFrom(ss.Context(), ServerInfoFromGrpc(*info), p, i.rcFunc)
			}
		}()
		err = handler(srv, ss) // is assignment necessary?
		return err
	}
}

func recoverFrom(ctx context.Context, si *RPCCallInfo, p interface{}, r RecoveryHandler) error {
	if r == nil {
		return defaultRecoveryHandler(ctx, si, p)
	}
	return r(ctx, si, p)
}

func defaultRecoveryHandler(ctx context.Context, si *RPCCallInfo, p interface{}) error {
	message := fmt.Sprintf("Call %s ", si.FullMethod)
	if c, ok := peer.FromContext(ctx); ok {
		message += fmt.Sprintf("from client %s ", c.Addr.String())
	}
	message += fmt.Sprintf("got error: %s", p)
	return status.Error(codes.Internal, message)
}
