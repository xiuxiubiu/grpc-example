package main

import (
	"context"
	"github.com/xiuxiubiu/grpc-example/simple-example/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type server struct {
	proto.GreeterServer
}

func (s *server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloResponse, error) {
	log.Printf("Received: %s", request.GetName())
	return &proto.HelloResponse{Message: "Hello " + request.Name}, nil
}

func main() {
	
	lis, err := net.Listen("tcp", ":8181")
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}
	
	s := grpc.NewServer()
	proto.RegisterGreeterServer(s, &server{})
	
	reflection.Register(s)
	
	log.Println("Serving grpc on 0.0.0.0:8181")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
