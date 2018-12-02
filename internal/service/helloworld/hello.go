//go:generate protoc -I ../../../proto/helloworld --go_out=plugins=grpc:../../../proto/helloworld  ../../../proto/helloworld/helloworld.proto

package helloworld

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hxzhao527/grpcdemo/proto/helloworld"
	"golang.org/x/net/context"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// Server is used to implement helloworld.GreeterServer.
type Server struct {
	mu          sync.Mutex
	count       map[string]int
	specialFlag int32
}

func NewServer() *Server {
	return &Server{count: make(map[string]int), specialFlag: 0}
}

// SayHello implements helloworld.GreeterServer
func (s *Server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}

// SayHelloOnce
func (s *Server) SayHelloOnce(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Track the number of times the user has been greeted.
	s.count[in.Name]++
	if s.count[in.Name] > 1 {
		st := status.New(codes.ResourceExhausted, "Request limit exceeded.")
		ds, err := st.WithDetails(
			&epb.QuotaFailure{
				Violations: []*epb.QuotaFailure_Violation{{
					Subject:     fmt.Sprintf("name:%s", in.Name),
					Description: "Limit one greeting per person",
				}},
			},
		)
		if err != nil {
			return nil, st.Err()
		}
		return nil, ds.Err()
	}
	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}

// TryPanic will return panic when it is called an the first time. And later it will act normally, just return empty.
func (s *Server) TryPanic(context.Context, *empty.Empty) (*empty.Empty, error) {
	if atomic.CompareAndSwapInt32(&s.specialFlag, 0, 1) {
		panic("just try to panic and recovery")
	}
	return &empty.Empty{}, nil
}

func (s *Server) Register(rpcServer *grpc.Server) {
	helloworld.RegisterHelloServer(rpcServer, s)
}

func (s *Server) Status() grpc_health_v1.HealthCheckResponse_ServingStatus {
	return grpc_health_v1.HealthCheckResponse_SERVING
}

func (s *Server) Name() string {
	// https://www.ietf.org/rfc/rfc2782.txt
	// alpha-numerics and dashes, DNS friendly
	return "helloworld-Hello"
}
