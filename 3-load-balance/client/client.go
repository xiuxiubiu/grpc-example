package main

import (
	"context"
	"fmt"
	"github.com/xiuxiubiu/grpc-example/load-balance/proto"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"time"
)

func main() {

	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"127.0.0.1:2379"}})
	if err != nil {
		fmt.Printf("clientv3 new err: %v", err)
		return
	}

	resolver.Register(NewBuilder(cli))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		"etcd:///test",
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "weighted_round_robin"}`),
	)

	if err != nil {
		fmt.Printf("grpc DialContext err: %v", err)
		return
	}

	bc := proto.NewBalanceClient(conn)

	for i := 0; i < 100; i++ {

		_, err := bc.Call(context.Background(), &proto.CallRequest{Num: int32(i)})
		if err != nil {
			fmt.Printf("NewBalanceClient call err: %v\n", err)
			return
		}

		fmt.Printf("call num: %d\n", i)
	}
}

type Builder struct {
	client *clientv3.Client
}

func NewBuilder(client *clientv3.Client) *Builder {
	return &Builder{
		client: client,
	}
}

func (b *Builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	r := &EtcdResolver{
		c:       b.client,
		target:  target.Endpoint,
		cc:      cc,
		address: map[string]struct{}{},
	}
	r.start()
	return r, nil
}

func (b *Builder) Scheme() string {
	return "etcd"
}

type EtcdResolver struct {
	c       *clientv3.Client
	target  string
	cc      resolver.ClientConn
	address map[string]struct{}
}

func (er *EtcdResolver) start() {

	resp, _ := er.c.Get(context.Background(), "/etcd3_naming", clientv3.WithPrefix())
	for _, kv := range resp.Kvs {
		er.address[string(kv.Value)] = struct{}{}
	}
	er.ResolveNow(resolver.ResolveNowOptions{})

	go func() {
		wc := er.c.Watch(context.Background(), "/etcd3_naming", clientv3.WithPrefix())
		for m := range wc {
			for _, ev := range m.Events {
				switch ev.Type {
				case mvccpb.PUT:
					er.address[string(ev.Kv.Value)] = struct{}{}
				case mvccpb.DELETE:
					delete(er.address, string(ev.Kv.Value))
				}
			}
			er.ResolveNow(resolver.ResolveNowOptions{})
		}
	}()
}

func (er *EtcdResolver) ResolveNow(options resolver.ResolveNowOptions) {
	addr := make([]resolver.Address, len(er.address))
	for k := range er.address {
		addr = append(addr, resolver.Address{Addr: k})
	}
	_ = er.cc.UpdateState(resolver.State{Addresses: addr})
}

func (er *EtcdResolver) Close() {

}
