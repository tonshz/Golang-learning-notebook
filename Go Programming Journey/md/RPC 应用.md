# Go 语言编程之旅(三)：RPC 应用

## 一、gRPC 和 Protobuf

### 1. gRPC

#### a. 什么是 RPC 

RPC 代指远程过程调用（Remote Procedure Call），它的调用包含了传输协议和编码（对象序列）协议等等，允许运行于一台计算机的程序调用另一台计算机的子程序，而开发人员无需额外地为这个交互作用编程，因此我们也常常称 RPC 调用，就像在进行本地函数调用一样方便。

#### b. 什么是 gRPC

gRPC 是一个是一个高性能、开源和通用的 RPC 框架，面向移动和基于 HTTP/2 设计。目前提供 C、Java 和 Go 语言等等版本，分别是：grpc、grpc-java、grpc-go，其中 C 版本支持 C、C++、Node.js、Python、Ruby、Objective-C、PHP 和 C# 支持。

gRPC 基于 HTTP/2 标准设计，带来诸如双向流、流控、头部压缩、单 TCP 连接上的多复用请求等特性。这些特性使得其在移动设备上表现更好，在一定的情况下更节省空间占用。

gRPC 的接口描述语言（Interface description language，缩写 IDL）使用的是 Protobuf，都是由 Google 开源的。

#### c. gRPC 调用模型

接下来看一个 gRPC 的最简调用模型。

![image](https://raw.githubusercontent.com/tonshz/test/master/img/202205212056184.jpeg)

+ 客户端（gRPC Stub）在程序中调用某方法，发起 RPC 调用。
+ 对请求信息使用 Protobuf 进行对象序列化压缩（IDL）。
+ 服务端（gRPC Server）接收到请求后，解码请求体，进行业务逻辑处理并返回。
+ 对响应结果使用 Protobuf 进行对象序列化压缩（IDL）。
+ 客户端接受到服务端响应，解码请求体。回调被调用的 A 方法，唤醒正在等待响应（阻塞）的客户端调用并返回响应结果。

### 2. Protobuf

**Protocol Buffers（Protobuf）是一种与语言、平台无关，可扩展的序列化结构化数据的数据描述语言，常常称其为 IDL，常用于通信协议，数据存储等等，相较于 JSON、XML，它更小、更快，因此也更受开发人员的青眯。**

#### a. 基本语法

```protobuf
syntax = "proto3";
package helloworld;
service Greeter {
    rpc SayHello (HelloRequest) returns (HelloReply) {}
}
message HelloRequest {
    string name = 1;
}
message HelloReply {
    string message = 1;
}
```

+ 第一行（非空的非注释行）声明使用 `proto3` 语法。如果不声明，将默认使用 `proto2` 语法。同时建议无论是用 v2 还是 v3 版本，都应当进行显式声明。而在版本上，目前主流推荐使用 v3 版本。
+ 定义名为 `Greeter` 的 RPC 服务（Service），其包含 RPC 方法 `SayHello`，入参为 `HelloRequest` 消息体（message），出参为 `HelloReply` 消息体。
+ 定义 `HelloRequest`、`HelloReply` 消息体，每一个消息体的字段包含三个属性：**类型、字段名称、字段编号**。在消息体的定义上，除类型以外均不可重复。

在编写完`.proto `文件后，一般会进行编译和生成对应语言的 proto 文件操作，这个时候 Protobuf 的编译器会根据选择的语言不同、调用的插件情况，生成相应语言的 Service Interface Code 和 Stubs。

#### b. 基本数据类型

在生成了对应语言的 proto 文件后，需要注意的是 protobuf 所生成出来的数据类型并非与原始的类型完全一致，因此需要有一个基本的了解，下面是列举了的一些常见的类型映射，如下表：

| .proto Type | C++ Type | Java Type  | Go Type | PHP Type       |
| :---------- | :------- | :--------- | :------ | :------------- |
| double      | double   | double     | float64 | float          |
| float       | float    | float      | float32 | float          |
| int32       | int32    | int        | int32   | integer        |
| int64       | int64    | long       | int64   | integer/string |
| uint32      | uint32   | int        | uint32  | integer        |
| uint64      | uint64   | long       | uint64  | integer/string |
| sint32      | int32    | int        | int32   | integer        |
| sint64      | int64    | long       | int64   | integer/string |
| fixed32     | uint32   | int        | uint32  | integer        |
| fixed64     | uint64   | long       | uint64  | integer/string |
| sfixed32    | int32    | int        | int32   | integer        |
| sfixed64    | int64    | long       | int64   | integer/string |
| bool        | bool     | boolean    | bool    | boolean        |
| string      | string   | String     | string  | string         |
| bytes       | string   | ByteString | []byte  | string         |

### 3. 思考 gRPC

#### a. gRPC 与 RESTful API 对比

| 特性       | gRPC                   | RESTful API          |
| :--------- | :--------------------- | :------------------- |
| 规范       | 必须.proto             | 可选 OpenAPI         |
| 协议       | HTTP/2                 | 任意版本的 HTTP 协议 |
| 有效载荷   | Protobuf（小、二进制） | JSON（大、易读）     |
| 浏览器支持 | 否（需要 grpc-web）    | 是                   |
| 流传输     | 客户端、服务端、双向   | 客户端、服务端       |
| 代码生成   | 是                     | OpenAPI+ 第三方工具  |

#### b. gRPC 优势

##### 性能

gRPC 使用的 IDL 是 Protobuf，Protobuf 在客户端和服务端上都能快速地进行序列化，并且序列化后的结果较小，能够有效地节省传输占用的数据大小。另外众多周知，gRPC 是基于 HTTP/2 协议进行设计的，有非常显著的优势。

另外常常会有人问，为什么是 Protobuf，为什么 gRPC 不用 JSON、XML 这类 IDL 呢，主要有如下原因：

- 在定义上更简单，更明了。
- 数据描述文件只需原来的 1/10 至 1/3。
- 解析速度是原来的 20 倍至 100 倍。
- 减少了二义性。
- 生成了更易使用的数据访问类。
- 序列化和反序列化速度快。
- 开发者本身在传输过程中并不需要过多的关注其内容。

##### 代码生成

在代码生成上，只需要一个 proto 文件就能够定义 gRPC 服务和消息体的约定，并且 gRPC 及其生态圈提供了大量的工具从 proto 文件中生成服务基类、消息体、客户端等等代码，**也就是客户端和服务端共用一个 proto 文件就可以了，保证了 IDL 的一致性且减少了重复工作。**

##### 流传输

gRPC 通过 HTTP/2 对流传输提供了大量的支持：

+ Unary RPC：一元 RPC
+ Server-side streaming RPC：服务端流式 RPC。
+ Client-side streaming RPC：客户端流式 RPC。
+ Bidirectional streaming RPC：双向流式 RPC。

##### 超时和取消

gRPC 允许客户端设置截止时间，若超出截止时间那么本次 RPC 请求将会被取消，与此同时服务端也会接收到取消动作的事件，因此客户端和服务端都可以在达到截止时间后进行取消事件的相关联动处理。

并且根据 Go 语言的上下文（context）的特性，截止时间的传递是可以一层层传递下去的，也就是可以通过一层层 gRPC 调用来进行上下文的传播截止日期和取消事件，有助于处理一些上下游的连锁问题等等场景。但是同时也会带来隐患，如果没有适当处理，第一层的上下文取消，可以把最后的调用也给取消掉，这在某些场景下可能是有问题的（需要根据实际业务场景判别）。

#### c. gRPC 缺点

##### 可读性

默认情况下 gRPC 使用 Protobuf 作为其 IDL，**Protobuf 序列化后本质上是二进制格式的数据**，并不可读，因此其可读性差，没法像 HTTP/1.1 那样直接目视调试，除非进行其它的特殊操作调整格式支持。

##### 浏览器支持

目前来讲，无法直接通过浏览器来调用 gRPC 服务，这意味着单从调试上来讲就没那么便捷了，更别提在其它的应用场景上了。

gRPC-Web 提供了一个 JavaScript 库，使浏览器客户端可以访问 gRPC 服务，但它也是有限的 gRPC 支持（对流传输的支持比较弱）。gRPC-Web 由两部分组成：一个支持浏览器的 JavaScript 客户端，以及服务器上的一个 gRPC-Web 代理。调用流程为：gRPC-Web 客户端调用代理，代理将根据 gRPC 请求转发到 gRPC 服务。

但总归是需要额外的组件进行支持的，因此对浏览器的支持是有限的。

##### 外部组件支持

gRPC 是基于 HTTP/2 设计的，HTTP/2 标准在 2015 年 5 月以 RFC 7540 正式发表，虽然已经过去了好几年，HTTP/3 也已经有了踪影，但目前为止各大外部组件对 gRPC 这类基于 HTTP/2 设计的组件支持仍然不够完美，甚至有少数暂时就完全不支持。与此同时，即使外部组件支持了，但其在社区上的相关资料也比较少，需要开发人员花费部分精力进行识别和研究，这是一个需要顾及的点。

## 二、Protobuf 的使用

### 1. 安装 Protobuf

#### a.  安装 protoc 编译器

在 gRPC 开发中，常常需要与 Protobuf 进行打交道，而在编写了`.proto `文件后，会需要到一个编译器，那就是 `protoc`，`protoc `是 Protobuf 的编译器，是用 C++ 所编写的，其主要功能是用于编译`.proto `文件。

接下来进行 `protoc `的安装。在 GitHub 上下载所需版本[Releases · protocolbuffers/protobuf (github.com)](https://github.com/protocolbuffers/protobuf/releases) ，配置好环境变量，并将 `protoc.exe`文件移至`GOPATH/bin`目录下即可。

![image-20220521222758801](https://raw.githubusercontent.com/tonshz/test/master/img/202205212227883.png)

```bash
$ protoc --version
libprotoc 3.20.1
```

#### b. protoc 插件安装

在上一步安装了 protoc 编译器，但是还是不够的，针对不同的语言，还需要不同的运行时的 protoc 插件，那么对应 Go 语言就是 protoc-gen-go 插件，接下来可以在命令行执行如下安装命令：

```bash
$ go get -u github.com/golang/protobuf/protoc-gen-go
```

### 2. 初始化 Demo 项目

在使用`go mod init`初始化目录结构后，新建`server、client、proto`目录，便于后续使用。

```lua
grpc-demo
├── go.mod
├── client
├── proto
└── server
```

### 3. 编译和生成 proto 文件

#### a. 创建 proto 文件

在项目的`proto`目录下新建`helloworld.proto`文件。

```protobuf
// 声明使用的语法
syntax = "proto3";
package helloworld;

// 定义名为 Greeter 的 RPC 服务
service Greeter {
    // 包含 RPC 方法 SayHello
    rpc SayHello (HelloRequest) returns (HelloReply) {}
}

// 入参
message HelloRequest {
    // 类型 字段名称 字段编号
    string name = 1;
}

// 出参
message HelloReply {
    string message = 1;
}
```

#### b. 生成 proto 文件

接下来在项目根目录下执行 `protoc`的相关命令来生成对应的 `pb.go`文件。

```bash
$ protoc --go_out=plugins=grpc:. ./proto/*.proto
```

+ `--go_out`：设置所生成 Go 代码的输出目录，该指令会加载`protoc-gen-go`插件达到生成 Go 代码的目的，生成的文件以`.pb.go`为文件后缀，在这里`:`（冒号）充当分隔符的作用，后跟命令所需要的参数集，在这里代表着要将所生成的 Go 代码输出到所指向`protoc`编译的当前目录。	
+ `plugins=pligin1+plugin2`：指定要加载的子插件列表，定义的 `proto`文件涉及了 RPC 服务，而默认是不会生成 RPC 代码，因此需要在 `go_out`中给出`plugins`参数传递给`protoc-gen-go`，告诉编译器，请支持 RPC（指定了内置的 grpc 插件）。

在执行这条命令后，就会生成此`proto`文件对应的`.pb.go`文件。

```bash
$ ls
helloworld.pb.go helloworld.proto
```

#### c. 生成的 .pb.go 文件

查看刚刚所生成的` helloworld.pb.go` 文件，`pb.go `文件是针对` proto `文件所生成的对应的 Go 语言代码，是实际在应用中会引用到的文件，代码如下：

```go
type HelloRequest struct {
    Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
    ...
}
func (m *HelloRequest) Reset()         { *m = HelloRequest{} }
func (m *HelloRequest) String() string { return proto.CompactTextString(m) }
func (*HelloRequest) ProtoMessage()    {}
func (*HelloRequest) Descriptor() ([]byte, []int) {
    return fileDescriptor_4d53fe9c48eadaad, []int{0}
}
func (m *HelloRequest) GetName() string {...}
```

在上述代码中，主要涉及针对 HelloRequest 类型，其包含了一组 Getters 方法，能够提供便捷的取值方式，并且处理了一些空指针取值的情况，还能够通过 Reset 方法来重置该参数。而该方法通过实现 ProtoMessage 方法，以此表示这是一个实现了 proto.Message 的接口。另外 HelloReply 类型也是类似的生成结果，因此不重复概述。

接下来看到.pb.go 文件的初始化方法，其中比较特殊的就是 fileDescriptor 的相关语句，如下：

```go
func init() {
    proto.RegisterType((*HelloRequest)(nil), "helloworld.HelloRequest")
    proto.RegisterType((*HelloReply)(nil), "helloworld.HelloReply")
}

// 注意此处
func init() { proto.RegisterFile("proto/helloworld.proto", fileDescriptor_4d53fe9c48eadaad) }

var fileDescriptor_4d53fe9c48eadaad = []byte{
    0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2b, 0x28, 0xca, 0x2f,
    ...
}
```

实际上看到的 `fileDescriptor_4d53fe9c48eadaad` 表示的是一个经过编译后的 proto 文件，是对 proto 文件的整体描述，其包含了 proto 文件名、引用（import）内容、包（package）名、选项设置、所有定义的消息体（message）、所有定义的枚举（enum）、所有定义的服务（ service）、所有定义的方法（rpc method）等等内容，可以认为就是整个 proto 文件的信息都能够取到。

同时在每一个 Message Type 中都包含了 Descriptor 方法，Descriptor 代指对一个消息体（message）定义的描述，而这一个方法则会在 fileDescriptor 中寻找属于自己 Message Field 所在的位置再进行返回，如下：

```go
func (*HelloRequest) Descriptor() ([]byte, []int) {
    return fileDescriptor_4d53fe9c48eadaad, []int{0}
}

func (*HelloReply) Descriptor() ([]byte, []int) {
    return fileDescriptor_4d53fe9c48eadaad, []int{1}
}
```

接下来再往下看可以看到` GreeterClient `接口，因为 Protobuf 是客户端和服务端可共用一份.proto 文件的，因此除了存在数据描述的信息以外，还会存在客户端和服务端的相关内部调用的接口约束和调用方式的实现，在后续多服务内部调用的时候会经常用到，如下：

```go
// GreeterClient is the client API for Greeter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GreeterClient interface {
   // 包含 RPC 方法 SayHello
   SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
}

type greeterClient struct {
   cc *grpc.ClientConn
}

func NewGreeterClient(cc *grpc.ClientConn) GreeterClient {
   return &greeterClient{cc}
}

func (c *greeterClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
   out := new(HelloReply)
   err := c.cc.Invoke(ctx, "/helloworld.Greeter/SayHello", in, out, opts...)
   if err != nil {
      return nil, err
   }
   return out, nil
}
```

### 4. 更多的类型支持

#### a. 通用类型

在 Protobuf 中一共支持 double、float、int32、int64、uint32、uint64、sint32、sint64、fixed32、fixed64、sfixed32、sfixed64、bool、string、bytes 类型，例如一开始使用的是字符串类型，当然也可以根据实际情况，修改成上述类型，例如：

```protobuf
message HelloRequest {
    bytes name = 1;
}
```

另外常常会遇到需要传递动态数组的情况，在 protobuf 中，可以使用 repeated 关键字，如果一个字段被声明为 repeated，那么该字段可以重复任意次（包括零次），重复值的顺序将保留在 protobuf 中，将重复字段视为动态大小的数组，如下：

```protobuf
message HelloRequest {
    repeated string name = 1;
}
```

#### b. 嵌套类型

嵌套类型，也就是字面意思，在 message 消息体中，又嵌套了其它的 message 消息体，一共有两种模式，如下：

第一种是将 World 消息体定义在 HelloRequest 消息体中，也就是其归属在消息体 HelloRequest 下，若要调用则需要使用 `HelloRequest.World` 的方式，外部才能引用成功。

```protobuf
message HelloRequest {
    message World {
        string name = 1;
    }
    repeated World worlds = 1;
}
```

第二种是将 World 消息体定义在外部，一般比较推荐使用这种方式，清晰、方便，如下：

````protobuf
message World {
    string name = 1;
}
message HelloRequest {
    repeated World worlds = 1;
}
````

#### c. Oneof

**如果希望消息体可以包含多个字段，但前提条件是最多同时只允许设置一个字段**，那么就可以使用 oneof 关键字来实现这个功能，如下：

```protobuf
message HelloRequest {
    oneof name {
        string nick_name = 1;
        string true_name = 2;
    }
}
```

#### d. Enum

枚举类型，限定所传入的字段值必须是预定义的值列表之一，如下：

```protobuf
enum NameType {
    NickName = 0;
    TrueName = 1;
}
message HelloRequest {
    string name = 1;
    NameType nameType = 2;
}
```

#### e. Map

map 类型，需要设置键和值的类型，map 类型不能用`repeated`关键字修饰。格式为 `map<key_type, value_type> map_field = N;`，示例如下：

```protobuf
message HelloRequest {
    map<string, string> names = 2;
}
```

## 三、gRPC 的使用

### 1. 安装

```bash
$ go get -u google.golang.org/grpc@v1.29.1
```

### 2. gRPC 的四种调用方式

在 gRPC 中，一共包含四种调用方式，分别是：

+ Unary RPC：一元 RPC。

+ Server-side streaming RPC：服务端流式 RPC。

+ Client-side streaming RPC：客户端流式 RPC。

+ Bidirectional streaming RPC：双向流式 RPC。

不同的调用方式往往代表着不同的应用场景，接下来将一同深入了解各个调用方式的实现和使用场景，在下述代码中，统一将项目下的 proto 引用名指定为 pb，并设置端口号都由外部传入，如下：

```bash
import (
    ...
    // 设置引用别名
    pb "github.com/go-programming-tour-book/grpc-demo/proto"
)
var port string
func init() {
    flag.StringVar(&port, "p", "8000", "启动端口号")
    flag.Parse()
}
```

下述的调用方法都是在 `server` 目录下的 `server.go` 和 `client` 目录的 `client.go `中完成，需要注意的该两个文件的 package 名称应该为 main（IDE 默认会创建与目录名一致的 package 名称），这样子 main 方法才能够被调用，并且在**本章中 proto 引用都会以引用别名 pb 来进行调用**。

另外在每个调用方式的 Proto 小节都会给出该类型 RPC 方法的 Proto 定义，请注意自行新增并在项目根目录执行重新编译生成语句，如下：

```bash
$ protoc --go_out=plugins=grpc:. ./proto/*.proto
```

#### a. Unary RPC：一元 RPC

一元 RPC，也就是是单次 RPC 调用，简单来讲就是客户端发起一次普通的 RPC 请求，响应，是最基础的调用类型，也是最常用的方式，大致如图：

![image](https://raw.githubusercontent.com/tonshz/test/master/img/202205212233573.png)

##### Proto

```protobuf
// 在service 内添加
rpc SayHello (HelloRequest) returns (HelloReply) {};
```

##### Server

```go
package main

import (
	pb "ch03/proto"
	"context"
	"google.golang.org/grpc"
	"net"
)

type GreeterServer struct{}

func (s *GreeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "hello.world"}, nil
}

func main() {
	// 设置监听端口
	port := "8080"
	// Server 端的抽象对象
	server := grpc.NewServer()
	// 注册到 gRPC Server
	pb.RegisterGreeterServer(server, &GreeterServer{})
	// 创建 Listen，监听 TCP 端口。
	lis, _ := net.Listen("tcp", ":"+port)
	// gRPC Server 开始 lis.Accept，直到 Stop 或 GracefulStop（优雅停止）。
	server.Serve(lis)
}
```

+ 创建 gRPC Server 对象，可以理解为它是 Server 端的抽象对象。
+ 将 GreeterServer（**其包含需要被调用的服务端接口**）注册到 gRPC Server 的内部注册中心。这样可以在接受到请求时，**通过内部的 “服务发现”，发现该服务端接口并转接进行逻辑处理。**
+ 创建 Listen，监听 TCP 端口。
+ gRPC Server 开始 lis.Accept，直到 Stop 或 GracefulStop（优雅停止）。

##### Client

```go
package main

import (
	pb "ch03/proto"
	"context"
	"google.golang.org/grpc"
	"log"
)

func main() {
	// 设置通信端口，与服务器端监听端口一致
	port := "8080"
	// 创建与给定目标（服务端）的连接句柄
	conn, _ := grpc.Dial(":"+port, grpc.WithInsecure())
	defer conn.Close()
	// 创建 Greeter 的客户端对象。
	client := pb.NewGreeterClient(conn)
	_ = SayHello(client)
}

func SayHello(client pb.GreeterClient) error {
	// 发送 RPC 请求，等待同步响应，得到回调后返回响应结果。
	resp, _ := client.SayHello(context.Background(), &pb.HelloRequest{Name: "eddycjy"})
	log.Printf("client.SayHello resp: %s", resp.Message)
	return nil
}

```

+ 创建与给定目标（服务端）的连接句柄。
+ 创建 Greeter 的客户端对象。
+ 发送 RPC 请求，等待同步响应，得到回调后返回响应结果。

先运行`server.go`后运行`client.go`（否则会报错）。

```bash
2022/05/21 22:40:26 client.SayHello resp: hello.world
```

#### b. Server-side streaming RPC：服务端流式 RPC

服务器端流式 RPC，也就是是单向流，并代指 Server 为 Stream，Client 为普通的一元 RPC 请求。

**简单来讲就是客户端发起一次普通的 RPC 请求，服务端通过流式响应多次发送数据集，客户端 Recv 接收数据集。**大致如图：

![image](https://raw.githubusercontent.com/tonshz/test/master/img/202205212248267.png)

##### Proto 

```protobuf
// stream 表示流，下述代码表示返回流 HelloReply
rpc SayList (HelloRequest) returns (stream HelloReply) {};
```

##### Server

```go
...

func (s *GreeterServer) SayList(r *pb.HelloRequest, stream pb.Greeter_SayListServer) error {
   for n := 0; n <= 6; n++ {
      _ = stream.Send(&pb.HelloReply{Message: "hello.list"})
   }
   return nil
}

...
```

在 Server 端，主要留意 `stream.Send` 方法，通过阅读源码，可得知是 protoc 在生成时，根据定义生成了各式各样符合标准的接口方法。最终再统一调度内部的 `SendMsg` 方法，该方法涉及以下过程:

+ 消息体（对象）序列化。
+ 压缩序列化后的消息体。
+ 对正在传输的消息体增加 5 个字节的 header（标志位）。
+ 判断压缩 + 序列化后的消息体总字节长度是否大于预设的 maxSendMessageSize（预设值为 `math.MaxInt32`），若超出则提示错误。
+ 写入给流的数据集。

##### Client

```go
...

func main() {
   // 设置通信端口，与服务器端监听端口一致
   port := "8080"
   // 创建与给定目标（服务端）的连接句柄
   conn, _ := grpc.Dial(":"+port, grpc.WithInsecure())
   defer conn.Close()
   // 创建 Greeter 的客户端对象。
   client := pb.NewGreeterClient(conn)
   _ = SayHello(client)

   // 添加 SayList 的调用
   _ = SayList(client, &pb.HelloRequest{Name: "eddycjy"})
}

...

func SayList(client pb.GreeterClient, r *pb.HelloRequest) error {
   stream, _ := client.SayList(context.Background(), r)
   for {
      resp, err := stream.Recv()
      if err == io.EOF {
         break
      }
      if err != nil {
         return err
      }
      log.Printf("resp: %v", resp)
   }
   return nil
}
```

在 Client 端，主要留意 `stream.Recv()` 方法，可以思考一下，什么情况下会出现 `io.EOF` ，又在什么情况下会出现错误信息呢？实际上 `stream.Recv` 方法，是对 `ClientStream.RecvMsg` 方法的封装，而 RecvMsg 方法会从流中读取完整的 gRPC 消息体，可得知：

- RecvMsg 是阻塞等待的。
- RecvMsg 当流成功/结束（调用了 Close）时，会返回 `io.EOF`。
- RecvMsg 当流出现任何错误时，流会被中止，错误信息会包含 RPC 错误码。而在 RecvMsg 中可能出现如下错误，例如：
  - io.EOF、io.ErrUnexpectedEOF
  - transport.ConnectionError
  - google.golang.org/grpc/codes（gRPC 的预定义错误码）

需要注意的是，默认的 MaxReceiveMessageSize 值为 `1024 *1024* 4`，若有特别需求，可以适当调整。

```bash
2022/05/21 23:23:38 client.SayHello resp: hello.world
2022/05/21 23:23:38 resp: message:"hello.list" 
2022/05/21 23:23:38 resp: message:"hello.list" 
2022/05/21 23:23:38 resp: message:"hello.list" 
2022/05/21 23:23:38 resp: message:"hello.list" 
2022/05/21 23:23:38 resp: message:"hello.list" 
2022/05/21 23:23:38 resp: message:"hello.list" 
2022/05/21 23:23:38 resp: message:"hello.list" 
```

#### c. Client-side streaming RPC：客户端流式 RPC

客户端流式 RPC，单向流，客户端通过流式发起**多次** RPC 请求给服务端，服务端发起**一次**响应给客户端，大致如图：

![image](https://raw.githubusercontent.com/tonshz/test/master/img/202205212303920.png)

##### Proto

```protobuf
rpc SayRecord(stream HelloRequest) returns (HelloReply) {};
```

##### Server

```go
...
func (s *GreeterServer) SayRecord(stream pb.Greeter_SayRecordServer) error {
   for {
      resp, err := stream.Recv()
      if err == io.EOF {
         return stream.SendAndClose(&pb.HelloReply{Message: "say.record"})
      }
      if err != nil {
         return err
      }
      log.Printf("resp: %v", resp)
   }
   return nil
}
...
```

可以发现在这段程序中，对每一个 Recv 都进行了处理，当发现 `io.EOF` (流关闭) 后，需要通过 `stream.SendAndClose` 方法将最终的响应结果发送给客户端，同时关闭正在另外一侧等待的 Recv。

##### Client

```go
...
func main() {
   // 设置通信端口，与服务器端监听端口一致
   port := "8080"
   // 创建与给定目标（服务端）的连接句柄
   conn, _ := grpc.Dial(":"+port, grpc.WithInsecure())
   defer conn.Close()
   // 创建 Greeter 的客户端对象。
   client := pb.NewGreeterClient(conn)
   _ = SayHello(client)

   // 添加 SayList 的调用
   _ = SayList(client, &pb.HelloRequest{Name: "eddycjy"})

   // 添加 SayRecord 的调用
   _ = SayRecord(client, &pb.HelloRequest{Name: "eddycjy"})
}

...
func SayRecord(client pb.GreeterClient, r *pb.HelloRequest) error {
   stream, _ := client.SayRecord(context.Background())
   for n := 0; n < 6; n++ {
      _ = stream.Send(r)
   }
   resp, _ := stream.CloseAndRecv()
   log.Printf("resp err: %v", resp)
   return nil
}
```

在 Server 端的 `stream.SendAndClose`，与 Client 端 `stream.CloseAndRecv` 是配套使用的方法。

```bash
# Server 端输出
2022/05/21 23:09:30 resp: name:"record" 
2022/05/21 23:09:30 resp: name:"record" 
2022/05/21 23:09:30 resp: name:"record" 
2022/05/21 23:09:30 resp: name:"record" 
2022/05/21 23:09:30 resp: name:"record" 
2022/05/21 23:09:30 resp: name:"record" 

# Client 端输出
2022/05/21 23:09:30 client.SayHello resp: hello.world
2022/05/21 23:09:30 resp: message:"hello.list" 
2022/05/21 23:09:30 resp: message:"hello.list" 
2022/05/21 23:09:30 resp: message:"hello.list" 
2022/05/21 23:09:30 resp: message:"hello.list" 
2022/05/21 23:09:30 resp: message:"hello.list" 
2022/05/21 23:09:30 resp: message:"hello.list" 
2022/05/21 23:09:30 resp: message:"hello.list" 
2022/05/21 23:09:30 resp Record: message:"say.record"  # record
```

#### d. Bidirectional streaming RPC：双向流式 RPC

双向流式 RPC，顾名思义是双向流，由客户端以流式的方式发起请求，服务端同样以流式的方式响应请求。

**首个请求一定是 Client 发起**，但具体交互方式（谁先谁后、一次发多少、响应多少、什么时候关闭）根据程序编写的方式来确定（可以结合协程）。

假设该双向流是**按顺序发送**的话，大致如图：

![image](https://raw.githubusercontent.com/tonshz/test/master/img/202205212318786.png)

##### Proto

```protobuf
rpc SayRoute(stream HelloRequest) returns (stream HelloReply) {};
```

##### Server

```go
func (s *GreeterServer) SayRoute(stream pb.Greeter_SayRouteServer) error {
    n := 0
    for {
        _ = stream.Send(&pb.HelloReply{Message: "say.route"})
        resp, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }
        n++
        log.Printf("resp: %v", resp)
    }
}
```

##### Client

```go
func SayRoute(client pb.GreeterClient, r *pb.HelloRequest) error {
    stream, _ := client.SayRoute(context.Background())
    for n := 0; n <= 6; n++ {
        _ = stream.Send(r)
        resp, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
        log.Printf("resp err: %v", resp)
    }
    _ = stream.CloseSend()
    return nil
}
```

````bash
# Server 端输出
2022/05/21 23:23:38 resp: name:"record" 
2022/05/21 23:23:38 resp: name:"record" 
2022/05/21 23:23:38 resp: name:"record" 
2022/05/21 23:23:38 resp: name:"record" 
2022/05/21 23:23:38 resp: name:"record" 
2022/05/21 23:23:38 resp: name:"record" 
2022/05/21 23:23:38 resp: name:"route"  # route 输出
2022/05/21 23:23:38 resp: name:"route" 
2022/05/21 23:23:38 resp: name:"route" 
2022/05/21 23:23:38 resp: name:"route" 
2022/05/21 23:23:38 resp: name:"route" 
2022/05/21 23:23:38 resp: name:"route" 
2022/05/21 23:23:38 resp: name:"route" 

# Client 端输出
2022/05/21 23:23:38 client.SayHello resp: hello.world
2022/05/21 23:23:38 resp: message:"hello.list" 
2022/05/21 23:23:38 resp: message:"hello.list" 
2022/05/21 23:23:38 resp: message:"hello.list" 
2022/05/21 23:23:38 resp: message:"hello.list" 
2022/05/21 23:23:38 resp: message:"hello.list" 
2022/05/21 23:23:38 resp: message:"hello.list" 
2022/05/21 23:23:38 resp: message:"hello.list" 
2022/05/21 23:23:38 resp Record: message:"say.record" 
2022/05/21 23:23:38 resp Route: message:"say.route"  # route 输出
2022/05/21 23:23:38 resp Route: message:"say.route" 
2022/05/21 23:23:38 resp Route: message:"say.route" 
2022/05/21 23:23:38 resp Route: message:"say.route" 
2022/05/21 23:23:38 resp Route: message:"say.route" 
2022/05/21 23:23:38 resp Route: message:"say.route" 
2022/05/21 23:23:38 resp Route: message:"say.route" 
````

### 3. Unary 和 Streaming RPC

#### a. 为什么不用 Unary RPC

StreamingRPC 为什么要存在呢，是 Unary RPC 有什么问题吗，通过模拟业务场景，可得知在使用 Unary RPC 时，有如下问题：

- 在一些业务场景下，数据包过大，可能会造成瞬时压力。
- 接收数据包时，需要所有数据包都接受成功且正确后，才能够回调响应，进行业务处理（无法客户端边发送，服务端边处理）。

#### b. 为什么使用 Streaming RPC

- 持续且大数据包场景。
- 实时交互场景。

#### c. 思考模拟场景

每天早上 6 点，都有一批百万级别的数据集要同从 A 同步到 B，在同步的时候，会做一系列操作（归档、数据分析、画像、日志等），这一次性涉及的数据量确实大。

在同步完成后，也有人马上会去查阅数据，为了新的一天筹备。也符合实时性。在仅允许使用 Unary 或 StreamingRPC 的情况下，两者相较下，这个场景下更适合使用 Streaming RPC。

### 4. 客户端与服务端如何交互

刚刚对 gRPC 的四种调用方式进行了探讨，但光会用还是不够的，知其然知其所然很重要，因此需要对 gRPC 的整体调用流转有一个基本印象，那么最简单的方式就是对 Client 端调用 Server 端进行抓包去剖析，看看整个过程中它都做了些什么事。

通过 WireShark 对回环地址进行抓包，选择`Adapter for loopback traffic capture`进行抓包。

![image-20220522001906944](https://raw.githubusercontent.com/tonshz/test/master/img/202205220019045.png)

查看抓包情况如下：

![image-20220521234901590](https://raw.githubusercontent.com/tonshz/test/master/img/202205212349755.png)

略加整理发现共有十二个行为，从上到下分别是 Magic、SETTINGS、HEADERS、DATA、SETTINGS、WINDOW_UPDATE、PING、HEADERS、DATA、HEADERS、WINDOW_UPDATE、PING 是比较重要的。

接下来将针对每个行为进行分析，下图中 2388 端口为客户端监听端口，8080 为服务器端监听端口。

#### a. 行为分析

##### Magic

![image-20220521234925260](https://raw.githubusercontent.com/tonshz/test/master/img/202205212349334.png)

Magic 帧的主要作用是建立 HTTP/2 请求的前言。在 HTTP/2 中，要求两端都要发送一个连接前言，作为对所使用协议的最终确认，并确定 HTTP/2 连接的初始设置，客户端和服务端各自发送不同的连接前言。

而上图中的 Magic 帧是客户端的前言之一，内容为 `PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n`，以确定启用 HTTP/2 连接。

##### SETTINGS

![image-20220521235203713](https://raw.githubusercontent.com/tonshz/test/master/img/202205212352789.png "客户端连接前言")

![image-20220521235511116](https://raw.githubusercontent.com/tonshz/test/master/img/202205212355195.png "服务器端连接前言")

SETTINGS 帧的主要作用是设置这一个连接的参数，作用域是整个连接而并非单一的流。

而上图的 SETTINGS 帧都是空 SETTINGS 帧，图一是客户端连接的前言（Magic 和 SETTINGS 帧分别组成连接前言）。图二是服务端的。另外从图中可以看到多个 SETTINGS 帧，这是为什么呢？是因为发送完连接前言后，客户端和服务端还需要有一步互动确认的动作。对应的就是带有 ACK 标识 SETTINGS 帧。

##### HEADERS

![image-20220521235801106](https://raw.githubusercontent.com/tonshz/test/master/img/202205212358227.png)

HEADERS 帧的主要作用是存储和传播 HTTP 的标头信息。可以看到 HEADERS 里有一些眼熟的信息，分别如下：

- method：POST
- scheme：http
- path：/proto.SearchService/SayList
- authority：:8080
- content-type：application/grpc
- user-agent：grpc-go/1.29.1

这些东西非常眼熟，其实都是 gRPC 的基础属性，实际上远远不止这些，只是设置了多少展示多少。例如像平时常见的 `grpc-timeout`、`grpc-encoding` 也是在这里设置的。

##### DATA

![image-20220522000031361](https://raw.githubusercontent.com/tonshz/test/master/img/202205220000458.png)

DATA 帧的主要作用是装填主体信息，是数据帧。而在上图中，可以很明显看到请求参数 gRPC 存储在里面。只需要了解到这一点就可以了。

##### HEADERS, DATA, HEADERS

![image-20220522001400115](https://raw.githubusercontent.com/tonshz/test/master/img/202205220014151.png)

![image-20220522000714274](https://raw.githubusercontent.com/tonshz/test/master/img/202205220007409.png "HEADERS 帧")

在上图中 HEADERS 帧比较简单，体现了 HTTP 响应状态和响应的内容格式。

![image-20220522001330483](https://raw.githubusercontent.com/tonshz/test/master/img/202205220013600.png "DATA 帧")

在上图中 DATA 帧主要承载了响应结果的数据集，图中的 gRPC Server 就是我们 RPC 方法的响应结果。

![image-20220522001535967](https://raw.githubusercontent.com/tonshz/test/master/img/202205220015102.png "HEADERS 帧")

在上图中 HEADERS 帧主要承载了 gRPC 的状态信息，对应图中的 `grpc-status` 和 `grpc-message` 就是本次 gRPC 调用状态的结果。

#### b. 其他步骤

##### WINDOW_UPDATE

主要作用是管理和流的窗口控制。通常情况下打开一个连接后，服务器和客户端会立即交换 SETTINGS 帧来确定流控制窗口的大小。默认情况下，该大小设置为约 65 KB，但可通过发出一个 WINDOW_UPDATE 帧为流控制设置不同的大小。

![image-20220522001709762](https://raw.githubusercontent.com/tonshz/test/master/img/202205220017881.png)

##### PING/PONG

主要作用是判断当前连接是否仍然可用，也常用于计算往返时间。其实也就是 PING/PONG。

### 5. 小结

在本章节中，对于 gRPC 的基本使用和交互原理进行了一个简单剖析，总结如下：

- gRPC 一共支持四种调用方式，分别是：
  - Unary RPC：一元 RPC。
  - Server-side streaming RPC：服务端流式 RPC。
  - Client-side streaming RPC：客户端流式 RPC。
  - Bidirectional streaming RPC：双向流式 RPC。
- gRPC 在建立连接之前，客户端/服务端都会发送连接前言（Magic+SETTINGS），确立协议和配置项。
- gRPC 在传输数据时，是会涉及滑动窗口（WINDOW_UPDATE）等流控策略的。
- 传播 gRPC 附加信息时，是基于 HEADERS 帧进行传播和设置；而具体的请求/响应数据是存储的 DATA 帧中的。
- gRPC 请求/响应结果会分为 HTTP 和 gRPC 状态响应（grpc-status、grpc-message）两种类型。
- 客户端发起 PING，服务端就会回应 PONG，反之亦可。

## 四、运行一个 gRPC 服务

在了解了 `gRPC` 和` Protobuf `的具体使用和情况后，将结合常见的应用场景，完成一个 `gRPC` 服务。而为了防止重复用工，这一个` gRPC` 服务将会直接通过 HTTP 调用上一章节的博客后端服务，以此来获得标签列表的业务数据，只需要把主要的精力集中在 `gRPC `服务相关联的知识上就可以了，同时后续的数个章节知识点的开展都会围绕着这个服务来进行。

### 1. 初始化项目

项目最终的目录结构如下：

```lua
tag-service
├── main.go
├── go.mod
├── go.sum
├── pkg
├── internal
├── proto
├── server
└── third_party
```

完成项目基础目录的创建后，在项目根目录执行`gRPC` 的安装命令。

```bash
$ go get -u google.golang.org/grpc@v1.29.1
```

### 2. 编译和生成 proto 文件

在 `GoLand`安装插件`Protobuf`并为插件配置路径为`proto`目录所在路径。

![image-20220522165047275](https://raw.githubusercontent.com/tonshz/test/master/img/202205221650416.png)

在正式的开始编写服务前，需要先编写对应的 RPC 方法所需的 proto 文件，这是日常要先做的事情之一，因此接下来开始进行公共 proto 的编写，在项目的 proto 目录下新建 `common.proto `文件，写入如下代码：

```protobuf
syntax = "proto3";

package proto;

message Pager{
    int64 page = 1;
    int64 page_size = 2;
    int64 total_rows = 3;
}
```

接着再编写获取标签列表的 RPC 方法，继续新建 `tag.proto `文件，写入如下代码：

```go
syntax = "proto3";

package proto;

import "proto/common.proto"; // 引入公共文件

service TagService {
    rpc GetTagList (GetTagListRequest) returns (GetTagListReply) {};
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

message GetTagListReply {
    repeated Tag list = 1;
    Pager pager = 2;
}
```

在上述` proto` 代码中，引入了公共文件 `common.proto`，并依据先前博客后端服务一致的数据结构定义了 RPC 方法，完成后就可以编译和生成 proto 文件，在项目根目录下执行如下命令：

```bash
$ protoc --go_out=plugins=grpc:. ./proto/*.proto
```

需要注意的一点是，在` tag.proto` 文件中 import 了 `common.proto`，因此在执行 `protoc` 命令生成时，如果只执行命令 `protoc --go_out=plugins=grpc:. ./proto/tag.proto` 是会存在问题的。

因此建议若所需生成的 proto 文件和所依赖的 proto 文件都在同一目录下，可以直接执行 `./proto/*.proto` 命令来解决，又或是指定所有含关联的 proto 引用 `./proto/common.proto ./proto/tag.proto` ，这样子就可以成功生成`.pb.go `文件，并且避免了很多的编译麻烦。

但若实在是存在多层级目录的情况，可以利用 `protoc` 命令的 `-I` 和 `M` 指令来进行特定处理。

### 3. 编写 gRPC 方法

#### a. 获取博客 API 的数据

由于数据源是第二章节的博客后端，因此需要编写一个最简单的 API SDK 去进行调用，在项目的 `pkg` 目录新建 `bapi` 目录，并创建 `api.go` 文件，写入如下代码：

```go
package bapi

import (
   "context"
   "encoding/json"
   "fmt"
   "io/ioutil"
   "net/http"
)

const (
   APP_KEY    = "admin"
   APP_SECRET = "go-learning"
)

type AccessToken struct {
   Token string `json:"token"`
}

type API struct {
   URL string
}

func NewAPI(url string) *API {
   return &API{URL: url}
}

// 获取所有 API 请求都需要带上的 token
func (a *API) getAccessToken(ctx context.Context) (string, error) {
   body, err := a.httpGet(ctx, fmt.Sprintf("%s?app_key=%s&app_secret=%s", "auth", APP_KEY, APP_SECRET))
   if err != nil {
      return "", err
   }

   var accessToken AccessToken
   _ = json.Unmarshal(body, &accessToken)
   return accessToken.Token, nil
}

// 统一的 HTTP GET 请求方法
func (a *API) httpGet(ctx context.Context, path string) ([]byte, error) {
   resp, err := http.Get(fmt.Sprintf("%s/%s", a.URL, path))
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   body, _ := ioutil.ReadAll(resp.Body)
   return body, nil
}

// 具体的获取标签列表的方法实现
func (a *API) GetTagList(ctx context.Context, name string) ([]byte, error) {
   // 获取AccessToken
   token, err := a.getAccessToken(ctx)
   if err != nil {
      return nil, err
   }

   body, err := a.httpGet(ctx, fmt.Sprintf("%s?token=%s&name=%s", "api/v1/tags", token, name))
   if err != nil {
      return nil, err
   }

   return body, nil
}
```

首先编写了两个主要方法，分别是 API SDK 统一的 HTTP GET 的请求方法，以及所有 API 请求都需要带上的 `AccessToken` 的获取，接下来就是具体的获取标签列表的方法编写。上述代码主要是实现从第二章的博客后端中获取 `AccessToken` 和完成各类数据源的接口编写，并不是本章节的重点，因此只进行了简单实现，若有兴趣可以进一步的实现 `AccessToken` 的缓存和刷新，以及多 HTTP Method 的接口调用等等。

#### b. 编写 gRPC Server

在完成了 API SDK 的编写后，在项目`server`目录下创建`tag.go`文件，针对获取标签列表的接口逻辑进行编写。

```go
package server

import (
	"ch03/pkg/bapi"
	pb "ch03/proto"
	"context"
	"encoding/json"
)

type TagServer struct {
}

func NewTagServer() *TagServer {
	return &TagServer{}
}

func (t *TagServer) GetTagList(ctx context.Context, r *pb.GetTagListRequest) (*pb.GetTagListResponse, error) {
	// localhost：不通过网卡传输，不受网络防火墙和网卡相关的限制。
	// 127.0.0.1：通过网卡传输，依赖网卡，并受到网卡和防火墙相关的限制。
	api := bapi.NewAPI("http://127.0.0.1:8000")
	body, err := api.GetTagList(ctx, r.GetName())
	if err != nil {
		return nil, err
	}

	tagList := pb.GetTagListResponse{}
	err = json.Unmarshal(body, &tagList)
	if err != nil {
		return nil, err
	}

	return &tagList, nil
}
```

在上述代码中，主要是指定了博客后端的服务地址（`http://127.0.0.1:8000`），然后调用 `GetTagList` 方法的 API，通过 HTTP 调用到第二章节所编写的博客后端服务获取标签列表数据，然后利用 `json.Unmarshal` 的特性，将其直接转换，并返回。

### 4. 编写启动文件

在项目的根目录下创建`main.go`文件。

```go
package main

import (
   pb "ch03/proto"
   "ch03/server"
   "google.golang.org/grpc"
   "log"
   "net"
)

func main() {
   s := grpc.NewServer()
   pb.RegisterTagServiceServer(s, server.NewTagServer())

   lis, err := net.Listen("tcp", ":8080")
   if err != nil {
      log.Fatalf("net.Listen err: %v", err)
   }

   err = s.Serve(lis)
   if err != nil {
      log.Fatalf("server.Serve err: %v", err)
   }
}
```

至此，一个简单的标签服务就完成了，它将承担整个篇章的研讨功能。接下来在项目根目录下执行 `go run main.go` 命令，启动这个服务，检查是否一切是否正常。

### 5. 调试 gRPC 接口

在服务启动后，除了要验证服务是否正常运行，还要调试或验证 RPC 方法是否运行正常，而 gRPC 是基于 HTTP/2 协议的，因此不像普通的 HTTP/1.1 接口可以直接通过 postman 或普通的 curl 进行调用。但目前开源社区也有一些方案，例如像` grpcurl`，`grpcurl` 是一个命令行工具，可让用户与 gRPC 服务器进行交互，安装命令如下：

```go
$ go get github.com/fullstorydev/grpcurl
$ go get github.com/fullstorydev/grpcurl/cmd/grpcurl
$ go get google.golang.org/grpc/credentials/oauth@v1.44.0
```

但使用该工具的前提是 gRPC Server 已经注册了反射服务，因此需要修改上述服务的启动文件，如下：

```go
package main

import (
   pb "ch03/proto"
   "ch03/server"
   "google.golang.org/grpc"
   "google.golang.org/grpc/reflection"
   "log"
   "net"
)

func main() {
   s := grpc.NewServer()
   pb.RegisterTagServiceServer(s, server.NewTagServer())
   // gRPC Serve 注册反射服务
   reflection.Register(s)

   lis, err := net.Listen("tcp", ":8080")
   if err != nil {
      log.Fatalf("net.Listen err: %v", err)
   }

   err = s.Serve(lis)
   if err != nil {
      log.Fatalf("server.Serve err: %v", err)
   }
}
```

`reflection `包是 gRPC 官方所提供的反射服务，在启动文件新增了 `reflection.Register` 方法的调用后，需要重新启动服务，反射服务才可用。

接下来就可以借助 `grpcurl `工具进行调试了，可以首先执行下述 list 命令：

```bash
$ grpcurl -plaintext localhost:8001 list
grpc.reflection.v1alpha.ServerReflection # 注册的反射方法
proto.TagService # 自定义的 RPC Service 方法

$ grpcurl -plaintext localhost:8001 list proto.TagService # 查看子类的 RPC 信息
proto.TagService.GetTagList

$ grpcurl localhost:8001 list # 不指定 -plaintext 会报错
Failed to dial target host "localhost:8001": read tcp [::1]:3307->[::1]:800
1: wsarecv: An existing connection was forcibly closed by the remote host.
```

一共指定了三个选项：

+ plaintext：`grpcurl`工具默认使用 TLS 认证（可通过 `-cert`和`-key`参数设置公钥和密钥），但由于服务是非 TLS 认证的，**因此需要通过指定这个选项来忽略 TLS 认证。**
+ localhost:8001：指定运行服务的 HOST
+ list：指定所执行的命令，list 子命令可以获取该服务的 RPC 方法列表信息。例如上述的输出结果，一共有两个方法，一个是注册的反射方法，一个是自定义的 RPC Service 方法，因此可以更进一步的执行命令 `grpcurl -plaintext localhost:8001 list proto.TagService` 查看其子类的 RPC 方法信息。

在了解该服务具体有什么 RPC 方法后，可以执行下述命令去调用 RPC 方法：

````bash
# 注意引号前要加 \ ，否则无法转换成 json，会报错
$ grpcurl -plaintext  -d '{\"name\": \"Java\"}' localhost:8001 proto.TagService.GetTagList
{
  "list": [
    {
      "id": "2",
      "name": "Java",
      "state": 1
    }
  ],
  "pager": {
    "page": "1",
    "pageSize": "10",
    "totalRows": "1"
  }
}

# 错误请求
$ grpcurl -plaintext -d '{"name":"Java"}' localhost:8001 proto.TagService.GetTagList 
Error invoking method "proto.TagService.GetTagList": error getting request data: invalid character 'n' looking for beginning of object key string
````

在这里使用到了 `grpcurl` 工具的`-d `选项，**其输入的内容必须为 JSON 格式**，该内容将被解析，最终以 `protobuf` 二进制格式传输到 gRPC Server，可以简单理解为 RPC 方法的入参信息，**也可以不传，不指定-d 选项即可。****==注意引号前要加 \ ，否则无法转换成 JSON ，会报错。==**

### 6. 错误处理

在项目的实际运行中，常常会有各种奇奇怪怪的问题触发，也就是要返回错误的情况，在这里可以将数据源，也就是博客的后端服务停掉，再利用工具重新请求。

```bash
$ grpcurl -plaintext localhost:8001 proto.TagService.GetTagList 
ERROR:
  Code: Unknown
  Message: Get "http://127.0.0.1:8000/auth?app_key=admin&app_secret=go-learning": dial tcp 127.0.0.1:8000: connectex: No connection could be made because the target machine actively refused it.
```

会发现其返回的字段分为两个，一个是 Code，另外一个是 Message，也就是对应着第三章提到 `grpc-status` 和 `grpc-message` 两个字段，它们共同代表着 gRPC 的整体调用情况。

#### a. gRPC 状态码

下表为官方给出的全部状态响应码。

| Code | Status              | Notes                                        |
| :--- | :------------------ | :------------------------------------------- |
| 0    | OK                  | 成功                                         |
| 1    | CANCELLED           | 该操作被调用方取消                           |
| 2    | UNKNOWN             | 未知错误                                     |
| 3    | INVALID_ARGUMENT    | 无效的参数                                   |
| 4    | DEADLINE_EXCEEDED   | 在操作完成之前超过了约定的最后期限。         |
| 5    | NOT_FOUND           | 找不到                                       |
| 6    | ALREADY_EXISTS      | 已经存在                                     |
| 7    | PERMISSION_DENIED   | 权限不足                                     |
| 8    | RESOURCE_EXHAUSTED  | 资源耗尽                                     |
| 9    | FAILED_PRECONDITION | 该操作被拒绝，因为未处于执行该操作所需的状态 |
| 10   | ABORTED             | 该操作被中止                                 |
| 11   | OUT_OF_RANGE        | 超出范围，尝试执行的操作超出了约定的有效范围 |
| 12   | UNIMPLEMENTED       | 未实现                                       |
| 13   | INTERNAL            | 内部错误                                     |
| 14   | UNAVAILABLE         | 该服务当前不可用。                           |
| 15   | DATA_LOSS           | 不可恢复的数据丢失或损坏。                   |

那么对应在刚刚的调用结果，状态码是 UNKNOWN，这是为什么呢，可以查看底层的处理源码，如下：

````go
func FromError(err error) (s *Status, ok bool) {
    ...
    if se, ok := err.(interface {
        GRPCStatus() *Status
    }); ok {
        return se.GRPCStatus(), true
    }
    return New(codes.Unknown, err.Error()), false
}
````

可以看到，实际上若不是 `GRPCStatus `类型的方法，都是默认返回 `codes.Unknown`，也就是未知。而目前的报错，实际上是直接返回 `return err` 的，需要自定义返回的话只需要遵循内部规范实现即可。

#### b. 错误码处理

在项目的 `pkg` 目录下新建` errcode `目录，并创建 `errcode.go `文件。

```go
package errcode

import "fmt"

type Error struct {
   code int
   msg  string
}

var _codes = map[int]string{}

func NewError(code int, msg string) *Error {
   if _, ok := _codes[code]; ok {
      panic(fmt.Sprintf("错误码 %d 已经存在，请更换一个", code))
   }

   _codes[code] = msg
   return &Error{code: code, msg: msg}
}

func (e *Error) Error() string {
   return fmt.Sprintf("错误码：%d, 错误信息:：%s", e.Code(), e.Msg())
}

func (e *Error) Code() int {
   return e.code
}
func (e *Error) Msg() string {
   return e.msg
}
```

接下来继续在目录下新建 `common_error.go` 文件，写入如下公共错误码：

```go
package errcode

var (
   Success          = NewError(0, "成功")
   Fail             = NewError(10000000, "内部错误")
   InvalidParams    = NewError(10000001, "无效参数")
   Unauthorized     = NewError(10000002, "认证错误")
   NotFound         = NewError(10000003, "没有找到")
   Unknown          = NewError(10000004, "未知")
   DeadlineExceeded = NewError(10000005, "超出最后截止期限")
   AccessDenied     = NewError(10000006, "访问被拒绝")
   LimitExceed      = NewError(10000007, "访问限制")
   MethodNotAllowed = NewError(10000008, "不支持该方法")
)
```

继续在目录下新建` rpc_error.go `文件，写入如下 RPC 相关的处理方法：

```go
package errcode

import (
   "google.golang.org/grpc/codes"
   "google.golang.org/grpc/status"
)

func TogRPCError(err *Error) error {
   s := status.New(ToRPCCode(err.code), err.Msg())
   return s.Err()
}

func ToRPCCode(code int) codes.Code {
   var statusCode codes.Code
   switch code {
   case Fail.Code():
      statusCode = codes.Internal
   case InvalidParams.Code():
      statusCode = codes.InvalidArgument
   case Unauthorized.Code():
      statusCode = codes.Unauthenticated
   case AccessDenied.Code():
      statusCode = codes.PermissionDenied
   case DeadlineExceeded.Code():
      statusCode = codes.DeadlineExceeded
   case NotFound.Code():
      statusCode = codes.NotFound
   case LimitExceed.Code():
      statusCode = codes.ResourceExhausted
   case MethodNotAllowed.Code():
      statusCode = codes.Unimplemented
   default:
      statusCode = codes.Unknown
   }
   
   return statusCode
}
```

#### c. 业务错误码

这个时候会发现，返回的错误最后都会被转换为 RPC 的错误信息，那原始的业务错误码，可以放在哪里呢，因为没有业务错误码，怎么知道错在具体哪个业务板块，面向用户的客户端又如何特殊处理呢？

那么实际上，在 gRPC 的状态消息中其一共包含三个属性，分别是错误代码、错误消息、错误详细信息，因此可以通过错误详细信息这个字段来实现这个功能，其 `googleapis` 的` status.pb.go `原型如下：

```go
type Status struct {
    Code      int32 `protobuf:"..."`
    Message   string `protobuf:"..."`
    Details   []*any.Any `protobuf:"..."`
    ...
}
```

因此只需要对应其下层属性，让其与应用程序的错误码机制产生映射关系即可，首先在 `common.proto` 中增加 `any.proto` 文件的引入和消息体 `Error` 的定义，将其作为应用程序的错误码原型，如下：

```protobuf
syntax = "proto3";

import "google/protobuf/any.proto";

package proto;

message Pager{
    int64 page = 1;
    int64 page_size = 2;
    int64 total_rows = 3;
}

message Error {
    int32 code = 1;
    string message  = 2;
    google.protobuf.Any detail = 3;
}
```

在 win10 上直接执行编译命令会报错。

```bash
$ protoc --go_out=plugins=grpc:. ./proto/*.proto
google/protobuf/any.proto: File not found.
proto/common.proto:3:1: Import "google/protobuf/any.proto" was not found or had errors.
proto/common.proto:16:5: "google.protobuf.Any" is not defined.
```

最简单的方法是将`any.proto`文件复制到项目目录下的`google/protobuf`下。

![image-20220522213910003](https://raw.githubusercontent.com/tonshz/test/master/img/202205222139037.png)

接着重新执行编译命令 `protoc --go_out=plugins=grpc:. ./proto/*.proto` ，再打开刚刚编写的 `rpc_error.go `文件，修改 `TogRPCError` 方法，新增 `Details` 属性，如下：

```go
...
func TogRPCError(err *Error) error {
	//s := status.New(ToRPCCode(err.code), err.Msg())
	// 新增 Details 属性
	s, _ := status.New(ToRPCCode(err.Code()), err.Msg()).WithDetails(&pb.Error{Code: int32(err.Code()), Message: err.Msg()})
	return s.Err()
}
...
```

这时候又有新的问题了，那就是服务自身，在处理 err 时，如何能够获取到错误类型呢，可以通过新增` FromError` 方法。而针对有的应用程序，除了希望把业务错误码放进 Details 中，还希望把其它信息也放进去的话，可以通过新增`ToRPCStatus`方法。

```go
...
type Status struct {
   *status.Status
}

func FromError(err error) *Status {
   s, _ := status.FromError(err)
   return &Status{s}
}

// 将其他信息也放入 Details 中
func ToRPCStatus(code int, msg string) *Status {
   s, _ := status.New(ToRPCCode(code), msg).WithDetails(&pb.Error{Code: int32(code), Message: msg})
   return &Status{s}
}

// 将业务错误码放进 Details 中
func TogRPCError(err *Error) error {
   //s := status.New(ToRPCCode(err.code), err.Msg())
   // 新增 Details 属性
   s, _ := status.New(ToRPCCode(err.Code()), err.Msg()).WithDetails(&pb.Error{Code: int32(err.Code()), Message: err.Msg()})
   return s.Err()
}
...
```

#### d. 验证

在项目的 errcode 目录下新建 `module_error.go `文件，写入模块的业务错误码。

```go
package errcode

var (
   ErrorGetTagListFail = NewError(20010001, "获取标签列表失败")
)
```

接下来修改 server 目录下的 `tag.go `文件中的 `GetTagList` 方法，将业务错误码填入（同时建议记录日志，可参见 HTTP 应用章节），如下：

```go
...

func (t *TagServer) GetTagList(ctx context.Context, r *pb.GetTagListRequest) (*pb.GetTagListResponse, error) {
	// localhost：不通过网卡传输，不受网络防火墙和网卡相关的限制。
	// 127.0.0.1：通过网卡传输，依赖网卡，并受到网卡和防火墙相关的限制。
	api := bapi.NewAPI("http://127.0.0.1:8000")
	body, err := api.GetTagList(ctx, r.GetName())
	if err != nil {
		// 填入业务错误码 
		return nil, errcode.TogRPCError(errcode.ErrorGetTagListFail)
	}

	tagList := pb.GetTagListResponse{}
	err = json.Unmarshal(body, &tagList)
	if err != nil {
		// 填入错误码 
		return nil, errcode.TogRPCError(errcode.Fail)
	}

	return &tagList, nil
}
```

这个时候还是保持该 RPC 服务的数据源（博客服务）的停止运行，并在添加完业务错误码后重新运行 RPC 服务，然后利用 `grpcurl `工具查看错误码是否达到预期结果。

```bash
$ grpcurl -plaintext localhost:8001 proto.TagService.GetTagList
ERROR:
  Code: Unknown
  Message: 获取标签列表失败
  Details:
  1)    {"@type":"type.googleapis.com/proto.Error","code":20010001,"message":"获取标签列表失败"}
```

那么外部客户端可通过 Details 属性知道业务错误码了，那内部客户端要如何使用呢，如下：

```go
err := errcode.TogRPCError(errcode.ErrorGetTagListFail)
sts := errcode.FromError(err)
details := sts.Details()
```

最终错误信息是以 RPC 所返回的 err 进行传递的，因此只需要利用先前编写的` FromError `方法，解析为 status，接着调用 Details 方法，进行 Error 的断言，就可以精确的获取到业务错误码了。

### 7. 为什么，是什么

#### a. 为什么可以转换

在编写 RPC 方法时，可以直接用 json 和 protobuf 所生成出来的结构体互相转换，这是为什么呢，为什么可以这么做，可以一起看看所生成的.pb.go 文件内容，如下：

```go
type GetTagListReply struct {
    List                 []*Tag   `protobuf:"... json:"list,omitempty"`
    Pager                *Pager   `protobuf:"... json:"pager,omitempty"`
}
```

实际上在 protoc 生成`.pb.go`文件时，会在所生成的结构体上打入 JSON Tag，**其默认的规则就是下划线命名**，因此可以通过该方式进行转换，而若是出现字段刚好不兼容的情况，也可以通过结构体转结构体的方式，最后达到这种效果。

#### b. 为什么零值不显示

在实际的运行调用中，有一个问题常常被初学者所问到，占比非常高，那就是为什么在调用过程中，有的数据没有展示出来，例如：name 为空字符串、state 为 0 的话，就不会在 RPC 返回的数据中展示出来，会发现其实是有规律的，他们都是零值，可以看到所生成的`.pb.go `文件内容，如下：

```go
type Tag struct {
    Id                   int64    `protobuf:"... json:"id,omitempty"`
    Name                 string   `protobuf:"... json:"name,omitempty"`
    State                uint32   `protobuf:"... json:"state,omitempty"`
    ...
}
```

在上小节也有提到，实际上所生成的结构体是有打 JSON Tag 的，它在所有的字段中都标明了 `omitempty` 属性，也就是当值为该类型的零值时将不会序列化该字段。

那么紧跟这个问题，就会出现第二个最常见的被提问的疑惑，那就是能不能解决这个”问题“，实际上这个并不是”问题“，因为这是 Protobuf 的规范，在官方手册的 JSON Mapping 小节明确指出**，如果字段在 Protobuf 中具有默认值，则默认情况下会在 JSON 编码数据中将其省略以节省空间。**

#### c. googleapis 是什么

googleapis 代指 Google API 的公共接口定义，在 Github 上搜索 googleapis 就可以找到对应的仓库了，不过需要注意的是由于 Go 具有不同的目录结构，因此很难在原始的 googleapis 存储库存储和生成 Go gRPC 源代码，因此 Go gRPC 实际使用的是 go-genproto 仓库，该仓库有如下两个主要使用来源：

+ `google/protobuf`：protobuf 和 ptypes 子目录中的代码是均从存储库派生的， protobuf 中的消息体用于描述 Protobuf 本身。 ptypes 下的消息体定义了常见的常见类型。

+ `googleapis/googleapis`：专门用于与 Google API 进行交互的类型。

## 五、gRPC 服务间的内调

在上一个章节中，运行了一个最基本的 gRPC 服务，那么在实际上的应用场景，服务是会有多个的，并且随着需求的迭代拆分重合，服务会越来越多，到上百个也是颇为常见的。因此在这么多的服务中，最常见的就是 gRPC 服务间的内调行为，再细化下来，其实就是客户端如何调用 gRPC 服务端的问题，那么在本章节将会使用客户端如何调用 gRPC 服务端和并对此做一个深入了解。

### 1. 进行 gRPC 调用

理论上在任何能够执行 Go 语言代码，且网络互通的地方都可以进行 gRPC 调用，它并不受限于必须在什么类型应用程序下才能够调用。接下来在项目下新建 `client` 目录，创建 `client.go `文件，编写一个示例来调用先前所编写的 gRPC 服务，如下代码：

```go
package main

import (
   pb "ch03/proto"
   "context"
   "google.golang.org/grpc"
   "google.golang.org/grpc/credentials/insecure"
   "log"
)

func main() {
   /*
      context.Background()
      返回一个非零的空上下文。它永远不会被取消，没有价值，也没有最后期限。
      它通常由主函数、初始化和测试使用，并作为传入请求的顶级上下文。
   */
   ctx := context.Background()
   clientConn, _ := GetClientConn(ctx, "localhost:8004", nil)
   defer clientConn.Close()

   // 初始化指定 RPC Proto Service 的客户端实例对象
   tagServiceClient := pb.NewTagServiceClient(clientConn)
   // 发起指定 RPC 方法的调用
   resp, _ := tagServiceClient.GetTagList(ctx, &pb.GetTagListRequest{Name: "Go"})

   log.Printf("resp: %v", resp)
}

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
   // grpc.WithInsecure() 已弃用，作用为跳过对服务器证书的验证，此时客户端和服务端会使用明文通信
   // 使用 WithTransportCredentials 和 insecure.NewCredentials() 代替
   //opts = append(opts, grpc.WithInsecure())
   // insecure.NewCredentials 返回一个禁用传输安全的凭据
   opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

   /*
      grpc.DialContext() 创建到给定目标的客户端连接。
      默认情况下，它是一个非阻塞拨号（该功能不会等待建立连接，并且连接发生在后台）。
      要使其成为阻塞拨号，请使用 WithBlock() 拨号选项。
   */
   return grpc.DialContext(ctx, target, opts...)
}
```

在上述 gRPC 调用的示例代码中，一共分为三大步，分别是：

- `grpc.DialContext`：创建给定目标的客户端连接，另外所要请求的服务端是非加密模式的，因此调用了 `grpc.WithInsecure` 方法禁用了此 ClientConn 的传输安全性验证。
- `pb.NewTagServiceClient`：初始化指定 RPC Proto Service 的客户端实例对象。
- `tagServiceClient.GetTagList`：发起指定 RPC 方法的调用。

### 2. grpc.Dial 做了什么

常常有的人会说在调用 `grpc.Dial` 或 `grpc.DialContext` 方法时，客户端就已经与服务端建立起了连接，但这对不对呢，这是需要细心思考的一个点，客户端真的是一调用 `Dial` 相关方法就马上建立了可用连接吗，一起尝试一下，示例代码：

```go
func main() {
    ctx := context.Background()
    clientConn, _ := GetClientConn(ctx, "localhost:8004", nil)
    defer clientConn.Close()
}
```

在上述代码中，只保留了创建给定目标的客户端连接的部分代码，然后执行该程序，接着马上查看抓包工具的情况下，竟然提示一个包都没有，那么这算真正连接了吗？

实际上，如果真的想在调用 `DialContext` 方法时就马上打通与服务端的连接，那么需要调用 `WithBlock` 方法来进行设置，那么它在发起拨号连接时就会阻塞等待连接完成，并且最终连接会到达 `Ready` 状态，这样子在此刻的连接才是正式可用的，代码如下：

```go
func main() {
    ctx := context.Background()
    clientConn, _ := GetClientConn(
        ctx,
        "localhost:8004",
        []grpc.DialOption{grpc.WithBlock()},
    )
    defer clientConn.Close()
}
```

再次进行抓包（使用 WireShark 抓包报错`tcp port numbers reused `），查看效果，如下：

![image-20220522210941862](https://raw.githubusercontent.com/tonshz/test/master/img/202205222109965.png)

#### a. 源码分析

那么在调用 `grpc.Dial` 或 `grpc.DialContext` 方法时，到底做了什么事情呢，为什么还要调用 `WithBlock` 方法那么“麻烦”，接下来一起看看正在调用时运行的 `goroutine `情况，如下：

![image](https://raw.githubusercontent.com/tonshz/test/master/img/202205222110346.jpeg)

可以看到有几个核心方法一直在等待/处理信号，通过分析底层源码可得知。涉及如下：

```go
func (ac *addrConn) connect()
func (ac *addrConn) resetTransport()
func (ac *addrConn) createTransport(addr resolver.Address, copts transport.ConnectOptions, connectDeadline time.Time)
func (ac *addrConn) getReadyTransport()
```

在这里主要分析所提示的 `resetTransport` 方法，看看都做了什么。核心代码如下：

```go
func (ac *addrConn) resetTransport() {
	ac.mu.Lock()
	if ac.state == connectivity.Shutdown {
		ac.mu.Unlock()
		return
	}

	addrs := ac.addrs
	backoffFor := ac.dopts.bs.Backoff(ac.backoffIdx)
	// This will be the duration that dial gets to finish.
	dialDuration := minConnectTimeout
	if ac.dopts.minConnectTimeout != nil {
		dialDuration = ac.dopts.minConnectTimeout()
	}

	if dialDuration < backoffFor {
		// Give dial more time as we keep failing to connect.
		dialDuration = backoffFor
	}
	// We can potentially spend all the time trying the first address, and
	// if the server accepts the connection and then hangs, the following
	// addresses will never be tried.
	//
	// The spec doesn't mention what should be done for multiple addresses.
	// https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md#proposed-backoff-algorithm
	connectDeadline := time.Now().Add(dialDuration)

	ac.updateConnectivityState(connectivity.Connecting, nil)
	ac.mu.Unlock()

	if err := ac.tryAllAddrs(addrs, connectDeadline); err != nil {
		ac.cc.resolveNow(resolver.ResolveNowOptions{})
		// After exhausting all addresses, the addrConn enters
		// TRANSIENT_FAILURE.
		ac.mu.Lock()
		if ac.state == connectivity.Shutdown {
			ac.mu.Unlock()
			return
		}
		ac.updateConnectivityState(connectivity.TransientFailure, err)

		// Backoff.
		b := ac.resetBackoff
		ac.mu.Unlock()

		timer := time.NewTimer(backoffFor)
		select {
		case <-timer.C:
			ac.mu.Lock()
			ac.backoffIdx++
			ac.mu.Unlock()
		case <-b:
			timer.Stop()
		case <-ac.ctx.Done():
			timer.Stop()
			return
		}

		ac.mu.Lock()
		if ac.state != connectivity.Shutdown {
			ac.updateConnectivityState(connectivity.Idle, err)
		}
		ac.mu.Unlock()
		return
	}
	// Success; reset backoff.
	ac.mu.Lock()
	ac.backoffIdx = 0
	ac.mu.Unlock()
}
```

通过上述代码可得知，在该方法中会不断地去尝试创建连接，若成功则结束。否则不断地根据 `Backoff` 算法的重试机制去尝试创建连接，直到成功为止。

#### b. 小结

单纯调用 `grpc.DialContext` 方法是异步建立连接的，并不会马上就成为可用连接了，仅处于 `Connecting` 状态（需要多久则取决于外部因素，例如：网络），正式要到达 `Ready` 状态，这个连接才算是真正的可用。

再回顾到前面的示例中，为什么抓包时一个包都抓不到，实际上连接立即建立了，但 main 结束的很快，因此可能刚建立就被销毁了，也可能还处于 `Connecting` 状态，没来得及产生具体的网络活动，自然也就抓取不到任何包了。

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

## 七、生成接口文档

所开发的接口有没有接口文档，接口文档有没有及时更新，是一个永恒的话题。在本章节将继续使用 Swagger 作为接口文档平台，但是与第二章不同，载体变成了` Protobuf`，`Protobuf `是强规范的，其本身就包含了字段名和字段类型等等信息，因此其会更加的简便。

在接下来的章节中，将继续基于同端口同 RPC 方法支持双流量（grpc-gateway 方案）的服务代码来进行开发和演示。

### 1. 安装和下载

#### a. 安装

针对后续 Swagger 接口文档的开发和使用，需要安装 `protoc `的插件 `protoc-gen-swagger`，它的作用是通过 proto 文件来生成 swagger 定义（`.swagger.json`），安装命令如下：

```bash
$ go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
```

#### b. 下载 Swagger UI 文件 

Swagger 提供可视化的接口管理平台，也就是 Swagger UI，首先需要到 `https://github.com/swagger-api/swagger-ui` 上将其源码压缩包下载下来，接着在项目的 `third_party` 目录下新建 `swagger-ui` 目录，将其 `dist` 目录下的所有资源文件拷贝到项目的 `third_party/swagger-ui` 目录中去。

### 2. 静态资源转换

在上一步中将 Swagger UI 的资源文件拷贝到了项目的` swagger-ui `目录中，但是这时候应用程序还不可以使用它，**需要使用 `go-bindata `库将其资源文件转换为 Go 代码**，便于后续使用，安装命令如下：

```bash
$ go get -u github.com/go-bindata/go-bindata/...
```

接下来在项目的 `pkg` 目录下新建 `swagger-ui` 目录，并在项目根目录执行下述转换命令：

```bash
$ go-bindata --nocompress -pkg swagger -o pkg/swagger/data.go third_party/swagger-ui/...
```

在执行完毕后，应当在项目的 `pkg/swagger` 目录下创建了 `data.go `文件。\

### 3. Swagger UI 处理和访问

为了让刚刚转换的静态资源代码能够让外部访问到，需要安装 `go-bindata-assetfs` 库，它能够结合 `net/http` 标准库和 `go-bindata` 所生成 `Swagger UI` 的 `Go` 代码两者来供外部访问，安装命令如下：

```bash
$ go get -u github.com/elazarl/go-bindata-assetfs/...
```

安装完成后，打开启动文件 `main.go`，修改 HTTP Server 相关的代码，如下：

```go
import (
   "ch03/pkg/swagger"
   ...
   assetfs "github.com/elazarl/go-bindata-assetfs"
   ...
)
...
// 启动 HTTP 服务器
func runHttpServer() *http.ServeMux {
	// 初始化一个 HTTP 请求多路复用器
	serveMux := http.NewServeMux()
	// HandleFunc() 为给定模式注册处理函数
	// 新增了一个 /ping 路由及其 Handler，可做心跳检测
	serveMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`pong`))
	})

    // 注意此处前后都有 /
	prefix := "/swagger-ui/"
	var fileServer = http.FileServer(&assetfs.AssetFS{
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
		Prefix:   "third_party/swagger-ui",
	})
	serveMux.Handle(prefix, http.StripPrefix(prefix, fileServer))
	return serveMux
}
```

在上述代码中，通过引用先前通过` go-bindata-assetfs `所生成的 `data.go` 文件中的资源信息，结合两者来对外提供 `swagger-ui` 的 Web UI 服务。需要注意的是，因为所生成的文件的原因，因此 `swagger.Asset` 和 `swagger.AssetDir` 的引用在 IDE 识别上（会标红）存着一定的问题，但实际程序运行是没有问题的，只需要通过命令行启动就可以了。

可通过以下方法解决 IDE 爆红的问题。

```properties
idea.max.intellisense.filesize=10000
```

![image-20220523204651453](https://raw.githubusercontent.com/tonshz/test/master/img/202205232047287.png)

![image-20220523204744666](https://raw.githubusercontent.com/tonshz/test/master/img/202205232047596.png)

重新运行该服务，通过浏览器访问 `http://127.0.0.1:8004/swagger-ui/`，查看结果如下：

![image-20220523205832161](https://raw.githubusercontent.com/tonshz/test/master/img/202205232058492.png)

以上看到的就是 Swagger UI 的默认展示界面，默认展示的是 Demo 示例（也就是输入框中的 swagger 地址），如果看到如上界面，那说明一切正常。

### 4. Swagger 描述文件生成和读取

既然 Swagger UI 已经能够看到了，那项目的接口文档又如何读取呢，其实刚刚在上一步，可以看到默认示例，读取的是一个` swagger.json `的远程地址，也就是只要本地服务中也有对应的 `swagger.json`，就能够展示项目服务的接口文档了。

因此先需要进行 swagger 定义文件的生成，在项目根目录下执行以下命令：

```bash
$ protoc -I C:\Users\zyc\protoc\include -I C:\Users\zyc\go\pkg\mod\github.com\grpc-ecosystem\grpc-gateway@v1.14.5\third_party\googleapis -I . --swagger_out=logtostderr=true:. ./proto/*.proto
```

执行完毕后会发现 proto 目录下会多出 `common.swagger.json` 和 `tag.swagger.json `两个文件，文件内容是对应的 API 描述信息。

接下来需要让浏览器能够访问到本地所生成的` swagger.json`，也就是需要有一个能够访问本地 proto 目录下的`.swagger.json` 的文件服务，继续修改 `main.go `文件，如下：

```go
// 启动 HTTP 服务器
func runHttpServer() *http.ServeMux {
   // 初始化一个 HTTP 请求多路复用器
   serveMux := http.NewServeMux()
   // HandleFunc() 为给定模式注册处理函数
   // 新增了一个 /ping 路由及其 Handler，可做心跳检测
   serveMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
      _, _ = w.Write([]byte(`pong`))
   })

   // 注意此处前后都有 /
   prefix := "/swagger-ui/"
   var fileServer = http.FileServer(&assetfs.AssetFS{
      Asset:    swagger.Asset,
      AssetDir: swagger.AssetDir,
      Prefix:   "third_party/swagger-ui",
   })
   serveMux.Handle(prefix, http.StripPrefix(prefix, fileServer))
    
   // 修改以便获取本地 proto 目录下的 .swagger.json 文件
   serveMux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request){
      if !strings.HasSuffix(r.URL.Path, "swagger.json"){
         http.NotFound(w, r)
         return
      }
      
      p := strings.TrimPrefix(r.URL.Path, "/swagger/")
      p = path.Join("proto", p)
      
      http.ServeFile(w, r, p)
   })
   return serveMux
}
```

重新运行该服务，通过浏览器访问 `http://127.0.0.1:8004/swagger/tag.swagger.json`，查看能够成功得到所生成的 API 描述信息。

### 5. 查看接口文档

接下来只需要把想要查看的 `swagger.json `的访问地址，填入输入框中，点击“Explore”就可以查看到对应的接口信息了，如下图：

![image-20220523211223671](https://raw.githubusercontent.com/tonshz/test/master/img/202205232112843.png)

## 八、拦截器介绍和实际使用

想在每一个 RPC 方法的前面或后面做某些操作，想针对某个业务模块的 RPC 方法进行统一的特殊处理，想对 RPC 方法进行鉴权校验，想对 RPC 方法进行上下文的超时控制，想对每个 RPC 方法的请求都做日志记录，怎么做呢？

这诸如类似的一切需求的答案，都在本章节将要介绍的拦截器（Interceptor）上，能够借助它实现许许多多的定制功能且不直接侵入业务代码。

### 1. 拦截器的类型

在 gRPC 中，根据拦截器拦截的 RPC 调用的类型，拦截器在分类上可以分为如下两种：

- 一元拦截器（Unary Interceptor）：拦截和处理一元 RPC 调用。
- 流拦截器（Stream Interceptor）：拦截和处理流式 RPC 调用。

虽然总的来说是只有两种拦截器分类，但是再细分下去，客户端和服务端每一个都有其自己的一元和流拦截器的具体类型。因此，gRPC 中也可以说总共有四种不同类型的拦截器。

### 2. 客户端和服务器端拦截器

#### a.  客户端

##### 一元拦截器

客户端的一元拦截器类型为 `UnaryClientInterceptor`，方法原型如下：

```go
type UnaryClientInterceptor func(
    ctx context.Context, 
    method string,
    req, 
    reply interface{},
    cc *ClientConn, 
    invoker UnaryInvoker, 
    opts ...CallOption
) error
```

一元拦截器的实现通常可以分为三个部分：预处理，调用 RPC 方法和后处理。其一共分为七个参数，分别是：RPC 上下文、所调用的方法、RPC 方法的请求参数和响应结果，客户端连接句柄、所调用的 RPC 方法以及调用的配置。

##### 流拦截器

客户端的流拦截器类型为 `StreamClientInterceptor`，方法原型如下：

```go
type StreamClientInterceptor func(
    ctx context.Context,
    desc *StreamDesc, 
    cc *ClientConn,
    method string, 
    streamer Streamer, 
    opts ...CallOption
) (ClientStream, error)
```

流拦截器的实现包括预处理和流操作拦截，并不能在事后进行 RPC 方法调用和后处理，而是拦截用户对流的操作。

#### b. 服务端

##### 一元拦截器

服务端的一元拦截器类型为 `UnaryServerInterceptor`，方法原型如下：

```go
type UnaryServerInterceptor func(
    ctx context.Context, 
    req interface{}, 
    info *UnaryServerInfo, 
    handler UnaryHandler
) (resp interface{}, err error)
```

其一共包含四个参数，分别是 RPC 上下文、RPC 方法的请求参数、RPC 方法的所有信息、RPC 方法本身。

##### 流拦截器

服务端的流拦截器类型为 `StreamServerInterceptor`，方法原型如下：

```go
type StreamServerInterceptor func(
    srv interface{}, 
    ss ServerStream, 
    info *StreamServerInfo, 
    handler StreamHandler
) error
```

### 3. 实现一个拦截器

在了解了 gRPC 拦截器的基本概念后，修改项目`main.go`文件，添加拦截器相关代码。

```go
...
func runGrpcServer() *grpc.Server {
   // 新增拦截器相关代码
   opts := []grpc.ServerOption{
      grpc.UnaryInterceptor(HelloInterceptor),
   }

   s := grpc.NewServer(opts...)
   pb.RegisterTagServiceServer(s, server.NewTagServer())
   reflection.Register(s)
   return s
}

// 实现一个简单的一元拦截器
func HelloInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
   handler grpc.UnaryHandler) (interface{}, error) {
   log.Println("hello")
   resp, err := handler(ctx, req)
   log.Println("goodbye")
   return resp, err
}
...
```

在上述代码中，除了实现一个简单的一元拦截器以外，还初次使用到了 `grpc.ServerOption`，gRPC Server 的相关属性都可以在此设置，例如：`credentials`、`keepalive` 等等参数。服务端拦截器也在此注册，但是需要以指定的类型进行封装，例如一元拦截器是使用 `grpc.UnaryInterceptor`。

在验证上，修改完毕后需要重新启动该服务，调用对应的 RPC 接口，查看控制台是否输出“你好”和“再见”两个字符串，若有则实现正确。

![image-20220523213321141](https://raw.githubusercontent.com/tonshz/test/master/img/202205232133169.png)

### 4. 能使用多少个拦截器

既然实现了一个拦截器，那么在实际的应用程序中，肯定是不止一个了，按常规来讲，既然支持了一个，支持多个拦截器的注册和使用应该不过分吧，再来试试，代码如下：

```go
func runGrpcServer() *grpc.Server {
   // 新增拦截器相关代码
   opts := []grpc.ServerOption{
      grpc.UnaryInterceptor(HelloInterceptor),
      grpc.UnaryInterceptor(WorldInterceptor),
   }

   s := grpc.NewServer(opts...)
   pb.RegisterTagServiceServer(s, server.NewTagServer())
   reflection.Register(s)
   return s
}

// 实现一个简单的一元拦截器
func HelloInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
   handler grpc.UnaryHandler) (interface{}, error) {
   log.Println("hello")
   resp, err := handler(ctx, req)
   log.Println("goodbye")
   return resp, err
}

func WorldInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
   handler grpc.UnaryHandler) (interface{}, error) {
   log.Println("world")
   resp, err := handler(ctx, req)
   log.Println("say hi")
   return resp, err
}
```

重新运行服务，查看输出结果，如下：

```bash
panic: The unary server interceptor was already set and may not be reset.
```

会发现启动服务就报错了，会提示“一元服务器拦截器已经设置，不能重置”，**也就是一种类型的拦截器只允许设置一个。**

### 5. 需要多个拦截器

虽然 grpc-go 官方只允许设置一个拦截器，但不代表只能”用”一个拦截器。

在实际使用上，常常会希望将不同功能设计为不同的拦截器，这个时候，除了自己实现一套多拦截器的逻辑（拦截器中调拦截器即可）以外，还可以直接使用 gRPC 应用生态（`grpc-ecosystem`）中的 `go-grpc-middleware `所提供的 `grpc.UnaryInterceptor` 和 `grpc.StreamInterceptor` 链式方法来达到这个目的。

#### a. 安装

```bash
$ go get -u github.com/grpc-ecosystem/go-grpc-middleware@v1.1.0
```

#### b. 使用

修改 gRPC Server 的相关代码，进行多拦截器的注册。

```go
func runGrpcServer() *grpc.Server {
   // 新增拦截器相关代码
   //opts := []grpc.ServerOption{
   // grpc.UnaryInterceptor(HelloInterceptor),
   //}
   // 进行多拦截器的注册
   opts := []grpc.ServerOption{
      grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
         HelloInterceptor,
         WorldInterceptor,
      )),
   }go

   s := grpc.NewServer(opts...)
   pb.RegisterTagServiceServer(s, server.NewTagServer())
   reflection.Register(s)
   return s
}
```

在 `grpc.UnaryInterceptor` 中嵌套 `grpc_middleware.ChainUnaryServer` 后重新启动服务，查看输出结果：

```bash
2022/05/23 21:39:31 hello # 第一个注册的拦截器
2022/05/23 21:39:31 world # 第二个注册的拦截器
2022/05/23 21:39:33 say hi # 第二个注册的拦截器
2022/05/23 21:39:33 goodbye # 第一个注册的拦截器
```

两个拦截器都调用成功，完成常规多拦截器的需求。

#### c. 实现原理

单单会用还是不行的，`go-grpc-middleware `是如何实现的呢？

```go
// ChainUnaryServer creates a single interceptor out of a chain of many interceptors.
//
// Execution is done in left-to-right order, including passing of context.
// For example ChainUnaryServer(one, two, three) will execute one before two before three, and three
// will see context changes of one and two.
// ChainUnaryServer 从许多拦截器链中创建一个拦截器。执行按从左到右的顺序完成，包括传递上下文。
// 例如 ChainUnaryServer(one, two, three) 会在 3 之前执行 1 之前 2 ，并且 3 会看到 1 和 2 的上下文变化。
func ChainUnaryServer(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
   n := len(interceptors)

   return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
      chainer := func(currentInter grpc.UnaryServerInterceptor, currentHandler grpc.UnaryHandler) grpc.UnaryHandler {
         return func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
            return currentInter(currentCtx, currentReq, info, currentHandler)
         }
      }

      chainedHandler := handler
      for i := n - 1; i >= 0; i-- {
         chainedHandler = chainer(interceptors[i], chainedHandler)
      }
go
      return chainedHandler(ctx, req)
   }
}
```

当拦截器数量大于 1 时，从 `interceptors[1]` 开始递归，每一个递归的拦截器 `interceptors[i]` 会不断地执行，最后才会去真正执行代表 RPC 方法的 `handler` 。

### 6. 服务端-常用拦截器

在项目的运行中，常常会有那么一些应用拦截器，是必须要有的，因此可以总结出来一套简单的而行之有效的“公共”拦截器，在本节将模拟实际的使用场景进来实现。

在项目的 `internal/middleware` 目录下新建存储服务端拦截器的` server_interceptor.go `文件，另外后续的服务端拦截器的相关注册行为也均在 `runGrpcServer `方法中进行处理，例如：

```bash
func runGrpcServer() *grpc.Server {
   opts := []grpc.ServerOption{
      grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
         middleware.XXXXX,
      )),
   }
   ...
}
```

#### a. 日志

在应用程序的运行中，常常需要一些信息来协助做问题的排查和追踪，因此日志信息的及时记录和处理是非常必要的，接下来将针对常见的访问日志和错误日志进行日志输出。在实际使用的过程中，可以将案例中的默认日志实例替换为应用中实际在使用的文件日志的模式（例如参考第二章的日志器）。

##### 访问日志

打开 `server_interceptor.go `文件，新增针对访问记录的日志拦截器，代码如下：

```go
package middleware

import (
   "context"
   "google.golang.org/grpc"
   "log"
   "time"
)

func AccessLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
   handler grpc.UnaryHandler) (interface{}, error) {
   requestLog := "access request log: method: %s, begin_time: %d, request: %v"
   beginTime := time.Now().Local().Unix()
   // FullMethod 是完整的 RPC 方法字符串，即 package.servicemethod。
   log.Printf(requestLog, info.FullMethod, beginTime, req)

   resp, err := handler(ctx, req)

   responseLog := "access response log: method: %s, begin_time: %d, end_time: %d, response: %v"
   endTime := time.Now().Local().Unix()
   log.Printf(responseLog, info.FullMethod, beginTime, endTime, resp)
   return resp, err
}
```

完成 `AccessLog `拦截器的编写后，将其注册到 gRPC Server 中去，然后重新启动服务进行验证，在进行验证时，也就是调用 RPC 方法时，会输出两条日志记录。

```go
func runGrpcServer() *grpc.Server {
   // 新增拦截器相关代码
   //opts := []grpc.ServerOption{
   // grpc.UnaryInterceptor(HelloInterceptor),
   //}
   // 进行多拦截器的注册
   opts := []grpc.ServerOption{
      grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
         middleware.AccessLog,
      )),
   }

   s := grpc.NewServer(opts...)
   pb.RegisterTagServiceServer(s, server.NewTagServer())
   reflection.Register(s)
   return s
}
```

这时候可能有的读者会疑惑，为什么要在该 RPC 方法执行的前后各输出一条类似但不完全一致的日志呢，这有什么用意，为什么不是直接在 RPC 方法执行完毕后输出一条就好了，会不会有些重复?

其实不然，如果仅仅只在 RPC 方法执行完毕后才输出、落地日志，那么可以来假设两个例子：

+ 这个 RPC 方法在执行遇到了一些意外情况，执行了很久，也不知道什么时候返回（无其它措施的情况下）。

+ 在执行过程中因极端情况出现了 OOM，RPC 方法未执行完毕，就被系统杀掉了。

这两个例子的情况可能会造成什么问题呢，一般来讲，会去看日志，基本是因为目前应用系统已经出现了问题，那么第一种情况，就非常常见，如果只打 RPC 方法执行完毕后的日志，单看日志，可能会压根就没有所需要的访问日志，因为它还在执行中；而第二种情况，就根本上也没有达到完成。

那么从结果上来讲，日志的部分缺失有可能会导致误判当前事故的原因，影响全链路追踪，需要花费更多的精力去排查。

##### 错误日志

打开 `server_interceptor.go `文件，新增普通错误记录的日志拦截器，代码如下：

```go
func ErrorLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
   handler grpc.UnaryHandler) (interface{}, error) {
   resp, err := handler(ctx, req)
   if err != nil {
      errLog := "error log: method: %s, code: %v, message: %v, details: %v"
      s := errcode.FromError(err)
      log.Printf(errLog, info.FullMethod, s.Code(), s.Err().Error(), s.Details())
   }
   return resp, err
}
```

同之前的拦截器一样，编写后进行注册，该拦截器可以针对所有 RPC 方法的错误返回进行记录，便于对 error 级别的错误进行统一规管和观察。

```bash
2022/05/23 21:56:42 access request log: method: /proto.TagService/GetTagList, begin_time: 1653314202, request: 
2022/05/23 21:56:44 error log: method: /proto.TagService/GetTagList, code: Unknown, message: rpc error: code = Unknown desc = 获取标签列表失败, details: [code:20010001 message:"\350\216\267\345\217\226\346\240\207\347\255\276\345\210\227\350\241\250\345\244\261\350\264\245" ]
2022/05/23 21:56:44 access response log: method: /proto.TagService/GetTagList, begin_time: 1653314202, end_time: 1653314204, response: <nil>
```

#### b. 异常捕获

接下来针对异常进行捕获处理，再开始编写拦截器之前，现在`GetTagList`方法中加入`panic`语句，模拟抛出异常的情况，再重新运行服务，观察调用结果。

```go
func (t *TagServer) GetTagList(ctx context.Context, r *pb.GetTagListRequest) (*pb.GetTagListResponse, error) {
   panic("抛出异常！")
   ...
}
```

```bash
2022/05/23 21:59:18 access request log: method: /proto.TagService/GetTagList, begin_time: 1653314358, request: 
panic: 抛出异常！

goroutine 53 [running]:
ch03/server.(*TagServer).GetTagList(0x70, {0x14, 0x70}, 0x14)
...
```

服务直接因为异常抛出而中断了，这意味着该服务无法再提供相应，因此，为了解决这个问题，需要新增一个自定义的异常捕获拦截器，修改 `middelware`目录下的`server_interceptor.go`文件。

```go
func Recovery(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
   handler grpc.UnaryHandler) (interface{}, error) {
   defer func() {
      if e := recover(); e != nil {
         recoveryLog := "recovery log: method: %s, message: %v, stack: %s"
         /*
            debug.Stack()
            Stack 返回调用它的 goroutine 的格式化堆栈跟踪。
            它使用足够大的缓冲区调用 runtime.Stack 来捕获整个跟踪。
         */
         log.Printf(recoveryLog, info.FullMethod, e, string(debug.Stack()[:]))
      }
   }()

   return handler(ctx, req)
}
```

同之前的拦截器一样，编写后进行注册，该拦截器可以针对所有 RPC 方法所抛出的异常进行捕抓和记录，确保不会因为未知的 `panic` 语句的执行导致整个服务中断，在实际项目的应用中，可以根据公司内的可观察性的技术栈情况，进行一些定制化的处理，那么它就会更加的完善。

```bash
2022/05/23 22:11:00 access request log: method: /proto.TagService/GetTagList, begin_time: 1653315060, request: 
2022/05/23 22:11:00 recovery log: method: /proto.TagService/GetTagList, message: 抛出异常！, stack: goroutine 27 [running]:
runtime/debug.Stack()
	C:/Users/zyc/sdk/go1.17.2/src/runtime/debug/stack.go:24 +0x65
ch03/internal/middleware.Recovery.func1()
...
2022/05/23 22:11:00 access response log: method: /proto.TagService/GetTagList, begin_time: 1653315060, end_time: 1653315060, response: <nil>
```

### 7. 客户端-常用拦截器

在项目的 `internal/middleware` 目录下新建存储客户端拦截器的` client_interceptor.go` 文件，针对一些常用场景编写一些客户端拦截器。另外后续的客户端拦截器相关注册行为是在调用 `grpc.Dial` 或 `grpc.DialContext` 前通过 `DialOption` 配置选项进行注册的，例如：

```go
    var opts []grpc.DialOption
    opts = append(opts, grpc.WithUnaryInterceptor(
        grpc_middleware.ChainUnaryClient(
            middleware.XXXXX(),
        ),
    ))
    opts = append(opts, grpc.WithStreamInterceptor(
        grpc_middleware.ChainStreamClient(
            middleware.XXXXX(),
        ),
    ))
    clientConn, err := grpc.DialContext(ctx, target, opts...)
    ...
```

#### a. 超时控制（上下文）

超时时间的设置和适当控制，是在微服务架构中非常重要的一个保全项。

假设一个应用场景，有多个服务，他们分别是 A、B、C、D，他们之间是最简单的关联依赖，也就是` A=>B=>C=>D`。在某一天，有一个需求上线了，修改的代码内容正好就是与服务 D 相关的，恰好这个需求就对应着一轮业务高峰的使用，但突然发现不知道为什么，服务 A、B、C、D 全部都出现了响应缓慢，整体来看，开始出现应用系统雪崩….这到底是怎么了？

从根本上来讲，是服务 D 出现了问题，所导致的这一系列上下游服务出现连锁反应，因为在服务调用中默认没有设置超时时间，或者所设置的超时时间过长，都会导致多服务下的整个调用链雪崩，导致非常严重的事故，因此任何调用的默认超时时间的设置是非常有必要的，在 gRPC 中更是强调 `TL;DR（Too long, Don’t read）`并建议始终设定截止日期。

因此在本节将针对 RPC 的内部调用设置默认的超时控制，在 `client_interceptor.go` 文件下，新增如下代码：

```go
package middleware

import (
   "context"
   "google.golang.org/grpc"
   "time"
)

func defaultContextTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
   var cancel context.CancelFunc
   // Deadline() 未设置截止日期时，返回 ok==false
   if _, ok := ctx.Deadline(); !ok {
      defaultTimeout := 60 * time.Second
      // context.WithTimeout() 设置默认超时时间 60s
      ctx, cancel = context.WithTimeout(ctx, defaultTimeout)
   }

   return ctx, cancel
}

// 一元调用的客户端拦截器
func UnaryContextTimeout() grpc.UnaryClientInterceptor {
   return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
      ctx, cancel := defaultContextTimeout(ctx)
      if cancel != nil {
         defer cancel()
      }
      return invoker(ctx, method, req, reply, cc, opts...)
   }
}

// 流式调用的客户端拦截器
func StreamContextTimeout() grpc.StreamClientInterceptor {
   return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
      ctx, cancel := defaultContextTimeout(ctx)
      if cancel != nil {
         defer cancel()
      }

      return streamer(ctx, desc, cc, method, opts...)
   }
}
```

在上述代码中，通过对传入的` context `调用 `ctx.Deadline` 方法进行检查，若未设置截止时间的话，其将会返回 `false`，那么就会对其调用 `context.WithTimeout` 方法设置默认超时时间为 60 秒（该超时时间设置是针对整条调用链路的，若需要另外调整，可在应用代码中再自行调整）。接下来分别对 gRPC 的一元调用和流式调用编写对应的客户端拦截器。在编写完拦截器后，在进行 RPC 内调前进行注册就可以生效了。

#### b. 重试操作

在整体的服务运行中，偶尔会出现一些“奇奇怪怪”的网络波动、流量限制、服务器资源突现异常（但很快下滑），需要稍后访问的情况，这时候常常需要采用一些退避策略，稍作等待后进行二次重试，确保应用程序的最终成功，因此对于 gRPC 客户端来讲，一个基本的重试是必要的。如果没有定制化需求的话，可以直接采用 gRPC 生态圈中的 `grpc_retry` 拦截器实现基本的重试功能，如下：

```go
    var opts []grpc.DialOption
    opts = append(opts, grpc.WithUnaryInterceptor(
        grpc_middleware.ChainUnaryClient(
            grpc_retry.UnaryClientInterceptor(
                grpc_retry.WithMax(2),
                grpc_retry.WithCodes(
                  codes.Unknown,
                  codes.Internal,
                  codes.DeadlineExceeded,
                ),
            ),
        ),
    ))
    ...
```

在上述 `grpc_retry` 拦截器中，设置了最大重试次数为 2 次，仅针对 gRPC 错误码为 `Unknown`、`Internal`、`DeadlineExceeded` 的情况。

这里需要注意的第一点是，它确定是否需要重试的维度是以错误码为标准，因此做服务重试的第一点，就是需要在设计微服务的应用时，明确其状态码的规则，确保多服务的状态码的标准是一致的（可通过基础框架、公共代码库等方式落地），另外第二点是要尽可能的保证接口设计是幂等的，保证允许重试，也不会造成灾难性的问题，例如：重复扣库存。

```go
// 设置客户端拦截器
opts := []grpc.DialOption{
   grpc.WithUnaryInterceptor(
      grpc_middleware.ChainUnaryClient(
         middleware.UnaryContextTimeout(),
         grpc_retry.UnaryClientInterceptor(
            grpc_retry.WithMax(2),
            grpc_retry.WithCodes(
               codes.Unknown,
               codes.Internal,
               codes.DeadlineExceeded,
            ),
         ),
      )),
   grpc.WithStreamInterceptor(
      grpc_middleware.ChainStreamClient(
         middleware.StreamContextTimeout())),
}

// 传入配置信息到 DialContext() , GetClientConn() 中调用了 DialContext()
clientConn, _ := GetClientConn(ctx, "localhost:8004", opts)
```

### 8. 实战演练

在刚刚的超时控制的拦截器中，完善了默认的超时控制，那本项目的系统中，有没有类似的风险，那当然是有的，而且还是在 Go 语言编程中非常经典的问题。在实现 gRPC Server 的` GetTagList `方法时，数据源是来自第二章的博客后端应用（blog-service），如下：

```go
func (t *TagServer) GetTagList(ctx context.Context, r *pb.GetTagListRequest) (*pb.GetTagListResponse, error) {
   // localhost：不通过网卡传输，不受网络防火墙和网卡相关的限制。
   // 127.0.0.1：通过网卡传输，依赖网卡，并受到网卡和防火墙相关的限制。
   api := bapi.NewAPI("http://127.0.0.1:8000")
   body, err := api.GetTagList(ctx, r.GetName())
   if err != nil {
      // 填入业务错误码
      return nil, errcode.TogRPCError(errcode.ErrorGetTagListFail)
   }

   tagList := pb.GetTagListResponse{}
   err = json.Unmarshal(body, &tagList)
   if err != nil {
      // 填入错误码
      return nil, errcode.TogRPCError(errcode.Fail)
   }

   return &tagList, nil
}
```

来模拟一下，假设这个博客后端的应用，出现了问题，假死，持续不返回，如下（通过休眠来模拟）：

```go
func (t Tag) List(c *gin.Context) {
    time.Sleep(time.Hour)
    ...
}
```

打开第二章的博客后端应用，对“获取标签列表接口”新增长时间的休眠，接着调用 gRPC Server 的接口，例如：当 gRPC Server 起的是 8000 端口，则调用 `http://127.0.0.1:8000/api/v1/tags`。

会看到 gRPC Server 中输出的访问日志如下：

```bash
2022/05/23 22:46:23 access request log: method: /proto.TagService/GetTagList, begin_time: 1653317183, request: 

// 没有下文，等了很久，response log 迟迟没有出现。
```

可以发现并没有响应结果，只输出了访问时间中的` request log`，这时候可以把这段请求一直挂着，无论多久它都不会返回，直至休眠时间结束。

考虑一下实际场景，一般会用到 HTTP API，基本上都是因为依赖第三方接口。那假设这个第三方接口，出现了问题，也就是接口响应极度缓慢，甚至假死，没有任何响应。但是，应用是正常的，那么流量就会不断地打进应用中，这就会形成一个恶性循环，阻塞等待的协程会越来越多，开销越来越大，最终就会导致上游服务出现问题，那么这个下游服务也会逐渐崩溃，最终形成连锁反应。

#### a. 原因

为什么一个休眠会带来那么大的问题呢，可以看看 `HTTP API SDK `的代码，如下：

```go
func (a *API) httpGet(ctx context.Context, path string) ([]byte, error) {
    resp, err := http.Get(fmt.Sprintf("%s/%s", a.URL, path))
    ...
}
```

默认使用的 `http.Get` 方法，其内部源码：

```go
func Get(url string) (resp *Response, err error) {
    return DefaultClient.Get(url)
}
```

实际上它使用的是标准库中预定的包全局变量` DefaultClient`，而` DefaultClient` 的 Timeout 的默认值是零值，相当于是 0，那么当 Timeout 值为 0 时，**默认认为是没有任何超时时间限制的，也就是会无限等待，**直至响应为止，这就是其出现问题的根本原因之一。

#### b. 解决方案

那么针对现在这个问题，有至少两种解决方法，分别是自定义` HTTPClient`，又或是通过超时控制来解决这个问题，如下：

```go
func (a *API) httpGet(ctx context.Context, path string) ([]byte, error) {
    // 自定义 HTTPClient                                                               
    resp, err := ctxhttp.Get(ctx, http.DefaultClient, fmt.Sprintf("%s/%s", a.URL, path))
    ...
}
```

将 `http.Get` 方法修改为 `ctxhttp.Get` 方法，将上下文（`ctx`）传入到该方法中，那么它就会受到上下文的超时控制了。但是这种方法有一个前提，那就是客户端在调用时需要将超时控制的拦截器注册进去，如下：

```go 
func main() {
    ctx := context.Background()
    clientConn, err := GetClientConn(ctx, "tag-service", []grpc.DialOption{grpc.WithUnaryInterceptor(
        grpc_middleware.ChainUnaryClient(middleware.UnaryContextTimeout()),
    )})
    ...
}
```

再次进行验证，如下：

```bash
2022/05/23 22:55:22 access request log: method: /proto.TagService/GetTagList, begin_time: 1653317722, request: name:"Go" 
# 客户端上下文中设置的超时时间为 60s
2022/05/23 22:56:22 error log: method: /proto.TagService/GetTagList, code: Unknown, message: rpc error: code = Unknown desc = 获取标签列表失败, details: [code:20010001 message:"\350\216\267\345\217\226\346\240\207\347\255\276\345\210\227\350\241\250\345\244\261\350\264\245" ]
2022/05/23 22:56:22 access response log: method: /proto.TagService/GetTagList, begin_time: 1653317722, end_time: 1653317782, response: <nil>
```

在到达截止时间后，客户端将自动断开，提示 `DeadlineExceeded`，那么从结果上来讲，当上游服务出现问题时，当前服务再去调用它，也不会受到过多的影响，因为通过超时时间进行了及时的止损，因此默认超时时间的设置和设置多少是非常有意义和考究的。

但是此时此刻，服务端本身可能还在无限的阻塞中，客户端断开的仅仅只是自己，因此服务端本身也建议设置默认的最大执行时间，以确保最大可用性和避免存在忘记设置超时控制的客户端所带来的最坏情况。

#### c. 如何发现

最简单的方式有两种，分别是通过日志和链路追踪发现，假设是上述提到的问题，在所打的访问日志中，它只会返回 `request log`，而不会返回` response log`。那如果是用分布式链路追踪系统，会非常明显的出现某个 Span 的调用链会耗时特别久，这就是一个危险的味道。更甚至可以通过对这些指标数据进行分析，当出现该类情况时，直接通过分析确定是否要报警、自愈，那将更妥当。

## 九、Metadata 和 RPC 自定义认证

#### 1. Metadata 介绍

在 HTTP/1.1 中，常常通过直接操纵 Header 来传递数据，而对于 gRPC 来讲，它基于 HTTP/2 协议，本质上也可是通过 Header 来进行传递，但一般不会直接的去操纵它，而是通过 gRPC 中的 metadata 来进行调用过程中的数据传递和操纵。但需要注意的是，metadata 的使用需要所使用的库进行支持，并不能像 HTTP/1.1 那样自行去 Header 去取。

在 gRPC 中，Metadata 实际上就是一个 map 结构，其原型如下：

```go
type MD map[string][]string
```

是一个字符串与字符串切片的映射结构。

#### a. 创建 metadata

在 `google.golang.org/grpc/metadata` 中分别提供了两个方法来创建 metadata，第一种是 `metadata.New` 方法，如下：

```go
metadata.New(map[string]string{"go": "programming", "tour": "book"})
```

使用 New 方法所创建的 metadata，将会直接被转换为对应的 MD 结构，参考结果如下：

```go
go:   []string{"programming"}
tour: []string{"book"}
```

第二种是 `metadata.Pairs` 方法，如下：

```go
metadata.Pairs(
    "go", "programming",
    "tour", "book",
    "go", "eddycjy",
)
```

使用 Pairs 方法所创建的 metadata，将会以奇数来配对，并且所有的 Key 都会被默认转为小写，若出现同名的 Key，将会追加到对应 Key 的切片（slice）上，参考结果如下：

```go
go:   []string{"programming", "eddycjy"}
tour: []string{"book"}
```

#### b. 设置/获取 metadata

```go
ctx := context.Background()
md := metadata.New(map[string]string{"go": "programming", "tour": "book"})

newCtx1 := metadata.NewIncomingContext(ctx, md) 
newCtx2 := metadata.NewOutgoingContext(ctx, md)
```

在 gRPC 中对于 metadata 进行了区别，分为了传入和传出用的 metadata，这是为了防止 metadata 从入站 RPC 转发到其出站 RPC 的情况（详见 issues #1148），针对此提供了两种方法来分别进行设置，如下：

- `NewIncomingContext`：创建一个附加了所传入的 md 新上下文，仅供自身的 gRPC 服务端内部使用。
- `NewOutgoingContext`：创建一个附加了传出 md 的新上下文，**可供外部的 gRPC 客户端、服务端使用。**

因此相对的在 metadata 的获取上，也区分了两种方法，分别是 FromIncomingContext 和 NewOutgoingContext，与设置的方法所相对应的含义，如下：

```go
md1, _ := metadata.FromIncomingContext(ctx)
md2, _ := metadata.FromOutgoingContext(ctx)
```

那么总的来说，这两种方法在实现上有没有什么区别呢，可以一起深入看看：

```go
type mdIncomingKey struct{}
type mdOutgoingKey struct{}

func NewIncomingContext(ctx context.Context, md MD) context.Context {
    return context.WithValue(ctx, mdIncomingKey{}, md)
}

func NewOutgoingContext(ctx context.Context, md MD) context.Context {
    return context.WithValue(ctx, mdOutgoingKey{}, rawMD{md: md})
}
```

实际上主要是在内部进行了 Key 的区分，以所指定的 Key 来读取相对应的 metadata，以防造成脏读，其在实现逻辑上本质上并没有太大的区别。另外可以看到，其对 Key 的设置，是用一个结构体去定义的，这是 Go 语言官方一直在推荐的写法。

#### c. 实际使用场景

在上面已经介绍了关键的 metadata 以及其相对的 `IncomingContext`、`OutgoingContext` 类别的相关方法，但在实际的使用中，仍然常常会有开发人员用错，然后出现了疑惑，最后无奈只能调试半天，才恍然大悟。

假设现在有一个 ServiceA 作为服务端，然后有一个 Client 去调用 ServiceA，想传入自定义的 metadata 信息，那该怎么写才合适，流程图如下：

![image](https://raw.githubusercontent.com/tonshz/test/master/img/202205242110157.jpeg)

在常规情况下，在 ServiceA 的服务端，应当使用 `metadata.FromIncomingContext` 方法进行读取，如下：

```go
func (t *TagServer) GetTagList(ctx context.Context, r *pb.GetTagListRequest) (*pb.GetTagListReply, error) {
    md, _ := metadata.FromIncomingContext(ctx)
    log.Printf("md: %+v", md)
    ...
}
```

而在 Client，应当使用 `metadata.AppendToOutgoingContext` 方法，如下：

```go
func main() {
    ctx := context.Background()
    newCtx := metadata.AppendToOutgoingContext(ctx, "eddycjy", "Go 语言编程之旅")
    clientConn, _ := GetClientConn(newCtx, ...)
    defer clientConn.Close()
    tagServiceClient := pb.NewTagServiceClient(clientConn)
    resp, _ := tagServiceClient.GetTagList(newCtx, &pb.GetTagListRequest{Name: "Go"})
    ...
}
```

这里需要注意一点，**在新增 metadata 信息时，务必使用 Append 类别的方法，否则如果直接 New 一个全新的 md，将会导致原有的 metadata 信息丢失**（除非确定希望得到这样的结果）。

### 2. Metadata 是如何传递的

在上小节中，已经知道 metadata 其实是存储在 context 之中的，那么 context 中的数据又是承载在哪里呢，继续对前面的 gRPC 调用例子进行调整，将已经传入 metadata 的 context 设置到对应的 RPC 方法调用上，代码如下：

```go
func main() {
    ctx := context.Background()
    md := metadata.New(map[string]string{"go": "programming", "tour": "book"})
    newCtx := metadata.NewOutgoingContext(ctx, md)
    clientConn, err := GetClientConn(newCtx, "localhost:8004", nil)
    if err != nil {
        log.Fatalf("err: %v", err)
    }
    defer clientConn.Close()
    tagServiceClient := pb.NewTagServiceClient(clientConn)
    resp, err := tagServiceClient.GetTagList(newCtx, &pb.GetTagListRequest{Name: "Go"})
    ...
}
...
```

再重新查看抓包工具的结果（在本地测试未抓包到内容，只有 tcp 协议的报文，没有 grpc 报文）：

![image](https://raw.githubusercontent.com/tonshz/test/master/img/202205242203826.jpeg)

显然，所传入的 `"go": "programming", "tour": "book"` 是在 Header 中进行传播的。

### 3. 对 RPC 方法做自定义认证

在实际需求中，有时候会需要对某些模块的 RPC 方法做特殊认证或校验，这时候可以利用 gRPC 所提供的 Token 接口，如下：

```go
type PerRPCCredentials interface {
    GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)
    RequireTransportSecurity() bool
}
```

在 gRPC 中所提供的 `PerRPCCredentials`，它就是本节的主角，是 gRPC 默认提供用于自定义认证 Token 的接口，它的作用是将所需的安全认证信息添加到每个 RPC 方法的上下文中。其包含两个接口方法，如下：

- GetRequestMetadata：获取当前请求认证所需的元数据（metadata）。
- RequireTransportSecurity：是否需要基于 TLS 认证进行安全传输。

#### a. 客户端

打开先前章节编写的 gRPC 调用的代码（也就是 gRPC 客户端的角色），那么在客户端的重点在于实现 `type PerRPCCredentials interface` 所需的接口方法，代码如下：

```go
package main

import (
	pb "ch03/proto"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type Auth struct {
	Appkey    string
	AppSecret string
}

func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_key": a.Appkey, "app_secret": a.AppSecret}, nil
}

func (a *Auth) RequireTransportSecurity() bool {
	return false
}
func main() {
	auth := Auth{
		Appkey:    "admin",
		AppSecret: "go-learning",
	}
	ctx := context.Background()
    // 调用 grpc.WithPerRPCCredentials() 进行注册
    // grpc.WithPerRPCCredentials() 设置凭据并在每个出站 RPC 上放置身份验证状态
	opts := []grpc.DialOption{grpc.WithPerRPCCredentials(&auth)}
	clientConn, _ := GetClientConn(ctx, "localhost:8004", opts)
	defer clientConn.Close()

	tagServiceClient := pb.NewTagServiceClient(clientConn)
	resp, _ := tagServiceClient.GetTagList(ctx, &pb.GetTagListRequest{Name: "Golang"})
	log.Printf("resp: %v", resp)
}

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	// grpc.WithInsecure() 已弃用，作用为跳过对服务器证书的验证，此时客户端和服务端会使用明文通信
	// 使用 WithTransportCredentials 和 insecure.NewCredentials() 代替
	//opts = append(opts, grp c.WithInsecure())
	// insecure.NewCredentials 返回一个禁用传输安全的凭据
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	/*
		grpc.DialContext() 创建到给定目标的客户端连接。
		默认情况下，它是一个非阻塞拨号（该功能不会等待建立连接，并且连接发生在后台）。
		要使其成为阻塞拨号，请使用 WithBlock() 拨号选项。
	*/
	return grpc.DialContext(ctx, target, opts...)
}
```

在上述代码中，声明了 Auth 结构体，并实现了所需的两个接口方法，最后在 `DialOption` 配置中调用 `grpc.WithPerRPCCredentials` 方法进行了注册。

#### b. 服务端

客户端的校验数据已经传过来了，接下来需要修改先前的服务端代码，对其进行 Token 校验，如下：

```go
package server

import (
   "ch03/pkg/bapi"
   "ch03/pkg/errcode"
   pb "ch03/proto"
   "context"
   "encoding/json"
   "google.golang.org/grpc/metadata"
   "log"
)

type TagServer struct {
   auth *Auth
}

type Auth struct{}

func (a *Auth) GetAppKey() string {
   return "admin"
}

func (a *Auth) GetAppSecret() string {
   return "go-learning"
}

func (a *Auth) Check(ctx context.Context) error {
   md, _ := metadata.FromIncomingContext(ctx)

   var appKey, appSecret string
   if value, ok := md["app_key"]; ok {
      appKey = value[0]
   }
   if value, ok := md["app_secret"]; ok {
      appSecret = value[0]
   }

   log.Printf("test: %s, %s", appKey, appSecret)
   if appKey != a.GetAppKey() || appSecret != a.GetAppSecret() {
      return errcode.TogRPCError(errcode.Unauthorized)
   }
   return nil
}

func NewTagServer() *TagServer {
   return &TagServer{}
}

func (t *TagServer) GetTagList(ctx context.Context, r *pb.GetTagListRequest) (*pb.GetTagListResponse, error) {
   if err := t.auth.Check(ctx); err != nil {
      return nil, err
   }

   // localhost：不通过网卡传输，不受网络防火墙和网卡相关的限制。
   // 127.0.0.1：通过网卡传输，依赖网卡，并受到网卡和防火墙相关的限制。
   api := bapi.NewAPI("http://127.0.0.1:8000")
   body, err := api.GetTagList(ctx, r.GetName())
   if err != nil {
      // 填入业务错误码
      return nil, errcode.TogRPCError(errcode.ErrorGetTagListFail)
   }

   tagList := pb.GetTagListResponse{}
   err = json.Unmarshal(body, &tagList)
   if err != nil {
      // 填入错误码
      return nil, errcode.TogRPCError(errcode.Fail)
   }

   return &tagList, nil
}
```

上述代码实际就是调用 `metadata.FromIncomingContext` 从上下文中获取 metadata，再在不同的 RPC 方法中进行认证检查就可以了。

```go
package bapi

import (
   "context"
   "encoding/json"
   "fmt"
   "golang.org/x/net/context/ctxhttp"
   "io/ioutil"
   "log"
   "net/http"
   "strings"
)

const (
   APP_KEY    = "admin"
   APP_SECRET = "go-learning"
)

type AccessToken struct {
   Token string `json:"token"`
}

type API struct {
   URL string
}
type AuthParams struct {
   AppKey    string `json:"app_key"`
   AppSecret string `json:"app_secret"`
}

func NewAPI(url string) *API {
   return &API{URL: url}
}

// 获取所有 API 请求都需要带上的 token
func (a *API) getAccessToken(ctx context.Context) (string, error) {
   body, err := a.httpPost(ctx, "/auth", APP_KEY, APP_SECRET)
   if err != nil {
      return "", err
   }

   var accessToken AccessToken
   _ = json.Unmarshal(body, &accessToken)
   return accessToken.Token, nil
}

// 统一的 HTTP GET 请求方法
func (a *API) httpGet(ctx context.Context, token string, path string) ([]byte, error) {
   // 自定义 HTTPClient
   req, _ := http.NewRequest("GET", a.URL+path, nil)
   req.Header.Set("token", token)
   resp, err := ctxhttp.Do(ctx, http.DefaultClient, req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   body, _ := ioutil.ReadAll(resp.Body)
   return body, nil
}

// 统一的 HTTP POST 请求方法
func (a *API) httpPost(ctx context.Context, path string, appKey string, appSecret string) ([]byte, error) {
   //resp, err := ctxhttp.Post(ctx, http.DefaultClient, a.URL+path, "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("app_key=%s&app_secret=%s", appKey, appSecret)))
   // 使用 json 传输数据
   reqParam, _ := json.Marshal(&AuthParams{appKey, appSecret})
   reqBody := strings.NewReader(string(reqParam))
   resp, err := ctxhttp.Post(ctx, http.DefaultClient, a.URL+path, "application/json", reqBody)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   body, _ := ioutil.ReadAll(resp.Body)
   return body, nil
}

// 具体的获取标签列表的方法实现
func (a *API) GetTagList(ctx context.Context, name string) ([]byte, error) {
   // 获取AccessToken
   token, err := a.getAccessToken(ctx)
   log.Printf("token: %s", token)
   if err != nil {
      return nil, err
   }

   body, err := a.httpGet(ctx, token, fmt.Sprintf("%s?name=%s", "/api/v1/tags", name))
   if err != nil {
      return nil, err
   }

   return body, nil
}
```

## 十、链路追踪

在前面的章节中，介绍了`gRPC`拦截器和 metadata 的使用，同时也提到了分布式链路追踪的重要性。在微服务体系中，在服务中注入链路追踪是非常重要的事情，一次 RPC 调用可能涉及多个服务，服务内有可能调用其他的 HTTP API、SQL、Redis 等。

![image-20220524234024466](https://raw.githubusercontent.com/tonshz/test/master/img/202205242340520.png)

在上图中，假设 B 服务出现了问题，它的调用链路是服务 A -> 服务 C  ->  服务 B -> 服务 D，还是服务 C  ->  服务 B -> 服务 D，又或是其他呢？不同的调用链路的业务场景是不一样的，入参的模式也不一样。故在微服务这种复杂分布式场景下，注入链路追踪是很有必要的。项目基于`gRPC+Jaeger`实现链路追踪。

### 1. 注入追踪信息

做链路追踪的基本条件是要注入追踪信息，而最简单的方法就是使用服务端和客户端拦截器组成完整的链路信息，具体如下：

+ 服务端拦截器：从 metadata 中提取链路信息，将其设置并追加到服务端的调用上下文中。也就是说，如果发现本次调用没有上一级的链路信息，那么它将会生成对应的父级信息，自己成为父级；如果发现本次调用存在上一级的链路信息，那么它将会根据上一级链路信息进行设置，成为其子级。
+ 客户端拦截器：从调用的上下文中提取链路信息，并将其作为 metadata 追加到 RPC 调用中。

### 2. 初始化 Jaeger

借助 OpenTracing API 和 Jaeger Client 实现与追踪系统的对接，在项目根目录中执行如下命令：

```bash
$ go get -u github.com/opentracing/opentracing-go
$ go get -u github.com/uber/jaeger-client-go
```

### 3. metadata 的读取和设置

在 OpenTracing  中， 可以指定 SpanContexts 的三种传输表示方式。

```go
type BuiltinFormat byte

const (
	Binary BuiltinFormat = iota

	TextMap

	HTTPHeaders
)

// TextMapWriter is the Inject() carrier for the TextMap builtin format. With
// it, the caller can encode a SpanContext for propagation as entries in a map
// of unicode strings.
type TextMapWriter interface {
	// Set a key:value pair to the carrier. Multiple calls to Set() for the
	// same key leads to undefined behavior.
	//
	// NOTE: The backing store for the TextMapWriter may contain data unrelated
	// to SpanContext. As such, Inject() and Extract() implementations that
	// call the TextMapWriter and TextMapReader interfaces must agree on a
	// prefix or other convention to distinguish their own key:value pairs.
	Set(key, val string)
}

// TextMapReader is the Extract() carrier for the TextMap builtin format. With it,
// the caller can decode a propagated SpanContext as entries in a map of
// unicode strings.
type TextMapReader interface {
	// ForeachKey returns TextMap contents via repeated calls to the `handler`
	// function. If any call to `handler` returns a non-nil error, ForeachKey
	// terminates and returns that error.
	//
	// NOTE: The backing store for the TextMapReader may contain data unrelated
	// to SpanContext. As such, Inject() and Extract() implementations that
	// call the TextMapWriter and TextMapReader interfaces must agree on a
	// prefix or other convention to distinguish their own key:value pairs.
	//
	// The "foreach" callback pattern reduces unnecessary copying in some cases
	// and also allows implementations to hold locks while the map is read.
	ForeachKey(handler func(key, val string) error) error
}
```

+ Binary: 不透明的二进制数据。
+ TextMap: 键值字符串对。
+ HTTPHeaders: HTTP Header 格式的字符串。

在项目的 pkg 目录下新建 metatext 目录，并创建 `metadata.go`文件。

```go
package metatext

import (
   "google.golang.org/grpc/metadata"
   "strings"
)

type MetadataTextMap struct {
   metadata.MD
}

func (m MetadataTextMap) ForeachKey(handler func(key, val string) error) error {
   for k, vs := range m.MD {
      for _, v := range vs {
         if err := handler(k, v); err != nil {
            return err
         }
      }
   }
   return nil
}

func (m MetadataTextMap) Set(key, val string) {
   key = strings.ToLower(key)
   m.MD[key] = append(m.MD[key], val)
}
```

在上述代码中，基于 TextMap 模式，对照实现了 metadata 的设置和读取方法。

### 4. 服务端

在服务端需要注册链路追踪的拦截器，用来生成和设置链路信息。

在项目根目录下新建`global`目录，并在其下新建`tracer.go`文件。

```go
package global

import "github.com/opentracing/opentracing-go"

var (
   Tracer opentracing.Tracer
)
```

修改 `middleware` 下的 `server_interceptor.go`文件，添加链路追踪的拦截器。

```go
func ServerTracing(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
   handler grpc.UnaryHandler) (interface{}, error) {
   md, ok := metadata.FromIncomingContext(ctx)
   if !ok {
      md = metadata.New(nil)
   }
   parentSpanContext, _ := global.Tracer.Extract(opentracing.TextMap,
      metatext.MetadataTextMap{md})
   spanOpts := []opentracing.StartSpanOption{
      opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
      ext.SpanKindRPCServer,
      ext.RPCServerOption(parentSpanContext),
   }

   span := global.Tracer.StartSpan(info.FullMethod, spanOpts...)
   defer span.Finish()
   ctx = opentracing.ContextWithSpan(ctx, span)
   return handler(ctx, req)
}
```

在上述代码中，首先通过读取 RPC 方法传入的上下文信息，可以解析出 metadata。然后从给定的载体中解码出 SpanContext 实例，并创建和设置本次跨度的标签信息。最后，根据当前的快读返回一个新的 `context.Context`，一般后续使用。

接下来把拦截器注册到服务中，代码如下:

```go
func runGrpcServer() *grpc.Server {
   // 新增拦截器相关代码
   //opts := []grpc.ServerOption{
   // grpc.UnaryInterceptor(HelloInterceptor),
   //}
   // 进行多拦截器的注册
   opts := []grpc.ServerOption{
      grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
         middleware.AccessLog,
         middleware.ErrorLog,
         middleware.Recovery,
         // 注册链路追踪拦截器
         middleware.ServerTracing,
      )),
   }

   s := grpc.NewServer(opts...)
   pb.RegisterTagServiceServer(s, server.NewTagServer())
   reflection.Register(s)
   return s
}
```

### 5. 客户端

在客户端，调用注册链路追踪中的拦截器，获取和设置链路信息，在 `client_interceptor.go`文件中，新增如下代码。

```go
package middleware

import (
   "ch03/global"
   "ch03/pkg/metatext"
   "context"
   "github.com/opentracing/opentracing-go"
   "github.com/opentracing/opentracing-go/ext"
   "google.golang.org/grpc"
   "google.golang.org/grpc/metadata"
   "time"
)

func defaultContextTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
   var cancel context.CancelFunc
   // Deadline() 未设置截止日期时，返回 ok==false
   if _, ok := ctx.Deadline(); !ok {
      defaultTimeout := 60 * time.Second
      // context.WithTimeout() 设置默认超时时间 60s
      ctx, cancel = context.WithTimeout(ctx, defaultTimeout)
   }

   return ctx, cancel
}

// 一元调用的客户端拦截器
func UnaryContextTimeout() grpc.UnaryClientInterceptor {
   return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
      ctx, cancel := defaultContextTimeout(ctx)
      if cancel != nil {
         defer cancel()
      }
      return invoker(ctx, method, req, reply, cc, opts...)
   }
}

// 流式调用的客户端拦截器
func StreamContextTimeout() grpc.StreamClientInterceptor {
   return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
      ctx, cancel := defaultContextTimeout(ctx)
      if cancel != nil {
         defer cancel()
      }

      return streamer(ctx, desc, cc, method, opts...)
   }
}

// 设置链路追踪拦截器
func ClientTracing() grpc.UnaryClientInterceptor {
   return func(ctx context.Context, method string, req, reply interface{},
      cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
      var parentCtx opentracing.SpanContext
      var spanOpts []opentracing.StartSpanOption
      var parentSpan = opentracing.SpanFromContext(ctx)
      
      if parentSpan != nil {
         parentCtx = parentSpan.Context()
         spanOpts = append(spanOpts, opentracing.ChildOf(parentCtx))
      }
      spanOpts = append(spanOpts, []opentracing.StartSpanOption{
         opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
         ext.SpanKindRPCClient,
      }...)
      
      span := global.Tracer.StartSpan(method, spanOpts...)
      defer span.Finish()
      
      md, ok := metadata.FromOutgoingContext(ctx)
      if !ok {
         md = metadata.New(nil)
      }
      _ = global.Tracer.Inject(span.Context(), opentracing.TextMap,
         metatext.MetadataTextMap{md})
      
      newCtx := opentracing.ContextWithSpan(metadata.NewOutgoingContext(ctx, md), span)
      return invoker(newCtx, method, req, reply, cc, opts...)
   }
}
```

在上述代码中，首先调用了 `opentracing.SpanFromContext()`来解析上下文信息，检查其是否包含上一级的跨度信息，若存在，则获取上一级的上下文信息，把它作为接下来本次跨度的父级，接下来就是常规的创建和设置本次跨度的标签信息，再对传出的 md 信息进行转换，把它设置到新的上下文信息中，以便后续再调用时使用。

在进行客户端内部调用时，只需将拦截器注册进去，就可以达到追踪的效果。

```go
package main

import (
   "ch03/internal/middleware"
   pb "ch03/proto"
   "context"
   grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
   grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
   "google.golang.org/grpc"
   "google.golang.org/grpc/codes"
   "google.golang.org/grpc/credentials/insecure"
   "log"
)

func main() {
   /*
      context.Background()
      返回一个非零的空上下文。它永远不会被取消，没有价值，也没有最后期限。
      它通常由主函数、初始化和测试使用，并作为传入请求的顶级上下文。
   */
   ctx := context.Background()

   // 设置客户端拦截器
   opts := []grpc.DialOption{
      grpc.WithUnaryInterceptor(
         grpc_middleware.ChainUnaryClient(
            middleware.UnaryContextTimeout(),
            grpc_retry.UnaryClientInterceptor(
               grpc_retry.WithMax(2),
               grpc_retry.WithCodes(
                  codes.Unknown,
                  codes.Internal,
                  codes.DeadlineExceeded,
               ),
            ),
            middleware.ClientTracing(),
         )),
      grpc.WithStreamInterceptor(
         grpc_middleware.ChainStreamClient(
            middleware.StreamContextTimeout())),
   }

   clientConn, _ := GetClientConn(ctx, "localhost:8004", opts)
   defer clientConn.Close()

   // 初始化指定 RPC Proto Service 的客户端实例对象
   tagServiceClient := pb.NewTagServiceClient(clientConn)
   // 发起指定 RPC 方法的调用
   resp, _ := tagServiceClient.GetTagList(ctx, &pb.GetTagListRequest{Name: "Go"})

   log.Printf("resp: %v", resp)
}

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
   // grpc.WithInsecure() 已弃用，作用为跳过对服务器证书的验证，此时客户端和服务端会使用明文通信
   // 使用 WithTransportCredentials 和 insecure.NewCredentials() 代替
   //opts = append(opts, grp c.WithInsecure())
   // insecure.NewCredentials 返回一个禁用传输安全的凭据
   opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

   /*
      grpc.DialContext() 创建到给定目标的客户端连接。
      默认情况下，它是一个非阻塞拨号（该功能不会等待建立连接，并且连接发生在后台）。
      要使其成为阻塞拨号，请使用 WithBlock() 拨号选项。
   */
   return grpc.DialContext(ctx, target, opts...)
}
```

### 6. 实现 HTTP 追踪

在获取标签列表的 RPC 方法中，数据源实际上调用的是前面章节中的博客后端接口，既然可以追踪 SQL ,当然也可以追踪 HTTP 调用。修改`bapi`目录下的`api.go`文件中的`http.Get`方法。

```go
// 统一的 HTTP GET 请求方法
func (a *API) httpGet(ctx context.Context, token string, path string) ([]byte, error) {
   // 自定义 HTTPClient
   req, _ := http.NewRequest("GET", a.URL+path, nil)
   req.Header.Set("token", token)

   span, newCtx := opentracing.StartSpanFromContext(
      ctx, "HTTP GET: "+a.URL,
      opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
   )
   span.SetTag("url", a.URL+path)
   _ = opentracing.GlobalTracer().Inject(
      span.Context(),
      opentracing.HTTPHeaders,
      opentracing.HTTPHeadersCarrier(req.Header),
   )

   resp, err := ctxhttp.Do(newCtx, http.DefaultClient, req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   defer span.Finish()

   // 读取消息主体，在实际封装中可以将其抽离
   body, _ := ioutil.ReadAll(resp.Body)
   return body, nil
}
```

在上述代码中可以发现，HTTP 追踪和 RPC 追踪的内部设置逻辑都是基于一个模式的。即首先创建并设置当前跨度的信息和标签内容，需传入上下文信息，以保证链路完整性；然后传入附带信息，并将它设置到对应的链路信息上，最后进行调用，并返回新的上下文，以便后续使用。

### 7. 验证

#### a. 外部调用

通过调用 gRPC 服务中的 grpc-gateway 的 HTTP 接口来进行快速检查，用浏览器访问`http://127.0.0.1:8004/api/v1/tags`后，查看 Jaeger 的链路结果。使用下列命令运行 Jaeger Web UI（http://localhost:16686/）。

```bash
$ docker run -d --name Jaeger -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p 5775:5775/udp -p 6831:6831/udp -p 6832:6832/udp -p 5778:5778 -p 16686:16686 -p 14268:14268 -p 9411:9411 jaegertracing/all-in-one:1.16
```

![image-20220526001751072](https://raw.githubusercontent.com/tonshz/test/master/img/202205260017251.png)

可以看到左侧为 HTTP 调用，右侧为 HTTP 调用对应的链路信息展示。从整体来看，可以看到 tour-service 服务发起的 HTTP 调用被识别为了 blog-service，由于在 blog-service 中也在进行链路追踪，故也可以看到一些其中的打点追踪行为。

需要在 `main.go`与`client.go`中均添加下列代码：

```go
func setupTracer() error {
   var err error
   jaegerTracer, _, err := tracer.NewJaegerTracer("article-service", "127.0.0.1:6831")
   if err != nil {
      return err
   }
   global.Tracer = jaegerTracer
   return nil
}

func init() {
	err := setupTracer()
	if err != nil {
		log.Fatalf("init.setupTracer err: %v", err)
	}
}
```

并在``pkg`目录下创建 `tracer`目录，在其下新增`tracer.go`。

```go
package tracer

import (
   "io"
   "time"

   "github.com/opentracing/opentracing-go"
   "github.com/uber/jaeger-client-go/config"
)

func NewJaegerTracer(serviceName, agentHostPort string) (opentracing.Tracer, io.Closer, error) {
   cfg := &config.Configuration{
      ServiceName: serviceName,
      Sampler: &config.SamplerConfig{
         Type:  "const",
         Param: 1,
      },
      Reporter: &config.ReporterConfig{
         LogSpans:            true,
         BufferFlushInterval: 1 * time.Second,
         LocalAgentHostPort:  agentHostPort,
      },
   }
   tracer, closer, err := cfg.NewTracer()
   if err != nil {
      return nil, nil, err
   }
   opentracing.SetGlobalTracer(tracer)
   return tracer, closer, nil
}
```

#### b. 服务内调

在 RPC 方法 `GetTagList`中，对先前启动的服务进行内部调用。

![image-20220526004812434](https://raw.githubusercontent.com/tonshz/test/master/img/202205260050966.png)

可以看到其一共有三个部分的链路调用：article-service 调用了两次 tour-service，而 tour-service 又调用了一次 blog-service，最后在 blog-service 中又调用了 SQL 来查询数据。

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

## 十二、实现自定义的 protoc 插件

在开发 gRPC +`Protobuf` 的相关服务时，用到了许多与 `protoc `相关的插件来实现各种功能。

| 插件名称                  | 对应的命令           |
| ------------------------- | -------------------- |
| `protoc-gen-go`           | `--go_out`           |
| `protoc-gen-grpc-gateway` | `--grpc-gateway_out` |
| `protoc-gen-swagger`      | `--swagger_out`      |

并非所有插件都是由 `gRPC`或是`Protobuf`官方人员开发的，开发人员可以定制化开发一个新的`protoc`插件。

### 1. 插件的内部逻辑

通常来说，在实现自定义 `protoc`插件之前，需要分析既有的插件，看看它是怎么做的。首先先来看官方提供的`protoc-gen-go`插件的内部逻辑。

#### a. protoc-gen-gogo/main.go

```go
package main

import (
	"github.com/gogo/protobuf/vanity/command"
)

func main() {
    // Read() => Generate() => Write()
	command.Write(command.Generate(command.Read()))
}
```

#### b. vanity/command/command.go

```go
package command

import (
   "fmt"
   "go/format"
   "io/ioutil"
   "os"
   "strings"

   _ "github.com/gogo/protobuf/plugin/compare"
   _ "github.com/gogo/protobuf/plugin/defaultcheck"
   _ "github.com/gogo/protobuf/plugin/description"
   _ "github.com/gogo/protobuf/plugin/embedcheck"
   _ "github.com/gogo/protobuf/plugin/enumstringer"
   _ "github.com/gogo/protobuf/plugin/equal"
   _ "github.com/gogo/protobuf/plugin/face"
   _ "github.com/gogo/protobuf/plugin/gostring"
   _ "github.com/gogo/protobuf/plugin/marshalto"
   _ "github.com/gogo/protobuf/plugin/oneofcheck"
   _ "github.com/gogo/protobuf/plugin/populate"
   _ "github.com/gogo/protobuf/plugin/size"
   _ "github.com/gogo/protobuf/plugin/stringer"
   "github.com/gogo/protobuf/plugin/testgen"
   _ "github.com/gogo/protobuf/plugin/union"
   _ "github.com/gogo/protobuf/plugin/unmarshal"
   "github.com/gogo/protobuf/proto"
   "github.com/gogo/protobuf/protoc-gen-gogo/generator"
   _ "github.com/gogo/protobuf/protoc-gen-gogo/grpc"
   plugin "github.com/gogo/protobuf/protoc-gen-gogo/plugin"
)

// No.1
func Read() *plugin.CodeGeneratorRequest {
   // 1.创建默认的代码生成器
   g := generator.New()
    
   // 2.从标准输出中读取所需的 CodeGeneratorRequest 信息
   data, err := ioutil.ReadAll(os.Stdin)
   if err != nil {
      g.Error(err, "reading input")
   }

   //3.序列化所读取的 CodeGeneratorRequest
   if err := proto.Unmarshal(data, g.Request); err != nil {
      g.Error(err, "parsing input proto")
   }

   // 4.检查待生成文件的源文件数量，若数量为 0 则不需要生成
   if len(g.Request.FileToGenerate) == 0 {
      g.Fail("no files to generate")
   }
   return g.Request
}

// filenameSuffix replaces the .pb.go at the end of each filename.
func GeneratePlugin(req *plugin.CodeGeneratorRequest, p generator.Plugin, filenameSuffix string) *plugin.CodeGeneratorResponse {
   g := generator.New()
   g.Request = req
   if len(g.Request.FileToGenerate) == 0 {
      g.Fail("no files to generate")
   }

   g.CommandLineParameters(g.Request.GetParameter())

   g.WrapTypes()
   g.SetPackageNames()
   g.BuildTypeNameMap()
   g.GeneratePlugin(p)

   for i := 0; i < len(g.Response.File); i++ {
      g.Response.File[i].Name = proto.String(
         strings.Replace(*g.Response.File[i].Name, ".pb.go", filenameSuffix, -1),
      )
   }
   if err := goformat(g.Response); err != nil {
      g.Error(err)
   }
   return g.Response
}

func goformat(resp *plugin.CodeGeneratorResponse) error {
   for i := 0; i < len(resp.File); i++ {
      formatted, err := format.Source([]byte(resp.File[i].GetContent()))
      if err != nil {
         return fmt.Errorf("go format error: %v", err)
      }
      fmts := string(formatted)
      resp.File[i].Content = &fmts
   }
   return nil
}

// No.2
func Generate(req *plugin.CodeGeneratorRequest) *plugin.CodeGeneratorResponse {
   // Begin by allocating a generator. The request and response structures are stored there
   // so we can do error handling easily - the response structure contains the field to
   // report failure.
   g := generator.New()
   g.Request = req

   // 5.将从标准输出中所读取到的 CodeGeneratorRequest 传递给 protoc 的代码生成器
   g.CommandLineParameters(g.Request.GetParameter())

   // Create a wrapped version of the Descriptors and EnumDescriptors that
   // point to the file that defines them.
   // 6.封装到 Generator 中的文件引用对象中
   g.WrapTypes()

   // 7.设置包名
   g.SetPackageNames()
    
   // 8.将类型名称映射为对象
   g.BuildTypeNameMap()

   // 9.生成所有的源文件
   g.GenerateAllFiles()

   if err := goformat(g.Response); err != nil {
      g.Error(err)
   }

   testReq := proto.Clone(req).(*plugin.CodeGeneratorRequest)

   testResp := GeneratePlugin(testReq, testgen.NewPlugin(), "pb_test.go")

   for i := 0; i < len(testResp.File); i++ {
      if strings.Contains(*testResp.File[i].Content, `//These tests are generated by github.com/gogo/protobuf/plugin/testgen`) {
         g.Response.File = append(g.Response.File, testResp.File[i])
      }
   }

   return g.Response
}


// No.3
func Write(resp *plugin.CodeGeneratorResponse) {
   g := generator.New()
   // Send back the results.
   // 10.将序列化后的 CodeGeneratorRequest 输出到标准输出
   data, err := proto.Marshal(resp)
   if err != nil {
      g.Error(err, "failed to marshal output proto")
   }
   _, err = os.Stdout.Write(data)
   if err != nil {
      g.Error(err, "failed to write output proto")
   }
}
```

从上述代码中可以看出，此处没有自定义插件的相关逻辑，在`grpc.go`文件中，即如下初始化代码。

```go
_ "github.com/gogo/protobuf/protoc-gen-gogo/grpc"
```

#### c. grpc.go

```go
...

func init() {
   generator.RegisterPlugin(new(grpc))
}

// grpc is an implementation of the Go protocol buffer compiler's
// plugin architecture.  It generates bindings for gRPC support.
type grpc struct {
   gen *generator.Generator
}

// Name returns the name of this plugin, "grpc".
func (g *grpc) Name() string {
   return "grpc"
}

...

// Init initializes the plugin.
func (g *grpc) Init(gen *generator.Generator) {
   g.gen = gen
}

...

// Generate generates code for the services in the given file.
func (g *grpc) Generate(file *generator.FileDescriptor) {
   if len(file.FileDescriptorProto.Service) == 0 {
      return
   }

   contextPkg = string(g.gen.AddImport(contextPkgPath))
   grpcPkg = string(g.gen.AddImport(grpcPkgPath))

   g.P("// Reference imports to suppress errors if they are not otherwise used.")
   g.P("var _ ", contextPkg, ".Context")
   g.P("var _ ", grpcPkg, ".ClientConn")
   g.P()

   // Assert version compatibility.
   g.P("// This is a compile-time assertion to ensure that this generated file")
   g.P("// is compatible with the grpc package it is being compiled against.")
   g.P("const _ = ", grpcPkg, ".SupportPackageIsVersion", generatedCodeVersion)
   g.P()

   for i, service := range file.FileDescriptorProto.Service {
      g.generateService(file, service, i)
   }
}

// GenerateImports generates the import declaration for this file.
func (g *grpc) GenerateImports(file *generator.FileDescriptor) {}

...
```

实际上，它将自定义的 grpc 插件注册到了 generator 中，这意味着只需将公共的 main 方法中的内容拷贝一份然后就可以根据实际需要实现自定义的 `protoc`插件了。

### 2. generator.Plugin 接口

在查看 grpc 插件时，会发现它实现了很多方法，只需重点关注 `generator.Plugin`的相关接口即可。

```go
type Plugin interface {
	// Name identifies the plugin.
	Name() string
	// Init is called once after data structures are built but before
	// code generation begins.
	Init(g *Generator)
	// Generate produces the code generated by the plugin for this file,
	// except for the imports, by calling the generator's methods P, In, and Out.
	Generate(file *FileDescriptor)
	// GenerateImports produces the import declarations for this file.
	// It is called after Generate.
	GenerateImports(file *FileDescriptor)
}
```

在实现自定义插件时，只需满足该`Plugin`接口的定义，就可以无缝的接入 protoc，此接口设计 4 个方法。

+ `Name`：插件的名称
+ `Init`：插件初始化动作
+ `Generate`：生成文件所需的具体代码
+ `GenerateImports`：生成文件所需的具体导入声明。

简单来说，只要包含了公共的 main 方法并实现了`generator.Plugin`接口的4个方法，就可以在`protoc`中使用自定义插件。由此可见，grpc 插件是根据`Plugin`接口的主题流程来流转的。

### 3. FileDescriptor 属性

既然是针对 `proto`文件来生成对应的 Go 代码，那么在生成文件时必然需要很多可用的细腻些，即在`proto`文件中定义的信息。查看所使用的`tag.proto`文件。

```go
syntax = "proto3"; // 1

package proto; // 2

import "proto/common.proto"; // 3
import "google/api/annotations.proto";

service TagService { // 4
    rpc GetTagList (GetTagListRequest) returns (GetTagListResponse) {
        option (google.api.http) = {
            get: "/api/v1/tags" // 5
        };
    };
}

message GetTagListRequest { // 6
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

可以看到在一个普通的`proto`文件中，至少包含了6大属性。

+ proto 文件的语法版本，一般为`proto3`或`proto2`
+ 包（package）名称
+ 引入的多个依赖包
+ 多个 RPC 服务（service）的定义
+ 多个 RPC 方法（rpc method）和对应入参、出参的消息体的定义
+ 其他多个消息体的定义

在根据 proto 文件生成对应代码时是需要用到这些信息的，否则无法进行具体的生成。Plugin 接口方法中提供的`FileDescriptor`属性就包含了这些信息。

```go
// FileDescriptor describes an protocol buffer descriptor file (.proto).
// It includes slices of all the messages and enums defined within it.
// Those slices are constructed by WrapTypes.
type FileDescriptor struct {
   *descriptor.FileDescriptorProto
   desc []*Descriptor          // All the messages defined in this file.
   enum []*EnumDescriptor      // All the enums defined in this file.
   ext  []*ExtensionDescriptor // All the top-level extensions defined in this file.
   imp  []*ImportedDescriptor  // All types defined in files publicly imported by this file.

   // Comments, stored as a map of path (comma-separated integers) to the comment.
   comments map[string]*descriptor.SourceCodeInfo_Location

   // The full list of symbols that are exported,
   // as a map from the exported object to its symbols.
   // This is used for supporting public imports.
   exported map[Object][]symbol

   importPath  GoImportPath  // Import path of this file's package.
   packageName GoPackageName // Name of this file's Go package.

   proto3 bool // whether to generate proto3 code for this file
}
```

`FileDescriptor`属性所包含的信息非常多，主要分为文件的描述信息、消息体的定义、枚举的定义、顶级扩展的定义、公开导入的文件中的所有类型的定义、注释、导出的符号的完整列表（作为从导出的对象到其符号的映射）、该文件包的导入路径、该文件的 Go 软件包的名称，以及是否以此文件生成 `proto3`代码。

其包含的属性层层嵌套，展示一些与本次生成 proto 文件相关的属性，以便后续的开发和使用。

#### a. proto 文件的描述（FileDescriptor）

| 字段名           | 类型                      | 含义                                        |
| ---------------- | ------------------------- | ------------------------------------------- |
| Name             | *string                   | proto 文件名                                |
| Package          | *string                   | proto 文件的包名称                          |
| Dependency       | []string                  | proto 文件中引用（import）的依赖包列表      |
| PublicDependency | []int32                   | Dependency 引用的依赖包列表中的索引         |
| WeakDependency   | []int32                   | 仅适用于 Google 内部迁移，不需要了解        |
| MessageType      | []*DescriptorProto        | 所有消息体（message）类型的定义             |
| EnumType         | []EnumDescriptorProto     | 所有枚举（enum）类型的定义                  |
| Service          | []*ServiceDescriptorProto | 所有服务（service）类型的定义               |
| Extension        | []*FiledDescriptorProto   | 所有扩展信息的定义                          |
| Options          | *FileOptions              | 文件选项                                    |
| SourceCodeInfo   | *SourceCodeInfo           | 源代码的相关信息                            |
| Syntax           | *string                   | proto 文件的语法版本，值为 proto2 或 proto3 |

#### b. Service 定义的描述（ServiceDescriptorProto）

| 字段名  | 类型                     | 含义                   |
| ------- | ------------------------ | ---------------------- |
| Name    | *string                  | 服务（service）名      |
| Method  | []*MethodDescriptorProto | 方法（rpc method）列表 |
| Options | *ServiceOptions          | 服务选项               |

#### c. Service 中 Method 定义的描述（MethodDescriptorProto）

| 字段名          | 类型           | 含义                                 |
| --------------- | -------------- | ------------------------------------ |
| Name            | *string        | 方法名称                             |
| InputType       | *string        | 方法的入参类型                       |
| OutputType      | *string        | 方法的出参类型                       |
| Options         | *MethodOptions | 方法选项                             |
| ClientStreaming | *bool          | 标识客户端是否流式传输多个客户端信息 |
| ServerStreaming | *bool          | 标识服务端是否流式传输多个服务信息   |

### 4. 实现一个简单的自定义插件

创建一个新的项目目录 `protoc-gen-go-tour`，对项目进行初始化，并新建对应的文件，代码如下。

```bash
$ go mod init protoc-gen-go-tour
```

最终目录结构如下。

![image-20220528165445360](https://raw.githubusercontent.com/tonshz/test/master/img/202205281654904.png)

下载所需的依赖模块，在项目根目录下执行：

```bash
$ go get -u github.com/golang/protobuf@v1.3.3
```

#### a. 公共 generator

打开`main.go`文件，写入固定的公共`generator`逻辑。

```go
package main

import (
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/generator"
)

func main() {
	// Begin by allocating a generator. The request and response structures are stored there
	// so we can do error handling easily - the response structure contains the field to
	// report failure.
	g := generator.New()

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		g.Error(err, "reading input")
	}

	if err := proto.Unmarshal(data, g.Request); err != nil {
		g.Error(err, "parsing input proto")
	}

	if len(g.Request.FileToGenerate) == 0 {
		g.Fail("no files to generate")
	}

	g.CommandLineParameters(g.Request.GetParameter())

	// Create a wrapped version of the Descriptors and EnumDescriptors that
	// point to the file that defines them.
	g.WrapTypes()

	g.SetPackageNames()
	g.BuildTypeNameMap()

	g.GenerateAllFiles()

	// Send back the results.
	data, err = proto.Marshal(g.Response)
	if err != nil {
		g.Error(err, "failed to marshal output proto")
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		g.Error(err, "failed to write output proto")
	}
}
```

#### b. 实现 tour 插件

在编写完 generator 的启动代码后，新建`tour.go`文件，写入自定义插件的代码。

```GO
package tour

import "github.com/golang/protobuf/protoc-gen-go/generator"

func init() {
	generator.RegisterPlugin(new(tour))
}

type tour struct {
	gen *generator.Generator
}

func (g *tour) Name() string {
	return "tour"
}

func (g *tour) Init(gen *generator.Generator) {
	g.gen = gen
}

func (g *tour) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	} 
}

func (g *tour) GenerateImports(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}
}
```

再在与`main.go`同级的`link_tour.go`文件中写入初始化代码。

```go
import _ "protoc-gen-go-tour/tour"
```

#### c. 验证

在编写玩插件后，是不能直接使用这个插件的，而是需要对其进行编译，然后将编译好的二进制文件移到对应的 bin 目录下，执行如下命令。

```bash
$ go build .
$ mv ./protoc-gen-go-tour.exe $HOME/go/bin # 此处命令在 WIn10 下执行
```

将`tag.proto`文件拷贝一份命名为`tour.proto`，在`ch03`的根目录下执行如下命令。

```bash
$ protoc -I C:\Users\zyc\protoc\include -I C:\Users\zyc\go\pkg\mod\github.com\grpc-ecosystem\grpc-gateway@v1.14.5\third_party\googleapis -I . --go-tour_out=plugins=tour:. ./proto/tour.proto
```

在上述命令中，新增的命令主要是`--go-tour_out=plugins=tour`，`--go-tour_out`会告诉 protoc 编译器去查找并使用名为`protoc-gen-tour`的插件，而`plugins=tour`则指定使用`protoc-gen-tour`插件中的 tour 子插件（允许在插件中自定义多个子插件）。

#### d. 源码分析

生成完毕后，即可在`proto`目录下看到新生成的`tour.pb.go`文件。

```go
// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/tour.proto

package proto

import (
   fmt "fmt"
   proto "github.com/golang/protobuf/proto"
   _ "google.golang.org/genproto/googleapis/api/annotations"
   math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type GetTagListRequest struct {
   ...
}

func (m *GetTagListRequest) Reset()         { *m = GetTagListRequest{} }
func (m *GetTagListRequest) String() string { return proto.CompactTextString(m) }
func (*GetTagListRequest) ProtoMessage()    {}
func (*GetTagListRequest) Descriptor() ([]byte, []int) {
   return fileDescriptor_1b8a0da1d14bccfa, []int{0}
}

func (m *GetTagListRequest) XXX_Unmarshal(b []byte) error {
   return xxx_messageInfo_GetTagListRequest.Unmarshal(m, b)
}
func (m *GetTagListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
   return xxx_messageInfo_GetTagListRequest.Marshal(b, m, deterministic)
}
func (m *GetTagListRequest) XXX_Merge(src proto.Message) {
   xxx_messageInfo_GetTagListRequest.Merge(m, src)
}
func (m *GetTagListRequest) XXX_Size() int {
   return xxx_messageInfo_GetTagListRequest.Size(m)
}
func (m *GetTagListRequest) XXX_DiscardUnknown() {
   xxx_messageInfo_GetTagListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetTagListRequest proto.InternalMessageInfo

func (m *GetTagListRequest) GetName() string {
   ...
}

func (m *GetTagListRequest) GetState() uint32 {
   ...
}

type Tag struct {
   ...
}

func (m *Tag) Reset()         { *m = Tag{} }
func (m *Tag) String() string { return proto.CompactTextString(m) }
func (*Tag) ProtoMessage()    {}
func (*Tag) Descriptor() ([]byte, []int) {
   return fileDescriptor_1b8a0da1d14bccfa, []int{1}
}

func (m *Tag) XXX_Unmarshal(b []byte) error {
   return xxx_messageInfo_Tag.Unmarshal(m, b)
}
func (m *Tag) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
   return xxx_messageInfo_Tag.Marshal(b, m, deterministic)
}
func (m *Tag) XXX_Merge(src proto.Message) {
   xxx_messageInfo_Tag.Merge(m, src)
}
func (m *Tag) XXX_Size() int {
   return xxx_messageInfo_Tag.Size(m)
}
func (m *Tag) XXX_DiscardUnknown() {
   xxx_messageInfo_Tag.DiscardUnknown(m)
}

var xxx_messageInfo_Tag proto.InternalMessageInfo

func (m *Tag) GetId() int64 {
   ...
}

func (m *Tag) GetName() string {
   ...
}

func (m *Tag) GetState() uint32 {
   ...
}

type GetTagListResponse struct {
   ...
}

func (m *GetTagListResponse) Reset()         { *m = GetTagListResponse{} }
func (m *GetTagListResponse) String() string { return proto.CompactTextString(m) }
func (*GetTagListResponse) ProtoMessage()    {}
func (*GetTagListResponse) Descriptor() ([]byte, []int) {
   return fileDescriptor_1b8a0da1d14bccfa, []int{2}
}

func (m *GetTagListResponse) XXX_Unmarshal(b []byte) error {
   return xxx_messageInfo_GetTagListResponse.Unmarshal(m, b)
}
func (m *GetTagListResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
   return xxx_messageInfo_GetTagListResponse.Marshal(b, m, deterministic)
}
func (m *GetTagListResponse) XXX_Merge(src proto.Message) {
   xxx_messageInfo_GetTagListResponse.Merge(m, src)
}
func (m *GetTagListResponse) XXX_Size() int {
   return xxx_messageInfo_GetTagListResponse.Size(m)
}
func (m *GetTagListResponse) XXX_DiscardUnknown() {
   xxx_messageInfo_GetTagListResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetTagListResponse proto.InternalMessageInfo

func (m *GetTagListResponse) GetList() []*Tag {
   ...
}

func (m *GetTagListResponse) GetPager() *Pager {
   ...
}

func init() {
   ...
}

func init() { proto.RegisterFile("proto/tour.proto", fileDescriptor_1b8a0da1d14bccfa) }

var fileDescriptor_1b8a0da1d14bccfa = []byte{
   ...
}
```

从上述代码中可以发现，即使没有编写任何新代码，没有输出任何新内容，只是实现了`generator.Plugin`接口中定义的方法，tour 插件也会按规范生成代码。

对`tag.pb.go`文件和`tour.pb.go`文件进行对比，可以发现后者生成的代码没有像前者一样生成 gRPC 相关的方法。`tour.pb.go`文件生成的代码都是针对`Protobuf`，即只生成了 proto 文件（`tour.proto`）和消息体（message）定义、配套方法。

`tour.pb.go`文件中的这些基础 IDL 方法均来源于`protoc-gen-go`插件。`generator`生成器直接使用的是`protoc-gen-go`插件，而`protoc-gen-go`插件的真正本体时 `Protobuf`这个 IDL 生成器。因此指定了`plugins=grpc`，所有会有 gRPC 相关的代码，即在将子插件更换为 tour 后，便只剩默认的与 `Protobuf`相关的 Go 代码的生成功能了。

### 5. 实现定制化的 gRPC 自定义插件

在本章的项目中，至少有两种情况是需要自定义插件的。第一种是基于 gRPC 插件实现一些定制化功能，第二种是基于 proto 文件生成一些扩展性的代码和功能（如在社区中很常见的`protoc-gen-grpc-gateway`和`protoc-gen-swagger`）。

本节基于官方 protoc 的 gRPC 插件实现一个最简单的定制化需求。

#### a. 确认需求

在业务开发过程中，有时候会有多租户的概念，即根据租户标识的不同，获取其对应的租户实例信息和数据。比如，根据租户标识，判定当前部署的环境（如灰度环境或私有环境），以便进行更多的精确调度。因此必须在调用中传播租户标识。

#### b. 解决方法

解决这个问题方法之一就是开发自定义插件，在调用 gRPC Client 时，就要求必须传入租户标识，并且可以在内部进行入参校验和节点发现。此方案能够较好的解决这个问题，并且对开发人员的直接侵入性较小。

#### c. 实现 tour 插件

首先拷贝一份官方的 gRPC 插件中的代码，原因有两个：第一，官方的 gRPC 插件中的代码是不可导出（无法外部引用）的；第二，需要对源代码进行较大的修改。

#####  拷贝插件模板

将`protoc-gen-go/grpc/grpc.go`中的代码完整的复制到项目`tour`目录下的`tour.go`文件中，并将其对应所有的标识修改为 tour。

```go
package tour
...
type tour struct {
   gen *generator.Generator
}

func (g *tour) Name() string {
	return "tour"
}
...
```

主要将包名称修改为 tour，原本的 grpc 结构体及其方法也都调整为了 tour，并且修改了插件的名称，即 Name 方法返回的结果。

##### 分析插件模板

在修改官方的 gRPC 插件之前，必须了解其原有逻辑，下面先来看看最核心的`Generate`方法。

```go
// Generate generates code for the services in the given file.
func (g *tour) Generate(file *generator.FileDescriptor) {
   if len(file.FileDescriptorProto.Service) == 0 {
      return
   }

   contextPkg = string(g.gen.AddImport(contextPkgPath))
   grpcPkg = string(g.gen.AddImport(grpcPkgPath))

   g.P("// Reference imports to suppress errors if they are not otherwise used.")
   g.P("var _ ", contextPkg, ".Context")
   g.P("var _ ", grpcPkg, ".ClientConnInterface")
   g.P()

   // Assert version compatibility.
   g.P("// This is a compile-time assertion to ensure that this generated file")
   g.P("// is compatible with the grpc package it is being compiled against.")
   g.P("const _ = ", grpcPkg, ".SupportPackageIsVersion", generatedCodeVersion)
   g.P()

   for i, service := range file.FileDescriptorProto.Service {
      g.generateService(file, service, i)
   }
}
```

从上述代码中可以看出，使用次数最多的时 P 方法，P 方法时插件中最常用的方法之一，其作用是将传入的参数打印到所需生成的文件输出上，这是一个相对原子的方法，并不会做过多的事情。

从逻辑上看，该方法主要是对 proto 文件的引入和版本信息进行了输出和定义，然后将最重要的服务（service）逻辑，通过`FileDescriptorProto.Service`（标识所有服务类型的定义）循环调用`generateService`方法进行转换和输出。

`generateService`为主要的生成处理方法，其处理 gRPC Client 生成的方法名为 `generateClientMethod`。

```go
func (g *tour) generateClientMethod(servName, fullServName, serviceDescVar string, method *pb.MethodDescriptorProto, descExpr string) {
   sname := fmt.Sprintf("/%s/%s", fullServName, method.GetName())
   methName := generator.CamelCase(method.GetName())
   inType := g.typeName(method.GetInputType())
   outType := g.typeName(method.GetOutputType())

   if method.GetOptions().GetDeprecated() {
      g.P(deprecationComment)
   }
   g.P("func (c *", unexport(servName), "Client) ", g.generateClientSignature(servName, method), "{")
   if !method.GetServerStreaming() && !method.GetClientStreaming() {
      g.P("out := new(", outType, ")")
      // TODO: Pass descExpr to Invoke.
      g.P(`err := c.cc.Invoke(ctx, "`, sname, `", in, out, opts...)`)
      g.P("if err != nil { return nil, err }")
      g.P("return out, nil")
      g.P("}")
      g.P()
      return
   }
   // Stream auxiliary types and methods： 流辅助类型和方法
   ...
}
```

在分析代码前，可以结合生成后的代码进行观察。

```go
func (c *tagServiceClient) GetTagList(ctx context.Context, in *GetTagListRequest, opts ...grpc.CallOption) (*GetTagListResponse, error) {
   out := new(GetTagListResponse)
   err := c.cc.Invoke(ctx, "/proto.TagService/GetTagList", in, out, opts...)
   if err != nil {
      return nil, err
   }
   return out, nil
}
```

实际上，通过对生成代码和生成后的代码进行对比，可以快速的知道各行代码都做了那些事情。

##### 二次开发

需要做的是对 gRPC Client 方法进行租户标识（orgcode）的获取和判断。若不存在，则直接返回响应的错误信息。因此需要编写获取和设置租户标识值的方法，编写判定租户标识值是否正确的方法，在`tour.go`文件中新增代码。

```go
// 需要在 pb.go 文件中新增获取和设置租户标识值的方法。下为模板代码
/*
   type orgCodeKey struct{}

   func (c *tagServiceClient) WithOrgCode(ctx context.Context, orgCode string) context.Context {
      return context.WithValue(ctx, orgCodeKey{}, orgCode)
   }

   func (c *tagServiceClient) OrgCode(ctx context.Context) string {
      return ctx.Value(orgCodeKey{}).(string)
   }
*/
func (g *tour) generateOrgCodeMethod() {
   g.P("type orgCodeKey struct{}")
   g.P()
   g.P("func (c *tagServiceClient) WithOrgCode(ctx context.Context, orgCode string) context.Context {")
   g.P("return context.WithValue(ctx, orgCodeKey{}, orgCode)")
   g.P("}")
   g.P()
   g.P("func (c *tagServiceClient) OrgCode(ctx context.Context) string {")
   g.P("return ctx.Value(orgCodeKey{}).(string)")
   g.P("}")
}
```

在上述代码中，声明了获取租户标识值得获取和设置方法，在生成具体的 gRPC Client 方法前，需要将`generateOrgCodeMethod`方法的调用添加到主流程中，并且将`WithOrgCode`的定义加入`Client interface`中，以便调用。

```go
// generateService generates all the code for the named service.
func (g *tour) generateService(file *generator.FileDescriptor, service *pb.ServiceDescriptorProto, index int) {
   ...
   for i, method := range service.Method {
      g.gen.PrintComments(fmt.Sprintf("%s,2,%d", path, i)) // 2 means method in a service.
      if method.GetOptions().GetDeprecated() {
         g.P("//")
         g.P(deprecationComment)
      }
      g.P(g.generateClientSignature(servName, method))
   }
   // 第一处：新增的 WithOrgCode 定义
   g.P("WithOrgCode(ctx context.Context, orgCode string) context.Context")
   g.P("}")
   g.P()

   // Client structure.
   g.P("type ", unexport(servName), "Client struct {")
   g.P("cc ", grpcPkg, ".ClientConnInterface")
   g.P("}")
   g.P()

   // NewClient factory.
   if deprecated {
      g.P(deprecationComment)
   }
   g.P("func New", servName, "Client (cc ", grpcPkg, ".ClientConnInterface) ", servName, "Client {")
   g.P("return &", unexport(servName), "Client{cc}")
   g.P("}")
   g.P()

   // 第二处：新增 OrgCode 定义
   g.generateOrgCodeMethod()
   ...
}
```

在完成了所需方法的声明和调用后，对`generateClientMethod`方法进行修改。

```go
func (g *tour) generateClientMethod(servName, fullServName, serviceDescVar string, method *pb.MethodDescriptorProto, descExpr string) {
   ...
   if !method.GetServerStreaming() && !method.GetClientStreaming() {
      // 添加租户逻辑
      g.P("orgCode, ok := c.OrgCode(ctx)")
      g.P(`if !ok || orgCode == "" {`)
      g.P(`return nil, errors.New("请调用 WithOrgCode 方法设置租户标识")`)
      g.P("}")
      // =======
      g.P("out := new(", outType, ")")
      // TODO: Pass descExpr to Invoke.
      g.P(`err := c.cc.Invoke(ctx, "`, sname, `", in, out, opts...)`)
      g.P("if err != nil { return nil, err }")
      g.P("return out, nil")
      g.P("}")
      g.P()
      return
   }
}
```

在生成 gRPC Client 方法时，调用了声明的对应方法，并新增了对租户标识的获取和判断，若出现不存在或值为空的情况（也可针对实际业务场景自行定制），则直接返回错误。

因为在返回错误时需要调用 errors 的标准库，而在原先的生成文件中并没有引用该标准库，因此需要在 `tour.go`中声明`errorsPkgPath`，并且在`Generate`方法中进行调用。

```go
...
const (
   contextPkgPath = "context"
   errorPkgPath   = "errors"
   grpcPkgPath    = "google.golang.org/grpc"
   codePkgPath    = "google.golang.org/grpc/codes"
   statusPkgPath  = "google.golang.org/grpc/status"
)

...
// Generate generates code for the services in the given file.
func (g *tour) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}

	// 新增 errors 标准库的引入
	_ = g.gen.AddImport(errorPkgPath)
	contextPkg = string(g.gen.AddImport(contextPkgPath))
	grpcPkg = string(g.gen.AddImport(grpcPkgPath))

	...
}
```

##### 验证生成

回到`protoc-gen-go-tour`项目的根目录，重新编译和移动，执行如下命令：

````bash
$ go build . 
#需要先删除目录下的文件，不然会报错 mv : 当文件已存在时，无法创建该文件。
$ mv ./protoc-gen-go-tour.exe $HOME/go/bin
````

在`ch03`目录下重新对`tag.proto`文件进行生成。

```bash
$ protoc -I C:\Users\zyc\protoc\include -I C:\Users\zyc\go\pkg\mod\github.com\grpc-ecosystem\grpc-gateway@v1.14.5\third_party\googleapis -I . --go-tour_out=plugins=tour:. ./proto/tour.proto
```

再次查看`tour.pb.go`文件。

```go
// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/tour.proto

package proto

import (
   context "context"
   // 1.导入 errors 标准库
   errors "errors"
   ...
)

...
// TagServiceClient is the client API for TagService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TagServiceClient interface {
   GetTagList(ctx context.Context, in *GetTagListRequest, opts ...grpc.CallOption) (*GetTagListResponse, error)
    // 2.新增 WithOrgCode() 定义
   WithOrgCode(ctx context.Context, orgCode string) context.Context
}

type tagServiceClient struct {
   cc grpc.ClientConnInterface
}

func NewTagServiceClient(cc grpc.ClientConnInterface) TagServiceClient {
   return &tagServiceClient{cc}
}

// 3.新增 OrgCode 定义
type orgCodeKey struct{}

// 4.设置租户标识
func (c *tagServiceClient) WithOrgCode(ctx context.Context, orgCode string) context.Context {
   return context.WithValue(ctx, orgCodeKey{}, orgCode)
}

// 5.获取租户标识
func (c *tagServiceClient) OrgCode(ctx context.Context) string {
   return ctx.Value(orgCodeKey{}).(string)
}
func (c *tagServiceClient) GetTagList(ctx context.Context, in *GetTagListRequest, opts ...grpc.CallOption) (*GetTagListResponse, error) {
   // 6.新增对租户标识的获取和判断
   orgCode, ok := c.OrgCode(ctx)
   if !ok || orgCode == "" {
      return nil, errors.New("请调用 WithOrgCode 方法设置租户标识")
   }
   
   out := new(GetTagListResponse)
   err := c.cc.Invoke(ctx, "/proto.TagService/GetTagList", in, out, opts...)
   if err != nil {
      return nil, err
   }
   return out, nil
}

...
```

从代码中可以看出，在生成的代码中包含了所需的获取和设置租户标识的方法，并且在 gRPC Client 的方法定义和接口声明中也包含了预定义的方法。

##### 验证功能

重新运行 `tag-service`服务，调用客户端，在不进行任何设置的情况下，可以看到如下报错：

```bash
$ go run client.go
tagServiceClient.GetTagList err: 请调用 WithOrgCode 方法设置租户标识
```

此时在`client.go`文件中调用`WithOrgCode`方法进行设置即可。

```go
// 初始化指定 RPC Proto Service 的客户端实例对象
tagServiceClient := pb.NewTagServiceClient(clientConn)
newCtx := tagServiceClient.WithOrgCode(ctx, "Go 语言学习")
// 发起指定 RPC 方法的调用
resp, _ := tagServiceClient.GetTagList(newCtx, &pb.GetTagListRequest{Name: "Go"})
```

如果是通过`grpc-gateway`所映射的 HTTP 接口进行访问也会受到影响。

```json
{
    "error": "请调用 WithOrgCode 方法设置租户标识",
    "code": 2,
    "message": "请调用 WithOrgCode 方法设置租户标识"
}
```

如果希望`grpc-gateway`也能访问，那么需要调整上下文的键，将`grpc-gateway`提供的 Header（`grpc-metadata-*`）转换成 metedata，再进行设置并访问即可。

## 十三、对 gRPC 接口进行版本管理

一般来说，项目会随着时间和需求不断的迭代发展和变更，原本的 gRPC 接口也会出现或大或小的改动，改动需要兼容先前的客户端，需要更好的管理 gRPC 接口版本。

### 1. 接口变更

一般在迭代过程中，最常见的接口 IDL（Protobuf）变更有以下几种：

+ 新增 service 和 rpc 方法
+ 新增 message 中的字段
+ 删除 message 中的字段
+ 修改 message 中的字段类型

这几种变更操作中，有些是可兼容性修改，有些是破坏性修改。

### 2. 可兼容性修改

+ 新增新的 service
+ 在原有的 service 中新增新的 RPC 方法
+ 在原有的请求 message 中新增新的字段
+ 在原有的响应 message 中新增新的字段

### 3. 破坏性修改

+ 修改原有的字段数据类型：修改原有的字段类型会使协议出现重大变更，即在客户端和新服务端交互过程中会出现反序列化报错，影响使用。
+ 修改原有的字段标识位：字段标识位是用来标识其在 Protobuf Payload 中的位置，一旦发生改变，就会出现找不到或者错位的情况。
+ 修改原有的字段名：不会产生直接影响，字段名仅仅在生成的代码中使用，因此不会给 Protobuf 生成的二进制数据带来什么影响。但是，不管在服务端还是在客户端，应用代码都会用到这个字段，但字段名发生改变时，一旦重新编译 Protobuf，就会导致应用无法正常提供服务。
+ 修改 message 原有的命名：不会产生直接影响，因为 message 的命名变更并不会导致协议出现重大变更，但是如果客户端升级了所使用的 Protobuf，那么 message 的命名变更将导致应用代码报错。
+ 删除 service 或 RPC 方法：会造成直接影响。
+ 删除原有的字段：会造成直接影响，

### 4. 设计 gRPC 接口

在变更 gRPC 接口时，要尽可能不影响现有的客户端。因为一旦影响了现有的客户端，则需要通知所有客户端提前进行测试，在同步服务端的更新和发布，这是一件很麻烦的事情。

遇到非常大的接口出参和入参变动，则往往会选择两种方式：一是编写新的 RPC 方法，而是在原有接口内对入参或出参进行转换，然后再内部将流量导向新的 RPC 方法。

### 5. 版本号管理

当进行大版本变更时，应尽可能让两者完全隔离，即不出现版本内互通的情况。

```protobuf
syntax = "proto3";

package proto.v1;

import "proto/common.proto";
import "google/api/annotations.proto";

service TagService {
	rpc GetTagList (GetTagListRequest) returns (GetTagListReply) {
		option (google.api.http) = {
			get: "/api/v1/tags"
		};
	}
}
```

可以在 package 上指定版本号，格式为`proto.v1.TagService`或`proto.v2.TagService`。假设存在 HTTP 路由，则可以将大版本号定义在路由中，格式为`/api/v1/tag`或`/api/v1/tag`。

也可以通过 Header 继续传播，但在 gRPC 下 ，这种模式并不够显性化，并且可能需要进行额外的特殊处理，需要根据实际应用场景进行探讨。

## 十四、常见问题讨论

### 1. Q&A

#### a. 当调用 grpc.Dial 时会连接服务端吗？

会，并且是异步连接的，连接状态为 Connecting（正在连接）状态，如果设置了`grpc.WithBlock`选项，就会阻塞等待（等待到达 Ready）。==**需要注意的是，当未设置`grpc.WithBlock`时，上下文超时控制将不会生效。**==

#### b. 再调用 ClientConn 时不执行 Close 语句会导致泄露吗？

有可能会。除非客户端不是常驻进程，即在应用结束时会被动的回收资源。如果是常驻进程，同时没有执行 Close 语句，则有可能会泄露。

#### c. 如果不做超时控制，会出现什么问题？

短时间内不会出现问题，但是会不断泄露，直至服务无法正常运行。当自身服务逐渐出现问题时，其会影响上下游的服务调用，因此默认对上下文进行超时控制非常重要。

#### d. 频繁创建 ClientConn 有问题吗？

![image-20220528203426200](https://raw.githubusercontent.com/tonshz/test/master/img/202205282051067.png)

但应用场景存在高频次同时生成/调用 ClientConn 时，会导致系统的文件句柄占用过多，需要调整代码。 

#### e. 客户端请求失败后会默认重试吗？

会不断的进行重试，直至上下文取消。一般采用 `backoff` 算法，默认的最大重试时间间隔时 120s。

#### f. 在 Kubernetes 中，gRPC 负载均衡有问题吗？

gRPC 的 RPC 协议是基于 HTTP/2 标准实现的，HTTP/2 的一大特性是，它不像 HTTP/1.1 在每次发出请求时都会重新建立一个新连接，而是会复用原有的连接。

当使用 k8s Service 做负载均衡时，会导致 `kube-proxy` 只有在连接建立时才会做负载均衡，而在这之后的同一个客户端的每一次 RPC 请求都会复用原有的连接，也就是说，实际上后续的每一次 RPC 强求都会被送到同一个服务端去处理，最终导致负载不均衡。

#### g. 用什么微服务框架？

在完整的微服务框架选型上，需要根据实际情况自行处理。

### 2. 小结

+ 客户端请求若使用 `grpc.Dial`方法，则默认是一异步建立连接，连接状态为 Connecting。
+ 客户端请求若需要同步建立连接，则调用`WithBlock`方法，连接状态 Ready。
+ 在特定场景下，如果不对`grpc.ClientConn`进行调控，则会导致文件句柄超出系统限制，影响正常使用。
+ 若内部调用完毕后，`grpc.ClientConn`不关闭连接，则有可能会导致`goroutine`和 `Memory`等出现问题。
+ 任何内/外调用如果不进行超时控制，都会导致泄露和客户端不断重试，最终导致上下游服务出现连环问题，非常危险。
+ 当选择 gRPC 负载均衡模式时，需要谨慎，确定其在那一层进行负载，若直接使用 Kubernetes Service 且不做调整，则很容易出现问题，导致负载不均衡。
+ 完整的微服务框架需要根据企业实际情况进行选型，需根据具体情况具体讨论。

-----------------------------

## 参考

+ [参考文章](https://golang2.eddycjy.com/)

+ [GitHub代码](https://github.com/go-programming-tour-book)









