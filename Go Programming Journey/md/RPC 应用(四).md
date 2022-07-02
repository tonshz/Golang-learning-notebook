# Go 语言编程之旅(三)：RPC 应用(四) 

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

---------------

## 参考

+ [参考文章](https://golang2.eddycjy.com/)

+ [GitHub代码](https://github.com/go-programming-tour-book)



