package main

import (
	"log"
	"net"

	_ "github.com/fastcampus-backend-golang/grpc-api-example/data"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/fastcampus-backend-golang/grpc-api-example/proto"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterStockServiceServer(srv, &stockService{})

	// reflection untuk debugging
	reflection.Register(srv)

	log.Printf("server listening at %v", lis.Addr())
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
