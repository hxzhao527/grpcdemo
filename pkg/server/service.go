package server

import (
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type ServiceRegister interface {
	Name() string
	Register(*grpc.Server)
}

type ServiceStatuser interface {
	Status() grpc_health_v1.HealthCheckResponse_ServingStatus
}

type Service interface {
	ServiceRegister
	ServiceStatuser
}

func GetServiceNameFromFullMethod(method string) string {
	ps := strings.SplitN(method, "/", 3)
	if len(ps) < 2 {
		return ""
	}
	return ps[1]
}
