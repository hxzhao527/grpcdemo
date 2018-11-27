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
	consulAgentService, err := consulClient.Agent().Services()
	if err != nil {
		t.Errorf("query form consul get err %s", err)
	}
	for sid, cfg := range consulAgentService {
		t.Logf("svc: %s, addr:%s", sid, net.JoinHostPort(cfg.Address, strconv.Itoa(cfg.Port)))
	}
}
