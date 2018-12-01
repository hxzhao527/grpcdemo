package server

import (
	"net"
	"strconv"
	"testing"
)
import consulApi "github.com/hashicorp/consul/api"

func TestWithConsulIntegration(t *testing.T) {
	conf := consulApi.DefaultConfig()
	conf.Address = "127.0.0.1:8500"

	consulClient, err := consulApi.NewClient(conf)
	if err != nil {
		t.Error(err)
	}
	consulAgentService, _, err := consulClient.Catalog().Service("helloworld-Hello", "", nil)
	if err != nil {
		t.Errorf("query form consul get err %s", err)
	}
	for _, svc := range consulAgentService {
		t.Logf("svc: %s, addr:%s", svc.ServiceName, net.JoinHostPort(svc.ServiceAddress, strconv.Itoa(svc.ServicePort)))
	}
}
