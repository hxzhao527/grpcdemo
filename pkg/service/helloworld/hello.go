//go:generate protoc -I ../../../proto/helloworld --go_out=plugins=grpc:../../../proto/helloworld  ../../../proto/helloworld/helloworld.proto

package helloworld

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"grpcdemo/proto/helloworld"
	"sync"

	"golang.org/x/net/context"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server is used to implement helloworld.GreeterServer.
type Server struct {
	mu    sync.Mutex
	count map[string]int
}

func NewServer() *Server {
	return &Server{count: make(map[string]int)}
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

func (s *Server) TryPanic(context.Context, *empty.Empty) (*empty.Empty, error) {
	panic("just try to panic and recovery")
}

func (s *Server) Register(rpcServer *grpc.Server) {
	helloworld.RegisterHelloServer(rpcServer, s)
}
