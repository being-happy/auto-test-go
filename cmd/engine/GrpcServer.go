// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"auto-test-go/api/service"
	"auto-test-go/contract"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

const (
	Network = "tcp"
	Port    = ":9080"
)

type GrpcServer struct {
}

func (GrpcServer) start() {
	address := os.Getenv("GRPC_SERVICE")
	if address == "" {
		address = Port
	}

	listener, err := net.Listen(Network, address)
	if err != nil {
		log.Fatalf("[GrpcServer] net.Listen err: %v", err)
	}

	log.Println("[GrpcServer] " + address + " net.Listing...")
	grpcServer := grpc.NewServer()
	contract.RegisterStreamCaseCommandListenerServer(grpcServer, &service.CaseCommandStreamService{})
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("[GrpcServer] grpcServer.Serve err: %v", err)
	}
}
