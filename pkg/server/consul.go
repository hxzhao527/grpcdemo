package server

import (
	consulApi "github.com/hashicorp/consul/api"
	"grpcdemo/pkg/util"
	"log"
)

type ConsulConfig consulApi.Config

func WithConsulIntegration(consulAddress string) RPCServerOption {
	conf := consulApi.DefaultConfig()
	conf.Address = consulAddress
	return func(srv *RPCServer) {
		var err error
		srv.consulClient, err = consulApi.NewClient(conf)
		if err != nil {
			// how to deal with the error, it is up to you
			// I like panic
			panic(err)
		}
	}
}

func (srv *RPCServer) registerWithConsul() {
	ipAddrStr := srv.host
	if len(srv.host) == 0 {
		// not perfect, you should pass svc-address(not server host) though config-file ot others.
		ipAddrStr = util.GetSelfIPAddress().String()
	}
	for name := range srv.grpcsvc {
		err := srv.consulClient.Agent().ServiceRegister(&consulApi.AgentServiceRegistration{Name: name, Address: ipAddrStr, Port: srv.port})
		if err != nil {
			log.Printf("register svc %s to consul get error: %s", name, err)
		}
	}
}

func (srv *RPCServer) deRegisterWithConsul() {
	for name := range srv.grpcsvc {
		err := srv.consulClient.Agent().ServiceDeregister(name)
		if err != nil {
			log.Printf("deregister svc %s to consul get error: %s", name, err)
		}
	}
}
