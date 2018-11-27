package client

import (
	"context"
	"log"
	"net"
	"strconv"
	"time"

	consulApi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

const (
	defaultFreq = 30 * time.Second
)

func NewConsulResovleBuilder() resolver.Builder {
	return consulBuilder{}
}

type consulBuilder struct {
}

func (cb consulBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	cr := consulResolver{cc: cc}
	var err error
	conf := consulApi.DefaultConfig()
	conf.Address = target.Endpoint

	if cr.consulClient, err = consulApi.NewClient(conf); err != nil {
		return nil, err
	}
	cr.ctx, cr.cancel = context.WithCancel(context.Background())

	cr.rn = make(chan struct{}, 1) // cache channel
	cr.t = time.NewTicker(defaultFreq)

	go cr.watcher()

	cr.ResolveNow(resolver.ResolveNowOption{}) // important

	return cr, nil
}

func (cb consulBuilder) Scheme() string {
	return "consul"
}

type consulResolver struct {
	consulClient *consulApi.Client

	cc resolver.ClientConn

	ctx    context.Context
	cancel context.CancelFunc

	t  *time.Ticker
	rn chan struct{}
}

func (cr consulResolver) ResolveNow(option resolver.ResolveNowOption) {
	select {
	case cr.rn <- struct{}{}:
	default:
	}
}
func (cr consulResolver) Close() {
	cr.cancel()
	close(cr.rn)
	cr.t.Stop()
}

func (cr consulResolver) watcher() {
	for {
		select {
		case <-cr.ctx.Done():
			return
		case <-cr.t.C:
		case <-cr.rn:
		}
		var address []resolver.Address
		consulAgentService, err := cr.consulClient.Agent().Services()
		if err != nil {
			log.Printf("query form consul get err %s", err)
			continue
		}
		for sid, cfg := range consulAgentService {
			address = append(address, resolver.Address{Addr: net.JoinHostPort(cfg.Address, strconv.Itoa(cfg.Port)), ServerName: sid, Type: resolver.Backend})
		}
		log.Printf("found service: %#v", address)
		cr.cc.NewAddress(address)
	}
}

func init() {
	resolver.Register(NewConsulResovleBuilder())
}
