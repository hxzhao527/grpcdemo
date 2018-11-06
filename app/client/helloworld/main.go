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
	"log"
	"time"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc/status"

	"grpcdemo/proto/helloworld"

	"golang.org/x/net/context"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "世界"
	caFilePath  = "assets/public.pem"
)

var (
	name = flag.String("name", defaultName, "name to contact the server")
	ssl  = flag.Bool("ssl", false, "whether TLS enabled")
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

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := helloworld.NewHelloClient(conn)

	sayHello(client, *name)

	sayHelloOnce(client, *name)
	sayHelloOnce(client, *name)
}
