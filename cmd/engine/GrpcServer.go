package main

import (
	"auto-test-go/api/service"
	"auto-test-go/contract"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
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
