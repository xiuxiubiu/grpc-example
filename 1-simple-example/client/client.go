package main

import (
	"context"
	"github.com/xiuxiubiu/grpc-example/simple-example/proto"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {

	conn, err := grpc.Dial("127.0.0.1:8181", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	
	defer func() {
		_ = conn.Close()
	}()
	
	client := proto.NewGreeterClient(conn)
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
	resp, err := client.SayHello(ctx, &proto.HelloRequest{Name: "gRPC"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	
	log.Printf("Greeting: %s", resp.GetMessage())
}
