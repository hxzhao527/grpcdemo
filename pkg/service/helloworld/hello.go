//go:generate protoc -I ../../../proto/helloworld --go_out=plugins=grpc:../../../proto/helloworld  ../../../proto/helloworld/helloworld.proto

package helloworld

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpcdemo/proto/helloworld"
	"sync"
	"sync/atomic"
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
