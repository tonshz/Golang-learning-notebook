# Go 语言编程之旅(三)：RPC 应用(五) 

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

------------------------

## 参考

+ [参考文章](https://golang2.eddycjy.com/)

+ [GitHub代码](https://github.com/go-programming-tour-book)

