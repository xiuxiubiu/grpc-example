### 简单的示例

#### 准备工作

1. 使用以下命令编译protocol插件
```
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

2. 添加protoc编译命令到PATH
```
$ export PATH="$PATH:$(go env GOPATH)/bin"
```

### proto生成go代码
```
$ cd proto
$ protoc --go_out=. --go_opt paths=source_relative \
    --go-grpc_out=. --go-grpc_opt paths=source_relative \
    hello.proto
```

### 编译运行示例

1. 编译运行服务端代码
```
$ cd server
$ go run server.go
2021/06/04 23:03:48 Serving grpc on 0.0.0.0:8181
```

2. 再另一个终端编译运行客户端代码
```
$ cd client
$ go run client.go
2021/06/04 23:04:21 Greeting: Hello gRPC
```