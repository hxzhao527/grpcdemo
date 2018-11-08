/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"flag"
	"google.golang.org/grpc/codes"
	"grpcdemo/pkg/client"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/empty"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc/status"

	"grpcdemo/proto/helloworld"

	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"golang.org/x/net/context"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "世界"
	caFilePath  = "assets/public.pem"
	authToken   = "grpcdemo"
)

var (
	name = flag.String("name", defaultName, "name to contact the server")
	ssl  = flag.Bool("ssl", false, "whether TLS enabled")
	auth = flag.Bool("auth", false, "whether oauth enabled")
)

func sayHello(client helloworld.HelloClient, name string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.SayHello(ctx, &helloworld.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}

func sayHelloOnce(client helloworld.HelloClient, name string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.SayHelloOnce(ctx, &helloworld.HelloRequest{Name: name})
	if err != nil {
		s := status.Convert(err)
		for _, d := range s.Details() {
			switch info := d.(type) {
			case *epb.QuotaFailure:
				log.Printf("Quota failure: %v", info)
			default:
				log.Printf("Unexpected type: %v", info)
			}
		}
		return
	}
	log.Printf("Greeting: %s", r.Message)
}

func tryPanic(client helloworld.HelloClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := client.TryPanic(ctx, &empty.Empty{})
	if err != nil {
		s := status.Convert(err)
		log.Printf("request get error: %v", s.Message())
		return
	}
	log.Println("request return normally")
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption

	if *ssl {
		creds, err := credentials.NewClientTLSFromFile(caFilePath, "")
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	if *auth {
		opts = append(opts, grpc.WithPerRPCCredentials(&client.AuthCreds{Token: authToken}))
	}

	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithMax(3),
		grpc_retry.WithCodes(codes.Internal),
	}
	opts = append(opts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)))

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	helloClient := helloworld.NewHelloClient(conn)

	sayHello(helloClient, *name)

	sayHelloOnce(helloClient, *name)
	sayHelloOnce(helloClient, *name)

	tryPanic(helloClient)
}
