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

	client "github.com/hxzhao527/grpcdemo/internal/client/helloworld"
	rpcClient "github.com/hxzhao527/grpcdemo/pkg/client"

	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
)

const (
	defaultName = "世界"
	caFilePath  = "assets/public.pem"
	authToken   = "grpcdemo"
)

var (
	target        = "localhost:50051"
	name          = flag.String("name", defaultName, "name to contact the server")
	ssl           = flag.Bool("ssl", false, "whether TLS enabled")
	auth          = flag.Bool("auth", false, "whether oauth enabled")
	consulAddress = flag.String("consul", "", "consul address to register svc")
)

func showErrorDetail(err error) {
	s := status.Convert(err)
	for _, d := range s.Details() {
		switch info := d.(type) {
		case *epb.QuotaFailure:
			log.Printf("Quota failure: %v", info)
		default:
			log.Printf("Unexpected type: %v", info)
		}
	}
}

func main() {
	flag.Parse()
	var opts []rpcClient.RPCClientOption

	if *ssl {
		creds, err := credentials.NewClientTLSFromFile(caFilePath, "")
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, rpcClient.WithGrpcDialOption(grpc.WithTransportCredentials(creds)))
	} else {
		opts = append(opts, rpcClient.WithGrpcDialOption(grpc.WithInsecure()))
	}

	if *auth {
		opts = append(opts, rpcClient.WithGrpcDialOption(grpc.WithPerRPCCredentials(&rpcClient.AuthCreds{Token: authToken})))
	}

	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithMax(3),
		grpc_retry.WithCodes(codes.Internal),
	}
	opts = append(opts, rpcClient.WithGrpcDialOption(grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...))))

	if len(*consulAddress) > 0 {
		opts = append(opts, client.WithConsul(*consulAddress))
	} else {
		opts = append(opts, rpcClient.WithTarget(target))
	}

	helloClient, err := client.NewClient(opts...)
	defer helloClient.Close()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(helloClient.SayHello(*name))

	if resp, err := helloClient.SayHelloOnce(*name); err != nil {
		// how to deal err is business logic
		showErrorDetail(err)
	} else {
		log.Println(helloClient.SayHello(resp))
	}

	helloClient.TryPanic()
}
