// server health
// you can add a health-svc referring to https://github.com/grpc/grpc-go/blob/master/health/server.go
// it fellows [GRPC Health Checking Protocol](https://github.com/grpc/grpc/blob/master/doc/health-checking.md)
//
// but you need to find a way to update svc-status properly. timer? when call?
package server

import (
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"sync"
	"time"
)

var healthSvcOnce = sync.Once{}

func (srv *RPCServer) EnableHealth() {
	healthSvcOnce.Do(func() {
		srv.healthSvc = health.NewServer()
		grpc_health_v1.RegisterHealthServer(srv.grpcsrv, srv.healthSvc)
	})
}

func (srv *RPCServer) initServiceStatus() {
	for name, _ := range srv.grpcsvc {
		srv.UpdateSericeStatus(name, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	}
}

func (srv *RPCServer) checkServiceStatus() {
	for name, svc := range srv.grpcsvc {
		srv.UpdateSericeStatus(name, svc.Status())
	}
}

// checkServiceStatusInterval will bring back unhealthy-svc
func (srv *RPCServer) checkServiceStatusInterval(dur time.Duration) {
	srv.healthCheckTimer = time.NewTicker(dur)
TIME_CHECK:
	for {
		select {
		case <-srv.healthCheckTimer.C:
			srv.checkServiceStatus()
		case <-srv.done:
			break TIME_CHECK
		}
	}
}

func (srv *RPCServer) UpdateSericeStatus(svc string, status grpc_health_v1.HealthCheckResponse_ServingStatus) {
	if _, ok := srv.grpcsvc[svc]; !ok {
		return
	}
	log.Printf("svc %s update status to %d \n", svc, status)
	srv.healthSvc.SetServingStatus(svc, status)
}
