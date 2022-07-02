# Go 语言编程之旅(三)：RPC 应用(三) 

## 六、提供 HTTP 接口

### 1. 支持其他协议

在完成了多个 gRPC 服务后，总会有遇到一个需求，那就是提供 HTTP 接口，又或者针对一个 RPC 方法，提供多种协议的支持，但为什么会出现这种情况呢？

这基本是由于以下几种可能性，第一：心跳、监控接口等等，第二：业务场景变化，同一个 RPC 方法需要针对多种协议的业务场景提供它的服务了，但是总不可能重现一个一模一样的，因此多协议支持就非常迫切了。

另外在前文讲过，gRPC 协议本质上是 HTTP/2 协议，如果该服务想要在同个端口适配两种协议流量的话，是需要进行特殊处理的。因此在接下来的内容里，就将主要讲解接触频率最高的 HTTP/1.1 接口的支持，和与其对应所延伸出来的多种方案和思考。

接下来的将分为三个大案例进行实操讲解，虽然每个案例的代码都是相对独立的，但在知识点上是相互关联的。

### 2. 另起端口监听 HTTP

那么第一种，也就是最基础的需求：实现 gRPC（HTTP/2）和 HTTP/1.1 的支持，允许分为两个端口来进行，修改` main.go` 文件，修改其启动逻辑，并分别实现 gRPC 和 HTTP/1.1 的运行逻辑，写入如下代码：

```go
var grpcPort string
var httpPort string

func init() {
   flag.StringVar(&grpcPort, "grpc_port", "8001", "gRPC 启动端口号")
   flag.StringVar(&httpPort, "http_port", "9001", "HTTP 启动端口号")

   flag.Parse()
}
```

首先将原本的 gRPC 服务启动端口，调整为 HTTP/1.1 和 gRPC 的端口号读取，接下来实现具体的服务启动逻辑，继续写入如下代码：

```go
package main

import (
   pb "ch03/proto"
   "ch03/server"
   "flag"
   "google.golang.org/grpc"
   "google.golang.org/grpc/reflection"
   "log"
   "net"
   "net/http"
)

var grpcPort string
var httpPort string

func init() {
   flag.StringVar(&grpcPort, "grpc_port", "8001", "gRPC 启动端口号")
   flag.StringVar(&httpPort, "http_port", "9001", "HTTP 启动端口号")

   flag.Parse()
}

func main() {

   //s := grpc.NewServer()
   //pb.RegisterTagServiceServer(s, server.NewTagServer())
   //// gRPC Serve 注册反射服务
   //reflection.Register(s)
   //
   //lis, err := net.Listen("tcp", ":8001")
   //if err != nil {
   // log.Fatalf("net.Listen err: %v", err)
   //}
   //
   //err = s.Serve(lis)
   //if err != nil {
   // log.Fatalf("server.Serve err: %v", err)
   //}

   // 重新编写启动逻辑
   errs := make(chan error)
   go func() {
      err := RunHttpServer(httpPort)
      if err != nil {
         errs <- err
      }
   }()

   go func() {
      err := RunGrpcServer(grpcPort)
      if err != nil {
         errs <- err
      }
   }()

   select {
   case err := <-errs:
      log.Fatalf("Run Server err: %v", err)
   }
}

// 启动 HTTP 服务器
func RunHttpServer(port string) error {
   // 初始化一个 HTTP 请求多路复用器
   serveMux := http.NewServeMux()
   // HandleFunc() 为给定模式注册处理函数
   // 新增了一个 /ping 路由及其 Handler，可做心跳检测
   serveMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
      _, _ = w.Write([]byte(`pong`))
   })

   return http.ListenAndServe(":"+port, serveMux)
}

func RunGrpcServer(port string) error {
   s := grpc.NewServer()
   pb.RegisterTagServiceServer(s, server.NewTagServer())
   reflection.Register(s)
   lis, err := net.Listen("tcp", ":"+port)
   if err != nil {
      return err
   }

   return s.Serve(lis)
}
```

在上述代码中，一共把服务启动分为了两个方法，分别是针对 HTTP 的 `RunHttpServer `方法，其作用是初始化一个新的 HTTP 多路复用器，并新增了一个 `/ping` 路由及其 Handler，可用于做基本的心跳检测。另外 gRPC 与之前一致，保持实现了 gRPC Server 的相关逻辑，仅是重新封装为`RunGrpcServer `方法。

在启动逻辑中，先专门声明了一个 `chan` 用于接收 goroutine 的 err 信息，接下来分别在 goroutine 中调用 `RunHttpServer` 和 `RunGrpcServer` 方法，那为什么要放到 goroutine 中去调用呢，是因为实际上监听` HTTP EndPoint `和 `gRPC EndPoint `是一个阻塞的行为。

而如果 `RunHttpServer` 或 `RunGrpcServer` 方法启动或运行出现了问题，会将 err 写入` chan `中，因此只需要利用 select 对其进行检测即可。

接下来进行验证，检查输出结果是否与预期一致，命令如下：

```bash
$ grpcurl -plaintext localhost:8001 proto.TagService.GetTagList
{
  "list": [
    {
      "id": "1",
      "name": "create_tag_test",
      "state": 1
    },
    {
      "id": "2",
      "name": "Java",
      "state": 1
    }
  ],
  "pager": {
    "page": "1",
    "pageSize": "10",
    "totalRows": "2"
  }
}

$ curl http://127.0.0.1:9001/ping


StatusCode        : 200
StatusDescription : OK
Content           : pong
RawContent        : HTTP/1.1 200 OK
                    Content-Length: 4
                    Content-Type: text/plain; charset=utf-8
                    Date: Sun, 22 May 2022 13:56:31 GMT

                    pong
Forms             : {}
Headers           : {[Content-Length, 4], [Content-Type, text/plain; chars
                    et=utf-8], [Date, Sun, 22 May 2022 13:56:31 GMT]}
Images            : {}
InputFields       : {}
Links             : {}
ParsedHtml        : mshtml.HTMLDocumentClass
RawContentLength  : 4
```

第一条命令输出获取标签列表的结果集，第二条命令应当输出 pong 字符串，**至此完成在一个应用程序中分别在不同端口监听 gRPC Server 和 HTTP Server 的功能。**

### 3. 在同端口号同时监听

在上小节完成了双端口监听不同的流量的需求，但是在一些使用或部署场景下，会比较麻烦，还要兼顾两个端口，这时候就会出现希望在一个端口上兼容多种协议的需求。

#### a. 介绍与安装

在 Go 语言中，可以使用第三方开源库 `cmux` 来实现多协议支持的功能，`cmux `是根据有效负载（payload）对连接进行多路复用（也就是匹配连接的头几个字节来进行区分当前连接的类型），可以在同一 TCP Listener 上提供 gRPC、SSH、HTTPS、HTTP、Go RPC 以及几乎所有其它协议的服务，是一个相对通用的方案。

**但需要注意的是，一个连接可以是 gRPC 或 HTTP，但不能同时是两者。**也就是说，假设客户端连接用于 gRPC 或 HTTP，但不会同时在同一连接上使用两者。

接下来在项目根目录下执行如下安装命令：

```bash
$ go get -u github.com/soheilhy/cmux@v0.1.4
```

#### b. 多协议支持

修改项目根目录下的启动文件`main.go`。

```go
var (
	// 添加端口号
	port string
)

func init() {
	// 添加端口号
	flag.StringVar(&port, "port", "8003", "启动端口号")
	flag.Parse()
}
```

首先调整了启动端口号的默认端口号，而由于是在同端口，因此调整回一个端口变量，接下来编写具体的 Listener 的实现逻辑，与上小节其实本质上是一样的内容，但重新拆分了 TCP、gRPC、HTTP 的逻辑，以便于连接多路复用器的使用，修改为如下代码:

```go
func RunTCPServer(port string) (net.Listener, error) {
	return net.Listen("tcp", ":"+port)
}

// 启动 HTTP 服务器
func RunHttpServer(port string) *http.Server {
	// 初始化一个 HTTP 请求多路复用器
	serveMux := http.NewServeMux()
	// HandleFunc() 为给定模式注册处理函数
	// 新增了一个 /ping 路由及其 Handler，可做心跳检测
	serveMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`pong`))
	})

	return &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
}

func RunGrpcServer(port string) *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)
	return s
}
```

接下来修改 `main.go`中的启动逻辑。

```go
func main() {

   // 重新编写启动逻辑
   l, err := RunTCPServer(port)
   if err != nil {
      log.Fatalf("Run TCP Server err: %v", err)
   }

   // 实例化一个新的连接多路复用器。
   m := cmux.New(l)
   grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))
   httpL := m.Match(cmux.HTTP1Fast())

   grpcS := RunGrpcServer()
   httpS := RunHttpServer(port)

   go grpcS.Serve(grpcL)
   go httpS.Serve(httpL)

   err = m.Serve()
   if err != nil {
      log.Fatalf("Run Serve err: %v", err)
   }
}
```

在上述代码中，需要注意是几点，第一个点是第一个初始化的就是 TCP Listener，因为是实际上 gRPC（HTTP/2）、HTTP/1.1 在网络分层上都是基于 TCP 协议的，第二个点是 content-type 的` application/grpc` 标识，在前文中曾经分析过 gRPC 的也有特定标志位，也就是 `application/grpc`，同样的` cmux` 也是基于这个标识去进行分流。

至此，基于 `cmux `实现的同端口支持多协议已经完成了，需要重新启动服务进行验证，确保` grpcurl `工具和利用 curl 调用 HTTP/1.1 接口响应正常。

```bash
$ grpcurl -plaintext localhost:8003 proto.TagService.GetTagList
{
  "list": [
    {
      "id": "1",
      "name": "create_tag_test",
      "state": 1
    },
    {
      "id": "2",
      "name": "Java",
      "state": 1
    }
  ],
  "pager": {
    "page": "1",
    "pageSize": "10",
    "totalRows": "2"
  }
}

$ curl http://127.0.0.1:8003/ping


StatusCode        : 200
StatusDescription : OK
Content           : pong
RawContent        : HTTP/1.1 200 OK
                    Content-Length: 4
                    Content-Type: text/plain; charset=utf-8
                    Date: Sun, 22 May 2022 14:12:08 GMT

                    pong
Forms             : {}
Headers           : {[Content-Length, 4], [Content-Type, text/plain; chars
                    et=utf-8], [Date, Sun, 22 May 2022 14:12:08 GMT]}
Images            : {}
InputFields       : {}
Links             : {}
ParsedHtml        : mshtml.HTMLDocumentClass
RawContentLength  : 4
```

### 4. 同端口同方法提供双流量支持

虽然做了很多的尝试，但需求方还是想要更直接的方式，需求方就想在应用里实现一个 RPC 方法对 gRPC（HTTP/2）和 HTTP/1.1 的双流量支持，而不是单单是像前面那几个章节一样，只是单纯的另起 HTTP Handler，经过深入交流，其实他们是想用 gRPC 作为内部 API 的通讯的同时也想对外提供 RESTful，又不想搞个转换网关，写两套又太繁琐不符合….

同时也有内部的开发人员反馈说，他们平时就想在本地/开发调试时直接调用接口做一下基础验证….不想每次还要调用一下` grpcurl` 工具，看一下 list，再填写入参，相较和直接用 Postman 这类工具（具有 Web UI），那可是繁琐多了…

那有没有其它办法呢，实际上是有的，目前开源社区中的 `grpc-gateway`，就可以实现这个功能，如下图（来源自官方图）：

![image](https://raw.githubusercontent.com/tonshz/test/master/img/202205222214371.jpeg)

`grpc-gateway` 是 `protoc` 的一个插件，它能够读取 `protobuf` 的服务定义，并生成一个反向代理服务器，将 RESTful JSON API 转换为 gRPC，它主要是根据` protobuf` 的服务定义中的` google.api.http `进行生成的。

**简单来讲，grpc-gateway 能够将 RESTful 转换为 gRPC 请求，实现同一个 RPC 方法提供 gRPC 协议和 HTTP/1.1 的双流量支持的需求。**

#### a. grpc-gateway 介绍与安装

需要安装` grpc-gateway` 的 `protoc-gen-grpc-gateway `插件，安装命令如下：

```bash
$ go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.14.5
```

#### b. Proto 文件的处理

##### Proto 文件修改和编译

那么针对 grpc-gateway 的使用，需要调整项目 proto 命令下的` tag.proto `文件，修改为如下：

```protobuf
syntax = "proto3";

package proto;

import "proto/common.proto";
import "google/api/annotations.proto";

service TagService {
    rpc GetTagList (GetTagListRequest) returns (GetTagListResponse) {
        option (google.api.http) = {
            get: "/api/v1/tags"
        };
    };
}

message GetTagListRequest {
    string name = 1;
    uint32 state = 2;
}

message Tag {
    int64 id = 1;
    string name = 2;
    uint32 state = 3;
}

message GetTagListResponse {
    repeated Tag list = 1;
    Pager pager = 2;
}
```

在 proto 文件中增加了 `google/api/annotations.proto` 文件的引入，并在对应的 RPC 方法中新增了针对 HTTP 路由的注解。接下来重新编译 proto 文件，在项目根目录执行如下命令。==**需要注意，`annotations.proto`需要通过`-I`指定所在路径。**==，否则会报错找不到该文件。

```bash
# 注意 -I 后的路径地址
$ protoc -I C:\Users\zyc\protoc\include -I C:\Users\zyc\go\pkg\mod\github.com\grpc-ecosystem\grpc-gateway@v1.14.5\third_party\googleapis -I . --grpc-gateway_out=logtostderr=true:. ./proto/*.proto
```

执行完毕后将生成 `tag.pb.gw.go` 文件，也就是目前 `proto` 目录下用`.pb.go `和`.pb.gw.go `两种文件，分别对应两类功能支持。

这里使用到了一个新的 `protoc` 命令选项 `-I` 参数，它的格式为：`-IPATH, --proto_path=PATH`，作用是指定 `import` 搜索的目录（也就是 Proto 文件中的 import 命令），可指定多个，如果不指定则默认当前工作目录。

另外在实际使用场景中，还有一个较常用的选项参数，`M` 参数，例如`protoc `的命令格式为：`Mfoo/bar.proto=quux/shme`，则在生成、编译 Proto 时将所指定的包名替换为所要求的名字（如：`foo/bar.proto` 编译后为包名为 `quux/shme`），更多的选项支持可执行 `protoc --help` 命令查看帮助文档。

##### annotations.proto 是什么

刚刚在` grpc-gateway` 的 proto 文件生成中用到了 `google/api/annotations.proto` 文件，实际上它是 `googleapis` 的产物，在前面的章节有介绍过。

另外可以结合 `grpc-gateway` 的` protoc `的生成命令来看，会发现它在` grpc-gateway `的仓库下的` third_party `目录也放了个 `googleapis`，因此在引用 annotations.proto 时，用的就是 grpc-gateway 下的，这样子可以保证其兼容性和稳定性（版本可控）。

那么 `annotations.proto` 文件到底是什么，又有什么用呢，一起看看它的文件内容，如下：

```protobuf
// Copyright (c) 2015, Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

syntax = "proto3";

package google.api;

import "google/api/http.proto";
import "google/protobuf/descriptor.proto";

option go_package = "google.golang.org/genproto/googleapis/api/annotations;annotations";
option java_multiple_files = true;
option java_outer_classname = "AnnotationsProto";
option java_package = "com.google.api";
option objc_class_prefix = "GAPI";

extend google.protobuf.MethodOptions {
  // See `HttpRule`.
  HttpRule http = 72295728;
}
```

查看核心使用的 `http.proto` 文件中的一部分内容，如下：

```protobuf
message HttpRule {
  // Selects methods to which this rule applies.
  //
  // Refer to [selector][google.api.DocumentationRule.selector] for syntax details.
  string selector = 1;

  // Determines the URL pattern is matched by this rules. This pattern can be
  // used with any of the {get|put|post|delete|patch} methods. A custom method
  // can be defined using the 'custom' field.
  oneof pattern {
    // Used for listing and getting information about resources.
    string get = 2;

    // Used for updating a resource.
    string put = 3;

    // Used for creating a resource.
    string post = 4;

    // Used for deleting a resource.
    string delete = 5;

    // Used for updating a resource.
    string patch = 6;

    // The custom pattern is used for specifying an HTTP method that is not
    // included in the `pattern` field, such as HEAD, or "*" to leave the
    // HTTP method unspecified for this rule. The wild-card rule is useful
    // for services that provide content to Web (HTML) clients.
    CustomHttpPattern custom = 8;
  }

  // The name of the request field whose value is mapped to the HTTP body, or
  // `*` for mapping all fields not captured by the path pattern to the HTTP
  // body. NOTE: the referred field must not be a repeated field and must be
  // present at the top-level of request message type.
  string body = 7;

  // Optional. The name of the response field whose value is mapped to the HTTP
  // body of response. Other response fields are ignored. When
  // not set, the response message will be used as HTTP body of response.
  string response_body = 12;

  // Additional HTTP bindings for the selector. Nested bindings must
  // not contain an `additional_bindings` field themselves (that is,
  // the nesting may only be one level deep).
  repeated HttpRule additional_bindings = 11;
}message HttpRule {
  string selector = 1;
  oneof pattern {
    string get = 2;
    string put = 3;
    string post = 4;
    string delete = 5;
    string patch = 6;
    CustomHttpPattern custom = 8;
  }
  string body = 7;
  string response_body = 12;
  repeated HttpRule additional_bindings = 11;
}
```

总的来说，主要是针对的 HTTP 转换提供支持，定义了 `Protobuf `所扩展的 HTTP Option，在 Proto 文件中可用于定义 API 服务的 HTTP 的相关配置，并且可以指定每一个 RPC 方法都映射到一个或多个 HTTP REST API 方法上。

因此如果没有引入 `annotations.proto` 文件和在 Proto 文件中填写相关 HTTP Option 的话，执行生成命令，不会报错，但也不会生成任何东西。

#### c. 服务逻辑实现

接下来开始实现基于` grpc-gateway `的在同端口下同 RPC 方法提供 gRPC（HTTP/2）和 HTTP/1.1 双流量的访问支持，修改启动文件 main.go，修改为如下代码：

```go
var port string

func init() {
   flag.StringVar(&port, "port", "8004", "启动端口号")
   flag.Parse()
}
```

##### 不同协议的分流

调整案例的服务启动端口号，然后继续在` main.go `中写入如下代码：

```go
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
   return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      // gRPC 和 HTTP/1.1 的流量区分
      if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
         grpcServer.ServeHTTP(w, r)
      } else {
         otherHandler.ServeHTTP(w, r)
      }
   }), &http2.Server{})
}
```

这是一个很核心的方法，重要的分流和设置一共有两个部分，如下：

- gRPC 和 HTTP/1.1 的流量区分：
  - 对 `ProtoMajor `进行判断，该字段代表客户端请求的版本号，客户端始终使用 HTTP/1.1 或 HTTP/2。
  - Header 头 Content-Type 的确定：grpc 的标志位 `application/grpc` 的确定。
- gRPC 服务的非加密模式的设置：关注代码中的`h2c`标识，`h2c` 标识允许通过明文 TCP 运行 HTTP/2 的协议，此标识符用于 HTTP/1.1 升级标头字段以及标识 HTTP/2 over TCP，而官方标准库 `golang.org/x/net/http2/h2c` 实现了 HTTP/2 的未加密模式，直接使用即可。

在整体的方法逻辑上来讲，可以看到关键之处在于调用了 `h2c.NewHandler` 方法进行了特殊处理，`h2c.NewHandler` 会返回一个 `http.handler`，其主要是在内部逻辑是拦截了所有 `h2c` 流量，然后根据不同的请求流量类型将其劫持并重定向到相应的 `Hander` 中去处理，最终以此达到同个端口上既提供 HTTP/1.1 又提供 HTTP/2 的功能了。

##### Server 实现

完成了不同协议的流量分发和处理后，需要实现其 Server 的具体逻辑，继续在 `main.go `文件中写入如下代码：

```go
package main

import (
   ...
   "github.com/grpc-ecosystem/grpc-gateway/runtime"
   ...
)

...

func RunServer(port string) error {
   httpMux := runHttpServer()
   grpcS := runGrpcServer()
   gatewayMux := runGrpcGatewayServer()

   httpMux.Handle("/", gatewayMux)
   return http.ListenAndServe(":"+port, grpcHandlerFunc(grpcS, httpMux))
}

// 启动 HTTP 服务器
func runHttpServer() *http.ServeMux {
   // 初始化一个 HTTP 请求多路复用器
   serveMux := http.NewServeMux()
   // HandleFunc() 为给定模式注册处理函数
   // 新增了一个 /ping 路由及其 Handler，可做心跳检测
   serveMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
      _, _ = w.Write([]byte(`pong`))
   })

   return serveMux
}

func runGrpcServer() *grpc.Server {
   s := grpc.NewServer()
   pb.RegisterTagServiceServer(s, server.NewTagServer())
   reflection.Register(s)
   return s
}

func runGrpcGatewayServer() *runtime.ServeMux {
   endpoint := "0.0.0.0:" + port
   gwmux := runtime.NewServeMux()
   dopts := []grpc.DialOption{grpc.WithInsecure()}

   _ = pb.RegisterTagServiceHandlerFromEndpoint(context.Background(), gwmux,
      endpoint, dopts)
   return gwmux
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
   return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      // gRPC 和 HTTP/1.1 的流量区分
      if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
         grpcServer.ServeHTTP(w, r)
      } else {
         otherHandler.ServeHTTP(w, r)
      }
   }), &http2.Server{})
}
```

在上述代码中，与先前的案例中主要差异在于`RunServer` 方法中的`grpc-gateway `相关联的注册，核心在于调用了 `RegisterTagServiceHandlerFromEndpoint` 方法去注册 `TagServiceHandler` 事件，**其内部会自动转换并拨号到 gRPC Endpoint，并在上下文结束后关闭连接。**

另外在注册 `TagServiceHandler` 事件时，在 `grpc.DialOption` 中通过设置 `grpc.WithInsecure` 指定了 Server 为非加密模式，否则程序在运行时将会出现问题，因为 gRPC Server/Client 在启动和调用时，必须明确其是否加密。

##### 运行与验证

修改`main.go`文件，调用`RunServer()`。

```go
func main() {

   err := RunServer(port)
   if err != nil {
      log.Fatalf("Run Serve err: %v", err)
   }
}
```

重新启动服务后进行 RPC 方法的验证，如下：

```bash
$ grpcurl -plaintext localhost:8004 proto.TagService.GetTagList
{
  "list": [
    {
      "id": "1",
      "name": "create_tag_test",
      "state": 1
    },
    {
      "id": "2",
      "name": "Java",
      "state": 1
    }
  ],
  "pager": {
    "page": "1",
    "pageSize": "10",
    "totalRows": "2"
  }
}
```

![image-20220522232439089](https://raw.githubusercontent.com/tonshz/test/master/img/202205222324162.png)

正确的情况下，都会返回响应数据，分别对应心跳检测、RPC 方法的 HTTP/1.1 和 RPC 方法的 gRPC（HTTP/2）的响应。

##### 自定义错误

在完成验证后，又想到，在 gRPC 中可以通过引用 `google.golang.org/grpc/status` 内的方法可以对 `grpc-status`、`grpc-message` 以及 `grpc-details` 详细进行定制（` errcode `包就是这么做的），但是` grpc-gateway `又怎么定制呢，它作为一个代理，会怎么提示错误信息呢，如下:

```json
{
    "error": "获取标签列表失败",
    "code": 2,
    "message": "获取标签列表失败",
    "details": [
        {
            "@type": "type.googleapis.com/proto.Error",
            "code": 20010001,
            "message": "获取标签列表失败"
        }
    ]
}
```

通过结果上来看，这是真真实实的把 grpc 错误给完整转换了过来，太直接了，这显然不利于浏览器端阅读，调用的客户端会不知道以什么为标准。

实际上，`grpc-status `的含义其实对应的是 HTTP 状态码，业务错误码对应着客户端所需的消息主体，因此需要对` grpc-gateway` 的错误进行定制，继续在` main.go `文件中写入如下代码：

```go
func grpcGatewayError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler,
   w http.ResponseWriter, _ *http.Request, err error) {
   s, ok := status.FromError(err)
   if !ok {
      s = status.New(codes.Unknown, err.Error())
   }

   httpError := httpError{Code: int32(s.Code()), Message: s.Message()}
   details := s.Details()
   for _, detail := range details {
      if v, ok := detail.(*pb.Error); ok {
         httpError.Code = v.Code
         httpError.Message = v.Message
      }
   }

   resp, _ := json.Marshal(httpError)
   w.Header().Set("Content-Type", marshaler.ContentType())
   w.WriteHeader(runtime.HTTPStatusFromCode(s.Code()))
   _, _ = w.Write(resp)
}
```

在上述代码中，针对所返回的 gRPC 错误进行了两次处理，将其转换为对应的 HTTP 状态码和对应的错误主体，以确保客户端能够根据 RESTful API 的标准来进行交互。

接下来只需要将为 `grpc-gateway `所定制的错误处理方法，注册到对应的地方就可以了，如下：

```go
func runGrpcGatewayServer() *runtime.ServeMux {
   endpoint := "0.0.0.0:" + port
   // 注册定制的错误处理方法
   runtime.HTTPError = grpcGatewayError
   gwmux := runtime.NewServeMux()
   dopts := []grpc.DialOption{grpc.WithInsecure()}

   _ = pb.RegisterTagServiceHandlerFromEndpoint(context.Background(), gwmux,
      endpoint, dopts)
   return gwmux
}
```

重启服务再进行验证，查看输出结果，可以看到所输出的 HTTP 状态码和消息主体都是正确的。

![image-20220522233934159](https://raw.githubusercontent.com/tonshz/test/master/img/202205222339206.png)

#### d. 原理

虽然在上面已经讲到了 gRPC（HTTP/2）和 HTTP/1.1 的分流是通过 Header 中的 Content-Type 和 `ProtoMajor` 标识来进行分流的，但是分流后的处理逻辑又是怎么样的呢，gRPC 要进行注册（`RegisterTagServiceServer`），`grpc-gateway `也要进行注册（`RegisterTagServiceHandlerFromEndpoint`），到底有什么用呢？

解铃还须系铃人，接下来将进行探索，看看 grpc-gateway 是如何实现的，对于开发人员来讲，最常触碰到的就是`.pb.gw.go` 的注册方法，如下：

```go
// RegisterTagServiceHandlerFromEndpoint is same as RegisterTagServiceHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterTagServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
   conn, err := grpc.Dial(endpoint, opts...)
   if err != nil {
      return err
   }
   defer func() {
      if err != nil {
         if cerr := conn.Close(); cerr != nil {
            grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
         }
         return
      }
      // 根据 context 上下文控制连接的关闭时间
      go func() {
         <-ctx.Done()
         if cerr := conn.Close(); cerr != nil {
            grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
         }
      }()
   }()
   
   // RegisterTagServiceHandler 将服务 TagService 的 http 处理程序注册到“mux”。处理程序通过“conn”将请求转发到 grpc 端点。
   return RegisterTagServiceHandler(ctx, mux, conn)
}
```

实际上在调用这类` RegisterXXXXHandlerFromEndpoint `注册方法时，主要是进行 gRPC 连接的创建和管控，它在内部就已经调用了 `grpc.Dial` 对 gRPC Server 进行拨号连接，并保持住了一个 Conn 便于后续的 HTTP/1/1 调用转发。另外在关闭连接的处理上，处理的也比较的稳健，统一都是放到 defer 中进行关闭，又或者根据 context 的上下文来控制连接的关闭时间。

接下来就是，确切的内部注册方法 `RegisterTagServiceHandler`，其实际上调用的是如下方法：

```go
// RegisterTagServiceHandlerClient registers the http handlers for service TagService
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "TagServiceClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "TagServiceClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "TagServiceClient" to call the correct interceptors.
func RegisterTagServiceHandlerClient(ctx context.Context, mux *runtime.ServeMux, client TagServiceClient) error {

   mux.Handle("GET", pattern_TagService_GetTagList_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
      // 根据外部所传入的上下文进行控制，超时关闭
      ctx, cancel := context.WithCancel(req.Context())
      defer cancel()
      // 根据所传入的 MIME 类型进行默认序列化
      inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
      rctx, err := runtime.AnnotateContext(ctx, mux, req)
      if err != nil {
         runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
         return
      }
      resp, md, err := request_TagService_GetTagList_0(rctx, inboundMarshaler, client, req, pathParams)
      // NewServerMetadataContext 使用 ServerMetadata 创建新的上下文
      ctx = runtime.NewServerMetadataContext(ctx, md)
      if err != nil {
         runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
         return
      }

      forward_TagService_GetTagList_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

   })

   return nil
}
```

该方法包含了整体的 HTTP/1.1 转换到 gRPC 的前置操作，至少包含了以下四大处理：

- 注册方法：会将当前 RPC 方法所预定义的 HTTP Endpoint（根据 proto 文件所生成的`.pb.gw.go `中所包含的信息）注册到外部所传入的 HTTP 多路复用器中，也就是对应`main.go`中`runGrpcGatewayServer()`使用的的 `runtime.NewServeMux` 方法所返回的 `gwmux`。
- 超时时间：会根据外部所传入的上下文进行控制。
- 请求/响应数据：根据所传入的 MIME 类型进行默认序列化，例如：`application/jsonpb`、`application/json`。另外其在实现上是一个` Marshaler`，也就是可以通过调用 `grpc-gateway `中的 `runtime.WithMarshalerOption` 方法来注册所需要的 MIME 类型及其对应的 `Marshaler`。
- Metadata（元数据）：会将 gRPC metadata 转换为 context 中，便于使用。元数据`request_TagService_GetTagList_0()`方法中进行处理。

### 5. 其他方案

那么除了在应用中实现诸如` grpc-gateway `这种应用代理以外，还有没有其它的外部方案呢？

外部方案，也就是外部组件，普遍是代指网关，目前` Envoy `有提供` gRPC-JSON transcoder `来支持 RESTful JSON API 客户端通过 HTTP/1.1 向 `Envoy` 发送请求并代理到 gRPC 服务。另外像是` APISIX `也有提供类似的功能，其目前也进入了 Apache 开始孵化，也值得关注。

实际上可以选择的方案并不是特别多，并且都不是以单一技术方案提供，均是作为网关中的其中一个功能提供的。

--------------------------

## 参考

+ [参考文章](https://golang2.eddycjy.com/)

+ [GitHub代码](https://github.com/go-programming-tour-book)

