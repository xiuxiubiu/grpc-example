package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/xiuxiubiu/grpc-example/load-balance/proto"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func main() {

	port := flag.Int("port", 8080, "port")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		panic(err)
	}

	if err := register(context.Background(), "hello", "127.0.0.1", *port, []string{"127.0.0.1:2379"}, 1); err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	proto.RegisterBalanceServer(s, &server{})
	reflection.Register(s)
	panic(s.Serve(lis))
}

type server struct {
	proto.BalanceServer
}

func (s *server) Call(ctx context.Context, request *proto.CallRequest) (*proto.CallResponse, error) {
	fmt.Printf("Receive Num: %d\n", request.Num)
	return &proto.CallResponse{Message: fmt.Sprintf("Message: %d", request.Num)}, nil
}

type registerInfo struct {
	client  *clientv3.Client
	leaseId clientv3.LeaseID
}

var ri registerInfo
var prefix = "etcd3_naming"

func register(ctx context.Context, name, host string, port int, endpoints []string, ttl int) error {

	ri = registerInfo{}

	serviceValue := fmt.Sprintf("%s:%d", host, port)
	serviceKey := fmt.Sprintf("/%s/%s/%s", prefix, name, serviceValue)

	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		fmt.Printf("etcd client v3 new err: %v", err)
		return err
	}
	fmt.Printf("etcd client new: %v", endpoints)

	ri.client = cli

	lease, err := cli.Grant(ctx, int64(ttl))
	if err != nil {
		fmt.Printf("etcd client v3 grant err: %v\n", err)
		return err
	}
	fmt.Printf("etcd grant lease: %v\n", lease)

	ri.leaseId = lease.ID

	if _, err = cli.Put(ctx, serviceKey, serviceValue, clientv3.WithLease(lease.ID)); err != nil {
		fmt.Printf("etcd client v3 put err: %v\n", err)
		return err
	}
	fmt.Printf("etcd put key: %v, value: %v\n", serviceKey, serviceValue)

	go func() {
		kc, err := cli.KeepAlive(ctx, lease.ID)
		if err != nil {
			fmt.Printf("etcd client v3 keep alive lease id: %s, err: %v\n", lease.String(), err)
			return
		}
		fmt.Printf("etcd keep alive %v\n", lease.ID)

		for {
			<-kc
		}
	}()

	return nil
}
