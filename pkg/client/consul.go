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

func NewConsulResolveBuilder() resolver.Builder {
	return consulBuilder{}
}

type consulBuilder struct {
}

func (cb consulBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	cr := consulResolver{cc: cc, svcName: target.Endpoint, consulAddr: target.Authority}
	var err error
	conf := consulApi.DefaultConfig()
	conf.Address = target.Authority

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
	consulAddr   string
	consulClient *consulApi.Client
	svcName      string

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
		svcs, _, err := cr.consulClient.Catalog().Service(cr.svcName, "", nil)
		if err != nil {
			log.Printf("query svc: %s from consul: %s got error %s", cr.svcName, cr.consulAddr, err)
			return
		}
		for _, svc := range svcs {
			address = append(address, resolver.Address{Addr: net.JoinHostPort(svc.ServiceAddress, strconv.Itoa(svc.ServicePort)), ServerName: svc.ServiceName})
		}
		log.Printf("found service: %#v", address)
		cr.cc.NewAddress(address)
	}
}

func init() {
	resolver.Register(NewConsulResolveBuilder())
}
