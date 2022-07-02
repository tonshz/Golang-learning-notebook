# Go 语言编程之旅(三)：RPC 应用(一) 

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

----------------------------------

## 参考

+ [参考文章](https://golang2.eddycjy.com/)

+ [GitHub代码](https://github.com/go-programming-tour-book)



