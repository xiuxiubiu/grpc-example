module github.com/xiuxiubiu/grpc-example/simple-example/client

go 1.16

replace github.com/xiuxiubiu/grpc-example/simple-example/proto => ../proto

require (
	github.com/xiuxiubiu/grpc-example/simple-example/proto v0.0.0-00010101000000-000000000000 // indirect
	google.golang.org/grpc v1.38.0 // indirect
)
