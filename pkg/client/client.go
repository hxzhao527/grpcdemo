package client

import "google.golang.org/grpc"

type RPCClient struct {
	target       string
	grpcDialOpts []grpc.DialOption
}

type RPCClientOption func(client *RPCClient)

func WithTarget(target string) RPCClientOption {
	return func(c *RPCClient) {
		c.target = target
	}
}
func WithGrpcDialOption(opts ...grpc.DialOption) RPCClientOption {
	return func(c *RPCClient) {
		c.grpcDialOpts = append(c.grpcDialOpts, opts...)
	}
}

func (rc *RPCClient) Dial() (*grpc.ClientConn, error) {
	return grpc.Dial(rc.target, rc.grpcDialOpts...)
}
