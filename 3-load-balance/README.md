## gRPC负载均衡

### 负载均衡介绍

---

gRPC的client和server是基于长链接的，所以根据链接的负载均衡是没有太大意义的，因此gRPC的负载均衡是基于每次调用的。

负载均衡根据实现所在位置的不同，可以分为三类：

* 集中式LB（Proxy Model）
  
    在服务和客户端之间提供一个独立的LB用来转发客户端请求到服务，LB会根据某种策略，比如轮询（Round-Robin）等，将请求转发到合适的server，已达到负载均衡目的。

    由于所有的请求和响应都要经过LB，所以会有额外的性能开销和响应延迟。当请求量过大时，LB更容易成为性能瓶颈，产生单点故障。

* 客户端LB（Balancing-aware Client）
    
    此方案将负载均衡的的逻辑放到客户端，客户端会保存一份服务列表，然后根据某种策略，如轮询、随机等去选择服务链接。

    因此客户端需要开发实现不同语言的负载均衡客户端，每次升级都需要重新发布不同语言版本的客户端。某些负载均衡算法需要服务端的健康和负载等信息，还需要客户端通过rpc通信去收集，开发维护不方便。

* 独立LB服务（External Load Balancing Service）

    客户端负载方案中客户端实现了简单的负载策略算法，如轮询、随机等，复杂的算法可以由独立的LB实现，客户端通过LB获取负载配置和服务列表，当服务不可用或负载较高，LB通知客户端更新服务列表。


### gRPC负载均衡设计

----

gRPC官方的负载均衡设计使用的客户端LB方案，并在gRPC的API中提供了命名解析和负载均衡的接口，只需要扩展实现这些接口即可。

![load balancing workflow](./load-balancing.png)

基本实现原理：

1. 客户端请求名称解析器，将名称解析为一个或多个服务端访问地址
2. 客户端实例化负载均衡策略，如果解析的地址是负载均衡服务器地址，则客户端使用`grpclb`策略，否则使用服务配置的策略
3. 负载均衡策略为每一个服务地址创建一个链接
4. 当有rpc请求时，负载均衡策略决定使用哪个链接发送请求

这种模式客户端会和多个服务端建立链接，gRPC的client connection其实是维护了一组sub connection，每个sub connection都会与服务端建立链接，然后客户端直接请求服务端，没有额外的性能开销。详情参考文档[load balancing in gRPC](https://github.com/grpc/grpc/blob/master/doc/load-balancing.md)


### [etcd](https://github.com/etcd-io/etcd) 服务发现

---

根据gRPC官方提供的设计思路，只需要再结合分布式一致性组件（Consul、Zookeeper、Etcd）作为服务发现，即可实现一套完整的服务发现和负载均衡方案。

示例默认使用`etcd`版本`3`作为服务发现方案，详情参考[官方文档](https://etcd.io)

示例默认`etcd`单机模式，使用默认端口`2379`

使用v3版本的[客户端](https://github.com/etcd-io/etcd/tree/main/client/v3)

涉及到的命令如下:

```shell
# 写入键
$ etcdctl put foo bar
OK

# 读取键
$ etcdctl get foo
foo
bar

# 观察指定前缀键的变化
$ etcdctl watch  --prefix foo
# 在另外一个终端: etcdctl put foo bar
PUT
foo
bar
# 在另外一个终端: etcdctl put fooz1 barz1
PUT
fooz1
barz1

# 授予租约，TTL为10秒
$ etcdctl lease grant 10
lease 32695410dcc0ca06 granted with TTL(10s)

# 附加键 foo 到租约32695410dcc0ca06
$ etcdctl put --lease=32695410dcc0ca06 foo bar
OK

# 维持租约
$ etcdctl lease keep-alive 32695410dcc0ca06
lease 32695410dcc0ca06 keepalived with TTL(10)
lease 32695410dcc0ca06 keepalived with TTL(10)
lease 32695410dcc0ca06 keepalived with TTL(10)
...
```

代码服务发现流程如下：

1. 服务端通过api注册前缀相同的键到etcd，值为服务地址
2. 服务端给注册的信息授予租约，并使用keepalive维持租约
3. 客户端通过prefix获取所有的注册服务并保存
4. 客户端通过watch命令观察所有指定前缀键的变化
5. 服务端修改地址信息，客户端观测到put指令，修改本地服务信息
6. 服务端挂掉未续约，则相关信息自动删除，客户端观测到delete指令，从本地删除对应服务信息


### 执行示例代码

---

```shell
# 启动etcd
$ etcd
[WARNING] Deprecated '--logger=capnslog' flag is set; use '--logger=zap' flag instead
2021-06-09 17:26:44.831428 I | etcdmain: etcd Version: 3.4.16
....

# 另一个终端启动服务
$ cd server
$ go run server.go -port 8181

# 另一个终端启动服务
$ cd server
$ go run server.go -port 8282

# 另一个终端启动客户端
$ cd client
$ go run client.go
```
