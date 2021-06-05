## gRPC-Gateway

[gRPC-Gateway](https://github.com/grpc-ecosystem/grpc-gateway) 是一款protocol buffers编译插件，它会读取protobuf服务定义，并生成一个代理服务，将RESTful HTTP API转换成gRPC。用来同时生成RESTful和gRPC风格的API。

  
### 安装protoc插件

---

使用 [构建约束](https://golang.org/cmd/go/#hdr-Build_constraints) 来声明构建依赖的包：
```go
// +build tools

package tools

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
```
运行`go mod tidy`处理依赖关系，然后安装：
```shell
$ go install \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
```
以上命令将会在`$GOPATH/bin/`目录下生成可执行文件：
* `protoc-gen-grpc-gateway`
* `protoc-gen-openapiv2`
* `protoc-gen-go`
* `protoc-gen-go-grpc`
  

### 使用 [buf](https://buf.build) 编译proto

---

buf配置文件`buf.yaml`
```yaml
version: v1beta1
name: buf.build/yourorg/myprotos
deps:
  - buf.build/beta/googleapis
```

buf编译配置`buf.gen.yaml`
```yaml
version: v1beta1
plugins:
  - name: go
    out: .
    opt:
      - paths=source_relative
  - name: go-grpc
    out: .
    opt:
      - paths=source_relative
  - name: grpc-gateway
    out: .
    opt:
      - paths=source_relative
```

更新proto依赖
```shell
$ buf beta mod update
```

编译proto文件：
```shell
$ buf generate
```
  

### 使用protoc编译proto

---

和buf不同，使用protoc编译时需要将依赖文件放到本地，依赖文件可以在 [googleapis repository](https://github.com/googleapis/googleapis) 
```shell
$ git clone https://github.com/googleapis/googleapis
$ cd googleapis
$ pwd
/xxx/googleapis
```

protoc编译命令
```shell
$ cd proto
$ protoc --proto_path=. --proto_path=/googleapis \
 --go_out=. --go_opt paths=source_relative \
 --go-grpc_out=. --go-grpc_opt paths=source_relative \
 --grpc-gateway_out=. --grpc-gateway_opt paths=source_relative \
 hello.proto 
```
`--proto_path`指定了去哪里寻找import和要编译的proto文件，将/googleapis替换成googleapis库的绝对路径


### 请求服务

---

运行服务
```shell
$ go fun server.go
2021/06/05 20:04:01 Serving grpc on 0.0.0.0:8181
2021/06/05 20:04:01 Serving gRPC-Gateway on http://0.0.0.0:8282
```

测试Http请求
```shell
$ curl -H "Context-Type:application/json" -X POST --data '{"name": "grpc-gateway"}' 127.0.0.1:8282/v1/example/sayHello
{"message":"Hello grpc-gateway"}
```

使用 [grpcui](https://github.com/fullstorydev/grpcui) 测试gRPC请求
```shell
$ grpcui -plaintext 127.0.0.1:8181
```