package temp

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	fmt.Printf("%#v\n", req)
	if c, ok := peer.FromContext(ctx); ok {
		fmt.Printf("%v\n", c.Addr.String())
	}
	fmt.Printf("%v\n", info.FullMethod)

	return handler(ctx, req)
}
