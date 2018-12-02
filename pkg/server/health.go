// server health
// you can add a health-svc referring to https://github.com/grpc/grpc-go/blob/master/health/server.go
// it fellows [GRPC Health Checking Protocol](https://github.com/grpc/grpc/blob/master/doc/health-checking.md)
//
// but you need to find a way to update svc-status properly. timer? when call?
package server

import (
	"log"
	"sync"
	"time"

	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var healthSvcOnce = sync.Once{}

func (srv *RPCServer) EnableHealth() {
	log.Println("server enable health")
	healthSvcOnce.Do(func() {
		srv.healthSvc = health.NewServer()
		grpc_health_v1.RegisterHealthServer(srv.grpcsrv, srv.healthSvc)
	})
}

func (srv *RPCServer) initServiceStatus() {
	for name := range srv.grpcsvc {
		srv.UpdateServiceStatus(name, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	}
}

func (srv *RPCServer) checkServiceStatus() {
	for name, svc := range srv.grpcsvc {
		srv.UpdateServiceStatus(name, svc.Status())
	}
}

// checkServiceStatusInterval will bring back unhealthy-svc
func (srv *RPCServer) checkServiceStatusInterval(dur time.Duration) {
	srv.healthCheckTimer = time.NewTicker(dur)
TIMECHECK:
	for {
		select {
		case <-srv.healthCheckTimer.C:
			srv.checkServiceStatus()
		case <-srv.done:
			break TIMECHECK
		}
	}
}

func (srv *RPCServer) UpdateServiceStatus(svc string, status grpc_health_v1.HealthCheckResponse_ServingStatus) {
	if srv.healthSvc == nil {
		log.Printf("server health not enabled, updating status will be ignored")
		return
	}
	if _, ok := srv.grpcsvc[svc]; !ok {
		return
	}
	log.Printf("svc %s update status to %d \n", svc, status)
	srv.healthSvc.SetServingStatus(svc, status)
}
