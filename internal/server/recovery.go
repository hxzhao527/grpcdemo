package server

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/health/grpc_health_v1"
	"grpcdemo/internal/util"
	"grpcdemo/pkg/server"
)

func RecoveryFunc(rpc *server.RPCServer) server.RecoveryHandler {
	return func(ctx context.Context, method string, p interface{}) (err error) {
		rpc.UpdateServiceStatus(util.GetServiceNameFromFullMethod(method), grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		return server.DefaultRecoveryHandler(ctx, method, p)
	}
}
