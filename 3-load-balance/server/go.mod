module github.com/xiuxiubiu/grpc-example/load-balance/server

go 1.16

require (
	github.com/xiuxiubiu/grpc-example/load-balance/proto v0.0.0-00010101000000-000000000000 // indirect
	go.etcd.io/etcd/client/v3 v3.5.0-rc.0
	google.golang.org/grpc v1.38.0
)

replace github.com/xiuxiubiu/grpc-example/load-balance/proto => ../proto
