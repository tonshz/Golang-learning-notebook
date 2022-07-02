# Go 语言编程之旅(二)：HTTP 应用(八) 

## 十、进行链路追踪

在完成上一届的应用中间件后，解决了一系列的问题，但在对接新接口时，可能会出现个别接口响应速度十分缓慢的情况。

项目在不断迭代之后，可能会涉及许许多多的接口，而这些接口很可能是分布式部署的，即存在这多份副本，又存在着相互调用，并且在各自的调用中还可能包含大量的 SQL、HTTP、Redis 以及应用的调用逻辑。如果对每一个都进行调用堆栈的日志输出，就会导致日志文件过大。

故为了更好的解决这个问题，需要采用分布式链路追踪系统，它能有效解决可观察性上的一部分问题，即多程序部署可在多环境下调用链路的“观察”。

企业级项目总是在不断地迭代中，让程序尽可能的支持横向扩展，且具备一定的可观察性，对排查问题很有帮助。

### 1. OpenTracing 规范

OpenTracing 规范的出现是为了解决不同供应商的分布式追踪系统 API 互不兼容的问题，它提供了一个标准的、与供应商无关的工具框架，可以认为它是一个接入层。下面从多个维度进行分析：

+ 从功能上：在 OpenTracing 规范中会提供一系列与供应商无关的 API。
+ 从系统上：它能够让开发人员更便捷的对接（新增或替换）追踪系统，只需简单的更改 Tracer 的配置就可以了。
+ 从语言上：OpenTracing 规范是跨语言的，不会涉及特定的某类语言的标准，它通过接口的设计概念去封装一系列 API 的相关功能。
+ 从标准上：OpenTracing 规范并不是官方标准，它的主体 Cloud Native Computing Foundation(CNCF)并不是官方的标准机构。
+ 总的来说，通过 OpenTracing  规范来对接追踪系统后，可以很方便的在不同的追踪系统中进行切换，它不会与具体的某一个供应商系统产生强捆绑关系。
+ 目前，市面上比较流行的追踪系统的思维模型均起源于 Google 的 *Dapper, a Large-Scale Distributed System Tracing Infrastructure* 论文(建议阅读)。OpenTracing  规范也不例外，它由一系列约定的术语概念知识，追踪系统中常见的3个术语含义如下表所示。

| 术语        | 含义       | 概述                                                         |
| ----------- | ---------- | ------------------------------------------------------------ |
| Trace       | 跟踪       | 一个 Trace 代表了一个事务或者流程在（分布式）系统中的执行过程 |
| Span        | 跨度       | 代表了一个事务中的每个工作单元，通常多个 Span 将会组成一个完整的 Trace |
| SpanContext | 跨度上下文 | 代表一个事务的相关跟踪信息，不同的 Span 会根据 OpenTracing  规范封装不同的属性，包含操作名称、开始时间和结束时间、标签信息、日志信息、上下文信息等。 |

### 2. Jaeger 的使用

`Jaeger` 是 `uber `开源的一个分布式链路追踪系统，收到了 `Google Dapper `和`OpenZipkin `的启发，目前由 CNCF 托管。它提供了分布式上下文传播，分布式交易监控、原因分析、服务依赖性分析、性能/延迟优化分析等等核心功能。

目前，市面上比较流行的分布式追踪系统都已经完全支持 OpenTracing   规范。

#### 安装 Jaeger 

Jaeger 官方提供了 all-in-one 的安装包，提供了 Docker 和已打包好的二进制文件，可直接通过 Docker 的方式安装并启动。

```bash
$ docker run -d --name Jaeger -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p 5775:5775/udp -p 6831:6831/udp -p 6832:6832/udp -p 5778:5778 -p 16686:16686 -p 14268:14268 -p 9411:9411 jaegertracing/all-in-one:1.16
```

通过上面的命令成功将 Jaeger 运行起来，在命令中应设立许多端口号，它们的作用如下表所示。

| 端口  | 协议 | 功能                                      |
| ----- | ---- | ----------------------------------------- |
| 5775  | UDP  | 以 compact 协议接受 `zipkin.thrift `数据  |
| 6831  | UDP  | 以 compact 协议接受 `jaeger .thrift `数据 |
| 6832  | UDP  | 以 binary 协议接受`jaeger .thrift `数据   |
| 5778  | HTTP | Jaeger 的服务配置端口                     |
| 16686 | HTTP | Jaeger 的 Web UI                          |
| 14268 | HTTP | 通过 Client 直接接收`jaeger .thrift `数据 |
| 9411  | HTTP | 兼容 `Zipkin `的 HTTP 端口                |

在确保命令运行没有问题后，只需打开浏览器，访问`localhost:16686`，就可以看到 Jaeger  的 Web UI 界面了，并且 Jaeger  的后端是使用 Go 语言开发的，如果有定制化需求，还可以进行二次开发。

![image-20220518213807417](https://raw.githubusercontent.com/tonshz/test/master/img/202205182138746.png)

### 3. 在应用中注入追踪

接下来将在应用程序中接入链路追踪的功能，最简单的需求就是每次调用时，都能够在链路追踪系统（本项目采用 Jaeger  ）上查看到对应的调用链信息。

#### a. 安装第三方库

安装两个第三方库，借助它们实现与追踪系统的对接，分别是 OpenTracing API 和 Jaeger  Client 的 Go 语言实现，命令如下：

```bash
$ go get -u github.com/opentracing/opentracing-go@v1.1.0
$ go get -u github.com/uber/jaeger-client-go@v2.22.1 
```

#### b. 编写 tracer

在`pkg`目录下新建`tracer`目录，并创建`tracer.go`。

```go
package tracer

import (
   "github.com/opentracing/opentracing-go"
   "github.com/uber/jaeger-client-go/config"
   "io"
   "time"
)

func NewJaegerTrace(serviceName, agentHostPort string) (opentracing.Tracer, io.Closer, error) {
   // config.Configuration 为 jaeger client 的配置项，主要设置应用的基本信息
   cfg := &config.Configuration{
      ServiceName: serviceName,
      Sampler: &config.SamplerConfig{
         Type: "const", // 固定采样，对所有数据都进行采样
         Param: 1, // 1 表示 true,0 表示 false
      },
      Reporter: &config.ReporterConfig{
         LogSpans: true, // 是否启用 LoggingReporter
         BufferFlushInterval: 1 * time.Second, // 刷新缓冲区的频率
         LocalAgentHostPort: agentHostPort, // 上报的 Agent 地址
      },
   }
   // cfg.NewTracer() 根据配置项初始化 Tracer 对象，返回 opentracing.Tracer，而不是特定供应商的追踪系统对象
   tracer, closer, err := cfg.NewTracer()
   if err != nil{
      return nil, nil, err
   }
   // opentracing.SetGlobalTracer() 设置全局的 Tracer 对象
   opentracing.SetGlobalTracer(tracer)
   return tracer, closer, nil
}
```

上述代码主要分为三部分：

+ `config.Configuration`:  为 jaeger client 的配置项，主要设置应用的基本信息，如 Sampler（固定采样，对所有数据都进行采样）、Reporter（是否启用 `LoggingReporter`、刷新缓冲区的频率、上报的 Agent 地址）等。
+ `cfg.NewTracer`: 根据配置项初始化 Tracer 对象，返回 `opentracing.Tracer`，而不是特定供应商的追踪系统对象。
+ `opentracing.SetGlobalTracer`: 设置全局的 Tracer 对象，根据实际情况设置即可。通常会统一使用一套追踪系统，因此该语句常常会被使用。

#### c. 初始化配置

在编写完 tracer 后，需要在 `global`目录下新建 `tracer.go`，新增如下全局对象。go

```go
package global

import "github.com/opentracing/opentracing-go"

var(
   Tracer opentracing.Tracer
)
```

并在 `main.go`中新增`setupTracer`方法的初始化逻辑。

```go
func init() {
   ...
   // 链路追踪初始化
   err = setupTracer()go
   if err != nil {
      log.Fatalf("init.setupTracer err: %v", err)
   }
}
...
func setupTracer() error {
	// 6831 端口 以 compact 协议接受 jaeger .thrift 数据
	jaegerTracer, _, err := tracer.NewJaegerTrace("blog_service", "127.0.0.1:6831")
	if err != nil {
		return err
	}
	global.Tracer = jaegerTracer
	return nil
}
```

上述代码的主要功能是调用先前编写的 tracer，并将其注入全局变量 Tracer 中，以便后续在中间件中使用，或是

在不同的自定义 Span 中打点使用。

#### d. 中间件

至此，tracer 的流程基本编写完成，剩下的问题便是如何将 gin 与 tracer 衔接起来，让每次的接口调用都能够被精确的上报到追踪系统中。可以在中间件中实现这个功能，使其成为一个标准。在 `internal/middleware`目录下新建 `tracer.go`。

```go
package middleware

import (
   "context"
   "demo/ch02/global"
   "github.com/gin-gonic/gin"
   "github.com/opentracing/opentracing-go"
   "github.com/opentracing/opentracing-go/ext"
)

// 返回入参为一个 gin 上下文的函数
func Tracing() func(c *gin.Context) {
   return func(c *gin.Context) {
      var newCtx context.Context
      var span opentracing.Span
      // Extract() 返回给定 `format` 和 `carrier` 的 SpanContext 实例。
      spanCtx, err := opentracing.GlobalTracer().Extract(
         // HTTPHeaders 将 SpanContexts 表示为 HTTP 标头字符串对。
         opentracing.HTTPHeaders,
         opentracing.HTTPHeadersCarrier(c.Request.Header))
      if err != nil {
         /*
            StartSpanFromContextWithTracer
            使用在上下文中找到的跨度作为 ChildOfRef 启动并返回一个带有 `operationName` 的跨度。
            如果不存在，它会创建一个根跨度。它还返回一个围绕返回的范围构建的 context.Context 对象。
         */
         span, newCtx = opentracing.StartSpanFromContextWithTracer(
            c.Request.Context(),
            global.Tracer,
            c.Request.URL.Path)
      } else {
         span, newCtx = opentracing.StartSpanFromContextWithTracer(
            c.Request.Context(),
            global.Tracer,
            // 相对路径
            c.Request.URL.Path,
            // ChildOf 返回一个指向依赖父跨度的 StartSpanOption
            opentracing.ChildOf(spanCtx),
            // Tag 可以作为 StartSpanOption 传递以将标签添加到新 Span，或者其 Set 方法可用于将标签应用于现有 Span
            /*
               StartSpanOptions 允许 Tracer.StartSpan() 调用者和实现者使用一种机制来覆盖开始时间戳、
               指定 Span 引用并在 Span 开始时使单个标签或多个标签可用。
            */
            opentracing.Tag{Key: string(ext.Component), Value: "HTTP"})
      }
      defer span.Finish()
      c.Request = c.Request.WithContext(newCtx)
      c.Next()
   }
}
```

修改`internal/routers`下的 `router.go`文件，在`NewRoutwe() `中新增中间件的注册逻辑：

```go
// 新增链路追踪中间件注册
r.Use(middleware.Tracing())
```

**需要注意的是，链路追踪中间件的注册应该在对所有路由方法调用之前生效，因此需要在路由注册行为之前注册即可。**

### 4. 验证跟踪情况

在完成上述所有步骤后，需要重新启动服务，访问一个通用接口，例如，获取鉴权 Token 的接口，再到 Jaeger Web UI（localhost:16686）上查看追踪情况。可以在 Service 处发现刚刚注册的项目名称`blog-service`，并且可以看到调用的鉴权 Token 接口的调用概览。

![image-20220518223242922](https://raw.githubusercontent.com/tonshz/test/master/img/202205182232236.png)

![image-20220518223302596](https://raw.githubusercontent.com/tonshz/test/master/img/202205182233726.png)

通过针对性的查询可以看到链路的详细分析，在多调用的情况下非常的直观、便捷，同时还可以进行一些内部工具的分析和报警。

![image-20220518223615965](https://raw.githubusercontent.com/tonshz/test/master/img/202205182236229.png)

### 5. 实现日志追踪

在可观察性的三大核心要素中，日志是必不可少的一个环节。但是日志量通常都很大，在出现问题时，去翻找日志过于麻烦。在有了链路追踪系统后，可以在记录日志的同时，将链路的 `SpanID` 和`traceID`也记录进去，这样就可以串联起该次请求的所有请求链路和日志信息情况，而且这个功能的实现并不难。只需在对应方法的第一个参数中传入上下文(context)，并在内部解析此上下文来获取联络信息即可。

#### a. 日志包含的信息

期望最终日志包含的所有信息如下:

```json
{
    "callers": ["xxxx"],
    "level": "info",
    "message": "access log: method: GET, status_code: 200, begin_time: ..., end_time: ...",
    "request": "xxx",
    "response": "xxx",
    "time": xxxx,
    "span_id": "xxxx",
    "trace_id": "xxxx"
}
```

#### b. 修改中间件

在`internal/middleware`下修改`tracer.go`文件。

```go
package middleware

import (
	"context"
	"demo/ch02/global"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
)

// 返回入参为一个 gin 上下文的函数
func Tracing() func(c *gin.Context) {
	return func(c *gin.Context) {
		....

		defer span.Finish()
		// 添加链路 TraceID 与 SpanID 的记录
		var traceID string
		var spanID string
		var spanContext = span.Context()
		// 此方法为类型断言
		// interface{}.(type) 使用类型断言断定某个接口是否是指定的类型
		// 该写法必须与switch case联合使用，case中列出实现该接口的类型。
		switch spanContext.(type) {
		case jaeger.SpanContext:
			// 将 spanContext 转换为 jaeger.SpanContext 类型并赋值给 jaegerContext
			// jaegerContext, flag := spanContext.(jaeger.SpanContext)
			// 其中 flag 为是否转换成功，该参数可省略
			jaegerContext := spanContext.(jaeger.SpanContext)
			traceID = jaegerContext.TraceID().String()
			spanID = jaegerContext.SpanID().String()
		}

        // 将俩个 ID 注册到上下文的元数据中
		c.Set("X-Trace-ID", traceID)
		c.Set("X-Span-ID", spanID)
		c.Request = c.Request.WithContext(newCtx)
		c.Next()
	}
}
```

在上述代码中，通过对` spanContext`的`jaeger.SpanContext`做断言，获取了 `SpanID`和`traceID`，并将其注册到上下文的元数据中。

```go
// 类型断言样例
func main() {
	var a interface{} = 10
	b := a.(int)
	fmt.Println("b: ", b)
	c, flag := a.(float64)
	fmt.Printf("c: %v, flag: %v\n", c, flag)
	d, flag := a.(int)
	fmt.Printf("d: %v, flag: %v", d, flag)
}
```

```bash
b:  10
c: 0, flag: false
d: 10, flag: true
```

#### c. 日志追踪

接下来修改日志文件，在公共库中，会默认将 context 作为函数的受参数传入，本项目为了加强实践性，故没有预置。在方法形参中新增上下文（context）的传入，并对日志中调用的`WithTrace`方法进行设置。修改`pkg/logger`目录下的`logger.go`文件。

```go
package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"runtime"
	"time"
)

...

// 设置日志上下文属性
func (l *Logger) WithContext(ctx context.Context) *Logger {
	ll := l.clone()
	ll.ctx = ctx
	return ll
}

// 从上下文中获取 trace_id 和 span_id
func (l *Logger) WithTrace() *Logger {
	ginCtx, ok := l.ctx.(*gin.Context)
	if ok {
		return l.WithFields(Fields{
			"trace_id": ginCtx.MustGet("X-Trace-ID"),
			"span_id":  ginCtx.MustGet("X-Span-ID"),
		})
	}
	return l
}
...

// 日志分级输出: Debug、Info、Warn、Error、Fatal、Panic
func (l *Logger) Debug(ctx context.Context, v ...interface{}) {
    // 添加了 WithContext(ctx) 与 WithTrace()
    // 此处 ctx 类型为 context.Context 而不是中间件中的 gin.Context
	l.WithContext(ctx).WithTrace().Output(LevelDebug, fmt.Sprint(v...))
}

func (l *Logger) Debugf(ctx context.Context, format string, v ...interface{}) {}
func (l *Logger) Info(ctx context.Context, v ...interface{}) {}
func (l *Logger) Infof(ctx context.Context, format string, v ...interface{}) {}
func (l *Logger) Warn(ctx context.Context, v ...interface{}) {}
func (l *Logger) Warnf(ctx context.Context, format string, v ...interface{}) {}
func (l *Logger) Error(ctx context.Context, v ...interface{}) {}
func (l *Logger) Errorf(ctx context.Context, format string, v ...interface{}) {}
func (l *Logger) Fatal(ctx context.Context, v ...interface{}) {}
func (l *Logger) Fatalf(ctx context.Context, format string, v ...interface{}) {}
func (l *Logger) Panic(ctx context.Context, v ...interface{}) {}
func (l *Logger) Panicf(ctx context.Context, format string, v ...interface{}) {}
```

在 `WithTrace()`中将存储在上下文中的 `SpanID`和`traceID`读取出来，然后写入到日志信息中。

另外，需要对入参校验和绑定的日志、路由处理方法中的日志、访问日志拦截器和异常抛出拦截器等进行修改，需要将 Context 传入日志方法中，如访问日志，修改`internal/middleware`目录下的`access_log.go`文件。

```go
...

func AccessLog() gin.HandlerFunc {
   return func(c *gin.Context) {
      ...
      // 将 Context 传入日志方法中
      global.Logger.WithFields(fields).Infof(c, "access log: method: %s, status_code: %d, begin_time: %d, end_time: %d",
         c.Request.Method,
         bodyWriter.Status(),
         beginTime,
         endTime,
      )
   }
}
```

这样写入日志后，默认就会带上链路信息，最终日志输出如下：

```bash
2022/05/18 23: 16: 21 {
    "callers": [
        "C:/Users/zyc/GolandProjects/demo/ch02/main.go: 31 main.init.0"
    ],
    "level": "error",
    "message": "app.BindAndValid errs: AppSecret为必填字段,AppKey为必填字段",
    "span_id": "3e48ca5a8986ab56",
    "time": 1652886981497630400,
    "trace_id": "3e48ca5a8986ab56"
}
```

#### d. 思考

日志分级中的 `context.Context` 类型与中间件中设置的上下文类型`gin.Context`不同。之所以能成功是因为`context.Context`是一个接口，而`gin.Context`实现了这个接口中定义的方法，因此在 Go 语言中认为两者时“
相同”的（多态）。

```go
// A Context carries a deadline, a cancellation signal, and other values across
// API boundaries.
//
// Context's methods may be called by multiple goroutines simultaneously.
type Context interface {
   // Deadline returns the time when work done on behalf of this context
   // should be canceled. Deadline returns ok==false when no deadline is
   // set. Successive calls to Deadline return the same results.
   Deadline() (deadline time.Time, ok bool)

   // Done returns a channel that's closed when work done on behalf of this
   // context should be canceled. Done may return nil if this context can
   // never be canceled. Successive calls to Done return the same value.
   // The close of the Done channel may happen asynchronously,
   // after the cancel function returns.
   //
   // WithCancel arranges for Done to be closed when cancel is called;
   // WithDeadline arranges for Done to be closed when the deadline
   // expires; WithTimeout arranges for Done to be closed when the timeout
   // elapses.
   //
   // Done is provided for use in select statements:
   //
   //  // Stream generates values with DoSomething and sends them to out
   //  // until DoSomething returns an error or ctx.Done is closed.
   //  func Stream(ctx context.Context, out chan<- Value) error {
   //     for {
   //        v, err := DoSomething(ctx)
   //        if err != nil {
   //           return err
   //        }
   //        select {
   //        case <-ctx.Done():
   //           return ctx.Err()
   //        case out <- v:
   //        }
   //     }
   //  }
   //
   // See https://blog.golang.org/pipelines for more examples of how to use
   // a Done channel for cancellation.
   Done() <-chan struct{}

   // If Done is not yet closed, Err returns nil.
   // If Done is closed, Err returns a non-nil error explaining why:
   // Canceled if the context was canceled
   // or DeadlineExceeded if the context's deadline passed.
   // After Err returns a non-nil error, successive calls to Err return the same error.
   Err() error

   // Value returns the value associated with this context for key, or nil
   // if no value is associated with key. Successive calls to Value with
   // the same key returns the same result.
   //
   // Use context values only for request-scoped data that transits
   // processes and API boundaries, not for passing optional parameters to
   // functions.
   //
   // A key identifies a specific value in a Context. Functions that wish
   // to store values in Context typically allocate a key in a global
   // variable then use that key as the argument to context.WithValue and
   // Context.Value. A key can be any type that supports equality;
   // packages should define keys as an unexported type to avoid
   // collisions.
   //
   // Packages that define a Context key should provide type-safe accessors
   // for the values stored using that key:
   //
   //     // Package user defines a User type that's stored in Contexts.
   //     package user
   //
   //     import "context"
   //
   //     // User is the type of value stored in the Contexts.
   //     type User struct {...}
   //
   //     // key is an unexported type for keys defined in this package.
   //     // This prevents collisions with keys defined in other packages.
   //     type key int
   //
   //     // userKey is the key for user.User values in Contexts. It is
   //     // unexported; clients use user.NewContext and user.FromContext
   //     // instead of using this key directly.
   //     var userKey key
   //
   //     // NewContext returns a new Context that carries value u.
   //     func NewContext(ctx context.Context, u *User) context.Context {
   //        return context.WithValue(ctx, userKey, u)
   //     }
   //
   //     // FromContext returns the User value stored in ctx, if any.
   //     func FromContext(ctx context.Context) (*User, bool) {
   //        u, ok := ctx.Value(userKey).(*User)
   //        return u, ok
   //     }
   Value(key interface{}) interface{}
}
```

### 6. 实现 SQL 追踪

既然有了链路追踪，那么 SQL 追踪也是必不可少的，因为 SQL 是大部分 Web 应用中的第一性能“杀手”。针对 SQL，由于项目使用的是 GORM，因此需要结合`Callback`和`Context`来实现链路中的 SQL Span “打点”。

首先在项目根目录下执行下述命令进行安装：

```bash
$ go get -u github.com/eddycjy/opentracing-gorm
```

修改`internal/model`目录下的`model.go`文件，新增 OpenTracing 相关的注册回调。

```go
package model

import (
   ...
   otgorm "github.com/eddycjy/opentracing-gorm"
   ...
)
...

// 新增 NewDBEngine()
func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
   ...

   // 设置空闲连接池中的最大连接数
   db.DB().SetMaxIdleConns(databaseSetting.MaxIdleConns)
   db.DB().SetMaxOpenConns(databaseSetting.MaxOpenConns)
   // 添加 OpenTracing 注册回调
   otgorm.AddGormCallbacks(db)
   return db, nil
}
...
```

每一条链路都相当于一次请求，而每一次请求的上下文都是不一样的，因此每一次请求都需要将上下文注册进来（官方建议上下文参数放在函数的第一个参数），以确保联络联通。

修改`internal/service`目录下的`service.go`，新增数据库连接实例的上下文信息注册。

```go
package service

import (
   "context"
   "demo/ch02/global"
   "demo/ch02/internal/dao"
   otgorm "github.com/eddycjy/opentracing-gorm"
)

type Service struct {
   ctx context.Context
   dao *dao.Dao
}

func New(ctx context.Context) Service {
   svc := Service{ctx: ctx}
   //svc.dao = dao.New(global.DBEngine)
   // 新增数据库连接实例的上下文信息注册
   svc.dao = dao.New(otgorm.WithContext(svc.ctx, global.DBEngine))
   return svc
}
```

在完成 Callback 和 Context 的注册和设置后，重新运行应用程序，可以看到默认输出了对应的 Callback 注册信息，至此就完成了 SQL  的默认追踪。

打开 Jaeger Web UI（localhost:16686）查看追踪情况。

![image-20220518233839101](https://raw.githubusercontent.com/tonshz/test/master/img/202205182338340.png)

在途中一共有三个 Span，分别是本次请求的 HTTP 路由`ParentSpan`以及两条 SQL Span，出现两条的原因是因为执行了两条不一样的 SQL，并且两者之间是有顺序的，一个在前一个在后。如果是并发执行，则可能会出现多条 Span 重叠在同一个区域的时间轴上。

---------------

## 参考

+ [参考文章](https://golang2.eddycjy.com/)

+ [GitHub代码](https://github.com/go-programming-tour-book)

