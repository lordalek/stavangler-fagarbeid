package main

import (
	"fmt"
	"log"
	"net"

	"github.com/lordalek/stavangler-fagarbeid/event-sourcing/pkg/order"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 8877))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	srv, err := order.NewServer()

	if err != nil {
		log.Fatalf("failed to create order server: %v", err)
	}

	order.RegisterOrderServer(grpcServer, srv)

	fmt.Printf("OrderServer registered, starting grpc server.\n")

	grpcServer.Serve(lis)
}
