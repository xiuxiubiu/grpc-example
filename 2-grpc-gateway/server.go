package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/xiuxiubiu/grpc-example/grpc-gateway/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
)

type server struct {
	proto.GreeterServer
}

func (s *server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloResponse, error) {
	log.Printf("Received: %s", request.Name)
	return &proto.HelloResponse{Message: "Hello " + request.Name}, nil
}

func main() {

	// gRPC server
	lis, err := net.Listen("tcp", ":8181")
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}
	
	s := grpc.NewServer()
	proto.RegisterGreeterServer(s, &server{})
	
	reflection.Register(s)
	
	log.Println("Serving grpc on 0.0.0.0:8181")
	go func() {
		log.Fatalln(s.Serve(lis))
	}()
	
	
	// HTTP server
	cli, err := grpc.DialContext(context.Background(), "0.0.0.0:8181", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	
	mux := runtime.NewServeMux()
	if err := proto.RegisterGreeterHandler(context.Background(), mux, cli); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	
	gs := http.Server{
		Addr: ":8282",
		Handler: mux,
	}
	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8282")
	log.Fatalln(gs.ListenAndServe())
}
