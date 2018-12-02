package helloworld

import (
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hxzhao527/grpcdemo/pkg/client"
	"github.com/hxzhao527/grpcdemo/proto/helloworld"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type ClientOption func(client *Client)

type Client struct {
	conn   *grpc.ClientConn
	client helloworld.HelloClient

	config *client.RPCClient
}

func WithConsul(consulAddr string) client.RPCClientOption {
	return client.WithTarget("consul://" + consulAddr + "/helloworld-Hello")
}

func NewClient(opts ...client.RPCClientOption) (*Client, error) {
	c := Client{config: &client.RPCClient{}}

	for _, opt := range opts {
		opt(c.config)
	}

	conn, err := c.config.Dial()
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, client: helloworld.NewHelloClient(conn)}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) SayHello(name string) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.client.SayHello(ctx, &helloworld.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	return r.GetMessage()
}

func (c *Client) SayHelloOnce(name string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.client.SayHelloOnce(ctx, &helloworld.HelloRequest{Name: name})
	if err != nil {
		return "", err
	}
	return r.GetMessage(), nil
}

func (c *Client) TryPanic() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := c.client.TryPanic(ctx, &empty.Empty{})
	if err != nil {
		s := status.Convert(err)
		log.Printf("request get error: %v", s.Message())
		return
	}
	log.Println("request return normally")
}
