package server

import (
	"github.com/hxzhao527/grpcdemo/pkg/server"
	"golang.org/x/net/context"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func RecoveryFunc(rpc *server.RPCServer) server.RecoveryHandler {
	return func(ctx context.Context, method string, p interface{}) (err error) {
		rpc.UpdateServiceStatus(server.GetServiceNameFromFullMethod(method), grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		return server.DefaultRecoveryHandler(ctx, method, p)
	}
}
