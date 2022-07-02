# Go 语言编程之旅(二)：HTTP 应用(七) 

## 九、应用中间件

针对不同的环境，应该进行一些特殊的调整，而往往这些都是有规律可依的，一些常用的应用中间件就可以妥善的解决这些问题，接下来将会去编写一些在项目中比较常见的应用中间件。

### 1.访问日志记录

在出问题时，常常会需要去查日志，那么除了查错误日志、业务日志以外，还有一个很重要的日志类别，就是访问日志，从功能上来讲，它最基本的会记录每一次请求的请求方法、方法调用开始时间、方法调用结束时间、方法响应结果、方法响应结果状态码，更进一步的话，会记录 `RequestId`、`TraceId`、`SpanId `等等附加属性，以此来达到日志链路追踪的效果，如下图：

![image](https://raw.githubusercontent.com/tonshz/test/master/img/202205162046559.jpeg)

但是在正式开始前，又会遇到一个问题，没办法非常直接的获取到方法所返回的响应主体，这时候需要巧妙利用 Go interface 的特性，实际上在写入流时，调用的是 `http.ResponseWriter`，如下：

```go
// A ResponseWriter interface is used by an HTTP handler to
// construct an HTTP response.
//
// A ResponseWriter may not be used after the Handler.ServeHTTP method
// has returned.
type ResponseWriter interface {
   // Header returns the header map that will be sent by
   // WriteHeader. The Header map also is the mechanism with which
   // Handlers can set HTTP trailers.
   //
   // Changing the header map after a call to WriteHeader (or
   // Write) has no effect unless the modified headers are
   // trailers.
   //
   // There are two ways to set Trailers. The preferred way is to
   // predeclare in the headers which trailers you will later
   // send by setting the "Trailer" header to the names of the
   // trailer keys which will come later. In this case, those
   // keys of the Header map are treated as if they were
   // trailers. See the example. The second way, for trailer
   // keys not known to the Handler until after the first Write,
   // is to prefix the Header map keys with the TrailerPrefix
   // constant value. See TrailerPrefix.
   //
   // To suppress automatic response headers (such as "Date"), set
   // their value to nil.
   Header() Header

   // Write writes the data to the connection as part of an HTTP reply.
   //
   // If WriteHeader has not yet been called, Write calls
   // WriteHeader(http.StatusOK) before writing the data. If the Header
   // does not contain a Content-Type line, Write adds a Content-Type set
   // to the result of passing the initial 512 bytes of written data to
   // DetectContentType. Additionally, if the total size of all written
   // data is under a few KB and there are no Flush calls, the
   // Content-Length header is added automatically.
   //
   // Depending on the HTTP protocol version and the client, calling
   // Write or WriteHeader may prevent future reads on the
   // Request.Body. For HTTP/1.x requests, handlers should read any
   // needed request body data before writing the response. Once the
   // headers have been flushed (due to either an explicit Flusher.Flush
   // call or writing enough data to trigger a flush), the request body
   // may be unavailable. For HTTP/2 requests, the Go HTTP server permits
   // handlers to continue to read the request body while concurrently
   // writing the response. However, such behavior may not be supported
   // by all HTTP/2 clients. Handlers should read before writing if
   // possible to maximize compatibility.
   Write([]byte) (int, error)

   // WriteHeader sends an HTTP response header with the provided
   // status code.
   //
   // If WriteHeader is not called explicitly, the first call to Write
   // will trigger an implicit WriteHeader(http.StatusOK).
   // Thus explicit calls to WriteHeader are mainly used to
   // send error codes.
   //
   // The provided code must be a valid HTTP 1xx-5xx status code.
   // Only one header may be written. Go does not currently
   // support sending user-defined 1xx informational headers,
   // with the exception of 100-continue response header that the
   // Server sends automatically when the Request.Body is read.
   WriteHeader(statusCode int)
}
```

所以只需要写一个针对访问日志的 Writer 结构体，实现特定的 Write 方法就可以解决无法直接取到方法响应主体的问题了。在 `internal/middleware`下新建 `access_log.go `文件，写入如下代码：

```go
package middleware

import (
   "bytes"
   "github.com/gin-gonic/gin"
)

// 针对访问日志的 Writer 结构体
type AccessLogWriter struct {
   gin.ResponseWriter
   body *bytes.Buffer
}

// 实现特定的 Write 方法
func (w AccessLogWriter) Write(p []byte) (int, error) {
   // 获取响应主体
   if n, err := w.body.Write(p); err != nil {
      return n, err
   }
   return w.ResponseWriter.Write(p)
}
```

在` AccessLogWriter `的 Write 方法中，实现了双写，可以直接通过` AccessLogWriter` 的 body 取到值，接下来继续编写访问日志的中间件，写入如下代码：

```go
func AccessLog() gin.HandlerFunc {
   return func(c *gin.Context) {
      // 初始化 AccessLogWriter
      bodyWriter := &AccessLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
      // 将其赋值给当前的 Writer 写入流
      c.Writer = bodyWriter
      beginTime := time.Now().Unix()
      c.Next()
      endTime := time.Now().Unix()
      fields := logger.Fields{
         "request":  c.Request.PostForm.Encode(),
         "response": bodyWriter.body.String(),
      }
      global.Logger.WithFields(fields).Infof(c, "access log: method: %s, status_code: %d, begin_time: %d, end_time: %d",
         c.Request.Method,
         bodyWriter.Status(),
         beginTime,
         endTime,
      )
   }
}
```

在` AccessLog `方法中，初始化了 `AccessLogWriter`，将其赋予给当前的 Writer 写入流（可理解为替换原有），并且通过指定方法得到所需的日志属性，最终写入到日志中去，其中涉及到了如下信息：

- method：当前的调用方法。
- request：当前的请求参数。
- response：当前的请求结果响应主体。
- status_code：当前的响应结果状态码。
- begin_time/end_time：调用方法的开始时间，调用方法结束的结束时间。

### 2. 异常捕获处理 

在异常造成的恐慌发生时，开发者一定不在现场，因为不能随时随地的盯着控制台，在常规手段下也不知道它几时有可能发生，因此对于异常的捕获和及时的告警通知是非常重要的，而发现这些可能性的手段有非常多，本次采取的是最简单的捕获和告警通知，如下图：

![image](https://raw.githubusercontent.com/tonshz/test/master/img/202205162058607.jpeg)

#### a. 自定义 Recovery

在前文中看到 gin 本身已经自带了一个 Recovery 中间件，但是在项目中，需要针对我们的公司内部情况或生态圈定制 Recovery 中间件，确保异常在被正常捕抓之余，要及时的被识别和处理，因此自定义一个 Recovery 中间件是非常有必要的，在 `internal/middleware`下新建 `recovery.go `文件。

```go
package middleware

import (
   "demo/ch02/global"
   "demo/ch02/pkg/app"
   "demo/ch02/pkg/errcode"
   "github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
   return func(c *gin.Context) {
      defer func() {
         if err := recover(); err != nil {
            global.Logger.WithCallersFrames().Errorf(c, "panic recover err: %v", err)
            app.NewResponse(c).ToErrorResponse(errcode.ServerError)
            c.Abort()
         }
      }()
      c.Next()
   }
}
```

#### b. 邮件报警处理

另外在实现 Recovery 的同时，需要实现一个简单的邮件报警功能，确保出现 Panic 后，在捕抓之余能够通过邮件报警来及时的通知到对应的负责人。

##### 安装

在项目目录下执行安装命令：

```bash
$ go get -u gopkg.in/gomail.v2
```

`Gomail `是一个用于发送电子邮件的简单又高效的第三方开源库，目前只支持使用 SMTP 服务器发送电子邮件，但是其 API 较为灵活，如果有其它的定制需求也可以轻易地借助其实现，这恰恰好符合\需求，因为目前只需要一个小而美的发送电子邮件的库就可以了。

##### 邮件工具库

在项目目录 `pkg` 下新建 email 目录并创建 `email.go `文件，需要针对发送电子邮件的行为进行一些封装，写入如下代码：

```go
package email

import (
   "crypto/tls"
   "gopkg.in/gomail.v2"
)

type Email struct {
   *SMTPInfo
}

// 定义 SMTPinfo 结构体用于传递发送邮件所必需的信息
type SMTPInfo struct {
   Host     string
   Port     int
   IsSSL    bool
   UserName string
   Password string
   From     string
}

func NewEmail(info *SMTPInfo) *Email {
   return &Email{SMTPInfo: info}
}

func (e *Email) SendMail(to []string, subject, body string) error {
   // gomail.NewMessage() 创建一个消息实例
   m := gomail.NewMessage()
   // 设置邮件的一些必要信息
   m.SetHeader("From", e.From) // 发件人
   m.SetHeader("To", to...) // 收件人 to...: 将 to 切片打散
   m.SetHeader("Subject", subject) // 邮件主题
   m.SetBody("text/html", body) // 邮件正文
   
   // gomail.NewDialer() 创建一个新的 SMTP 拨号实例，设置对应的拨号信息用于连接 SMTP 服务器
   dialer := gomail.NewDialer(e.Host, e.Port, e.UserName, e.Password)
   dialer.TLSConfig = &tls.Config{InsecureSkipVerify: e.IsSSL}
   // DialAndSend() 打开与 SMTP 服务器的连接并发送电子邮件
   return dialer.DialAndSend(m)
}
```

在上述代码中，定义了` SMTPInfo `结构体用于传递发送邮箱所必需的信息，而在` SendMail `方法中，首先调用 `NewMessage `方法创建一个消息实例，可以用于设置邮件的一些必要信息，分别是：

- 发件人（From）
- 收件人（To）
- 邮件主题（Subject）
- 邮件正文（Body）

在完成消息实例的基本信息设置后，调用 `NewDialer `方法创建一个新的 SMTP 拨号实例，设置对应的拨号信息用于连接 SMTP 服务器，最后再调用 `DialAndSend `方法打开与 SMTP 服务器的连接并发送电子邮件。

#####  初始化配置信息

```lua
config/config.yaml => pkg/setting/section.go => global/setting.go => main.go
```

本次要做的发送电子邮件的行为，实际上可以理解是与一个 SMTP 服务进行交互，那么除了自建 SMTP 服务器以外，可以使用目前市面上常见的邮件提供商，它们也是有提供 SMTP 服务的，修改项目`config`目录下的配置文件 `config.yaml`，新增如下 Email 的配置项：

```yaml
# Email 初始化配置
Email:
  Host: smtp.163.com
  Port: 465
  UserName: dove_zyc@163.com
  Password: TYTJRVWRXGAHGTBW # 获取的 SMTP 密码
  IsSSL: true
  From: dove_zyc@163.com
  To:
    - 1103592040@qq.com
```

接下来在项目目录 `pkg/setting` 的 `section.go` 文件中，新增对应的 Email 配置项，如下：

```go
// Email 结构体
type EmailSettingS struct {
   Host     string
   Port     int
   UserName string
   Password string
   IsSSL    bool
   From     string
   To       []string
}
```

并在在项目目录 `global` 的 `setting.go `文件中，新增 Email 对应的配置全局对象，如下：

```go
var (
   ...
   EmailSetting    *setting.EmailSettingS
)
```

最后就是在项目根目录的` main.go `文件的 `setupSetting `方法中，新增 Email 配置项的读取和映射，如下：

```go
...

func setupSetting() error {
   ...
   err = settings.ReadSection("Email", &global.EmailSetting)
   if err != nil{
      return err
   }
   ...
}
...
```

##### 编写中间件

在项目目录 `internal/middleware` 下创建` recovery.go` 文件。

```go
package middleware

import (
   "demo/ch02/global"
   "demo/ch02/pkg/app"
   "demo/ch02/pkg/email"
   "demo/ch02/pkg/errcode"
   "fmt"
   "github.com/gin-gonic/gin"
   "time"
)

func Recovery() gin.HandlerFunc {
   // 转换为邮件结构体
   defailtMailer := email.NewEmail(&email.SMTPInfo{
      Host:     global.EmailSetting.Host,
      Port:     global.EmailSetting.Port,
      IsSSL:    global.EmailSetting.IsSSL,
      UserName: global.EmailSetting.UserName,
      Password: global.EmailSetting.Password,
      From:     global.EmailSetting.From,
   })
   return func(c *gin.Context) {
      defer func() {
         if err := recover(); err != nil {
            global.Logger.WithCallersFrames().Errorf(c, "panic recover err: %v", err)
            // 发送邮件
            err := defailtMailer.SendMail(
               global.EmailSetting.To,
               fmt.Sprintf("异常抛出，发生时间: %d", time.Now().Unix()),
               fmt.Sprintf("错误信息: %v", err),
            )
            if err != nil {
               global.Logger.Panicf(c, "mail.SendMail err: %v", err)
            }
            
            app.NewResponse(c).ToErrorResponse(errcode.ServerError)
            c.Abort()
         }
      }()
      c.Next()
   }
}
```

在本项目中，由于 Mailer（发信人） 是固定的，因此直接将其定义为了` defailtMailer`，接着在捕获到异常后调用 `SendMail `方法进行预警邮件发送，效果如下：

![image-20220516214122122](https://raw.githubusercontent.com/tonshz/test/master/img/202205162141164.png)

具体的邮件模板可以根据实际情况进行定制。

### 3. 服务信息存储

平时经常会需要在进程内上下文设置一些内部信息，例如是应用名称和应用版本号这类基本信息，也可以是业务属性的信息存储，例如是根据不同的租户号获取不同的数据库实例对象，这时候就需要有一个统一的地方处理，如下图：

![image](https://raw.githubusercontent.com/tonshz/test/master/img/202205162144750.jpeg)

在`internal/middleware` 目录下新建 `app_info.go` 文件。

```go
package middleware

import "github.com/gin-gonic/gin"

func AppInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("app_name", "blog-service")
		c.Set("app_version", "1.0.0")
		/*
			Next 应该只在中间件内部使用。
			它执行调用处理程序内链中的待处理处理程序。
			请参阅 GitHub 中的示例
		*/
		c.Next()
	}
}
```

在上述代码中需要用到 `gin.Context `所提供的 setter 和 getter，在 gin 中称为元数据管理（Metadata Management），大致如下：

```go
/************************************/
/******** METADATA MANAGEMENT********/
/************************************/

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes  c.Keys if it was not used previously.
func (c *Context) Set(key string, value interface{}) {
   c.mu.Lock()
   if c.Keys == nil {
      c.Keys = make(map[string]interface{})
   }

   c.Keys[key] = value
   c.mu.Unlock()
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (c *Context) Get(key string) (value interface{}, exists bool) {
   c.mu.RLock()
   value, exists = c.Keys[key]
   c.mu.RUnlock()
   return
}

func (c *Context) MustGet(key string) interface{} {...}
func (c *Context) GetString(key string) (s string) {...}
func (c *Context) GetBool(key string) (b bool) {...}
func (c *Context) GetInt(key string) (i int) {...}
func (c *Context) GetInt64(key string) (i64 int64) {...}
func (c *Context) GetFloat64(key string) (f64 float64) {...}
func (c *Context) GetTime(key string) (t time.Time) {...}
func (c *Context) GetDuration(key string) (d time.Duration) {...}
func (c *Context) GetStringSlice(key string) (ss []string) {...}
func (c *Context) GetStringMap(key string) (sm map[string]interface{}) {...}
func (c *Context) GetStringMapString(key string) (sms map[string]string) {...}
func (c *Context) GetStringMapStringSlice(key string) (smss map[string][]string) {...}
```

实际上可以看到在 gin 中的 metadata，其实就是利用内部实现的 `gin.Context `中的 Keys 进行存储的，并配套了多种类型的获取和设置方法，相当的方便。另外可以注意到在默认的 Get 和 Set 方法中，传入和返回的都是 interface 类型，实际在业务属性的初始化逻辑处理中，可以通过对返回的 interface 进行类型断言，就可以获取到所需要的类型了。**在需要使用设置的应用信息时可通过`c.Get("app_name")`获取信息，若存在方法返回`blog-service true`，不存在则返回`nil false`。**

### 4. 接口限流控制

在应用程序的运行过程中，会不断地有新的客户端进行访问，而有时候会突然出现流量高峰（例如：营销活动），如果不及时进行削峰，资源整体又跟不上，那就很有可能会造成事故，因此常常会才有多种手段进行限流削峰，而针对应用接口进行限流控制就是其中一种方法，如下图：

![image](https://raw.githubusercontent.com/tonshz/test/master/img/202205162202769.jpeg)

#### a. 安装

```bash
$ go get -u github.com/juju/ratelimit@v1.0.1
```

`Ratelimit `提供了一个简单又高效的令牌桶实现，能够提供大量的方法帮助开发者实现限流器的逻辑。

#### b. 限流控制

##### LimiterIface

在 `pkg/limiter` 目录下新建`limiter.go` 文件。

```go
package limiter

import (
   "github.com/gin-gonic/gin"
   "github.com/juju/ratelimit"
   "time"
)

// 声明了 LimiterIface 接口，用于定义当前限流器所必须要的方法
// 由于限流器的策略不同，需要定义通用的接口以保证接口的设计
type LimiterIface interface {
   // 获取对应限流器的键值对名称
   Key(c *gin.Context) string
   // 获取令牌桶
   GetBucket(key string) (*ratelimit.Bucket, bool)
   // 新增多个令牌桶
   AddBucket(rules ...LimiterBucketRule) LimiterIface
}

type Limiter struct {
   // 存储令牌桶和键值对名称的映射关系
   limiterBuckets map[string]*ratelimit.Bucket
}

// 定义 LimiterBucketRule 结构体用于存储令牌桶的规则属性
type LimiterBucketRule struct {
	// 自定义键值对名称
	Key string
	// 间隔多长时间释放 N 个令牌
	FillInterval time.Duration
	// 令牌桶的容量
	Capacity int64
	// 每次到达间隔时间后所释放的具体令牌数量
	Quantum int64
}
```

在上述代码中，声明了` LimiterIface `接口，用于定义当前限流器所必须要的方法。

之所以这样做的原因是因为限流器是存在多种实现的，可能某一类接口需要限流器 A，另外一类接口需要限流器 B，所采用的策略不是完全一致的，因此需要声明` LimiterIface` 这类通用接口，保证其接口的设计，初步的在 `Iface` 接口中，一共声明了三个方法，如下：

- `Key`：获取对应的限流器的键值对名称。
- `GetBucket`：获取令牌桶。
- `AddBuckets`：新增多个令牌桶。

同时定义 `Limiter `结构体用于存储令牌桶与键值对名称的映射关系，并定义 `LimiterBucketRule` 结构体用于存储令牌桶的一些相应规则属性，如下：

- `Key`：自定义键值对名称。
- `FillInterval`：间隔多久时间释放 N 个令牌。
- `Capacity`：令牌桶的容量。
- `Quantum`：每次到达间隔时间后所释放的具体令牌数量。

至此就完成了一个` Limter `最基本的属性定义了，接下来将针对不同的情况实现这个项目中的限流器。

##### MethodLimiter

第一个编写的简单限流器的主要功能是针对路由进行限流，因为在项目中，可能只需要对某一部分的接口进行流量调控，在 `pkg/limiter` 目录下新建 `method_limiter.go `文件，写入如下代码：

```go
package limiter

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"strings"
)

type MethodLimiter struct {
	*Limiter
}

func NewMethodLimiter() LimiterIface {
	return MethodLimiter{
		Limiter: &Limiter{limiterBuckets: make(map[string]*ratelimit.Bucket)},
	}
}

// 实现 LimiterIface 接口定义的三个通用方法
// 根据 RequestURI 切割出核心路由作为键值对名称
func (l MethodLimiter) Key(c *gin.Context) string {
	uri := c.Request.RequestURI
	// strings.Index() 返回子串的第一个实例，若不存在则返回 -1
	index := strings.Index(uri, "?")
	if index == -1 {
		return uri
	}

	return uri[:index]
}

// 获取 Bucket 方法实现
func (l MethodLimiter) GetBucket(key string) (*ratelimit.Bucket, bool) {
	bucket, ok := l.limiterBuckets[key]
	return bucket, ok
}

// 设置 Bucket 方法实现
func (l MethodLimiter) AddBucket(rules ...LimiterBucketRule) LimiterIface {
	for _, rule := range rules {
		if _, ok := l.limiterBuckets[rule.Key]; !ok {
			l.limiterBuckets[rule.Key] = ratelimit.NewBucketWithQuantum(rule.FillInterval, rule.Capacity, rule.Quantum)
		}
	}

	return l
}

```

在上述代码中，针对 `LimiterIface` 接口实现了`MethodLimiter` 限流器，**主要逻辑是在 Key 方法中根据 `RequestURI `切割出核心路由作为键值对名称，**并在` GetBucket` 和` AddBuckets` 进行获取和设置 Bucket 的对应逻辑。

#### c. 编写中间件

在完成了限流器的逻辑编写后，在`internal/middleware` 目录下新建 `limiter.go` 文件，将整体的限流器与对应的中间件逻辑串联起来。

```go
package middleware

import (
   "demo/ch02/pkg/app"
   "demo/ch02/pkg/errcode"
   "demo/ch02/pkg/limiter"
   "github.com/gin-gonic/gin"
)

// 入参为 LimiterIface 接口类型，这样只要符合该接口类型的具体限流器实现都可以传入使用
func RateLimiter(l limiter.LimiterIface) gin.HandlerFunc {
   return func(c *gin.Context) {
      key := l.Key(c)
      if bucket, ok := l.GetBucket(key); ok {
         /*
            TakeAvailable()
            占用存储桶中立即可用的令牌的数量。
            返回值为删除的令牌数，
            如果没有可用的令牌，则返回零。它不会阻塞。
         */
         count := bucket.TakeAvailable(1)
         if count == 0 {
            response := app.NewResponse(c)
            response.ToErrorResponse(errcode.TooManyRequests)
            c.Abort()
            return
         }
      }
      c.Next()
   }
}
```

在` RateLimiter `中间件中，需要注意的是入参应该为 `LimiterIface `接口类型，这样子的话只要符合该接口类型的具体限流器实现都可以传入并使用，另外比较重要的就是` TakeAvailable` 方法，它会占用存储桶中立即可用的令牌的数量，返回值为删除的令牌数，如**果没有可用的令牌，将会返回 0，也就是已经超出配额了**，因此这时候将返回` errcode.TooManyRequest` 状态告诉客户端需要减缓并控制请求速度。

![image-20220516231510575](https://raw.githubusercontent.com/tonshz/test/master/img/202205162315612.png)

### 5. 统一超时控制

在应用程序的运行中，常常会遇到一个头疼的问题，调用链如果是`应用 A => 应用 B =>应用 C`，那如果应用 C 出现了问题，在没有任何约束的情况下持续调用，就会导致应用 A、B、C 均出现问题，也就是很常见的上下游应用的互相影响，导致连环反应，最终使得整个集群应用出现一定规模的不可用，如下图：

![image](https://raw.githubusercontent.com/tonshz/test/master/img/202205162233545.jpeg)

为了规避这种情况，最简单也是最基本的一个约束点，那就是统一的在应用程序中针对所有请求都进行一个最基本的超时时间控制，如下图：

![image](https://raw.githubusercontent.com/tonshz/test/master/img/202205162233926.jpeg)

因此可以编写一个上下文超时控制的中间件来实现这个需求，在`internal/middleware` 目录下新建 `context_timeout.go` 文件。

```go
package middleware

import (
   "context"
   "github.com/gin-gonic/gin"
   "time"
)

func ContextTimeout(t time.Duration) func(c *gin.Context) {
   return func(c *gin.Context) {
      /*
         before c.Next()
      */
      // 使用 context.WithTimeout() 设置当前 context 的超时时间
      // 返回值 Context, CancelFunc，即 cancel 值是一个取消函数
      ctx, cancel := context.WithTimeout(c.Request.Context(), t)
      defer cancel()
      // 将设置了超时时间的 ctx 赋值给当前的 context
      c.Request = c.Request.WithContext(ctx)
      // c.Next() 使用后会使得在其之后的代码执行在 main() 之后
      c.Next()
      /*
         after c.Next()
         此处的代码在 main() 中代码处理完成后再执行
         即执行顺序为： before => main => after
         若不使用 c.Next() 执行顺序为： before => after => main
      */
   }
}
```

-------------------

#### ==c.Next()==

使用 `c.Next()`后程序流转过程：`before => main => after`

```go
// 注意此处的注释
func ContextTimeout(t time.Duration) func(c *gin.Context) {
   return func(c *gin.Context) {
      /*
         before c.Next()
      */
      // 使用 context.WithTimeout() 设置当前 context 的超时时间
      // 返回值 Context, CancelFunc，即 cancel 值是一个取消函数
      ctx, cancel := context.WithTimeout(c.Request.Context(), t)
      defer cancel()
      // 将设置了超时时间的 ctx 赋值给当前的 context
      c.Request = c.Request.WithContext(ctx)
      // c.Next() 使用后会使得在其之后的代码执行在 main() 之后
      c.Next()
      /*
         after c.Next()
         此处的代码在 main() 中代码处理完成后再执行
         即执行顺序为： before => main => after
         若不使用 c.Next() 执行顺序为： before => after => main
      */
   }
}
```

--------------

在上述代码中，调用了` context.WithTimeout`方法设置当前 context 的超时时间，并重新赋予给了 `gin.Context`，这样子在当前请求运行到指定的时间后，在使用了该 context 的运行流程就会针对 context 所提供的超时时间进行处理，并在指定的时间进行取消行为。效果如下：

```go
_, err := ctxhttp.Get(c.Request.Context(), http.DefaultClient, "https://www.google.com/")
if err != nil {
    log.Fatalf("ctxhttp.Get err: %v", err)
}
```

需要传入设置了超时的`c.Request.Context()`，在验证时可以将默认超时时间调短。

```bash
ctxhttp.Get err: context deadline exceeded
exit status 1
```

最后由于已经到达了截止时间，因此返回 `context deadline exceeded` 错误提示信息。另外这里还需要注意，**如果在进行多应用/服务的调用时**，把父级的上下文信息（`ctx`）不断地传递下去，那么在统计超时控制的中间件中所设置的超时时间，其实是针对整条链路的，而不是针对单单每一条，如果需要针对额外的链路进行超时时间的调整，那么只需要调用像 `context.WithTimeout` 等方法对父级 `ctx `进行设置，然后取得子级` ctx`，再进行新的上下文传递就可以了。==**单个服务该超时时间不生效，需要存在多个服务之间的调用。**==

### 6. 注册中间件

在完成一连串的通用中间件编写后，修改 `internal/routers` 下的 `router.go` 文件，修改注册应用中间件的逻辑。

```go
package routers

import (
	...
)

// 对认证接口进行限流
var methodLimiters = limiter.NewMethodLimiter().AddBucket(limiter.LimiterBucketRule{
	Key:          "/auth",          // 对接受的 auth 请求进行限流
	FillInterval: 10 * time.Second, // 对 10s 内接受的请求进行限流
	Capacity:     10,               // 10s 内最多接受 10 个
	Quantum:      10,               // 每过 10s 将令牌桶的数量减少 10 #{Quantum}
})

func NewRouter() *gin.Engine {
	// 或者 r := gin.Default() 也可
	r := gin.New()
	// 根据不同的部署环境进行了应用中间件的设置
	// 使用了自定义的 Logger 和 Recovery 后就不需要使用 gin 原生提供的组件了
	if global.ServerSetting.RunMode == "debug" {
		// 在本地开发环境中，可能没有全应用生态圈，故作特殊处理
		r.Use(gin.Logger())
		// 在注册顺序上需要注意，类似 Recovery 这类的中间件应当尽可能早的注册
		r.Use(gin.Recovery())
	} else {
		r.Use(middleware.AccessLog())
		r.Use(middleware.Recovery())
	}

	// 新增限流控制中间件的注册
	r.Use(middleware.RateLimiter(methodLimiters))
	// 新增统一超时控制中间件注册：DefaultContextTimeout 参数需要先在 yml 文件中设置
	r.Use(middleware.ContextTimeout(global.AppSetting.DefaultContextTimeout))
	// 新增中间件 Translations 的注册
	r.Use(middleware.Translations())
    // 新增应用信息中间件注册
	r.Use(middleware.AppInfo())

	...
	apiv1.Use(middleware.JWT())
	// 上面花括号是代表中间的语句属于一个空间内，不受外界干扰，可去掉
	{
		...
	}
	return r
}
```

在上述代码中，根据不同的部署环境（`RunMode`）进行了应用中间件的设置，因为实际上在使用了自定义的 `Logger` 和 `Recovery `后，就没有必要使用 gin 原有所提供的了，而在本地开发环境中，可能没有齐全应用生态圈，因此需要进行特殊处理。另外在常规项目中，自定义的中间件不仅包含了基本的功能，还包含了很多定制化的功能，**同时在注册顺序上也注意，Recovery 这类应用中间件应当尽可能的早注册，这根据实际所要应用的中间件情况进行顺序定制就可以了。**

原代码中 `middleware.ContextTimeout` 是写死的，可以对其进行配置化（映射配置和秒数初始化），将超时的时间配置调整到配置文件中，而不是在代码中硬编码，最终结果应当如下：

```go
r.Use(middleware.ContextTimeout(global.AppSetting.DefaultContextTimeout))
```

这样子的话，以后修改超时的时间就只需要通过修改配置文件就可以解决了，不需要人为的修改代码，甚至可以不需要开发人员的直接参与，让运维同事确认后直接修改即可。

----------------

## 参考

+ [参考文章](https://golang2.eddycjy.com/)

+ [GitHub代码](https://github.com/go-programming-tour-book)



