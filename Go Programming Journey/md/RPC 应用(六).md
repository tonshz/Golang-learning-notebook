# Go 语言编程之旅(三)：RPC 应用(六) 

## 十一、gRPC 服务注册和发现

### 1. 服务注册和发现

在分布式体系中，一个经常被提及的词就是服务注册和发现。

服务，有可能部署在多个不同的网络环境中，有可能是一个、两个、三个或更多个，在遇到特殊情况时，还会动态扩容和缩容，变得更多或更少。当其中一个服务突然在未知时刻宕机时，这个服务就不应该再被访问到了，它就应该被下线。

![image-20220526203629842](https://raw.githubusercontent.com/tonshz/test/master/img/202205262036894.png)

如果指定 IP 地址，即指定端口的访问方式，除非原先就是访问外部的负载组件，否则这种多个服务频繁发生网络位置变更和 调整的情况，根本没有办法处理。因为在现在复杂的网络环境下，服务端并不会固定在某一个网络环境中，客户端自然也就没有办法固定地址请求了，此时客户端需要动态获取实例地址，然后再次进行请求，为了解决这些问题，服务注册和发现这一概念也就应运而生了。

服务注册和发现这一概念，本质上就是增加多种角色对服务信息进行协调处理，常见的角色如下：

+ 注册中心：承担对服务信息进行注册、协调、管理等工作。
+ 服务提供者（服务端）：保罗特定端口，并提供一个到多个的服务允许外部访问。
+ 服务消费者（调用者）：服务提供者的调用方。

假设服务注册模式为自行上报，则服务注册者在启动服务时，会将自己的服务信息（如 IP 地址、端口号、版本号等）注册到注册中心。服务消费者在进行调用时，会以约定的命名标识（如服务名）到注册中心查询，发现当前哪些具体的服务可以调用。注册中心再根据约定的负载均衡算法进行调度，最终请求到服务提供者。

另外，当服务提供者出现问题时，或是当定期的心跳检测发现服务提供者无法正确响应时，这个出现问题的服务就会被下线，并被标识为不可用。即再启动时上报注册中心进行注册，把检测到出现问题的服务下线，通过这种方式来维护服务注册和发现的基本模型。

以上就是服务注册和发现的一个基本思路，除此之外，服务注册和发现还有其他很多实现的思路和方案。

### 2. gRPC 负载均衡策略

#### a. 正确的负载均衡、

一般来说，希望 gRPC 的负载均衡是基于每个调用进行，而不是基于每个连接进行。即使所有的调用都来自同一个 gRPC 客户端，仍然希望在所有服务端中实现负载均衡，即把所有的调用平均的分配到每个服务端上，而不是固定的某一个服务端。

![image-20220526204822596](https://raw.githubusercontent.com/tonshz/test/master/img/202205262048647.png)

上图中有两类客户端，分别是基于连接的负载和基于调用的负载。

客户端 B  是基于连接的负载，它把所有的请求都分配到了服务 A #1上，进而导致服务 A #1流量过大，资源开销过大，而服务 A #2、服务 A #3和服务 A #4则完全没有流量，由此可见，这种负载均衡的方式是失败的，会带来很严重的问题。

客户端 C 是基于调用的负载，它会在每次调用时根据具体的负载均衡策略选择最优的服务进行访问。除此之外，也可根据实际情况进行定制，以保证服务间的负载均衡。

#### b. 常见的负载均衡类型

##### 客户端负载

客户端负载是指再调用时，由客户端到注册中心对服务提供者进行查询，并获取所需的服务清单，服务清单中包括各个服务的石基信息（如 IP 地址、端口号、集群命名空间等），由客户端使用特定的负载均衡策略（如轮询）再服务清单中选择一个或多个服务进行调用。

客户端负载的优点是高性能、去中心化，并且不需要借助独立的外部负载均衡组件，缺点是实现成本比较高，因为对不同语言的客户端需要实现各自对应的 SDK 及其负载均衡策略，并且可能需要针对当时获得的服务清单进行过期反补处理等。

![image-20220526213019286](https://raw.githubusercontent.com/tonshz/test/master/img/202205262130332.png)

##### 服务端负载

服务端负载，又被称为代理模式，指在服务端侧搭设独立的负载均衡器，负载均衡器再根据给定的目标名称（如服务名）找到适合调用的服务实例。因此它具备负载均衡和反向代理两项功能。

服务端负载的优点是简单、透明，客户端不需要知道背后的逻辑，只需按给定的目标名称调用、访问即可，由服务端测管理负载、均衡策略及代理；缺点是外部的负载均衡器理论上可能成为性能瓶颈，会受到负载均衡器的吞吐率的影响，并且和客户端负载相比，有可能出现更高的网络延迟。同时，必须要保持高可用，因为它是整个系统的关键节点，一旦出现问人体，影响非常大。

![image-20220526213508716](https://raw.githubusercontent.com/tonshz/test/master/img/202205262135760.png)

#### c. 官方设计思路

gRPC 官方并没有直接给出 gRPC 服务发现和负载均衡相关的具体功能代码，而是在其官方文档 load-balancing 中进行了详细介绍，说明了所期望的实现思路，并在 gRPC API 中提供了各类应用的接口，以便外部扩展。

![image-20220526214814647](https://raw.githubusercontent.com/tonshz/test/master/img/202205262148701.png)

+ 在进行 gRPC 调用时，gRPC 客户端会向名称解析器（Name Resolver）发出服务端名称（即服务名）的名称解析请求，名称解析器会将服务名解析成一个或多个 IP 地址，每个 IP 地址都会标识它是服务端地址，还是负载均衡地址，以及客户端使用的负载均衡策略（如 round_robin 和 `grpclb`）。

+ 客户端实例化负载均衡策略，如果 gRPC 客户端获取的地址是负载均衡器地址，那么客户端将使用 `grpclb` 策略，否则使用服务配置请求的负载均衡策略。如果服务配置未请求负载均衡策略，则客户端默认选择第一个可用的服务端地址。

+ 负载均衡策略会为每个服务器地址创建一个子通道。

+ 对于每一个请求都由负载均衡策略决定将其发送到哪个子通道（即哪个 gRPC 服务端）。

以上就是 gRPC 官方提供的一个基本设计思路，简单来讲，核心如下：

+ 客户端根据服务名发起请求。
+ 名称解析器解析服务名并返回。
+ 客户端根据服务端类型选择相应的策略。
+ 最后根据不同的策略进行实际调用。

### 3. 实现服务注册和发现

下面分别基于 `etcd` 和 `consul`实现 gRPC 服务注册和发现，需要安装并运行 `etcd`和`consul`。 

安装完毕后，启动 `etcd server`。在启动 `etcd`时需要指定 `ETCDCTL_API` 环境变量来设定当前 `etcd`所使用的 API 版本（2或3），若不指定，则可能会出现调用问题，导致不可用。

```bash
$ cd C:\Users\zyc\Downloads\Compressed\etcd-v3.5.4-windows-amd64
# 文档中建议API version的版本设为3
$ set ETCDCTL_API=3 # win 10 命令，linux 中为 ETCDCTL_API=3 && ./etcd
```

#### a. 基于 etcd 的实现

##### 安装 etcd sdk

安装`etcd client sdk`以便应用程序调用 `etcd`的相关 API。在项目根目录下执行下列命令。

```bash
$ go get google.golang.org/grpc@v1.26.0
$ go get github.com/coreos/etcd/clientv3@v3.3.18
```

如果在拉取过程中出现`/go-systemd`模块的相关报错，则可尝试在`go.mod`文件中添加 replace 来解决这个问题。

```bash
replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
```

##### 服务端

修改前面编写的服务端代码，增加 `etcd sdk`和服务信息的注册行为。

```go
...

type httpError struct {
   Code    int32  `json:"code,omitempty"`
   Message string `json:"message,omitempty"`
}

var port string

// 服务提供者的唯一标识
const SERVICE_NAME = "tag-service"

...

func main() {

   err := RunServer(port)
   if err != nil {
      log.Fatalf("Run Serve err: %v", err)
   }
}
...

func RunServer(port string) error {
   httpMux := runHttpServer()
   grpcS := runGrpcServer()
   gatewayMux := runGrpcGatewayServer()
   httpMux.Handle("/", gatewayMux)

   // 创建 etcd sdk 实例
   etcdClient, err := clientv3.New(clientv3.Config{
      Endpoints:   []string{"http://localhost:2379"},
      DialTimeout: time.Second * 60,
   })
   if err != nil {
      return err
   }
   defer etcdClient.Close()

   // go-learning 会失败 go-learning-test 可以
   target := fmt.Sprintf("/etcdv3://go-learning-test/grpc/%s", SERVICE_NAME)
   // 调用官方提供的 grpcproxy.Register() 进行服务信息注册
   grpcproxy.Register(etcdClient, target, ":"+port, 60)

   return http.ListenAndServe(":"+port, grpcHandlerFunc(grpcS, httpMux))
}
...
```

##### ==`/etcdv3://go-learning-test/gprc/%s`中使用`go-learning`会失败`go-learning-test`则不存在问题，原因不详==

此处新增了一个常量 `SERVICE_NAME`用于表示当前服务名是什么。`SERVICE_NAME`是服务提供者的唯一标识，也常被用作各项服务信息的协调凭证，因此非常重要。

接下里创建`etcd sdk`的实例，链接地址为`http://localhost:2379`，可根据实际情况进行调整，然后调用官方提供的`grpcproxy.Register()`方法进行服务信息注册即可。

修改后，服务端在启动时就可以进行服务信息的注册。

需要注意的是，在`grpcproxy.Register()`设置的租约时间为 60s，如果服务一直运行，那么租约会不断进行续约（定时维持）。如果该服务已经失效，那么服务将在租约到期时被删除。

##### 客户端

修改客户端代码，让其通过`etcd`获取服务提供者的列表，然后通过负载均衡算法选择要请求的服务，进行 RPC 请求。

```go
...

func main() {
   ...
   clientConn, _ := GetClientConn(ctx, "tag-service", opts)
   defer clientConn.Close()

   ...

}

func GetClientConn(ctx context.Context, serviceName string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
   config := clientv3.Config{
      Endpoints: []string{"http://localhost:2379"},
      DialTimeout: time.Second*60,
   }
   cli, err := clientv3.New(config)
   if err != nil {
      return nil, err
   }
   
   r := &naming.GRPCResolver{Client: cli}
   target := fmt.Sprintf("/etcdv3://go-learning-test/gprc/%s", serviceName)
   
   // grpc.WithInsecure() 已弃用，作用为跳过对服务器证书的验证，此时客户端和服务端会使用明文通信
   // 使用 WithTransportCredentials 和 insecure.NewCredentials() 代替
   //opts = append(opts, grp c.WithInsecure())
   // insecure.NewCredentials 返回一个禁用传输安全的凭据
   opts = append(opts, grpc.WithInsecure(),
      grpc.WithBalancer(grpc.RoundRobin(r)), grpc.WithBlock())

   /*
      grpc.DialContext() 创建到给定目标的客户端连接。
      默认情况下，它是一个非阻塞拨号（该功能不会等待建立连接，并且连接发生在后台）。
      要使其成为阻塞拨号，请使用 WithBlock() 拨号选项。
   */
   return grpc.DialContext(ctx, target, opts...)
}

...
```

客户端主要改造的是`GetClientConn`，在入参上新增了服务名。因为在进行服务内调时，只需通过服务名即可到注册中心（此处为 `etcd`）获取响应的服务信息。

接着在方法中调整请求表示位（与服务端要求一致），并通过`DialOption`设置负载均衡器（`etcd`提供的 naming）和连接建立的要求（必须达到 Ready 状态才返回）。

##### 验证

在改造完服务端和客户端后，手动进行验证。首先启动多个服务端。

```go
$ go run main.go --port 8004
2022-05-26 23:33:18.099033 I | grpcproxy: registered ":8004" with 60-second lease

$ go run main.go --port 8005
2022-05-26 23:33:24.522681 I | grpcproxy: registered ":8005" with 60-second lease

$ go run main.go --port 8006
2022-05-26 23:33:33.137958 I | grpcproxy: registered ":8006" with 60-second lease
```

接着对客户端进行多次调用，通过查看拦截器中输出的日志可以发现，三个服务端都会接受请求，即使其中一个出现问题（被终止），客户端也不会出现调用失败的情况，同时，另外一个不同端口的新服务也能被请求到。

当所有服务都出现问题时，客户端的调用会被阻塞，直到有新的服务启动并接收响应才会完成调用流程。

```bash
# server
2022-05-27 00:30:07.751905 I | grpcproxy: registered ":8004" with 60-second lease
2022-05-27 00:30:29.679176 I | access request log: method: /proto.TagService/GetTagList, begin_time: 1653582629, request: name:"Golang" 
2022-05-27 00:30:29.679176 I | test: admin, go-learning
2022-05-27 00:30:29.693694 I | token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBfa2V5IjoiMjEyMzJmMjk3YTU3YTVhNzQzODk0YTBlNGE4MDFmYzMiLCJhcHBfc2VjcmV0IjoiNjgyYjU1NGRiYmQ5NGE3NDQ0NDU5NDJlOGMyZDk3Y2YiLCJleHAiOjE2NTM1ODk4MjksImlzcyI6ImJsb2ctc2VydmljZSJ9.XjK6jL4A8ZGdizBG7L3nAX_D2QMKXCfnpzPt6Vi8sxQ
2022-05-27 00:30:29.701402 I | access response log: method: /proto.TagService/GetTagList, begin_time: 1653582629, end_time: 1653582629, response: list:<id:3 name:"Golang" state:1 > pager:<page:1 page_size:10 total_rows:1 > 

# client
2022-05-27 00:30:29.702461 I | resp: list:<id:3 name:"Golang" state:1 > pager:<page:1 page_size:10 total_rows:1 > 
```

#### b. 源码分析

实际在 gRPC 中，官方提供了两个接口（Resolver 接口和 Watcher 接口）用于对外部组件进行扩展。

```go
// Resolver creates a Watcher for a target to track its resolution changes.
//
// Deprecated: please use package resolver.
type Resolver interface {
	// Resolve creates a Watcher for target.
	Resolve(target string) (Watcher, error)
}

// Watcher watches for the updates on the specified target.
//
// Deprecated: please use package resolver.
type Watcher interface {
	// Next blocks until an update or error happens. It may return one or more
	// updates. The first call should get the full set of the results. It should
	// return an error if and only if Watcher cannot recover.
	Next() ([]*Update, error)
	// Close closes the Watcher.
	Close()
}
```

Resolver 接口的主要作用是解析目标创建的观察程序，跟踪其信息的改变。Watcher  接口的主要作用是监视制定目标上的变更。

因此客户端中调用的`github.com/coreos/etcd/clientv3/naming`依赖其实就是 `etcd SDK`对官方两个接口的具体实现。

### 4. 其他方案

目前可实现负载均衡的技术方案有许多，比如在客户端负载上，可以基于 consul 实现，而在服务端负载上则更多，如 Nginx、Kubernetes、`Istio`或` Linkerd `等外部组件实现类似的负载均衡或服务发现和注册的相关功能。

---------------------

## 参考

+ [参考文章](https://golang2.eddycjy.com/)

+ [GitHub代码](https://github.com/go-programming-tour-book)

