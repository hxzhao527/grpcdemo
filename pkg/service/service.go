package service

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type ServiceRegister interface {
	Register(*grpc.Server)
}

type ServiceStatuser interface {
	Status() grpc_health_v1.HealthCheckResponse_ServingStatus
}

type Service interface {
	ServiceRegister
	ServiceStatuser
}
