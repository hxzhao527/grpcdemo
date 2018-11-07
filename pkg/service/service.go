package service

import "google.golang.org/grpc"

type ServiceRegister interface {
	Register(*grpc.Server)
}
