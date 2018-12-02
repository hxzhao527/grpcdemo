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

// Package main implements a simple gRPC client that demonstrates how to use gRPC-Go libraries
// to perform unary, client streaming, server streaming and full duplex RPCs.
//
// It interacts with the route guide service whose definition can be found in routeguide/route_guide.proto.
package main

import (
	"flag"
	"log"

	client "github.com/hxzhao527/grpcdemo/internal/client/routeguide"
	rpcClient "github.com/hxzhao527/grpcdemo/pkg/client"
	"github.com/hxzhao527/grpcdemo/proto/routeguide"
	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

const (
	serverAddr = "localhost:50051"
	caFilePath = "assets/public.pem"
	authToken  = "grpcdemo"
)

var (
	ssl  = flag.Bool("ssl", false, "whether TLS enabled")
	auth = flag.Bool("auth", false, "whether oauth enabled")
)

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

	opts = append(opts, rpcClient.WithTarget(serverAddr))

	routeClient, err := client.NewClient(opts...)
	if err != nil {
		log.Fatal(err)
	}

	// Looking for a valid feature
	routeClient.PrintFeature(&routeguide.Point{Latitude: 409146138, Longitude: -746188906})

	// Feature missing.
	routeClient.PrintFeature(&routeguide.Point{Latitude: 0, Longitude: 0})

	// Looking for features between 40, -75 and 42, -73.
	routeClient.PrintFeatures(&routeguide.Rectangle{
		Lo: &routeguide.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &routeguide.Point{Latitude: 420000000, Longitude: -730000000},
	})

	// RecordRoute
	routeClient.RunRecordRoute()

	// RouteChat
	routeClient.RunRouteChat()
}
