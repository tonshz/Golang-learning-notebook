# Go 语言编程之旅(三)：RPC 应用(二) 

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

---------------

## 参考

+ [参考文章](https://golang2.eddycjy.com/)

+ [GitHub代码](https://github.com/go-programming-tour-book)



