# Go 语言编程之旅(二)：HTTP 应用

## 一、开启博客之路

### 1. 快速启动

```go
package main

import (
   "github.com/gin-gonic/gin"
   "log"
)

func main() {
   // 声明一个默认路由
   r := gin.Default()
   r.GET("/ping", func(c *gin.Context){
      c.JSON(200, gin.H{"message": "pong"})
   })
   err := r.Run()
   if err != nil {
      log.Fatalln(err)
   }
}
```

```bash
$ go run main.go
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

# /ping 路由注册成功
[GIN-debug] GET    /ping                     --> main.main.func1 (3 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Environment variable PORT is undefined. Using port :8080 by default
[GIN-debug] Listening and serving HTTP on :8080

```

启动服务后，输出的运行信息主要分为四大块：

- 默认 Engine 实例：当前默认使用了官方所提供的 Logger 和 Recovery 中间件创建了 Engine 实例。
- 运行模式：当前为调试模式，并建议若在生产环境时切换为发布模式。
- 路由注册：注册了 `GET /ping` 的路由，并输出其调用方法的方法名。
- 运行信息：本次启动时监听 8080 端口，由于没有设置端口号等信息，因此默认为 8080。

### 2. 分析

![image](https://raw.githubusercontent.com/tonshz/test/master/img/202205102056391.jpeg)

#### a. gin.Default

```go
// Default returns an Engine instance with the Logger and Recovery middleware already attached.
func Default() *Engine {
   debugPrintWARNINGDefault() // 此处会检查 Go 版本是否达到了 gin 的最低要求
   engine := New()
   engine.Use(Logger(), Recovery())
   return engine
}
```

可以通过调用`gin.Default()`来创建默认的 Engine 实例，它会在初始化阶段引入 Logger 和 Recovery 中间件，从而保障应用程序的基本运行。

+ Logger: 输出请求日志，并标准化日志格式
+ Recovery: 异常捕获，针对每次请求进行 recovery 处理，防止出现 panic 导致服务崩溃，并标准化异常日志格式。

另外在调用 `debugPrintWARNINGDefault()` 方法时，会检查 Go 版本是否达到 gin 的最低要求，再进行调试的日志 `[WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.` 的输出，用于提醒开发人员框架内部已经默认检查和集成了缺省值。

#### b. gin.New

```go
// New returns a new blank Engine instance without any middleware attached.
// By default the configuration is:
// - RedirectTrailingSlash:  true
// - RedirectFixedPath:      false
// - HandleMethodNotAllowed: false
// - ForwardedByClientIP:    true
// - UseRawPath:             false
// - UnescapePathValues:     true
func New() *Engine {
   debugPrintWARNINGNew()
   engine := &Engine{
      RouterGroup: RouterGroup{
         Handlers: nil,
         basePath: "/",
         root:     true,
      },
      FuncMap:                template.FuncMap{},
      RedirectTrailingSlash:  true,
      RedirectFixedPath:      false,
      HandleMethodNotAllowed: false,
      ForwardedByClientIP:    true,
      RemoteIPHeaders:        []string{"X-Forwarded-For", "X-Real-IP"},
      TrustedPlatform:        defaultPlatform,
      UseRawPath:             false,
      RemoveExtraSlash:       false,
      UnescapePathValues:     true,
      MaxMultipartMemory:     defaultMultipartMemory,
      trees:                  make(methodTrees, 0, 9),
      delims:                 render.Delims{Left: "{{", Right: "}}"},
      secureJSONPrefix:       "while(1);",
      trustedProxies:         []string{"0.0.0.0/0"},
      trustedCIDRs:           defaultTrustedCIDRs,
   }
   engine.RouterGroup.engine = engine
   engine.pool.New = func() interface{} {
      return engine.allocateContext()
   }
   return engine
}
```

`gin.New()`会进行 Engine 示例的初始化动作并返回，在初始化时主要会设置以下参数：

- `RouterGroup`：路由组，所有的路由规则都由 `*RouterGroup` 所属的方法进行管理，在 gin 中和 Engine 实例形成一个重要的关联组件。
- `RedirectTrailingSlash`：是否自动重定向，如果启用了，在无法匹配当前路由的情况下，**则自动重定向到带有或不带斜杠的处理程序**。例如：当外部请求了 `/tour/` 路由，但当前并没有注册该路由规则，只有 `/tour` 的路由规则时，将会在内部进行判定，若是 HTTP GET 请求，将会通过 HTTP Code 301 重定向到 `/tour` 的处理程序去，但若是其他类型的 HTTP 请求，那么将会是以 HTTP Code 307 重定向，通过指定的 HTTP 状态码重定向到 `/tour` 路由的处理程序去。
- `RedirectFixedPath`：是否尝试修复当前请求路径，也就是在开启的情况下，gin 会尽可能的帮你找到一个相似的路由规则并在内部重定向过去，主要是对当前的请求路径进行格式清除（删除多余的斜杠）和不区分大小写的路由查找等。
- `HandleMethodNotAllowed`：判断当前路由是否允许调用其他方法，如果当前请求无法路由，则返回 Method Not Allowed（HTTP Code 405）的响应结果。如果无法路由，也不支持重定向其他方法，则交由 `NotFound Hander `进行处理。
- `ForwardedByClientIP`：如果开启，则尽可能的返回真实的客户端 IP，先从 `X-Forwarded-For` 取值，如果没有再从 `X-Real-Ip`。
- `UseRawPath`：如果开启，则会使用 `url.RawPath` 来获取请求参数，不开启则还是按 `url.Path `去获取。
- `UnescapePathValues`：是否对路径值进行转义处理。
- `MaxMultipartMemory`：相对应 `http.Request ParseMultipartForm` 方法，用于控制最大的文件上传大小。
- `trees`：多个压缩字典树（Radix Tree），每个树都对应着一种 HTTP Method。你可以理解为，每当你添加一个新路由规则时，就会往 HTTP Method 对应的那个树里新增一个 node 节点，以此形成关联关系。
- `delims`：用于 HTML 模板的左右定界符。

总的来讲，Engine 实例就像引擎一样，与整个应用的运行、路由、对象、模板等管理和调度都有关联，另外通过上述的解析，可以发现其实 gin 在初始化默认已经做了很多事情，可以说是既定了一些默认运行基础。

#### c. r.GET

```go
// GET is a shortcut for router.Handle("GET", path, handle).
func (group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) IRoutes {
   return group.handle(http.MethodGet, relativePath, handlers)
}

func (group *RouterGroup) handle(httpMethod, relativePath string, handlers HandlersChain) IRoutes {
	absolutePath := group.calculateAbsolutePath(relativePath) // 计算路由的绝对路径
	handlers = group.combineHandlers(handlers) // 合并 Handler
	group.engine.addRoute(httpMethod, absolutePath, handlers) // 追加路由规则
	return group.returnObj()
}
```

- 计算路由的绝对路径，也就是 `group.basePath` 与定义的路由路径组装，那么 group 又是什么东西呢，实际上在 gin 中存在组别路由的概念，这个知识点在后续实战中会使用到。
- 合并现有和新注册的 Handler，并创建一个函数链 HandlersChain。
- 将当前注册的路由规则（含 HTTP Method、Path、Handlers）追加到对应的树中。

这类方法主要是针对路由的各类计算和注册行为，并输出路由注册的调试信息，如运行时的路由信息：

```bash
[GIN-debug] GET    /ping                     --> main.main.func1 (3 handlers) # 注意此处为3 handlers
```

明明只注册了 `/ping` 这一条路由而已，是不是应该是1个 Handler。其实不然，看看上述创建函数链 HandlersChain 的详细步骤，就知道为什么了，如下:

```go
func (group *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
   finalSize := len(group.Handlers) + len(handlers)
   if finalSize >= int(abortIndex) {
      panic("too many handlers")
   }
   mergedHandlers := make(HandlersChain, finalSize)
   copy(mergedHandlers, group.Handlers) // 优先级高于外部传入的 handlers
   copy(mergedHandlers[len(group.Handlers):], handlers)
   return mergedHandlers
}

// Use adds middleware to the group, see example code in GitHub.
func (group *RouterGroup) Use(middleware ...HandlerFunc) IRoutes {
	group.Handlers = append(group.Handlers, middleware...)
	return group.returnObj()
}
```

可以看到在 `combineHandlers` 方法中，最终函数链 HandlersChain 的是由 `group.Handlers` 和外部传入的 `handlers` 组成的，从拷贝的顺序来看，`group.Handlers` 的优先级高于外部传入的 `handlers`。

那么以此再结合 `Use` 方法来看，很显然是在 `gin.Default` 方法中注册的中间件影响了这个结果，因为中间件也属于 `group.Handlers` 的一部分，也就是在调用 `gin.Use`，就已经注册进去了。

```go
engine.Use(Logger(), Recovery())
```

所注册的路由加上内部默认设置的两个中间件，最终使得显示的结果为 **3 handlers**。

#### d. r.Run

```go
// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) Run(addr ...string) (err error) {
   defer func() { debugPrintError(err) }()

   if engine.isUnsafeTrustedProxies() {
      debugPrint("[WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.\n" +
         "Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.")
   }

   address := resolveAddress(addr)
   debugPrint("Listening and serving HTTP on %s\n", address)
   err = http.ListenAndServe(address, engine) // 将 Engine 实例作为 Handle 注册进去
   return
}
```

该方法会通过解析地址，再调用 `http.ListenAndServe` 将 Engine 实例作为 `Handler` 注册进去，然后启动服务，开始对外提供 HTTP 服务。

这里值得关注的是，为什么 Engine 实例能够传进去呢，明明形参要求的是 `Handler` 接口类型。**这是因为在 Go 语言中如果某个结构体实现了 interface 定义声明的那些方法，那么就可以认为这个结构体实现了 interface（Go 多态）。**

那么在 gin 中，Engine 结构体实现了 `ServeHTTP` 方法的，符合 `http.Handler` 接口标准，代码如下：

```go
// ServeHTTP conforms to the http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
   c := engine.pool.Get().(*Context)
   c.writermem.reset(w)
   c.Request = req
   c.reset()

   engine.handleHTTPRequest(c) // 处理外部的 HTTP 请求

   engine.pool.Put(c)
}
```

- 从 `sync.Pool` 对象池中获取一个上下文对象。
- 重新初始化取出来的上下文对象。
- 处理外部的 HTTP 请求。
- 处理完毕，将取出的上下文对象返回给对象池。

在这里上下文的池化主要是为了防止频繁反复生成上下文对象，相对的提高性能，并且针对 gin 本身的处理逻辑进行二次封装处理。

## 二、进行项目设计

### 1. 目录结构

项目标准目录结构如下：

```go
blog-service(ch02)
├── configs
├── docs
├── global
├── internal
│   ├── dao
│   ├── middleware
│   ├── model
│   ├── routers
│   └── service
├── pkg
├── storage
├── scripts
└── third_party
```

- configs：配置文件。
- docs：文档集合。
- global：全局变量。
- internal：内部模块。
  - dao：数据访问层（Database Access Object），所有与数据相关的操作都会在 dao 层进行，例如 MySQL、Elasticsearch 等。
  - middleware：HTTP 中间件。
  - model：模型层，用于存放 model 对象。
  - routers：路由相关逻辑处理。
  - service：项目核心业务逻辑。
- pkg：项目相关的模块包。
- storage：项目生成的临时文件。
- scripts：各类构建，安装，分析等操作的脚本。
- third_party：第三方的资源工具，例如 Swagger UI。

### 2. 数据库

![image](https://raw.githubusercontent.com/tonshz/test/master/img/202205102158488.jpeg)

三个表中的公共字段如下：

```mysql
  `created_on` int(10) unsigned DEFAULT '0' COMMENT '创建时间',
  `created_by` varchar(100) DEFAULT '' COMMENT '创建人',
  `modified_on` int(10) unsigned DEFAULT '0' COMMENT '修改时间',
  `modified_by` varchar(100) DEFAULT '' COMMENT '修改人',
  `deleted_on` int(10) unsigned DEFAULT '0' COMMENT '删除时间',
  `is_del` tinyint(3) unsigned DEFAULT '0' COMMENT '是否删除 0 为未删除、1 为已删除',
```

#### a. 创建标签表

```mysql
CREATE TABLE `blog_tag` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) DEFAULT '' COMMENT '标签名称',
  `created_on` int(10) unsigned DEFAULT '0' COMMENT '创建时间',
  `created_by` varchar(100) DEFAULT '' COMMENT '创建人',
  `modified_on` int(10) unsigned DEFAULT '0' COMMENT '修改时间',
  `modified_by` varchar(100) DEFAULT '' COMMENT '修改人',
  `deleted_on` int(10) unsigned DEFAULT '0' COMMENT '删除时间',
  `is_del` tinyint(3) unsigned DEFAULT '0' COMMENT '是否删除 0 为未删除、1 为已删除',
  `state` tinyint(3) unsigned DEFAULT '1' COMMENT '状态 0 为禁用、1 为启用',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='标签管理';
```

创建标签表，表字段主要为标签的名称、状态以及公共字段。

#### b. 创建文章表

```mysql
CREATE TABLE `blog_article` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(100) DEFAULT '' COMMENT '文章标题',
  `desc` varchar(255) DEFAULT '' COMMENT '文章简述',
  `cover_image_url` varchar(255) DEFAULT '' COMMENT '封面图片地址',
  `content` longtext COMMENT '文章内容',
  `created_on` int(10) unsigned DEFAULT '0' COMMENT '创建时间',
  `created_by` varchar(100) DEFAULT '' COMMENT '创建人',
  `modified_on` int(10) unsigned DEFAULT '0' COMMENT '修改时间',
  `modified_by` varchar(100) DEFAULT '' COMMENT '修改人',
  `deleted_on` int(10) unsigned DEFAULT '0' COMMENT '删除时间',
  `is_del` tinyint(3) unsigned DEFAULT '0' COMMENT '是否删除 0 为未删除、1 为已删除',
  `state` tinyint(3) unsigned DEFAULT '1' COMMENT '状态 0 为禁用、1 为启用',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章管理';
```

创建文章表，表字段主要为文章的标题、封面图、内容概述以及公共字段。

#### c. 创建文章标签关联表

```mysql
CREATE TABLE `blog_article_tag` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `article_id` int(11) NOT NULL COMMENT '文章 ID',
  `tag_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '标签 ID',
  `created_on` int(10) unsigned DEFAULT '0' COMMENT '创建时间',
  `created_by` varchar(100) DEFAULT '' COMMENT '创建人',
  `modified_on` int(10) unsigned DEFAULT '0' COMMENT '修改时间',
  `modified_by` varchar(100) DEFAULT '' COMMENT '修改人',
  `deleted_on` int(10) unsigned DEFAULT '0' COMMENT '删除时间',
  `is_del` tinyint(3) unsigned DEFAULT '0' COMMENT '是否删除 0 为未删除、1 为已删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章标签关联';
```

创建文章标签关联表，这个表主要用于记录文章和标签之间的 1:N 的关联关系。

### 3. 创建 model

#### a. 创建公共 model

在 `internal/model` 目录下创建 model.go 文件，写入如下代码：

```go
type Model struct {
    ID         uint32 `gorm:"primary_key" json:"id"`
    CreatedBy  string `json:"created_by"`
    ModifiedBy string `json:"modified_by"`
    CreatedOn  uint32 `json:"created_on"`
    ModifiedOn uint32 `json:"modified_on"`
    DeletedOn  uint32 `json:"deleted_on"`
    IsDel      uint8  `json:"is_del"`
}
```

#### b. 创建标签 model

在 `internal/model` 目录下创建 tag.go 文件，写入如下代码：

```go
type Tag struct {
    *Model // 继承
    Name  string `json:"name"`
    State uint8  `json:"state"`
}
func (t Tag) TableName() string {
    return "blog_tag"
}
```

#### c. 创建文章 model

在 `internal/model` 目录下创建 article.go 文件，写入如下代码：

```go
package model

type Article struct {
	*Model // 继承公共 Model 结构体中的属性
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	State         uint8  `json:"state"`
}

func (a Article) TableName() string {
	return "blog_article"
}
```

#### d. 创建文章标签 model

在 `internal/model` 目录下创建 article_tag.go 文件，写入如下代码：

```go
type ArticleTag struct {
    *Model // 继承
    TagID     uint32 `json:"tag_id"`
    ArticleID uint32 `json:"article_id"`
}
func (a ArticleTag) TableName() string {
    return "blog_article_tag"
}
```

### 4. 路由

在完成数据库的设计后，需要对业务模块的管理接口进行设计，这一块的核心是增删改查的 RESTful API 设计和编写，在 RESTful API 中 HTTP 方法对应的行为动作分别如下：

- GET：读取/检索动作。
- POST：新增/新建动作。
- PUT：更新动作，用于更新一个完整的资源，要求为幂等。
- PATCH：更新动作，用于更新某一个资源的一个组成部分，也就是只需要更新该资源的某一项，就应该使用 PATCH 而不是 PUT，可以不幂等。
- DELETE：删除动作。

在后续实现中可以根据 RESTful API 的基本规范针对业务模块设计路由规则，从业务角度来划分多个管理接口。

#### a. 标签管理

| 功能         | HTTP 方法 | 路径      |
| :----------- | :-------- | :-------- |
| 新增标签     | POST      | /tags     |
| 删除指定标签 | DELETE    | /tags/:id |
| 更新指定标签 | PUT       | /tags/:id |
| 获取标签列表 | GET       | /tags     |

#### b. 文章管理

| 功能         | HTTP 方法 | 路径          |
| :----------- | :-------- | :------------ |
| 新增文章     | POST      | /articles     |
| 删除指定文章 | DELETE    | /articles/:id |
| 更新指定文章 | PUT       | /articles/:id |
| 获取指定文章 | GET       | /articles/:id |
| 获取文章列表 | GET       | /articles     |

#### c. 路由管理

在确定了业务接口设计后，需要对业务接口进行一个基础编码，确定其方法原型，在`internal/routers`目录下新建`router.go`文件。

```go
package routers

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine {
	// 或者 r := gin.Default() 也可
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// 使用路由组设置访问路由的统一前缀 e.g. /api/v1
	// 此处定义了一个路由组 /api/v1
	apiv1 := r.Group("/api/v1")
	{
		apiv1.POST("/tags")
		apiv1.DELETE("/tags/:id")
		apiv1.PUT("/tags/:id")
		apiv1.PATCH("/tags/:id/state")
		apiv1.GET("/tags")
		apiv1.POST("/articles")
		apiv1.DELETE("/articles/:id")
		apiv1.PUT("/articles/:id")
		apiv1.PATCH("/articles/:id/state")
		apiv1.GET("/articles/:id")
		apiv1.GET("/articles")
	}
	return r
}
```

### 5. 处理程序

接下来编写对应路由的处理方法，在项目目录下新建 `internal/routers/api/v1` 文件夹，并新建`tag.go`（标签）和` article.go`（文章）文件.。

#### a. tag.go

```go
package v1

import "github.com/gin-gonic/gin"

type Tag struct {}
func NewTag() Tag {
   return Tag{}
}

func (t Tag) Get(c *gin.Context) {}
func (t Tag) List(c *gin.Context) {}
func (t Tag) Create(c *gin.Context) {}
func (t Tag) Update(c *gin.Context) {}
func (t Tag) Delete(c *gin.Context) {}
```

#### b. article.go

```go
package v1

import "github.com/gin-gonic/gin"

type Article struct{}
func NewArticle() Article {
   return Article{}
}

func (a Article) Get(c *gin.Context) {}
func (a Article) List(c *gin.Context) {}
func (a Article) Create(c *gin.Context) {}
func (a Article) Update(c *gin.Context) {}
func (a Article) Delete(c *gin.Context) {}
```

#### c. 路由管理

在编写好路由的 Handler 方法后，只需要将其注册到对应的路由规则上就好了，打开项目目录下 `internal/routers` 的 `router.go `文件，修改如下：

```go
package routers

import (
   v1 "demo/ch02/internal/routers/api/v1"
   "github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
   // 或者 r := gin.Default() 也可
   r := gin.New()
   r.Use(gin.Logger())
   r.Use(gin.Recovery())

   article := v1.NewArticle()
   tag := v1.NewTag()
   // 使用路由组设置访问路由的统一前缀 e.g. /api/v1
   // 此处定义了一个路由组 /api/v1
   apiv1 := r.Group("/api/v1")
   {
      // 将实现的 Handler 方法注册到对应的路由规则上
      apiv1.POST("/tags", tag.Create)
      apiv1.DELETE("/tags/:id", tag.Delete)
      apiv1.PUT("/tags/:id", tag.Update)
      apiv1.PATCH("/tags/:id/state", tag.Update)
      apiv1.GET("/tags", tag.List)

      apiv1.POST("/articles", article.Create)
      apiv1.DELETE("/articles/:id", article.Delete)
      apiv1.PUT("/articles/:id", article.Update)
      apiv1.PATCH("/articles/:id/state", article.Update)
      apiv1.GET("/articles/:id", article.Get)
      apiv1.GET("/articles", article.List)
   }
   return r
}
```

### 6. 启动接入

在完成了模型、路由的代码编写后，修改 main.go 文件，把它改造为这个项目的启动文件。

```go
package main

import (
   "demo/ch02/internal/routers"
   "log"
   "net/http"
   "time"
)

func main() {
   // 不再使用默认路由而使用项目下自定义的路由
   // router := gin.Default()
   router := routers.NewRouter()
   // 自定义 http.Server
   s := &http.Server{
      Addr:           ":8080", // 设置监听端口
      Handler:        router, // 设置处理程序
      ReadTimeout:    10 * time.Second, // 允许读取最大时间
      WriteTimeout:   10 * time.Second, // 允许写入最大时间
      MaxHeaderBytes: 1 << 20, // 请求头最大字节数
   }
   // 调用 ListenAndServe() 监听
   if err := s.ListenAndServe(); err != nil{
      log.Fatalf("监听失败：%v", err)
   }
}
```

通过自定义 `http.Server`，设置了监听的 TCP Endpoint、处理的程序、允许读取/写入的最大时间、请求头的最大字节数等基础参数，最后调用 `ListenAndServe` 方法开始监听。

### 7. 验证

```bash
$ go run main.go
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)
# 路由正常注册，测试调用接口返回正常
[GIN-debug] POST   /api/v1/tags              --> demo/ch02/internal/routers/api/v1.Tag.Create-fm (3 handlers)
[GIN-debug] DELETE /api/v1/tags/:id          --> demo/ch02/internal/routers/api/v1.Tag.Delete-fm (3 handlers)
[GIN-debug] PUT    /api/v1/tags/:id          --> demo/ch02/internal/routers/api/v1.Tag.Update-fm (3 handlers)
[GIN-debug] PATCH  /api/v1/tags/:id/state    --> demo/ch02/internal/routers/api/v1.Tag.Update-fm (3 handlers)
[GIN-debug] GET    /api/v1/tags              --> demo/ch02/internal/routers/api/v1.Tag.List-fm (3 handlers)
[GIN-debug] POST   /api/v1/articles          --> demo/ch02/internal/routers/api/v1.Article.Create-fm (3 handlers)
[GIN-debug] DELETE /api/v1/articles/:id      --> demo/ch02/internal/routers/api/v1.Article.Delete-fm (3 handlers)
[GIN-debug] PUT    /api/v1/articles/:id      --> demo/ch02/internal/routers/api/v1.Article.Update-fm (3 handlers)
[GIN-debug] PATCH  /api/v1/articles/:id/state --> demo/ch02/internal/routers/api/v1.Article.Update-fm (3 handlers)
[GIN-debug] GET    /api/v1/articles/:id      -->  demo/ch02/internal/routers/api/v1.Article.Get-fm (3 handlers)
[GIN-debug] GET    /api/v1/articles          --> demo/ch02/internal/routers/api/v1.Article.List-fm (3 handlers)
```

## 三、编写公共组件

### 1. 错误码标准化

#### a. 公共错误码

在项目目录下的 `pkg/errcode` 目录新建` common_code.go `文件，用于预定义项目中的一些公共错误码，便于引导和规范使用。

```go
package errcode

var (
   Success                   = NewError(0, "成功")
   ServerError               = NewError(10000000, "服务内部错误")
   InvalidParams             = NewError(10000001, "入参错误")
   NotFound                  = NewError(10000002, "找不到")
   UnauthorizedAuthNotExist  = NewError(10000003, "鉴权失败，找不到对应的 AppKey 和 AppSecret")
   UnauthorizedTokenError    = NewError(10000004, "鉴权失败，Token 错误")
   UnauthorizedTokenTimeout  = NewError(10000005, "鉴权失败，Token 超时")
   UnauthorizedTokenGenerate = NewError(10000006, "鉴权失败，Token 生成失败")
   TooManyRequests           = NewError(10000007, "请求过多")
)
```

#### b. 错误处理

在项目目录下的 `pkg/errcode` 目录新建` errcode.go `文件，编写常用的一些错误处理公共方法，标准化错误输出。

```go
package errcode

import (
   "fmt"
   "net/http"
)

type Error struct {
   // 首字母小写表示私有变量
   code int `json:"code"`
   msg string `json:"msg"`
   details []string `json:"details"`
}

// 全局变量
var codes = map[int]string{}
func NewError(code int, msg string) *Error {
   if _, ok := codes[code]; ok {
      // 使用 go 1.18版本 由于IDE版本较低会报错
      panic(fmt.Sprintf("错误码 %d 已经存在，请更换一个", code))
   }
   codes[code] = msg
   return &Error{code: code, msg: msg}
}

// 实现 error 接口中的方法
func (e *Error) Error() string {
   return fmt.Sprintf("错误码：%d, 错误信息:：%s", e.Code(), e.Msg())
}

func (e *Error) Code() int {
   return e.code
}

func (e *Error) Msg() string {
   return e.msg
}

func (e *Error) Msgf(args []interface{}) string {
   return fmt.Sprintf(e.msg, args...)
}

func (e *Error) Details() []string {
   return e.details
}

func (e *Error) WithDetails(details ...string) *Error {
   newError := *e
   newError.details = []string{}
   for _, d := range details {
      newError.details = append(newError.details, d)
   }
   return &newError
}

func (e *Error) StatusCode() int {
   switch e.Code() {
   case Success.Code():
      return http.StatusOK
   case ServerError.Code():
      return http.StatusInternalServerError
   case InvalidParams.Code():
      return http.StatusBadRequest
   case UnauthorizedAuthNotExist.Code():
      // 满足fallthrough上的某个 case 条件后会强制执行其后的所有 case 代码，不包含 default 部分
      fallthrough
   case UnauthorizedTokenError.Code():
      fallthrough
   case UnauthorizedTokenGenerate.Code():
      fallthrough
   case UnauthorizedTokenTimeout.Code():
      return http.StatusUnauthorized
   case TooManyRequests.Code():
      return http.StatusTooManyRequestsgo
   }
   return http.StatusInternalServerError
}
```

在错误码方法的编写中，声明了 `Error` 结构体用于表示错误的响应结果，并利用 `codes` 作为全局错误码的存储载体，便于查看当前注册情况，并在调用 `NewError` 创建新的 `Error` 实例的同时进行排重的校验。

另外相对特殊的是 `StatusCode` 方法，它主要用于针对一些特定错误码进行状态码的转换，因为不同的内部错误码在 HTTP 状态码中都代表着不同的意义，需要将其区分开来，便于客户端以及监控/报警等系统的识别和监听。

### 2. 配置管理: viper

在应用程序的运行生命周期中，最直接的关系之一就是应用的配置读取和更新。它的一举一动都有可能影响应用程序的改变，其分别包含如下行为：

![image](https://raw.githubusercontent.com/tonshz/test/master/img/202205112108661.jpeg)

- 在启动时：可以进行一些基础应用属性、连接第三方实例（MySQL、NoSQL）等等的初始化行为。
- 在运行中：可以监听文件或其他存储载体的变更来实现热更新配置的效果，例如：在发现有变更的话，就对原有配置值进行修改，以此达到相关联的一个效果。如果更深入业务使用的话，还可以通过配置的热更新，达到功能灰度的效果，这也是一个比较常见的场景。

另外，配置组件是会根据实际情况去选型的，一般大多为文件配置或配置中心的模式，在本次博客后端中的配置管理使用最常见的文件配置作为选型。

#### a. 安装

为了完成文件配置的读取，需要借助第三方开源库 `viper`，在项目根目录下执行以下安装命令（-u 表示如果本地存在旧包则进行更新）：

```bash
$ go get -u github.com/spf13/viper@v1.4.0
```

Viper 是适用于 Go 应用程序的完整配置解决方案，是目前 Go 语言中比较流行的文件配置解决方案，它支持处理各种不同类型的配置需求和配置格式。

#### b.配置文件

在项目目录下的 `configs` 目录新建 `config.yaml` 文件，写入以下配置：

```yaml
Server: # 服务配置
  RunMode: debug # 设置 gin 的运行模式
  HttpPort: 8000  # 服务端口号
  ReadTimeout: 60
  WriteTimeout: 60
App: # 应用配置
  DefaultPageSize: 10
  MaxPageSize: 100
  LogSavePath: storage/logs # 默认应用日志存储位置
  LogFileName: app # 默认应用日志名称
  LogFileExt: .log # 默认应用日志文件后缀名
Database: # 数据库配置
  DBType: mysql
  Username: root  # 数据库账号
  Password: root  # 数据库密码
  Host: 127.0.0.1:3306
  DBName: ch02 # 数据库名称
  TablePrefix: blog_ # 表名称前缀
  Charset: utf8
  ParseTime: True
  MaxIdleConns: 10
  MaxOpenConns: 30
```

配置文件中，分别针对如下内容进行了默认配置：

- Server：服务配置，设置 gin 的运行模式、默认的 HTTP 监听端口、允许读取和写入的最大持续时间。
- App：应用配置，设置默认每页数量、所允许的最大每页数量以及默认的应用日志存储路径。
- Database：数据库配置，主要是连接实例所必需的基础参数。

#### c. 编写组件

在完成了配置文件的确定和编写后，需要针对读取配置的行为进行封装，便于应用程序的使用，在项目目录下的 `pkg/setting` 目录下新建 `setting.go `文件。

```go
package setting

import "github.com/spf13/viper"

type Setting struct {
   vp *viper.Viper
}

// 用于初始化项目的基本配置
func NewSetting() (*Setting, error) {
   vp := viper.New()
   vp.SetConfigName("config") // 设置配置文件名称
   // 设置配置文件相对路径， viper 允许多个配置路径，可以不断调用 AddConfigPath()
   vp.AddConfigPath("configs/") 
   vp.SetConfigType("yaml") // 设置配置文件类型
   
   err := vp.ReadInConfig()
   if err != nil {
      return nil, err
   }
   return &Setting{vp}, nil
}
```

在其中编写了 `NewSetting` 方法，用于初始化本项目的配置的基础属性，设定配置文件的名称为 `config`，配置类型为 `yaml`，并且设置其配置路径为相对路径 `configs/`，以此确保在项目目录下执行运行时能够成功启动。

另外 viper 是允许设置多个配置路径的，这样子可以尽可能的尝试解决路径查找的问题，也就是可以不断地调用 `AddConfigPath` 方法。

接下来新建` section.go` 文件，用于声明配置属性的结构体并编写读取区段配置的配置方法。

```go
package setting

import "time"

// 声明配置属性的结构体
// 服务配置结构体
type ServerSettingS struct {
   RunMode      string
   HttpPort     string
   ReadTimeout  time.Duration
   WriteTimeout time.Duration
}

// 应用配置结构体
type AppSettingS struct {
   DefaultPageSize int
   MaxPageSize     int
   LogSavePath     string
   LogFileName     string
   LogFileExt      string
}

// 数据库配置结构体
type DatabaseSettingS struct {
   DBType       string
   UserName     string
   Password     string
   Host         string
   DBName       string
   TablePrefix  string
   Charset      string
   ParseTime    bool
   MaxIdleConns int
   MaxOpenConns int
}

// 读取相应配置的配置方法
func (s *Setting) ReadSection(k string, v interface{}) error {
    // 将配置文件 按照父节点(k)读取到相应的struct(v)中
   err := s.vp.UnmarshalKey(k, v)
   if err != nil {
      return err
   }
   return nil
}
```

#### d. 包全局变量

在读取了文件的配置信息后，还是不够的，因为需要将配置信息和应用程序关联起来，才能够去使用它，在项目目录下的 `global` 目录下新建 `setting.go` 文件。

```go
package global

import "demo/ch02/pkg/setting"

var (
   ServerSetting   *setting.ServerSettingS
   AppSetting      *setting.AppSettingSgo
   DatabaseSetting *setting.DatabaseSettingS
)
```

针对最初预估的三个区段配置，进行了全局变量的声明，便于在接下来的步骤将其关联起来，并且提供给应用程序内部调用。另外全局变量的初始化，是会随着应用程序的不断演进不断改变的，因此并不是一成不变，也就是这里展示的并不一定是最终的结果。

#### e. 初始化配置读取

在完成了所有的预备行为后，回到项目根目录下的 main.go 文件，修改代码。

```go
package main

import (
   "demo/ch02/global"
   "demo/ch02/internal/routers"
   "demo/ch02/pkg/setting"
   "log"
   "net/http"
   "time"
)

// Go 中的执行顺序: 全局变量初始化 =>init() => main()
// 在 main() 之前自动执行，进行初始化操作
func init() {
   err := setupSetting()
   if err != nil {
      log.Fatalf("init.setupSetting err: %v", err)
   }
}

func main() {
   ...
}

func setupSetting() error {
   settings, err := setting.NewSetting()
   if err != nil {
      return err
   }
   err = settings.ReadSection("Server", &global.ServerSetting)
   if err != nil {
      return err
   }
   err = settings.ReadSection("App", &global.AppSetting)
   if err != nil {
      return err
   }
   err = settings.ReadSection("Database", &global.DatabaseSetting)
   if err != nil {
      return err
   }
   // global.ServerSetting.ReadTimeout *=1000，将秒转换成毫秒
   global.ServerSetting.ReadTimeout *= time.Second
   global.ServerSetting.WriteTimeout *= time.Second
   return nil
}
```

新增了一个 `init` 方法，在 Go 语言中，`init` 方法常用于应用程序内的一些初始化操作，它在 `main` 方法之前自动执行，它的执行顺序是：`全局变量初始化 => init 方法 =>main 方法`，但并不是建议滥用，因为如果 `init` 过多，可能会迷失在各个库的 `init` 方法中，会非常麻烦。

而在此处的应用程序中，该 `init` 方法主要作用是进行应用程序的初始化流程控制，整个应用代码里也只会有一个 `init` 方法，因此在这里调用初始化配置的方法，达到配置文件内容映射到应用配置结构体的作用。

#### f. 修改服务端配置

接下来只需要在启动文件 main.go 中把已经映射好的配置和 gin 的运行模式进行设置，这样的话，在程序重新启动时后就可以生效。

```go
func main() {
   // 使用映射好的配置设置 gin 的运行模式: debug
   gin.SetMode(global.ServerSetting.RunMode)
   // 不再使用默认路由而使用项目下自定义的路由
   // router := gin.Default()
   router := routers.NewRouter()
   // 自定义 http.Server
   s := &http.Server{
      Addr:           ":8080", // 设置监听端口
      Handler:        router, // 设置处理程序
      ReadTimeout:    10 * time.Second, // 允许读取最大时间
      WriteTimeout:   10 * time.Second, // 允许写入最大时间
      MaxHeaderBytes: 1 << 20, // 请求头最大字节数
   }
   // 调用 ListenAndServe() 监听
   if err := s.ListenAndServe(); err != nil{
      log.Fatalf("监听失败：%v", err)
   }
}
```

#### g. 验证

在`main()`中添加输出语句来查看配置是否真正的映射到配置结构体上了。

```go
log.Printf("%+v\n%+v\n%+v\n", global.ServerSetting, global.AppSetting, global.DatabaseSetting)
```

```bash
2022/05/11 21:44:24 
&{RunMode:debug HttpPort:8000 ReadTimeout:1m0s WriteTimeout:1m0s}
&{DefaultPageSize:10 MaxPageSize:100 LogSavePath:storage/logs LogFileName:app LogFileExt:.log}
&{DBType:mysql UserName:root Password:root Host:127.0.0.1:3306 DBName:ch02 TablePrefix:blog_ Charset:utf8 ParseTime:true MaxIdleConns:10 MaxOpenConns:30}
```

### 3. 数据库连接: gorm

#### a. 安装

本项目中数据库相关的数据操作将使用第三方的开源库 `gorm`，它是目前 Go 语言中最流行的 ORM 库（从 Github Star 来看），同时它也是一个功能齐全且对开发人员友好的 ORM 库，目前在 Github 上相当的活跃，具有一定的保障，安装命令如下：

```bash
$ go get -u github.com/jinzhu/gorm@v1.9.12
```

#### b. 编写组件

在项目目录 `internal/model` 下的 model.go 文件，新增` NewDBEngine()`方法。

```go
package model

import (
	"demo/ch02/global"
	"demo/ch02/pkg/setting"
	"fmt"
	"github.com/jinzhu/gorm"
    // 引入 MySQL 驱动库进行初始化，见 gorm.Opem() 源码注释
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Model struct {
	ID         uint32 `gorm:"primary_key" json:"id"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	CreatedOn  uint32 `json:"created_on"`
	ModifiedOn uint32 `json:"modified_on"`
	DeletedOn  uint32 `json:"deleted_on"`
	IsDel      uint8  `json:"is_del"`
}

// 新增 NewDBEngine()
func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	// gorm.Open() 初始化一个 MySQL 连接，首先需要导入驱动
	db, err := gorm.Open(databaseSetting.DBType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime,
	))
	if err != nil {
		return nil, err
	}
	if global.ServerSetting.RunMode == "debug" {
		// 显示日志输出
		db.LogMode(true)
	}
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(databaseSetting.MaxIdleConns)
	db.DB().SetMaxOpenConns(databaseSetting.MaxOpenConns)
	return db, nil
}
```

其中的`gorm.Open()`源码处注释如下：

```go
// Open initialize a new db connection, need to import driver first, e.g:
//
//     import _ "github.com/go-sql-driver/mysql"
//     func main() {
//       db, err := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
//     }
// GORM has wrapped some drivers, for easier to remember driver's import path, so you could import the mysql driver with
//    import _ "github.com/jinzhu/gorm/dialects/mysql"
//    // import _ "github.com/jinzhu/gorm/dialects/postgres"
//    // import _ "github.com/jinzhu/gorm/dialects/sqlite"
//    // import _ "github.com/jinzhu/gorm/dialects/mssql"
```

通过上述代码，编写了一个针对创建 DB 实例的 `NewDBEngine` 方法，同时增加了 gorm 开源库的引入和 MySQL 驱动库 `github.com/jinzhu/gorm/dialects/mysql` 的初始化（不同类型的 DBType 需要引入不同的驱动库，否则会存在问题）。

#### c. 包全局变量

在项目目录下的 `global` 目录，新增` db.go `文件。

```go
package global

import "github.com/jinzhu/gorm"

var (
   // DBEngine 变量包含了当前数据库连接的信息
   DBEngine *gorm.DB
)
```

#### d. 初始化

在项目目录下的 `main.go `文件，新增` setupDBEngine() `方法初始化。

```go
package main

import (
   "demo/ch02/global"
   "demo/ch02/internal/model"
   "demo/ch02/internal/routers"
   "demo/ch02/pkg/setting"
   "github.com/gin-gonic/gin"
   "log"
   "net/http"
   "time"
)

// Go 中的执行顺序: 全局变量初始化 =>init() => main()
// 在 main() 之前自动执行，进行初始化操作
func init() {
   // 获取初始化配置
   err := setupSetting()
   if err != nil {
      log.Fatalf("init.setupSetting err: %v", err)
   }
   // 数据库初始化
   err = setupDBEngine()
   if err != nil {
      log.Fatalf("init.setupDBEngine err: %v", err)
   }
}

func main() {
   ...
}

func setupSetting() error {
   ...
}

func setupDBEngine() error {
   var err error
   // 初始化数据库连接信息，注意此处不是 := 而是 =，使用前者会导致在其他包中调用该变量时值为 nil
   global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
   if err != nil {
      return err
   }
   return nil
}
```

如果把` global.DBEngine`的初始化语句写成：`global.DBEngine, err := model.NewDBEngine(global.DatabaseSetting)`，会无法正常使用。因为 `:=` 会重新声明并创建了左侧的新局部变量，因此在其它包中调用 `global.DBEngine` 变量时，它仍然是 `nil`，在项目中无法正常使用，因为没有赋值到真正需要赋值的包全局变量 `global.DBEngine` 上，而只是一个局部变量。

### 4. 日志导入: lumberjack

上述应用代码中都是直接使用 Go 标准库 log 来进行的日志输出，这其实是有些问题的，因为在一个项目中，日志需要标准化的记录一些的公共信息，例如：代码调用堆栈、请求链路 ID、公共的业务属性字段等等，而直接输出标准库的日志的话，并不具备这些数据，也不够灵活。

日志的信息的齐全与否在排查和调试问题中是非常重要的一环，因此在应用程序中需要有一个标准的日志组件会进行统一处理和输出。

#### a. 安装

```bash
$ go get -u gopkg.in/natefinch/lumberjack.v2
```

拉取日志组件内要使用到的第三方的开源库 `lumberjack`，它的核心功能是将日志写入滚动文件中，该库支持设置所允许单日志文件的最大占用空间、最大生存周期、允许保留的最多旧文件数，如果出现超出设置项的情况，就会对日志文件进行滚动处理。

使用这个库，主要是为了减免一些文件操作类的代码编写，把核心逻辑摆在日志标准化处理上。

#### b. 编写组件

在项目目录下的 `pkg/` 目录新建 `logger` 目录，并创建 logger.go 文件。

```go
package logger

import (
   "context"
   "encoding/json"
   "fmt"
   "io"
   "log"
   "runtime"
   "time"
)

// 日志分级代码
type Level int8
type Fields map[string]interface{}
const (
   LevelDebug Level = iota
   LevelInfo
   LevelWarn
   LevelError
   LevelFatal
   LevelPanic
)

func (l Level) String() string {
   switch l {
   case LevelDebug:
      return "debug"
   case LevelInfo:
      return "info"
   case LevelWarn:
      return "warn"
   case LevelError:
      return "error"
   case LevelFatal:
      return "fatal"
   case LevelPanic:
      return "panic"
   }
   return ""
}
//============================

// 日志标准化，编写具体的方法去进行日志的实例初始化和标准化参数绑定
type Logger struct {
   newLogger *log.Logger
   ctx       context.Context
   fields    Fields
   callers   []string
}
func NewLogger(w io.Writer, prefix string, flag int) *Logger {
   l := log.New(w, prefix, flag)
   return &Logger{newLogger: l}
}
func (l *Logger) clone() *Logger {
   nl := *l
   return &nl
}

// 设置日志公共字段
func (l *Logger) WithFields(f Fields) *Logger {
   ll := l.clone()
   if ll.fields == nil {
      ll.fields = make(Fields)
   }
   for k, v := range f {
      ll.fields[k] = v
   }
   return ll
}

// 设置日志上下文属性
func (l *Logger) WithContext(ctx context.Context) *Logger {
   ll := l.clone()
   ll.ctx = ctx
   return ll
}

// 设置当前某一层调用栈的信息（程序计数器、文件信息、行号）
func (l *Logger) WithCaller(skip int) *Logger {
   ll := l.clone()
   pc, file, line, ok := runtime.Caller(skip)
   if ok {
      f := runtime.FuncForPC(pc)
      ll.callers = []string{fmt.Sprintf("%s: %d %s", file, line, f.Name())}
   }
   return ll
}

// 设置当前的整个调用栈信息
func (l *Logger) WithCallersFrames() *Logger {
   maxCallerDepth := 25
   minCallerDepth := 1
   callers := []string{}
   pcs := make([]uintptr, maxCallerDepth)
   depth := runtime.Callers(minCallerDepth, pcs)
   frames := runtime.CallersFrames(pcs[:depth])
   for frame, more := frames.Next(); more; frame, more = frames.Next() {
      callers = append(callers, fmt.Sprintf("%s: %d %s", frame.File, frame.Line, frame.Function))
      if !more {
         break
      }
   }
   ll := l.clone()
   ll.callers = callers
   return ll
}
//=====================

// 日志格式化输出
func (l *Logger) JSONFormat(level Level, message string) map[string]interface{} {
   data := make(Fields, len(l.fields)+4)
   data["level"] = level.String()
   data["time"] = time.Now().Local().UnixNano()
   data["message"] = message
   data["callers"] = l.callers
   if len(l.fields) > 0 {
      for k, v := range l.fields {
         if _, ok := data[k]; !ok {
            data[k] = v
         }
      }
   }
   return data
}
func (l *Logger) Output(level Level, message string) {
   body, _ := json.Marshal(l.JSONFormat(level, message))
   content := string(body)
   switch level {
   case LevelDebug:
      l.newLogger.Print(content)
   case LevelInfo:
      l.newLogger.Print(content)
   case LevelWarn:
      l.newLogger.Print(content)
   case LevelError:
      l.newLogger.Print(content)
   case LevelFatal:
      l.newLogger.Fatal(content)
   case LevelPanic:
      l.newLogger.Panic(content)
   }
}
//================================

// 日志分级输出: Debug、Info、Warn、Error、Fatal、Panic
// Debug 级别输出
func (l *Logger) Debug(v ...interface{}) {
   l.Output(LevelDebug, fmt.Sprint(v...))
}
func (l *Logger) Debugf(format string, v ...interface{}) {
   l.Output(LevelDebug, fmt.Sprintf(format, v...))
}

// Info 级别输出
func (l *Logger) Info(v ...interface{}) {
   l.Output(LevelInfo, fmt.Sprint(v...))
}
func (l *Logger) Infof(format string, v ...interface{}) {
   l.Output(LevelInfo, fmt.Sprintf(format, v...))
}

// Warn 级别输出
func (l *Logger) Warn(v ...interface{}) {
   l.Output(LevelWarn, fmt.Sprint(v...))
}
func (l *Logger) Warnf(format string, v ...interfgoace{}) {
   l.Output(LevelWarn, fmt.Sprintf(format, v...))
}

// Error 级别输出
func (l *Logger) Error(v ...interface{}) {
   l.Output(LevelError, fmt.Sprint(v...))
}
func (l *Logger) Errorf(format string, v ...interface{}) {
   l.Output(LevelError, fmt.Sprintf(format, v...))gp
}

// Fatal 级别输出
func (l *Logger) Fatal(v ...interface{}) {
   l.Output(LevelFatal, fmt.Sprint(v...))
}
func (l *Logger) Fatalf(format string, v ...interface{}) {
   l.Output(LevelFatal, fmt.Sprintf(format, v...))
}

// Panic 级别输出
func (l *Logger) Panic(v ...interface{}) {
   l.Output(LevelPanic, fmt.Sprint(v...))
}
func (l *Logger) Panicf(format string, v ...interface{}) {
   l.Output(LevelPanic, fmt.Sprintf(format, v...))
}
```

##### 日志分级

预先定义了应用日志的 Level 和 Fields 的具体类型，并且分为了 Debug、Info、Warn、Error、Fatal、Panic 六个日志等级，便于在不同的使用场景中记录不同级别的日志。

##### 日志标准化

编写具体的方法去进行日志的实例初始化和标准化参数绑定。

- WithLevel：设置日志等级。
- WithFields：设置日志公共字段。
- WithContext：设置日志上下文属性。
- WithCaller：设置当前某一层调用栈的信息（程序计数器、文件信息、行号）。
- WithCallersFrames：设置当前的整个调用栈信息。

##### 日志格式化输出

编写日志内容的格式化和日志输出动作的相关方法。

##### 日志分级输出

根据先前定义的日志分级，编写对应的日志输出的外部方法。上述代码中仅展示了 Info、Fatal 级别的日志方法，这里主要是根据 Debug、Info、Warn、Error、Fatal、Panic 六个日志等级编写对应的方法，可自行完善，除了方法名以及 WithLevel 设置的不一样，其他均为一致的代码。

#### c. 包全局变量

定义一个 Logger 对象便于应用程序使用，因此需要修改项目目录下的 `global/setting.go` 文件。在包全局变量中新增了 Logger 对象，用于日志组件的初始化。

```go
package global

import (
   "demo/ch02/pkg/logger"
   "demo/ch02/pkg/setting"
)

var (
   ServerSetting   *setting.ServerSettingS
   AppSetting      *setting.AppSettingS
   DatabaseSetting *setting.DatabaseSettingS
   Logger           *logger.Logger
)
```

#### d. 初始化

修改启动文件，也就是项目目录下的 main.go 文件，新增对刚刚定义的 Logger 对象的初始化。

```go
package main

import (
	"demo/ch02/global"
	"demo/ch02/internal/model"
	"demo/ch02/internal/routers"
	"demo/ch02/pkg/logger"
	"demo/ch02/pkg/setting"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"time"
)

// Go 中的执行顺序: 全局变量初始化 =>init() => main()
// 在 main() 之前自动执行，进行初始化操作
func init() {
	// 获取初始化配置
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
	// 数据库初始化
	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}
	// 日志初始化
	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}
}

func main() {
	...
}

func setupSetting() error {
	...
}

func setupDBEngine() error {
	...
}

func setupLogger() error {
	// 使用了 lumberjack 作为日志库的 io.Writer
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		// 设置生成日志文件存储的相对位置与文件名
		Filename: global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt,
		MaxSize:   600, // 设置最大占用空间
		MaxAge:    10, // 设置日志文件最大生存周期
		LocalTime: true, // 设置日志文件名的时间格式为本地时间
	}, "", log.LstdFlags).WithCaller(2)
	return nil
}

```

#### e. 验证

在`main()`中添加语句，并且在 `storage`目录下创建`logs`目录，日志内容添加格式是追加。

```go
global.Logger.Infof("%s: go-programming-tour-book/%s", "eddycjy", "blog-service")
```

接着可以查看项目目录下的 `storage/logs/app.log`，看看日志文件是否正常创建且写入了预期的日志记录，大致内容如下。

```bash
2022/05/11 22:51:58 {"callers":["C:/Users/zyc/GolandProjects/demo/ch02/main.go: 30 main.init.0"],"level":"info","message":"index: go-logger-test/ch02","time":1652280718874076900}
```

### 5. 响应处理

在应用程序中，与客户端对接的常常是服务端的接口，那客户端是怎么知道这一次的接口调用结果是怎么样的呢？一般来讲，主要是通过对返回的 HTTP 状态码和接口返回的响应结果进行判断，而判断的依据则是事先按规范定义好的响应结果。因此在这一小节将编写统一处理接口返回的响应处理方法，与错误码标准化是相对应的。

#### a.  类型转换

在项目目录下的 `pkg/convert` 目录下新建` convert.go` 文件。

```go
package convert

import "strconv"

type StrTo string
func (s StrTo) String() string {
   return string(s)
}
func (s StrTo) Int() (int, error) {
   v, err := strconv.Atoi(s.String())
   return v, err
}

// 强制转换为 int，不输出错误信息
func (s StrTo) MustInt() int {
   v, _ := s.Int()
   return v
}

// uint32: 32位无符号整数
func (s StrTo) UInt32() (uint32, error) {
   v, err := strconv.Atoi(s.String())
   return uint32(v), err
}

// 强制转换为 uint32, 不输出错误信息
func (s StrTo) MustUInt32() uint32 {
   v, _ := s.UInt32()
   return v
}
```

#### b. 分页处理

在项目目录下的 `pkg/app` 目录下新建` pagination.go `文件。

```go
package app

import (
   "demo/ch02/global"
   "demo/ch02/pkg/convert"
   "github.com/gin-gonic/gin"
)

func GetPage(c *gin.Context) int {
   // 获取页数
   page := convert.StrTo(c.Query("page")).MustInt()
   if page <= 0 {
      return 1
   }
   return page
}

func GetPageSize(c *gin.Context) int {
   // 获取每页大小
   pageSize := convert.StrTo(c.Query("page_size")).MustInt()
   if pageSize <= 0 {
      return global.AppSetting.DefaultPageSize
   }
   if pageSize > global.AppSetting.MaxPageSize {
      return global.AppSetting.MaxPageSize
   }
   return pageSize
}

func GetPageOffset(page, pageSize int) int {
   result := 0
   if page > 0 {
      result = (page - 1) * pageSize
   }
   return result
}
```

#### c. 响应处理

在项目目录下的 `pkg/app` 目录下新建 `app.go` 文件。

```go
package app

import (
   "demo/ch02/pkg/errcode"
   "github.com/gin-gonic/gin"
   "net/http"
)

type Response struct {
   Ctx *gin.Context
}
type Pager struct {
   Page int `json:"page"`
   PageSize int `json:"page_size"`
   TotalRows int `json:"total_rows"`
}
func NewResponse(ctx *gin.Context) *Response {
   return &Response{Ctx: ctx}
}

// 成功响应处理
func (r *Response) ToResponse(data interface{}) {
   if data == nil {
      data = gin.H{}
   }
   r.Ctx.JSON(http.StatusOK, data)
}

// 列表响应处理
func (r *Response) ToResponseList(list interface{}, totalRows int) {
   r.Ctx.JSON(http.StatusOK, gin.H{
      "list": list,
      "pager": Pager{
         Page:      GetPage(r.Ctx),
         PageSize:  GetPageSize(r.Ctx),
         TotalRows: totalRows,
      },
   })
}

// 错误响应处理
func (r *Response) ToErrorResponse(err *errcode.Error) {
   // gin.H 是 gin 框架中对 map[string]interface{} 的缩写
   response := gin.H{"code": err.Code(), "msg": err.Msg()}
   details := err.Details()
   if len(details) > 0 {
      response["details"] = details
   }
   // 在返回体中将结构体序列化为 JSON
   r.Ctx.JSON(err.StatusCode(), response)
}
```

#### d. 验证

随意使用项目中的一个接口方法，调用其对应的方法。

```go
func (a Article) Get(c *gin.Context) {
    app.NewResponse(c).ToErrorResponse(errcode.ServerError)
    return
}
```

```bash
$ curl -v http://127.0.0.1:8080/api/v1/articles/1
...
< HTTP/1.1 500 Internal Server Error
{"code":10000000,"msg":"服务内部错误"}
```

从响应结果上看，可以知道本次接口的调用结果的 HTTP 状态码为 500，响应消息体为约定的错误体，符合要求。

## 四、生成接口文档

### 1. 安装 Swagger

Swagger 相关的工具集会根据 OpenAPI 规范去生成各式各类的与接口相关联的内容，常见的流程是编写注解 =>调用生成库=>生成标准描述文件 =>生成/导入到对应的 Swagger 工具。使用系列命令安装 Go 对应的开源 Swagger 库。

```bash
# 推荐使用 go install, go get 将被遗弃
$ go get -u github.com/swaggo/swag/cmd/swag@v1.6.5 
$ go get -u github.com/swaggo/gin-swagger@v1.2.0 
$ go get -u github.com/swaggo/files
$ go get -u github.com/alecthomas/template
```

验证是否安装成功，如下：

```bash
$ swag -v
swag version v1.6.5
```

如果命令行提示寻找不到 swag 文件，可以检查一下对应的 bin 目录是否已经加入到环境变量 PATH 中。

### 2.  写入注解

在完成了 Swagger 关联库的安装后，需要针对项目里的 API 接口进行注解的编写，以便于后续在进行生成时能够正确的运行，接下来将使用到如下注解：

| 注解     | 描述                                                         |
| :------- | :----------------------------------------------------------- |
| @Summary | 摘要                                                         |
| @Produce | API 可以产生的 MIME 类型的列表，MIME 类型你可以简单的理解为**响应类型**，例如：json、xml、html 等等 |
| @Param   | 参数格式，从左到右分别为：参数名、入参类型、数据类型、是否必填、注释 |
| @Success | 响应成功，从左到右分别为：状态码、参数类型、数据类型、注释   |
| @Failure | 响应失败，从左到右分别为：状态码、参数类型、数据类型、注释   |
| @Router  | 路由，从左到右分别为：路由地址，HTTP 方法                    |

#### a. API

进入项目目录下的 `internal/routers/api/v1` 目录，修改 tag.go 文件。

##### tag.go

```go
package v1

import (
   "github.com/gin-gonic/gin"
)

type Tag struct {}
func NewTag() Tag {
   return Tag{}
}

func (t Tag) Get(c *gin.Context) {}

// @Summary 获取多个标签
// @Produce  json
// @Param name query string false "标签名称" maxlength(100)
// @Param state query int false "状态" Enums(0, 1) default(1)
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} model.Tag "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags [get]
func (t Tag) List(c *gin.Context) {}

// @Summary 新增标签
// @Produce  json
// @Param name body string true "标签名称" minlength(3) maxlength(100)
// @Param state body int false "状态" Enums(0, 1) default(1)
// @Param created_by body string true "创建者" minlength(3) maxlength(100)
// @Success 200 {object} model.Tag "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags [post]
func (t Tag) Create(c *gin.Context) {}

// @Summary 更新标签
// @Produce  json
// @Param id path int true "标签 ID"
// @Param name body string false "标签名称" minlength(3) maxlength(100)
// @Param state body int false "状态" Enums(0, 1) default(1)
// @Param modified_by body string true "修改者" minlength(3) maxlength(100)
// @Success 200 {array} model.Tag "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags/{id} [put]
func (t Tag) Update(c *gin.Context) {}

// @Summary 删除标签
// @Produce  json
// @Param id path int true "标签 ID"
// @Success 200 {string} string "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags/{id} [delete]
func (t Tag) Delete(c *gin.Context) {}
```

##### article.go

```go
package v1

import (
   "github.com/gin-gonic/gin"
)

type Article struct{}
func NewArticle() Article {
   return Article{}
}

// @Summary 获取单个文章
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} model.Article "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/articles/{id} [get]
func (a Article) Get(c *gin.Context) {
}

// @Summary 获取多个文章
// @Produce json
// @Param name query string false "文章名称"
// @Param tag_id query int false "标签ID"
// @Param state query int false "状态"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} model.Article "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/articles [get]
func (a Article) List(c *gin.Context) {}

// @Summary 创建文章
// @Produce json
// @Param tag_id body string true "标签ID"
// @Param title body string true "文章标题"
// @Param desc body string false "文章简述"
// @Param cover_image_url body string true "封面图片地址"
// @Param content body string true "文章内容"
// @Param created_by body int true "创建者"
// @Param state body int false "状态"
// @Success 200 {object} model.Article "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/articles [post]
func (a Article) Create(c *gin.Context) {}

// @Summary 更新文章
// @Produce json
// @Param tag_id body string false "标签ID"
// @Param title body string false "文章标题"
// @Param desc body string false "文章简述"
// @Param cover_image_url body string false "封面图片地址"
// @Param content body string false "文章内容"
// @Param modified_by body string true "修改者"
// @Success 200 {object} model.Article "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/articles/{id} [put]
func (a Article) Update(c *gin.Context) {}


// @Summary 删除文章
// @Produce  json
// @Param id path int true "文章ID"
// @Success 200 {string} string "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/articles/{id} [delete]
func (a Article) Delete(c *gin.Context) {}
```

#### b. main

但存在多个项目可针对项目写注解以便区分项目。

```go
// @title 博客系统
// @version 1.0
// @description Go 语言项目实战学习
// @termsOfService https://github.com/go-programming-tour-book
func main() {
    ...
}
```

### 3. 生成

但编写完所有 Swagger 注解后，在项目根目录下执行如下命令：

```bash
$ swag init
2022/05/12 21:00:09 Generate swagger docs....
2022/05/12 21:00:09 Generate general API Info, search dir:./
2022/05/12 21:00:09 Generating errcode.Error
2022/05/12 21:00:09 Generating model.Article
2022/05/12 21:00:09 Generating model.Tag
2022/05/12 21:00:09 create docs.go at  docs/docs.go
2022/05/12 21:00:09 create swagger.json at  docs/swagger.json
2022/05/12 21:00:09 create swagger.yaml at  docs/swagger.yaml
```

命令执行完毕后，可以在 docs 文件夹下看到 docs.go、swagger.json、swagger.yaml 三个文件。

### 4. 路由

项目中的注解编写完成，通过 `swag init`命令将 Swagger API 所需要的文件生成了，通过在 routers 中进行默认初始化和注册对应的路由即可访问到项目的接口文档，修改`internal/routers` 目录下的`router.go`文件。

```go
package routers

import (
   v1 "demo/ch02/internal/routers/api/v1"
   "github.com/gin-gonic/gin"
   ginSwagger "github.com/swaggo/gin-swagger"
   // 注意此处导包，不需要设置别名
   "github.com/swaggo/gin-swagger/swaggerFiles"
   // 初始化 docs 包
   _ "demo/ch02/docs"
)

func NewRouter() *gin.Engine {
   // 或者 r := gin.Default() 也可
   r := gin.New()
   r.Use(gin.Logger())
   r.Use(gin.Recovery())
   // 手动指定当前应用所启动的 swagger/doc.json 路径
   //url := ginSwagger.URL("http://127.0.0.1:8000/swagger/doc.json")
   //r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
   // 注册一个针对 swagger 的路由，默认指向当前应用所启动的域名下的 swagger/doc.json 路径
   r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
   ...
   return r
}
```

从表面上来看，主要做了两件事，分别是初始化 docs 包和注册一个针对 swagger 的路由，而在初始化 docs 包后，其 swagger.json 将会默认指向当前应用所启动的域名下的 swagger/doc.json 路径，如果有额外需求，可进行手动指定，如下：

```go
  url := ginSwagger.URL("http://127.0.0.1:8000/swagger/doc.json")
  r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
```

#### 5. 查看接口文档

在完成了上述设置后，启动服务端，在浏览器中访问 Swagger 的地址 `localhost:8000/swagger/index.html`。

![image-20220512212318905](https://raw.githubusercontent.com/tonshz/test/master/img/202205122123289.png)

上述图片中的 Swagger 文档展示主要分为三个部分，分别是项目主体信息、接口路由信息、模型信息。

### 6. 生成 Swagger 文档的原因

通过`swag init`命令生成的文件如下：

```bash
docs
├── docs.go
├── swagger.json
└── swagger.yaml
```

#### a. 初始化 docs

第一步，初始化 docs 包，相对应的其实是 `docs.go`文件。

```go
// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2022-05-12 21:00:09.5873562 +0800 CST m=+0.094954501

package docs

import (
   "bytes"
   "encoding/json"
   "strings"

   "github.com/alecthomas/template"
   "github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "https://github.com/go-programming-tour-book",
        "contact": {},
        "license": {},
        "version": "{{.Version}}"
    },
    ...
}`

type swaggerInfo struct {
   Version     string
   Host        string
   BasePath    string
   Schemes     []string
   Title       string
   Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
   Version:     "1.0",
   Host:        "",
   BasePath:    "",
   Schemes:     []string{},
   Title:       "博客系统",
   Description: "Go 语言项目实战学习",
}

type s struct{}

// 实现 Swagger 接口: Swagger is a interface to read swagger document.
func (s *s) ReadDoc() string {
   sInfo := SwaggerInfo
   sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

   t, err := template.New("swagger_info").Funcs(template.FuncMap{
      "marshal": func(v interface{}) string {
         a, _ := json.Marshal(v)
         return string(a)
      },
   }).Parse(doc)
   if err != nil {
      return doc
   }

   var tpl bytes.Buffer
   if err := t.Execute(&tpl, sInfo); err != nil {
      return doc
   }

   return tpl.String()
}

// 默认执行 init()
func init() {
    // 注册 swag.Name 值为 "swagger"
   swag.Register(swag.Name, &s{})
}
```

通过对源码的分析，可以得知实质上在初始化 docs 包时，会默认执行 init 方法，而在 init 方法中，会注册相关方法，主体逻辑是 swag 会在生成时去检索项目下的注解信息，然后将项目信息和接口路由信息按规范生成到包全局变量 doc 中去。紧接着会在 ReadDoc () 方法中做一些 template 的模板映射等工作，完善 doc 的输出。

#### b. 注册路由

在项目中通过调用 `gin.Swagger.WrapHandle(swaggerFiles.Handler)`关联注解数据源。

```go
// WrapHandler wraps `http.Handler` into `gin.HandlerFunc`.
func WrapHandler(h *webdav.Handler, confs ...func(c *Config)) gin.HandlerFunc {
	defaultConfig := &Config{
		URL: "doc.json",
	}

	for _, c := range confs {
		c(defaultConfig)
	}

	return CustomWrapHandler(defaultConfig, h)
}
```

实际上在调用 WrapHandler 后，swag 内部会将其默认调用的 URL 设置为 `doc.json`。

```go
// CustomWrapHandler wraps `http.Handler` into `gin.HandlerFunc`
func CustomWrapHandler(config *Config, h *webdav.Handler) gin.HandlerFunc {
   //create a template with name
   t := template.New("swagger_index.html")
   index, _ := t.Parse(swagger_index_templ)

   var rexp = regexp.MustCompile(`(.*)(index\.html|doc\.json|favicon-16x16\.png|favicon-32x32\.png|/oauth2-redirect\.html|swagger-ui\.css|swagger-ui\.css\.map|swagger-ui\.js|swagger-ui\.js\.map|swagger-ui-bundle\.js|swagger-ui-bundle\.js\.map|swagger-ui-standalone-preset\.js|swagger-ui-standalone-preset\.js\.map)[\?|.]*`)

   return func(c *gin.Context) {

      type swaggerUIBundle struct {
         URL string
      }

      var matches []string
      if matches = rexp.FindStringSubmatch(c.Request.RequestURI); len(matches) != 3 {
         c.Status(404)
         c.Writer.Write([]byte("404 page not found"))
         return
      }
      path := matches[2]
      prefix := matches[1]
      h.Prefix = prefix

      if strings.HasSuffix(path, ".html") {
         c.Header("Content-Type", "text/html; charset=utf-8")
      } else if strings.HasSuffix(path, ".css") {
         c.Header("Content-Type", "text/css; charset=utf-8")
      } else if strings.HasSuffix(path, ".js") {
         c.Header("Content-Type", "application/javascript")
      } else if strings.HasSuffix(path, ".json") {
         c.Header("Content-Type", "application/json")
      }

      switch path {
      case "index.html":
         index.Execute(c.Writer, &swaggerUIBundle{
            URL: config.URL,
         })
      case "doc.json":
         doc, err := swag.ReadDoc()
         if err != nil {
            panic(err)
         }
         c.Writer.Write([]byte(doc))
         return
      default:
         h.ServeHTTP(c.Writer, c.Request)
      }
   }
}
```

在 CustomWrapHandler 方法中，有一段 switch case 的逻辑。在第一个 case 中，处理是的 `index.html`，通过 `http://127.0.0.1:8000/swagger/index.html` 访问到 Swagger 文档的，对应的便是这里的逻辑。

在第二个 case 中，表明`doc.json`相当于一个内部标识，会去读取生成的 Swagger 注解，先前在访问的 Swagger 文档的顶部的文本框中 Explore 默认的就是 doc.json（也可以填写外部地址，只要输出的是对应的 Swagger 注解）。

![image-20220512214254789](https://raw.githubusercontent.com/tonshz/test/master/img/202205122142824.png)

### 7. 问题

在编写成功响应时，直接调用 model 作为数据类型。

```go
// @Success 200 {object} model.Tag "成功"
```

这样写的话，就会有一个问题，如果有 model.Tag 以外的字段，例如分页，那就无法展示了。更接近实践来讲，在编码中常常会遇到某个对象内中的某一个字段是 interface，这个字段的类型是不定的，也就是公共结构体，那注解又应该怎么写呢，如下情况：

```go
type Test struct {
    UserName string
    Content  interface{}
}
```

官方给出的建议很简单，就是定义一个针对 Swagger 的对象，专门用于 Swagger 接口文档展示，我们在 `internal/model` 的` tag.go `和 `article.go `文件中，新增如下代码：

```go
// tag.go
type TagSwagger struct {
    List  []*Tag
    Pager *app.Pager
}
// article.go
type ArticleSwagger struct {
    List  []*Article
    Pager *app.Pager
}
```

同时修改接口方法中对应的注解信息。

```go
// @Success 200 {object} model.TagSwagger "成功"
```

再在项目根目录下执行`swag init`,文件生成成功后再重新启动服务端，就可以看到最新的效果。

![old](https://raw.githubusercontent.com/tonshz/test/master/img/202205122150280.png "old")

![new](https://raw.githubusercontent.com/tonshz/test/master/img/202205122151839.png "new")

## 五、为接口做参数校验

### 1. 安装

在本项目中将使用开源项目 `go-playground/validator`作为基础库，它是一个基于标签来对结构体和字段进行值验证的验证器。由于是使用的是 gin 框架，其内部的模型绑定和验证默认使用的是该库来进行参数绑定和校验，所有使用起来相对方便。

在项目根目录下执行命令进行安装。

```bash
$ go get -u github.com/go-playground/validator/v10
```

### 2. 业务接口校验

接下来将正式开始对接口的入参进行校验规则的编写，**也就是将校验规则写在对应的结构体的字段标签上**，常见的标签含义如下：

| 标签     | 含义                      |
| :------- | :------------------------ |
| required | 必填                      |
| gt       | 大于                      |
| gte      | 大于等于                  |
| lt       | 小于                      |
| lte      | 小于等于                  |
| min      | 最小值                    |
| max      | 最大值                    |
| oneof    | 参数集内的其中之一        |
| len      | 长度要求与 len 给定的一致 |

#### a. 标签接口

在`internal/service` 目录新建 `tag.go` 文件，针对入参校验增加绑定/验证结构体，在路由方法前写入如下代码：

```go
package service

type CountTagRequest struct {
   Name  string `form:"name" binding:"max=100"`
   State uint8 `form:"state,default=1" binding:"oneof=0 1"`
}
type TagListRequest struct {
   Name  string `form:"name" binding:"max=100"`
   State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}
type CreateTagRequest struct {
   Name      string `form:"name" binding:"required,min=3,max=100"`
   CreatedBy string `form:"created_by" binding:"required,min=3,max=100"`
   State     uint8  `form:"state,default=1" binding:"oneof=0 1"`
}
type UpdateTagRequest struct {
   ID         uint32 `form:"id" binding:"required,gte=1"`
   Name       string `form:"name" binding:"min=3,max=100"`
   State      uint8  `form:"state" binding:"required,oneof=0 1"`
   ModifiedBy string `form:"modified_by" binding:"required,min=3,max=100"`
}
type DeleteTagRequest struct {
   ID uint32 `form:"id" binding:"required,gte=1"`
}
```

在上述代码中，主要针对业务接口中定义的的增删改查和统计行为进行了 Request 结构体编写，而在结构体中，应用到了两个 tag 标签，分别是 **form 和 binding**，它们分别代表着表单的映射字段名和入参校验的规则内容，**其主要功能是实现参数绑定和参数检验。**

#### b. 文章接口

接下来在项目的 `internal/service` 目录下新建 `article.go `文件，针对入参校验增加绑定/验证结构体。这块与标签模块的验证规则差不多，主要是必填，长度最小、最大的限制，以及要求参数值必须在某个集合内的其中之一。

```go
package service

type ArticleRequest struct {
   ID    uint32 `form:"id" binding:"required,gte=1"`
   State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type ArticleListRequest struct {
   TagID uint32 `form:"tag_id" binding:"gte=1"`
   State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type CreateArticleRequest struct {
   TagID         uint32 `form:"tag_id" binding:"required,gte=1"`
   Title         string `form:"title" binding:"required,min=2,max=100"`
   Desc          string `form:"desc" binding:"required,min=2,max=255"`
   Content       string `form:"content" binding:"required,min=2,max=4294967295"`
   CoverImageUrl string `form:"cover_image_url" binding:"required,url"`
   CreatedBy     string `form:"created_by" binding:"required,min=2,max=100"`
   State         uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type UpdateArticleRequest struct {
   ID            uint32 `form:"id" binding:"required,gte=1"`
   TagID         uint32 `form:"tag_id" binding:"required,gte=1"`
   Title         string `form:"title" binding:"min=2,max=100"`
   Desc          string `form:"desc" binding:"min=2,max=255"`
   Content       string `form:"content" binding:"min=2,max=4294967295"`
   CoverImageUrl string `form:"cover_image_url" binding:"url"`
   ModifiedBy    string `form:"modified_by" binding:"required,min=2,max=100"`
   State         uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type DeleteArticleRequest struct {
   ID uint32 `form:"id" binding:"required,gte=1"`
}
```

### 3. 国际化处理

#### a. 编写中间件

`go-playground/validator `默认的错误信息是英文，但项目实际使用中错误信息不一定是英文。如果是最简单的国际化需求，可以通过中间件配合语言包的方式去实现这个功能，在`internal/middleware` 目录下新建 `translations.go `文件，用于编写针对 validator 的语言包翻译的相关功能。

```go
package middleware

import (
   "github.com/gin-gonic/gin"
   "github.com/gin-gonic/gin/binding"
   // 多语言包
   "github.com/go-playground/locales/en"
   "github.com/go-playground/locales/zh"
   "github.com/go-playground/locales/zh_Hant_TW"
   // 通用翻译器
   "github.com/go-playground/universal-translator"
   validator "github.com/go-playground/validator/v10"
   // validator 的翻译器
   en_translations "github.com/go-playground/validator/v10/translations/en"
   zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

func Translations() gin.HandlerFunc {
   return func(c *gin.Context) {
      uni := ut.New(en.New(), zh.New(), zh_Hant_TW.New())
      // 通过 GetHeader 方法获取约定的 header 参数 locale,用于辨别当前请求的语言类别是en还是zh
      // 如果有其他语言环境要求，也可以继续引入其他语言类别，go-playground/locales 基本上都支持。
      locale := c.GetHeader("locale")
      // 对应语言的 Translator
      trans, _ := uni.GetTranslator(locale)
      // 验证器
      v, ok := binding.Validator.Engine().(*validator.Validate)
      if ok {
         switch locale {
         // 调用 RegisterDefaultTranslations 方法将验证器和对应语言类型的 Translator 注册进来
         case "zh":
            _ = zh_translations.RegisterDefaultTranslations(v, trans)
            break
         case "en":
            _ = en_translations.RegisterDefaultTranslations(v, trans)
            break
         default:
            _ = zh_translations.RegisterDefaultTranslations(v, trans)
            break
         }
         // 将 Translator 存储到全局上下文中，便于后续翻译时使用go
         c.Set("trans", trans)
      }
      c.Next()
   }
}
```

#### b. 注册中间件

在项目的 `internal/routers` 目录下的` router.go `文件，新增中间件 Translations 的注册。

```go
func NewRouter() *gin.Engine {
   // 或者 r := gin.Default() 也可
   r := gin.New()
   r.Use(gin.Logger())
   r.Use(gin.Recovery())
   // 新增中间件 Translations 的注册
   r.Use(middleware.Translations())
   ...
}
```

至此就完成了在项目中的自定义验证器注册、验证器初始化、错误提示多语言的功能支持了。

### 4. 接口校验

在项目下的 `pkg/app` 目录新建 form.go 文件。

```go
package app

import (
   "github.com/gin-gonic/gin"
   ut "github.com/go-playground/universal-translator"
   val "github.com/go-playground/validator/v10"
   "strings"
)

// 实现 error 接口
type ValidError struct {
   Key     string
   Message string
}
type ValidErrors []*ValidError
func (v *ValidError) Error() string {
   return v.Message
}
func (v ValidErrors) Error() string {
   return strings.Join(v.Errors(), ",")
}
func (v ValidErrors) Errors() []string {
   var errs []string
   for _, err := range v {
      errs = append(errs, err.Error())
   }
   return errs
}

func BindAndValid(c *gin.Context, v interface{}) (bool, ValidErrors) {
   var errs ValidErrors
   // 通过 ShouldBind 进行参数绑定和入参校验
   // ShouldBind 能够基于请求的不同，自动提取JSON、form表单和QueryString类型的数据，并把值绑定到指定的结构体对象
   // 此方法从上下文中获取传入的方法入参并进行绑定
   err := c.ShouldBind(v)
   if err != nil {
      v := c.Value("trans")
      trans, _ := v.(ut.Translator)
      verrs, ok := err.(val.ValidationErrors)
      if !ok {
         return false, errs
      }
      // 通过中间件 Translations 设置的 Translator 来对错误消息体进行翻译
      for key, value := range verrs.Translate(trans) {
         errs = append(errs, &ValidError{
            Key:     key,go
            Message: value,
         })
      }
      return false, errs
   }
   return true, nil
}
```

在上述代码中，主要是针对入参校验的方法进行了二次封装，在 BindAndValid 方法中，**通过 `ShouldBind` 进行参数绑定和入参校验**，当发生错误后，再通过上一步在中间件 Translations 设置的 Translator 来对错误消息体进行具体的翻译行为。

##### c.ShouldBind(v) 方法执行前

![image-20220514235035166](https://raw.githubusercontent.com/tonshz/test/master/img/202205142350218.png)

##### c.ShouldBind(v) 方法执行后

![image-20220514235050171](https://raw.githubusercontent.com/tonshz/test/master/img/202205142350215.png)

**使用传入的参数通过上下文为 v 赋值。`gin.Context`是 gin 最重要的部分。例如，它允许我们在中间件之间传递变量、管理流程、验证请求的 JSON 并呈现 JSON 响应。**

另外声明了 ValidError 相关的结构体和类型，对这块不熟悉的读者可能会疑惑为什么要实现其对应的 Error 方法呢，简单来看看标准库中 errors 的相关代码。

```go
func New(text string) error {
    return &errorString{text}
}
type errorString struct {
    s string
}
func (e *errorString) Error() string {
    return e.s
}
```

标准库 errors 的 New 方法实现非常简单，errorString 是一个结构体，内含一个 s 字符串，也只有一个 Error 方法，就可以认定为 error 类型，这是为什么呢？这一切的关键都在于 error 接口的定义。

```go
type error interface {
    Error() string
}
```

**在 Go 语言中，如果一个类型实现了某个 interface 中的所有方法，那么编译器就会认为该类型实现了此 interface，==它们是”一样“的==。**

### 5. 验证

在项目的 `internal/routers/api/v1` 下的` tag.go `文件中修改获取多个标签的 List 接口，用于验证 validator 是否正常。

```go
func (t Tag) List(c *gin.Context) {
	// 设置入参格式
	param := struct {
		Name  string `form:"name" binding:"max=100"`
		State uint8  `form:"state,default=1" binding:"oneof=0 1"`
	}{}
	// 创建响应
	response := app.NewResponse(c)
	// 进行入参校验
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf("app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	response.ToResponse(gin.H{})
	return
}
```

```bash
$ curl -X GET http://127.0.0.1:8000/api/v1/tags\?state\=6
{"code":10000001,"details":["State 必须是[0 1]中的一个"],"msg":"入参错误"}
```

另外还需要注意到 TagListRequest 的校验规则里其实并没有 required，**因此它的校验规则应该是有才校验，没有该入参的话，是默认无校验的**，也就是没有 state 参数，也应该可以正常请求。

```bash
$ curl -X GET http://127.0.0.1:8000/api/v1/tags          
{}
```

在 Response 中调用的是 `gin.H` 作为返回结果集，因此该输出结果正确。

### 6. 小结

在本章节中，介绍了 gin 框架如何通过 validator 进行参数校验，而在一些定制化场景中，常常需要自定义验证器，可以通过实现`bingding.Validator`接口的方式来替换自身的 validator。

```go
// binding/binding.go
type StructValidator interface {
    ValidateStruct(interface{}) error
    Engine() interface{}
}
func setupValidator() error {
    // 将你所自定义的 validator 写入
    binding.Validator = global.Validator
    return nil
}
```

也就是说如果有定制化需求，也完全可以自己实现一个验证器，效仿前面的模式，就可以替代 gin 框架原本的 validator 使用了。

## 六、模块开发：标签管理

在初步完成了业务接口的入参校验的逻辑处理后，接下来进入正式的业务模块的业务逻辑开发，在本章节将完成标签模块的接口代码编写，涉及的接口如下：

| 功能         | HTTP 方法 | 路径      |
| :----------- | :-------- | :-------- |
| 新增标签     | POST      | /tags     |
| 删除指定标签 | DELETE    | /tags/:id |
| 更新指定标签 | PUT       | /tags/:id |
| 获取标签列表 | GET       | /tags     |

### 1. 新建 model 方法

首先需要针对标签表进行处理，修改 `internal/model` 目录下的` tag.go `文件，针对标签模块的模型操作进行封装，并且只与实体产生关系，代码如下：

```go
package model

import (
   "demo/ch02/pkg/app"
   "github.com/jinzhu/gorm"
)

type Tag struct {
   *Model
   Name  string `json:"name"`
   State uint8  `json:"state"`
}

// tag.go
type TagSwagger struct {
   List  []*Tag
   Pager *app.Pager
}

func (t Tag) TableName() string {
   return "blog_tag"
}

// 使用 db *grom.DB 作为函数首参数传入
func (t Tag) Count(db *gorm.DB) (int, error) {
   var count int
   if t.Name != "" {
      // Where 设置筛选条件，接受 map、struct、string作为条件
      db = db.Where("name = ?", t.Name)
   }
   db = db.Where("state = ?", t.State)
   // Model 指定运行 DB 操作的模型实例，默认解析该结构体的名字为表名
   // Count 统计行为，用于统计模型的记录数
   if err := db.Model(&t).Where("is_del = ?", 0).Count(&count).Error; err != nil {
      return 0, err
   }
   return count, nil
}

func (t Tag) List(db *gorm.DB, pageOffset, pageSize int) ([]*Tag, error) {
   var tags []*Tag
   var err error
   if pageOffset >= 0 && pageSize > 0 {
      // Offset 偏移量，用于指定开始返回记录之前要跳过的记录数
      // Limit 限制检索的记录数
      db = db.Offset(pageOffset).Limit(pageSize)
   }
   if t.Name != "" {
      db = db.Where("name = ?", t.Name)
   }
   db = db.Where("state = ?", t.State)
   // Find 有两个参数，out 是数据接收者，where 是查询条件，可以代替 Where 来传入条件
   // err = e.g. db.Find(&tags, "is_del = 0").Error
   if err = db.Where("is_del = ?", 0).Find(&tags).Error; err != nil {
      return nil, err
   }
   return tags, nil
}

func (t Tag) Create(db *gorm.DB) error {
   return db.Create(&t).Error
}

func (t Tag) Update(db *gorm.DB) error {
   // Update 更新所选字段
   return db.Model(&Tag{}).Where("id = ? AND is_del = ?", t.ID, 0).Update(t).Error
}

func (t Tag) Delete(db *gorm.DB) error {
   // Delete 删除数据
   return db.Where("id = ? AND is_del = ?", t.Model.ID, 0).Delete(&t).Error
}
```

需要注意的是，在上述代码中，采用的是将`db *gorm.DB`作为函数首参数传入的方式，在实际开发中也可以基于结构体传入。

### 2. 处理 model 回调

在编写 model 代码时，并没有针对公共字段  created_on、modified_on、deleted_on、is_del 进行处理。可以通过设置 `model callback`的方式实现公共字段的处理，本项目使用的 ORM 是 GORM，其本身是提供回调支持的，因此可以根据自己的需要自定义 GORM 的回调操作，而在 GORM 中，可以分别进行如下的回调相关行为。

- 注册一个新的回调。
- 删除现有的回调。
- 替换现有的回调。
- 注册回调的先后顺序。

在本项目中使用到的“替换现有的回调”这一行为，修改`internal/model` 目录下的 model.go 文件，准备开始编写 model 的回调代码，下述所新增的回调代码均写入在` NewDBEngine `方法后。

```go
func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {}
func updateTimeStampForCreateCallback(scope *gorm.Scope) {}
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {}
func deleteCallback(scope *gorm.Scope) {}
func addExtraSpaceIfExist(str string) string {}
```

#### a. 新增行为的回调

```go
// 更新时间戳的新增行为的回调
// 当对数据库执行任何操作时，Scope 包含当前操作的信息
// Scope 允许复用通用的逻辑
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
   if !scope.HasError() {
      // 获取当前时间
      nowTime := time.Now().Unix()
      // scope.FiledByName() 获取当前是否包含所需字段
      if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
         // 若创建时间为空，则设置创建时间为当前时间
         // 通过判断 Filed.IsBlank 的值，可得知该字段值是否为空
         if createTimeField.IsBlank {
            // 通过 Filed.Set() 为字段赋值
            _ = createTimeField.Set(nowTime)
         }
      }
      if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
         // 若修改时间为空，则设置修改时间为当前时间
         if modifyTimeField.IsBlank {
            _ = modifyTimeField.Set(nowTime)
         }
      }
   }
}
```

- 通过调用 `scope.FieldByName` 方法，获取当前是否包含所需的字段。
- 通过判断 `Field.IsBlank` 的值，可以得知该字段的值是否为空。
- 若为空，则会调用 `Field.Set` 方法给该字段设置值，入参类型为 interface{}，内部也就是通过反射进行一系列操作赋值。

#### b. 更新行为的回调

```go
// 更新时间戳的更新行为的回调
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
   // scope.Get() 获取当前设置了标识 gorm:update_column 的字段属性
   if _, ok := scope.Get("gorm:update_column"); !ok {
      // 若不存在，即未自定义设置 update_column
      // 设置默认字段 ModifiedOn 的值为当前时间戳
      _ = scope.SetColumn("ModifiedOn", time.Now().Unix())
   }
}
```

- 通过调用 `scope.Get("gorm:update_column")` 去获取当前设置了标识 `gorm:update_column` 的字段属性。
- 若不存在，也就是没有自定义设置 `update_column`，那么将会在更新回调内设置默认字段` ModifiedOn `的值为当前的时间戳。

#### c. 删除行为的回调

```go
func deleteCallback(scope *gorm.Scope) {
   if !scope.HasError() {
      var extraOption string
      // 通过 scope.Get() 获取当前设置了标识 gorm:delete_option 的字段属性
      if str, ok := scope.Get("gorm:delete_option"); ok {
         extraOption = fmt.Sprint(str)
      }

      // 判断是否存在 DeletedOn 与 IsDel 字段
      deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")
      isDelField, hasIsDelField := scope.FieldByName("IsDel")

      // 若存在执行 UPDATE 进行软删除
      if !scope.Search.Unscoped && hasDeletedOnField && hasIsDelField {
         now := time.Now().Unix()
         // scope.Raw() 设置原始 sql
         scope.Raw(fmt.Sprintf(
            "UPDATE %v SET %v=%v,%v=%v%v%v",
            // scope.QuotedTableName() 获取当前所引用的表名
            scope.QuotedTableName(),
            // Quote() 使用引用字符串对数据库进行转义
            scope.Quote(deletedOnField.DBName),
            // AddToVars() 添加 value 作为 sql 的 vars，用于防止 SQL 注入
            scope.AddToVars(now),
            scope.Quote(isDelField.DBName),
            scope.AddToVars(1),
            // scope.CombinedConditionSql() 返回组合条件的 sql
            addExtraSpaceIfExist(scope.CombinedConditionSql()),
            addExtraSpaceIfExist(extraOption),
         )).Exec() // scope.Exec() 执行生成的sql
      } else {
         // 否则执行 DELETE 进行硬删除
         scope.Raw(fmt.Sprintf(
            "DELETE FROM %v%v%v",
            // 获取表明
            scope.QuotedTableName(),
            addExtraSpaceIfExist(scope.CombinedConditionSql()),
            addExtraSpaceIfExist(extraOption),
         )).Exec()
      }
   }
}
```

```go
// CombinedConditionSql return combined condition sql
func (scope *Scope) CombinedConditionSql() string {
   joinSQL := scope.joinsSQL()
   whereSQL := scope.whereSQL()
   if scope.Search.raw {
      whereSQL = strings.TrimSuffix(strings.TrimPrefix(whereSQL, "WHERE ("), ")")
   }
   return joinSQL + whereSQL + scope.groupSQL() +
      scope.havingSQL() + scope.orderSQL() + scope.limitAndOffsetSQL()
}
```

- 通过调用 `scope.Get("gorm:delete_option")` 去获取当前设置了标识 `gorm:delete_option` 的字段属性。
- 判断是否存在 `DeletedOn` 和 `IsDel` 字段，若存在则调整为执行 UPDATE 操作进行软删除（修改` DeletedOn` 和` IsDel` 的值），否则执行 DELETE 进行硬删除。
- 调用 `scope.QuotedTableName` 方法获取当前所引用的表名，并调用一系列方法针对 SQL 语句的组成部分进行处理和转移，最后在完成一些所需参数设置后调用 `scope.CombinedConditionSql` 方法完成 组合条件 SQL 语句的组装。

#### d. 注册回调行为

```go
package model

import (
	"demo/ch02/global"
	"demo/ch02/pkg/setting"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"

	// 引入 MYSQL驱动库进行初始化
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// 公共字段
type Model struct {
	ID         uint32 `gorm:"primary_key" json:"id"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	CreatedOn  uint32 `json:"created_on"`
	ModifiedOn uint32 `json:"modified_on"`
	DeletedOn  uint32 `json:"deleted_on"`
	IsDel      uint8  `json:"is_del"`
}

// 新增 NewDBEngine()
func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	// gorm.Open() 初始化一个 MySQL 连接，首先需要导入驱动
	db, err := gorm.Open(databaseSetting.DBType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime,
	))
	if err != nil {
		return nil, err
	}
	if global.ServerSetting.RunMode == "debug" {
		// 显示日志输出
		db.LogMode(true)
	}
	// 使用单表操作
	db.SingularTable(true)

	// 注册回调行为
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)

	// 设置空闲连接池中的最大连接数
	db.DB().SetMaxIdleConns(databaseSetting.MaxIdleConns)
	// 设置数据库的最大打开连接数。
	db.DB().SetMaxOpenConns(databaseSetting.MaxOpenConns)
	return db, nil
}

// 回调注册方法实现
func updateTimeStampForCreateCallback(scope *gorm.Scope) {...}
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {...}
func deleteCallback(scope *gorm.Scope) {...}
func addExtraSpaceIfExist(str string) string {...}
```

在最后回到 `NewDBEngine` 方法中，针对上述写的三个 Callback 方法进行回调注册，才能够让应用程序真正的使用上，至此，公共字段处理就完成了。

### 3. 新建 dao 方法

在项目的 `internal/dao` 目录下新建 `dao.go` 文件。

##### dao.go

```go
package dao

import "github.com/jinzhu/gorm"

type Dao struct {
   engine *gorm.DB
}
func New(engine *gorm.DB) *Dao {
   return &Dao{engine: engine}
}
```

接下来在同层级下新建 `tag.go` 文件，用于处理标签模块的 dao 操作。

##### tag.go

```go
package dao

import (
   "demo/ch02/internal/model"
   "demo/ch02/pkg/app"
)

func (d *Dao) CountTag(name string, state uint8) (int, error) {
   tag := model.Tag{Name: name, State: state}
   return tag.Count(d.engine)
}

func (d *Dao) GetTagList(name string, state uint8, page, pageSize int) ([]*model.Tag, error) {
   tag := model.Tag{Name: name, State: state}
   pageOffset := app.GetPageOffset(page, pageSize)
   return tag.List(d.engine, pageOffset, pageSize)
}

func (d *Dao) CreateTag(name string, state uint8, createdBy string) error {
   tag := model.Tag{
      Name:  name,
      State: state,
      Model: &model.Model{CreatedBy: createdBy},
   }
   return tag.Create(d.engine)
}

func (d *Dao) UpdateTag(id uint32, name string, state uint8, modifiedBy string) error {
   tag := model.Tag{
      Name:  name,
      State: state,
      Model: &model.Model{ID: id, ModifiedBy: modifiedBy},
   }
   return tag.Update(d.engine)
}

func (d *Dao) DeleteTag(id uint32) error {
   tag := model.Tag{Model: &model.Model{ID: id}}
   return tag.Delete(d.engine)
}
```

在 dao 层进行了数据访问对象的封装，并对针对业务所需字段进行了处理。

### 4. 新建 service 方法

在项目的 `internal/service` 目录下新建 `service.go` 文件。

##### service.go

```go
package service

import (
   "context"
   "demo/ch02/global"
   "demo/ch02/internal/dao"
)

type Service struct {
   ctx context.Context
   dao *dao.Dao
}
func New(ctx context.Context) Service {
   svc := Service{ctx: ctx}
   svc.dao = dao.New(global.DBEngine)
   return svc
}
```

修改同层级下的`tag.go`，用于处理标签模块的业务逻辑。

##### tag.go

```go
package service

import (
   "demo/ch02/internal/model"
   "demo/ch02/pkg/app"
)

// 设置方法的请求结构体和参数校验规则
type CountTagRequest struct {
   Name  string `form:"name" binding:"max=100"`
   State uint8 `form:"state,default=1" binding:"oneof=0 1"`
}
type TagListRequest struct {
   Name  string `form:"name" binding:"max=100"`
   State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}
type CreateTagRequest struct {
   Name      string `form:"name" binding:"required,min=3,max=100"`
   CreatedBy string `form:"created_by" binding:"required,min=3,max=100"`
   State     uint8  `form:"state,default=1" binding:"oneof=0 1"`
}
type UpdateTagRequest struct {
   ID         uint32 `form:"id" binding:"required,gte=1"`
   Name       string `form:"name" binding:"min=3,max=100"`
   State      uint8  `form:"state" binding:"required,oneof=0 1"`
   ModifiedBy string `form:"modified_by" binding:"required,min=3,max=100"`
}
type DeleteTagRequest struct {
   ID uint32 `form:"id" binding:"required,gte=1"`
}

func (svc *Service) CountTag(param *CountTagRequest) (int, error) {
   return svc.dao.CountTag(param.Name, param.State)
}
func (svc *Service) GetTagList(param *TagListRequest, pager *app.Pager) ([]*model.Tag, error) {
   return svc.dao.GetTagList(param.Name, param.State, pager.Page, pager.PageSize)
}
func (svc *Service) CreateTag(param *CreateTagRequest) error {
   return svc.dao.CreateTag(param.Name, param.State, param.CreatedBy)
}
func (svc *Service) UpdateTag(param *UpdateTagRequest) error {
   return svc.dao.UgopdateTag(param.ID, param.Name, param.State, param.ModifiedBy)
}
func (svc *Service) DeleteTag(param *DeleteTagRequest) error {
   return svc.dao.DeleteTag(param.ID)
}
```

在上述代码中，主要是定义了 Request 结构体作为接口入参的基准，而本项目由于并不会太复杂，所以直接放在了 service 层中便于使用，若后续业务不断增长，程序越来越复杂，service 也冗杂了，可以考虑将抽离一层接口校验层，便于解耦逻辑。

另外还在 service 中进行了一些简单的逻辑封装，在应用分层中，service 层主要是针对业务逻辑的封装，如果有一些业务聚合和处理可以在该层进行编码，同时也能较好的隔离上下两层的逻辑。

### 5. 新增业务错误码

在项目的 `pkg/errcode` 下新建 `module_code.go` 文件，针对标签模块，写入错误代码。

```go
package errcode

var (
   ErrorGetTagListFail = NewError(20010001, "获取标签列表失败")
   ErrorCreateTagFail  = NewError(20010002, "创建标签失败")
   ErrorUpdateTagFail  = NewError(20010003, "更新标签失败")
   ErrorDeleteTagFail  = NewError(20010004, "删除标签失败")
   ErrorCountTagFail   = NewError(20010005, "统计标签失败")
)
```

### 6. 新增路由方法

修改 `internal/routers/api/v1` 项目目录下的` tag.go `文件

```go
func (t Tag) List(c *gin.Context) {
   // 设置入参格式与参数校验规则
   param := service.TagListRequest{}
   // 初始化响应
   response := app.NewResponse(c)
   // 进行入参校验
   valid, errs := app.BindAndValid(c, &param)
   if !valid {
      global.Logger.Errorf("app.BindAndValid errs: %v", errs)
      response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
      return
   }

   svc := service.New(c.Request.Context())
   pager := app.Pager{Page: app.GetPage(c), PageSize: app.GetPageSize(c)}
   // 获取标签总数
   totalRows, err := svc.CountTag(&service.CountTagRequest{Name: param.Name, State: param.State})
   if err != nil {
      global.Logger.Errorf("svc.CountTag err: %v", err)
      response.ToErrorResponse(errcode.ErrorCountTagFail)
      return
   }
   // 获取标签列表
   tags, err := svc.GetTagList(&param, &pager)
   if err != nil {
      global.Logger.Errorf("svc.GetTagList err: %v", err)
      response.ToErrorResponse(errcode.ErrorGetTagListFail)
      return
   }
   
   // 序列化结果集
   response.ToResponseList(tags, totalRows)
   return
}
```

在上述代码中，完成了获取标签列表接口的处理方法，在方法中完成了入参校验和绑定、获取标签总数、获取标签列表、 序列化结果集等四大功能板块的逻辑串联和日志、错误处理。需要注意的是方法实现中的入参校验和绑定的处理代码基本都差不多，`tag.go`全部代码如下：

```go
package v1

import (
	"demo/ch02/global"
	"demo/ch02/internal/service"
	"demo/ch02/pkg/app"
	"demo/ch02/pkg/convert"
	"demo/ch02/pkg/errcode"
	"github.com/gin-gonic/gin"
)

type Tag struct {}
func NewTag() Tag {
	return Tag{}
}

func (t Tag) Get(c *gin.Context) {}

// @Summary 获取多个标签
// @Produce  json
// @Param name query string false "标签名称" maxlength(100)
// @Param state query int false "状态" Enums(0, 1) default(1)
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} model.TagSwagger "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags [get]
func (t Tag) List(c *gin.Context) {
	// 设置入参格式与参数校验规则
	param := service.TagListRequest{}
	// 初始化响应
	response := app.NewResponse(c)
	// 进行入参校验
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf("app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	pager := app.Pager{Page: app.GetPage(c), PageSize: app.GetPageSize(c)}
	// 获取标签总数
	totalRows, err := svc.CountTag(&service.CountTagRequest{Name: param.Name, State: param.State})
	if err != nil {
		global.Logger.Errorf("svc.CountTag err: %v", err)
		response.ToErrorResponse(errcode.ErrorCountTagFail)
		return
	}
	// 获取标签列表
	tags, err := svc.GetTagList(&param, &pager)
	if err != nil {
		global.Logger.Errorf("svc.GetTagList err: %v", err)
		response.ToErrorResponse(errcode.ErrorGetTagListFail)
		return
	}

	// 序列化结果集
	response.ToResponseList(tags, totalRows)
	return
}

// @Summary 新增标签
// @Produce  json
// @Param name body string true "标签名称" minlength(3) maxlength(100)
// @Param state body int false "状态" Enums(0, 1) default(1)
// @Param created_by body string true "创建者" minlength(3) maxlength(100)
// @Success 200 {object} model.Tag "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags [post]
func (t Tag) Create(c *gin.Context) {
	// 入参校验与绑定
	param := service.CreateTagRequest{}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf("app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.CreateTag(&param)
	if err != nil {
		global.Logger.Errorf("svc.CreateTag err: %v", err)
		response.ToErrorResponse(errcode.ErrorCreateTagFail)
		return
	}

	response.ToResponse(gin.H{})
	return
}

// @Summary 更新标签
// @Produce  json
// @Param id path int true "标签 ID"
// @Param name body string false "标签名称" minlength(3) maxlength(100)
// @Param state body int false "状态" Enums(0, 1) default(1)
// @Param modified_by body string true "修改者" minlength(3) maxlength(100)
// @Success 200 {array} model.Tag "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags/{id} [put]
func (t Tag) Update(c *gin.Context) {
	param := service.UpdateTagRequest{
		// 将 string 类型转换为 uint32
		ID: convert.StrTo(c.Param("id")).MustUInt32(),
	}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf("app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.UpdateTag(&param)
	if err != nil {
		global.Logger.Errorf("svc.UpdateTag err: %v", err)
		response.ToErrorResponse(errcode.ErrorUpdateTagFail)
		return
	}

	response.ToResponse(gin.H{})
	return
}

// @Summary 删除标签
// @Produce  json
// @Param id path int true "标签 ID"
// @Success 200 {string} string "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/tags/{id} [delete]
func (t Tag) Delete(c *gin.Context) {
	param := service.DeleteTagRequest{ID: convert.StrTo(c.Param("id")).MustUInt32()}
	response := app.NewResponse(c)
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.Logger.Errorf("app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	svc := service.New(c.Request.Context())
	err := svc.DeleteTag(&param)
	if err != nil {
		global.Logger.Errorf("svc.DeleteTag err: %v", err)
		response.ToErrorResponse(errcode.ErrorDeleteTagFail)
		return
	}

	response.ToResponse(gin.H{})
	return
}
```

### 7. 验证接口

启动服务，对标签模块的接口进行验证，请注意，验证示例中的 `{id}`，代指占位符，也就是填写实际调用中希望处理的标签 ID 即可。

#### a. 新增标签

使用 postman 进行接口测试，使用 `json` 作为入参无法成功创建标签，会报错入参错误，原因不明。

![image-20220514181031590](https://raw.githubusercontent.com/tonshz/test/master/img/202205141810834.png)

若不满足参数校验规则则会报错，例如设置 name 值为 Go。

```json
{
    "code": 10000001,
    "details": [
        "Name长度必须至少为3个字符"
    ],
    "msg": "入参错误"
}
```

#### b. 获取标签列表

![image-20220514181652124](https://raw.githubusercontent.com/tonshz/test/master/img/202205141816180.png)

```json
{
    "list": [
        {
            "id": 1,
            "created_by": "test",
            "modified_by": "",
            "created_on": 1652522978,
            "modified_on": 1652522978,
            "deleted_on": 0,
            "is_del": 0,
            "name": "create_tag_test",
            "state": 1
        },
        {
            "id": 2,
            "created_by": "test",
            "modified_by": "",
            "created_on": 1652523217,
            "modified_on": 1652523217,
            "deleted_on": 0,
            "is_del": 0,
            "name": "Java",
            "state": 1
        }
    ],
    "pager": {
        "page": 1, 
        "page_size": 2,
        "total_rows": 3
    }
}
```

修改 page 参数为 2。

```json
{
    "list": [
        {
            "id": 3,
            "created_by": "test",
            "modified_by": "",
            "created_on": 1652523235,
            "modified_on": 1652523235,
            "deleted_on": 0,
            "is_del": 0,
            "name": "Golang",
            "state": 1
        }
    ],
    "pager": {
        "page": 2,
        "page_size": 2,
        "total_rows": 3
    }
}
```

#### c. 修改标签

此处 postman 失败，不清楚原因，与参数校验相关，始终报错 `"Name长度必须至少为3个字符"`。

```bash
$ curl -X PUT http://127.0.0.1:8000/api/v1/tags/{id} -F state=0 -F modified_by=eddycjy
{}
```

#### d. 删除标签

![image-20220514185320956](https://raw.githubusercontent.com/tonshz/test/master/img/202205141853990.png)

删除标签后，数据库内容更新。

![image-20220514185303772](https://raw.githubusercontent.com/tonshz/test/master/img/202205141853815.png)

### 8. 发现问题：零值未更新

在完成了接口的检验后，还需要确定一下数据库内的数据变更是否正确。在经过一系列的对比后，发现在调用修改标签的接口时，通过接口入参，是希望将 id 为 1 的标签状态修改为 0，但是在对比后发现数据库内它的状态值仍然是 1，而且 SQL 语句内也没有出现 state 字段的设置，控制台输出的 SQL 语句如下：

```bash
UPDATE `blog_tag` SET `id` = 1, `modified_by` = 'eddycjy', `modified_on` = xxxxx  WHERE `blog_tag`.`id` = 1
```

**原因是只要字段是零值的情况下，GORM 就不会对该字段进行变更。**实际上，这有一个概念上的问题，先入为主的认为它一定会变更，其实是不对的，因为在程序中使用的是 struct 的方式进行更新操作，而在 GORM 中使用 struct 类型传入进行更新时，**GORM 是不会对值为零值的字段进行变更**。更根本的原因是因为在识别这个结构体中的这个字段值时，**很难判定是真的是零值，还是外部传入恰好是该类型的零值**，GORM 在这块并没有过多的去做特殊识别。

### 9. 解决问题

修改项目的 `internal/model` 目录下的` tag.go `文件里的` Update `方法。

```go
// 修改传入零值后，数据库未发生变化的问题
func (t Tag) Update(db *gorm.DB, values interface{}) error {
   // Update 更新所选字段
   if err := db.Model(t).Where("id = ? AND is_del = ?", t.ID, 0).Updates(values).Error; err != nil {
      return err
   }
   return nil
}
```

修改项目的 `internal/dao` 目录下的`tag.go `文件里的 `UpdateTag `方法。

```go
func (d *Dao) UpdateTag(id uint32, name string, state uint8, modifiedBy string) error {
   tag := model.Tag{
      Model: &model.Model{ID: id},
   }
   values := map[string]interface{}{
      "state":       state,
      "modified_by": modifiedBy,
   }
   if name != "" {
      values["name"] = name
   }
   return tag.Update(d.engine, values)
}
```

重新运行程序，请求修改标签接口，检查数据是否正常修改，在正确的情况下，该 id 为 1 的标签，modified_by 为 test，modified_on 应修改为当前时间戳，state 为 0。

### 10. 文章管理模块

参见 Github 代码。

## 七、上传图片和文件服务

实现文章的封面图片上传并用文件服务对外提供静态文件的访问服务，在上传图片后，可以通过约定的地址访问到对应的图片资源。

### 1. 新增配置

修改`configs/config.yaml` 配置文件，添加上传相关的配置。

```yaml
App: # 应用配置
  ...
  # 添加上传相关配置
  UploadSavePath: storage/uploads # 上传文件的最终保存目录
  UploadServerUrl: http://127.0.0.1:8000/static # 上传文件后用于展示的文件服务地址
  UploadImageMaxSize: 5  # 上传文件所允许的最大空间(MB)
  UploadImageAllowExts: # 上传文件所允许的文件后缀
    - .jpg
    - .jpeg
    - .png
```

一共新增了四项上传文件所必须的配置项，分别代表的作用如下：

- `UploadSavePath`：上传文件的最终保存目录。
- `UploadServerUrl`：上传文件后的用于展示的文件服务地址。
- `UploadImageMaxSize`：上传文件所允许的最大空间大小（MB）。
- `UploadImageAllowExts`：上传文件所允许的文件后缀。

接下来需要在对应的配置结构体上新增上传相关属性，修改`pkg/setting/section.go`。

```go
// 应用配置结构体
type AppSettingS struct {
   DefaultPageSize      int
   MaxPageSize          int
   LogSavePath          string
   LogFileName          string   
   LogFileExt           string
   UploadSavePath       string
   UploadServerUrl      string
   UploadImageMaxSize   int
   UploadImageAllowExts []string
}
```

### 2. 上传文件

接下来需要编写一个上传文件的工具库，它的主要作用是针对上传文件时的一些相关处理。在项目的 `pkg` 目录下新建 `util` 目录，并创建` md5.go `文件。

```go
package util

import (
   "crypto/md5"
   "encoding/hex"
)

// 将上传后的文件名进行格式化，将文件名 MD5 编码后在进行写入
func EncodeMD5(value string) string {
   m := md5.New()
   m.Write([]byte(value))
   return hex.EncodeToString(m.Sum(nil))
}
```

该方法用于针对上传后的文件名格式化，简单来讲，将文件名 MD5 后再进行写入，防止直接把原始名称就暴露出去了。接下来在项目的 `pkg/upload` 目录下新建`file.go `文件。

```go
package upload

import (
	"demo/ch02/global"
	"demo/ch02/pkg/util"
	"path"
	"strings"
)

// 定义 FileType 为 int 的别名
type FileType int

// 使用 FileType 作为类别表示的基础类型
const TypeImage FileType = iota + 1 // TypeImage = 1

// 返回经过加密处理后的文件名
func GetFileName(name string) string {
	ext := GetFileExt(name)
	// 返回没有后缀的文件名
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)
	return fileName + ext
}

func GetFileExt(name string) string {
	// 调用 path.Ext() 进行循环查找
	return path.Ext(name)
}

func GetSavePath() string {
	// 返回配置中的文件保存目录
	return global.AppSetting.UploadSavePath
}
```

在上述代码中，用到了两个比较常见的语法，首先是定义了 `FileType `为` int `的类型别名，并且利用` FileType `作为类别标识的基础类型，并 iota 作为了它的初始值。

实际上，在 Go 语言中 iota 相当于是一个 `const `的常量计数器，也可以理解为枚举值，第一个声明的 iota 的值为 0，在新的一行被使用时，它的值都会自动递增。

当然了，也可以像代码中那样，在初始的第一个声明时进行手动加一，那么它将会从 1 开始递增。其本质上是为了后续有其它的需求，能标准化的进行处理，例如：

```go
const (
    TypeImage FileType = iota + 1 // 1
    TypeExcel // 2
    TypeTxt // 3
)
```

如果未来需要支持其他的上传文件类型修改就很方便了，手工定义1，2，3，4不是可取的做法。

另外还声明了三个文件相关的方法，其作用分别如下：

- `GetFileName`：获取文件名称，先是通过获取文件后缀并筛出原始文件名进行 MD5 加密，最后返回经过加密处理后的文件名。
- `GetFileExt`：获取文件后缀，主要是通过调用 `path.Ext` 方法进行循环查找”.“符号，最后通过切片索引返回对应的文化后缀名称。
- `GetSavePath`：获取文件保存地址，这里直接返回配置中的文件保存目录即可，也便于后续的调整。

在完成了文件相关参数获取的方法后，接下来需要编写检查文件的相关方法，因为需要确保在文件写入时它已经达到了必备条件，否则要给出对应的标准错误提示，继续在文件内新增如下代码：

```go
// 检测路径是否存在
func CheckSavePath(dst string) bool {
   // 利用 os.Stat() 方法所返回的 error 值与系统所定义的 oserror.ErrNotExist 是否相等
   _, err := os.Stat(dst)
   return os.IsNotExist(err)
}

// 检测文件后缀是否满足设置条件
func CheckContainExt(t FileType, name string) bool {
   ext := GetFileExt(name)
   // 同一转换为大写进行匹配
   ext = strings.ToUpper(ext)
   switch t {
   case TypeImage:
      // 与配置文件中设置的允许的文件后缀名进行比较
      for _, allowExt := range global.AppSetting.UploadImageAllowExts {
         if strings.ToUpper(allowExt) == ext {
            return true
         }
      }
   }
   return false
}

// 检测最大大小是否超出最大限制
func CheckMaxSize(t FileType, f multipart.File) bool {
   content, _ := ioutil.ReadAll(f)
   size := len(content)
   switch t {
   case TypeImage:
      if size >= global.AppSetting.UploadImageMaxSize*1024*1024 {
         return true
      }
   }
   return false
}

// 检测文件权限是否足够
func CheckPermission(dst string) bool {
   // 与 CheckSavePath() 类似，与 oserror.ErrPermission 判断是否相等
   _, err := os.Stat(dst)
   return os.IsPermission(err)
}
```

- `CheckSavePath`：检查保存目录是否存在，通过调用 `os.Stat` 方法获取文件的描述信息 `FileInfo`，并调用 `os.IsNotExist` 方法进行判断，**其原理是利用 `os.Stat` 方法所返回的 error 值与系统中所定义的 `oserror.ErrNotExist` 进行判断，以此达到校验效果。**
- `CheckPermission`：检查文件权限是否足够，与 `CheckSavePath` 方法原理一致，是利用 `oserror.ErrPermission` 进行判断。
- `CheckContainExt`：检查文件后缀是否包含在约定的后缀配置项中，需要的是所上传的文件的后缀有可能是大写、小写、大小写等，因此需要调用 `strings.ToUpper` 方法统一转为大写（固定的格式）来进行匹配。
- `CheckMaxSize`：检查文件大小是否超出最大大小限制。

在完成检查文件的一些必要操作后，就可以编写涉及文件写入/创建的相关操作，继续在文件内新增如下代码：

```go
// 创建上传文件时所使用的保存目录
func CreateSavePath(dst string, perm os.FileMode) error {
   // os.MkdirAll 会根据传入的 os.FileMode 权限位递归创所需的所有目录结构
   // 若目录已存在则不会进行任何操作，直接返回 nil
   err := os.MkdirAll(dst, perm)
   if err != nil {
      return err
   }
   return nil
}

// 保存所上传的文件
func SaveFile(file *multipart.FileHeader, dst string) error {
   // 通过 file.Open 打开源地址的文件
   src, err := file.Open()
   if err != nil {
      return err
   }
   defer src.Close()
   // 通过 os.Create 创建目标地址的文件
   out, err := os.Create(dst)
   if err != nil {
      return err
   }
   defer out.Close()
   // 结合 io.Copy 实现两者之间的文件内容拷贝
   _, err = io.Copy(out, src)
   return err
}
```

- `CreateSavePath`：创建在上传文件时所使用的保存目录，在方法内部调用的 `os.MkdirAll` 方法，该方法将会以传入的 `os.FileMode` 权限位去递归创建所需的所有目录结构，若涉及的目录均已存在，则不会进行任何操作，直接返回 nil。
- `SaveFile`：保存所上传的文件，该方法主要是通过调用 `os.Create` 方法创建目标地址的文件，再通过 `file.Open` 方法打开源地址的文件，结合 `io.Copy` 方法实现两者之间的文件内容拷贝。

### 3. 新建 service 方法

将上一步所编写的上传文件工具库与具体的业务接口结合起来，在项目下的 `internal/service` 目录新建 `upload.go `文件。

```go
package service

import (
   "demo/ch02/global"
   "demo/ch02/pkg/upload"
   "errors"
   "mime/multipart"
   "os"
)

type FileInfo struct {
   Name      string
   AccessUrl string
}

func (svc *Service) UploadFile(fileType upload.FileType, file multipart.File, fileHeader *multipart.FileHeader) (*FileInfo, error) {
   fileName := upload.GetFileName(fileHeader.Filename)
   if !upload.CheckContainExt(fileType, fileName) {
      return nil, errors.New("file suffix is not supported.")
   }
   if upload.CheckMaxSize(fileType, file) {
      return nil, errors.New("exceeded maximum file limit.")
   }
   uploadSavePath := upload.GetSavePath()
   if upload.CheckSavePath(uploadSavePath) {
      if err := upload.CreateSavePath(uploadSavePath, os.ModePerm); err != nil {
         return nil, errors.New("failed to create save directory.")
      }
   }
   if upload.CheckPermission(uploadSavePath) {
      return nil, errors.New("insufficient file permissions.")
   }
   dst := uploadSavePath + "/" + fileName
   if err := upload.SaveFile(fileHeader, dst); err != nil {
      return nil, err
   }
   accessUrl := global.AppSetting.UploadServerUrl + "/" + fileName
   return &FileInfo{Name: fileName, AccessUrl: accessUrl}, nil
}
```

在 `UploadFile`方法中，主要是通过获取文件所需的基本信息，接着对其进行业务所需的文件检查（文件大小是否符合需求、文件后缀是否达到要求），并且判断在写入文件前对否具备必要的写入条件（目录是否存在、权限是否足够），最后在检查通过后再进行真正的写入文件操作。

### 4. 新增业务错误码

在项目的 `pkg/errcode` 下的` module_code.go `文件，针对上传模块，新增如下错误代码：

```go
package errcode

var (
   ...
   ErrorUploadFileFail    = NewError(20030001, "上传文件失败")
)
```

### 5. 新增路由方法

接下来需要编写上传文件的路由方法，将整套上传逻辑给串联起来，在项目的 `internal/routers` 目录下新建 `upload.go `文件。

```go
package routers

import (
   "demo/ch02/global"
   "demo/ch02/internal/service"
   "demo/ch02/pkg/app"
   "demo/ch02/pkg/convert"
   "demo/ch02/pkg/errcode"
   "demo/ch02/pkg/upload"
   "github.com/gin-gonic/gin"
)

type Upload struct{}

func NewUpload() Upload {
   return Upload{}
}

func (u Upload) UploadFile(c *gin.Context) {
   response := app.NewResponse(c)
   // 参数获取与检测
   // 通过 c.Request.FormFile() 读取入参 file 字段的上传文件信息
   file, fileHeader, err := c.Request.FormFile("file")
   if err != nil {
      response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
      return
   }
   // 使用入参 type 字段作为上传文件类型的确认依据
   fileType := convert.StrTo(c.PostForm("type")).MustInt()
   if fileHeader == nil || fileType <= 0 {
      response.ToErrorResponse(errcode.InvalidParams)
      return
   }
   
   // 调用 service 方法完成文件上传、文件保存，并返回文件展示地址
   svc := service.New(c.Request.Context())
   fileInfo, err := svc.UploadFile(upload.FileType(fileType), file, fileHeader)
   if err != nil {
      global.Logger.Errorf(c, "svc.UploadFile err: %v", err)
      response.ToErrorResponse(errcode.ErrorUploadFileFail.WithDetails(err.Error()))
      return
   }
   
   // 返回文件展示地址
   response.ToResponse(gin.H{
      "file_access_url": fileInfo.AccessUrl,
   })
}
```

在上述代码中，通过 `c.Request.FormFile` 读取入参 file 字段的上传文件信息，并利用入参 type 字段作为所上传文件类型的确立依据（也可以通过解析上传文件后缀来确定文件类型），最后通过入参检查后进行 `svc.UploadFile`的调用，完成上传和文件保存，返回文件的展示地址。

至此，业务接口的编写就完成了，下一步需要添加路由，让外部能够访问到该接口，依旧是在 `internal/routers` 目录下的` router.go `文件，在其中新增上传文件的对应路由，如下：

```go
package routers

import (
   ...
)

func NewRouter() *gin.Engine {
   ...
   // 添加上传文件的对应路由
   upload := api.NewUpload()
   r.POST("/upload/file", upload.UploadFile)
   // 使用路由组设置访问路由的统一前缀 e.g. /api/v1
   // 此处定义了一个路由组 /api/v1
   apiv1 := r.Group("/api/v1")
   {
      ...
   }
   return r
}
```

新增了` POST `方法的 `/upload/file` 路由，并调用其 `upload.UploadFile `方法来提供接口的方法响应，至此整体的路由到业务接口的联通就完成了。

### 6. 验证接口

检查接口返回是否与期望的一致，主体是由 `UploadServerUrl `与加密后的文件名称相结合。

![image-20220515203547801](https://raw.githubusercontent.com/tonshz/test/master/img/202205152036850.png)

![image-20220515203604648](https://raw.githubusercontent.com/tonshz/test/master/img/202205152036683.png)

### 7. 文件服务

在进行接口的返回结果校验时，会发现上小节中` file_access_url `这个地址无法访问到对应的文件资源，检查文件资源也确实存在 `storage/uploads` 目录下。

**实际上是需要设置文件服务去提供静态资源的访问**，才能实现让外部请求本项目` HTTP Server` 时同时提供静态资源的访问，实际上在 gin 中实现 File Server 是非常简单的，我们需要在 `NewRouter `方法中，新增如下路由：

```go
package routers

import (
	...
)

func NewRouter() *gin.Engine {
	...
	// 添加上传文件的对应路由
	upload := api.NewUpload()
	r.POST("/upload/file", upload.UploadFile)
	// 提供静态资源的访问
	// Static 只能展示文件，StaticFS 可以连目录也展示
	// http.Dir() 实现了 FileSystem接口，利用本地目录实现一个文件系统（FileSystem）
	r.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath))
	// 使用路由组设置访问路由的统一前缀 e.g. /api/v1
	// 此处定义了一个路由组 /api/v1
	apiv1 := r.Group("/api/v1")
	{
		...
	}
	return r
}
```

新增 `StaticFS` 路由完毕后，重新重启应用程序，再次访问 file_access_url 所输出的地址就可以查看到刚刚上传的静态文件了。

+ `router.Static()` : 指定某个目录为静态资源目录，可直接访问这个目录下的资源，`url `要具体到资源名称。

+ `router.StaticFS()` : 比前面一个多了个功能，当目录下不存 index.html 文件时，会列出该目录下的所有文件，**可以自定义文件系统。**

+ `router.StaticFile()`:  指定某个具体的文件作为静态资源访问。

```go
func main() {
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.StaticFS("/more_static", http.Dir("my_file_system"))
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}
```

### 8. 原理

设置一个` r.StaticFS` 的路由，就可以拥有一个文件服务，并且能够提供静态资源的访问。既然能够读取到文件的展示，那么就是在访问 `$HOST/static` 时，应用程序会读取到 `storage/uploads` 下的文件。可以看看 `StaticFS` 方法到底做了什么事，方法原型如下：

```go
// StaticFS works just like `Static()` but a custom `http.FileSystem` can be used instead.
// Gin by default user: gin.Dir()
func (group *RouterGroup) StaticFS(relativePath string, fs http.FileSystem) IRoutes {
   if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
      // 提供静态文件夹时不能使用 URL 参数
      panic("URL parameters can not be used when serving a static folder")
   }
   handler := group.createStaticHandler(relativePath, fs)
   urlPattern := path.Join(relativePath, "/*filepath")

   // Register GET and HEAD handlers
   group.GET(urlPattern, handler)
   group.HEAD(urlPattern, handler)
   return group.returnObj()
}
```

首先可以看到在暴露的 URL 中程序禁止了“*”和“:”符号的使用，然后通过 `createStaticHandler` 创建了静态文件服务，其实质最终调用的还是 `fileServer.ServeHTTP` 和对应的处理逻辑，如下：

```go
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
   absolutePath := group.calculateAbsolutePath(relativePath)
   fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

   return func(c *Context) {
      if _, noListing := fs.(*onlyFilesFS); noListing {
         c.Writer.WriteHeader(http.StatusNotFound)
      }

      file := c.Param("filepath")
      // 检查文件是否存在以及我们是否有权访问它
      // Check if file exists and/or if we have permission to access it
      f, err := fs.Open(file)
      if err != nil {
         c.Writer.WriteHeader(http.StatusNotFound)
         c.handlers = group.engine.noRoute
         // Reset index
         c.index = -1
         return
      }
      f.Close()

      // 最终实质调用
      fileServer.ServeHTTP(c.Writer, c.Request)
   }
}
```

在 `createStaticHandler` 方法中，需要留意下 `http.StripPrefix` 方法的调用，实际上在静态文件服务中很常见，它主要作用是从请求 URL 的路径中删除给定的前缀，然后返回一个 Handler。

另外在 `StaticFS `方法中看到 `urlPattern := path.Join(relativePath, "/*filepath")` 的代码块，`/*filepath` 通过语义可得知它是路由的处理逻辑，而 gin 的路由是基于 `httprouter` 的，通过查阅文档可以得到如下信息：

```lua
Pattern: /src/*filepath
 /src/                     match
 /src/somefile.go          match
 /src/subdir/somefile.go   match
```

**简单来讲，`*filepath` 将会匹配所有文件路径，但是前提是 `*filepath` 标识符必须在 Pattern 的最后。**

## 八、对接口进行访问控制

在完成了相关的业务接口的开发后，还有一个问题，这些 API 接口，没有鉴权功能，也就是所有知道地址的人都可以请求该项目的 API 接口和 Swagger 文档，甚至有可能会被网络上的端口扫描器扫描到后滥用，这非常的不安全，怎么办呢。实际上，应该要考虑做纵深防御，对 API 接口进行访问控制。

目前市场上比较常见的两种 API 访问控制方案，分别是 OAuth 2.0 和 JWT(JSON Web Token)，但实际上这两者并不能直接的进行对比，因为它们是两个完全不同的东西，对应的应用场景也不一样，可以先大致了解，如下：

- OAuth 2.0：**OAuth 2.0 是一种授权框架**，本质上是一个授权的行业标准协议，提供了一整套的授权机制的指导标准，常用于使用第三方登陆的情况，像是在网站登录时，会有提供其它第三方站点（例如用微信、QQ、Github 账号）关联登陆的，往往就是用 OAuth 2.0 的标准去实现的。并且 OAuth 2.0 会相对重一些，常常还会授予第三方应用去获取到对应账号的个人基本信息等等。在实现 OAuth 2.0 时可以将 JWT 作为一种认证机制使用。
- JWT：**JWT 是一种认证协议**，与 OAuth 2.0 完全不同，它常用于前后端分离的情况，能够非常便捷的给 API 接口提供安全鉴权，因此在本章节采用的就是 JWT 的方式，来实现 API 访问控制功能。

### 1. JWT 是什么

JSON Web Token（JWT）是一个开放标准（RFC7519），它定义了一种紧凑且自包含的方式，用于在各方之间作为 JSON 对象安全地传输信息。 由于此信息是经过数字签名的，因此可以被验证和信任。 可以使用使用 RSA 或 ECDSA 的公用/专用密钥对对 JWT 进行签名，其格式如下：

![image](https://raw.githubusercontent.com/tonshz/test/master/img/202205152113517.jpeg)

JSON Web 令牌（JWT）是由紧凑的形式三部分组成，这些部分由点 “.“ 分隔，组成为 `”xxxxx.yyyyy.zzzzz“ `的格式，三个部分分别代表的意义如下：

- Header：头部。
- Payload：有效载荷。
- Signature：签名。

#### a. Header

Header（头部）通常由两部分组成，**分别是令牌的类型和所使用的签名算法**（HMAC SHA256、RSA 等），其会组成一个 JSON 对象用于描述其元数据，例如：

```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

在上述 JSON 中 `alg` 字段表示所使用的签名算法，默认是 HMAC SHA256（HS256），而 type 字段表示所使用的令牌类型，使用的 JWT 令牌类型，在最后会对上面的 JSON 对象进行` base64UrlEncode `算法进行转换成为 JWT 的第一部分。

#### b. Payload

Payload（有效负载）也是一个 JSON 对象，**主要存储在 JWT 中实际传输的数据**，如下：

```json
{
  "sub": "1234567890",
  "name": "John Doe",
  "admin": true
}
```

- `aud（Audience）`：受众，也就是接受 JWT 的一方。
- `exp（ExpiresAt）`：所签发的 JWT 过期时间，过期时间必须大于签发时间。
- `jti（JWT Id）`：JWT 的唯一标识。
- `iat（IssuedAt）`：签发时间。
- `iss（Issuer）`：JWT 的签发者。
- `nbf（Not Before）`：JWT 的生效时间，如果未到这个时间则为不可用。
- `sub（Subject）`：主题。

同样也会对该 JSON 对象进行 base64UrlEncode 算法将其转换为 JWT Token 的第二部分。

这时候需要注意一个问题点，也就是 JWT 在转换时用的 base64UrlEncode 算法，也就是它是可逆的，因此一些敏感信息不要放到 JWT 中，若有特殊情况一定要放，**也应当进行一定的加密处理。**

#### c. Signature

Signature（签名）部分是对前面两个部分组合（Header+Payload）进行约定算法和规则的签名，**而签名将会用于校验消息在整个过程中有没有被篡改**，并且对有使用私钥进行签名的令牌，它还可以验证 JWT 的发送者是否它的真实身份。

在签名的生成上，在应用程序指定了密钥（secret）后，会使用传入的指定签名算法（默认是 HMAC SHA256），然后通过下述的签名方式来完成 Signature（签名）部分的生成，如下：

```lua
HMACSHA256(
  base64UrlEncode(header) + "." +
  base64UrlEncode(payload),
  secret)
```

可以看出 JWT 的第三部分是由 Header、Payload 以及 Secret 的算法组成而成的，因此它最终可达到用于校验消息是否被篡改的作用之一，因为如果一旦被篡改，Signature 就会无法对上。

#### d. Base64UrlEncode

实际上 Base64UrlEncode 是 Base64 算法的变种，为什么要变呢，原因是在实际开发过程中经常可以看到 JWT 令牌会被放入 Header 或 Query Param 中（也就是 URL）。

而在 URL 中，一些个别字符是有特殊意义的，例如：“+”、“/”、“=” 等等，因此在 Base64UrlEncode 算法中，会对其进行替换，例如：“+” 替换为 “-”、“/” 替换成 “_”、“=” 会被进行忽略处理，以此来保证 JWT 令牌的在 URL 中的可用性和准确性。

### 2. JWT 的使用场景

通常会先在内部约定好 JWT 令牌的交流方式，像是存储在 Header、Query Param、Cookie、Session 都有，但最常见的是存储在 Header 中。然后服务端提供一个获取 JWT 令牌的接口方法，返回而客户端去使用，在客户端请求其余的接口时需要带上所签发的 JWT 令牌，然后服务端接口也会到约定位置上获取 JWT 令牌来进行鉴权处理，以此流程来鉴定是否合法。

### 3. 安装 JWT

拉取 `jwt-go`，该库提供了 JWT 的 Go 实现，能够便捷的提供 JWT 支持，不需要自己去实现。

```bash
$ go get -u github.com/dgrijalva/jwt-go@v3.2.0
```

### 4. 配置 JWT

#### a. 创建认证表

在介绍 JWT 和其使用场景时，了解了实际上需要一个服务端的接口来提供 JWT 令牌的签发，并且可以将自定义的私有信息存入其中，那么必然需要一个地方来存储签发的凭证，否则谁来都签发，似乎不大符合实际的业务需求，因此要创建一个新的数据表，用于存储签发的认证信息，表 SQL 语句如下：

```sql
CREATE TABLE `blog_auth` (
                             `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                             `app_key` varchar(20) DEFAULT '' COMMENT 'Key',
                             `app_secret` varchar(50) DEFAULT '' COMMENT 'Secret',
                             `created_on` int(10) unsigned DEFAULT '0' COMMENT '创建时间',
                             `created_by` varchar(100) DEFAULT '' COMMENT '创建人',
                             `modified_on` int(10) unsigned DEFAULT '0' COMMENT '修改时间',
                             `modified_by` varchar(100) DEFAULT '' COMMENT '修改人',
                             `deleted_on` int(10) unsigned DEFAULT '0' COMMENT '删除时间',
                             `is_del` tinyint(3) unsigned DEFAULT '0' COMMENT '是否删除 0 为未删除、1 为已删除',
                              PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='认证管理';
```

上述表 SQL 语句的主要作用是创建了一张名为` blog_auth` 的表，其核心是 `app_key `和 `app_secret` 字段，用于签发的认证信息，接下来默认插入一条认证的 SQL 语句（也可以做一个接口），便于认证接口的后续使用，插入的 SQL 语句如下：

```sql
INSERT INTO `ch02`.`blog_auth`(`id`, `app_key`, `app_secret`, `created_on`, `created_by`, `modified_on`, `modified_by`, `deleted_on`, `is_del`) VALUES (1, 'admin', 'go-learning', 0, 'test', 0, '', 0, 0);
```

该条语句的主要作用是新增了一条`app_key` 为 admin以及 `app_secret `为 `go-learning`的数据。

#### b. 新建 model 对象

接下来打开项目的 `internal/model` 目录下的` auth.go` 文件，写入对应刚刚新增的` blog_auth `表的数据模型，如下：

```go
package model

type Auth struct {
   *Model
   AppKey    string `json:"app_key"`
   AppSecret string `json:"app_secret"`
}

func (a Auth) TableName() string{
   return "blog_auth"
}
```

#### c. 初始化配置

接下来需要针对 JWT 的一些相关配置进行设置，修改项目的 `configs/config.yaml` 配置文件，写入新的配置项，如下：

```yaml
# JWT 初始化配置
JWT:
  Secret: admin
  Issuer: blog-service
  Expire: 7200
```

然后对 JWT 的配置进行初始化操作，修改项目的启动文件` main.go`，修改其 `setupSetting` 方法，如下：

```go
func setupSetting() error {
	...
	err = settings.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		return err
	}

	global.JWTSetting.Expire *= time.Second
	...
}
```

在上述配置中，设置了 JWT 令牌的 Secret（密钥）为 `admin`，签发者（Issuer）是 `blog-service`，有效时间（Expire）为 7200 秒，这里需要注意的是 Secret 千万不要暴露给外部，只能有服务端知道，否则是可以解密出来的，非常危险。

### 5. 处理 JWT 令牌

虽然 `jwt-go `库能够帮助开发者快捷的处理 JWT 令牌相关的行为，但是还是需要根据项目特性对其进行设计的，简单来讲，就是组合其提供的 API，设计鉴权场景。

在 `pkg/app` 并创建` jwt.go` 文件，写入第一部分的代码：

```go
package app

import (
   "demo/ch02/global"
   "github.com/dgrijalva/jwt-go"
)

type Claims struct {
   AppKey    string `json:"app_key"`
   AppSecret string `json:"app_secret"`
   jwt.StandardClaims
}

func GetJWTSecret() []byte {
   return []byte(global.JWTSetting.Secret)
}
```

这块主要涉及 JWT 的一些基本属性，第一个是` GetJWTSecret` 方法，用于获取该项目的 JWT Secret，目前是直接使用配置所配置的 Secret，第二个是 Claims 结构体，分为两大块，第一块是项目嵌入的 `AppKey` 和 `AppSecret`，用于自定义的认证信息，第二块是 `jwt.StandardClaims` 结构体，它是` jwt-go` 库中预定义的，也是 JWT 的规范，其涉及字段如下：

```go
// Structured version of Claims Section, as referenced at
// https://tools.ietf.org/html/rfc7519#section-4.1
// See examples for how to use this with your own claim types
type StandardClaims struct {
   Audience  string `json:"aud,omitempty"`
   ExpiresAt int64  `json:"exp,omitempty"`
   Id        string `json:"jti,omitempty"`
   IssuedAt  int64  `json:"iat,omitempty"`
   Issuer    string `json:"iss,omitempty"`
   NotBefore int64  `json:"nbf,omitempty"`
   Subject   string `json:"sub,omitempty"`
}
```

它对应的其实是本章节中 Payload 的相关字段，这些字段都是非强制性但官方建议使用的预定义权利要求，能够提供一组有用的，可互操作的约定。

接下来在 `jwt.go`中写入第二部分代码。

```go
// 生成 JWT
func GenerateToken(appKey, appSecret string) (string, error) {
   nowTime := time.Now()
   expireTime := nowTime.Add(global.JWTSetting.Expire)
   claims := Claims{
      AppKey:    util.EncodeMD5(appKey),
      AppSecret: util.EncodeMD5(appSecret),
      StandardClaims: jwt.StandardClaims{
         ExpiresAt: expireTime.Unix(),
         Issuer:    global.JWTSetting.Issuer,
      },
   }
   // 根据 Claims 结构体创建 Token 实例，jwt.NewWithClaims() 包含两个形参
   // SigningMethod，其包含 SigningMethodHS256、SigningMethodHS384、SigningMethodHS512 三种 crypto.Hash 加密算法的方案
   // 第二个参数为 Claims 主要用于传递用户所预定义的一些权限要求，方便后续的加密、校验
   tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
   // SignedString() 生成签名后的 token 字符串
   token, err := tokenClaims.SignedString(GetJWTSecret())
   return token, err
}
```

在 `GenerateToken` 方法中，它承担了整个流程中比较重要的职责，也就是生成 JWT Token 的行为，主体的函数流程逻辑是根据客户端传入的` AppKey `和 `AppSecret `以及在项目配置中所设置的签发者（Issuer）和过期时间（`ExpiresAt`），根据指定的算法生成签名后的 Token。这其中涉及两个的内部方法，如下：

- `jwt.NewWithClaims`：根据 Claims 结构体创建 Token 实例，它一共包含两个形参，第一个参数是 `SigningMethod`，其包含 SigningMethodHS256、SigningMethodHS384、SigningMethodHS512 三种 `crypto.Hash `加密算法的方案。第二个参数是 Claims，主要是用于传递用户所预定义的一些权利要求，便于后续的加密、校验等行为。
- `tokenClaims.SignedString`：生成签名字符串，根据所传入 Secret 不同，进行签名并返回标准的 Token。

接下来继续在` jwt.go` 文件中写入第三部分代码，如下：

```go
// 解析和校验 Token
func ParseToken(token string) (*Claims, error) {
   // jwt.ParseWithClaims() 用于解析鉴权的声明，最终返回 *Token
   tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
      return GetJWTSecret(), nil
   })
   if err != nil {
      return nil, err
   }
   if tokenClaims != nil {
      // Token.Valid 当转换与核实 token 时填充该值
      if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
         return claims, nilgo
      }
   }
   return nil, err
}
```

```go
// A JWT Token.  Different fields will be used depending on whether you're
// creating or parsing/verifying a token.
type Token struct {
   Raw       string                 // The raw token.  Populated when you Parse a token
   Method    SigningMethod          // The signing method used or to be used
   Header    map[string]interface{} // The first segment of the token
   Claims    Claims                 // The second segment of the token
   Signature string                 // The third segment of the token.  Populated when you Parse a tokengo
   Valid     bool                   // Is the token valid?  Populated when you Parse/Verify a token
}
```

在 `ParseToken `方法中，它主要的功能是解析和校验 Token，承担着与 `GenerateToken` 相对的功能，其函数流程主要是解析传入的 Token，然后根据 Claims 的相关属性要求进行校验。这其中涉及两个的内部方法，如下：

- `ParseWithClaims`：用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回 `*Token`。
- `Valid`：验证基于时间的声明，例如：过期时间（`ExpiresAt`）、签发者（Issuer）、生效时间（Not Before），需要注意的是，如果在令牌中没有任何声明，仍然会被认为是有效的。

至此完成了 JWT 令牌的生成、解析、校验的方法编写，在后续的应用中间件中对其进行调用，使其能够在应用程序中将一整套的动作给串联起来。

### 6. 获取 JWT 令牌

#### a. 新建 model 方法

修改`internal/model` 下的 `auth.go `文件。

```go
// 通过传入的 app_key 和 app_secret 获取认证信息
func (a Auth) Get(db *gorm.DB) (Auth, error) {
   var auth Auth
   db = db.Where("app_key = ? AND app_secret = ? AND is_del = ?", a.AppKey, a.AppSecret, 0)
   err := db.First(&auth).Error
   if err != nil && err != gorm.ErrRecordNotFound {
      return auth, err
   }
   return auth, nil
}
```

上述方法主要是用于服务端在获取客户端所传入的 app_key 和 app_secret 后，根据所传入的认证信息进行获取，以此判别是否真的存在这一条数据。

#### b. 新建 dao 方法

在 `internal/dao` 下新建`auth.go `文件，并编写针对获取认证信息的方法。

```go
package dao

import "demo/ch02/internal/model"

func (d *Dao) GetAuth(appKey, appSecret string) (model.Auth, error) {
   auth := model.Auth{AppKey: appKey, AppSecret: appSecret}
   return auth.Get(d.engine)
}
```

#### c. 新建 service 方法

在 `internal/service` 下新建`auth.go `文件，针对一些相应的基本逻辑进行处理。

```go
package service

import "errors"

type AuthRequest struct {
   AppKey    string `form:"app_key" binding:"required"`
   AppSecret string `form:"app_secret" binding:"required"`
}

func (svc *Service) CheckAuth(param *AuthRequest) error {
   auth, err := svc.dao.GetAuth(param.AppKey, param.AppSecret)
   if err != nil {
      return err
   }
   if auth.ID > 0 {
      return nil
   }
   return errors.New("auth info does not exist.")
}
```

在上述代码中，声明了 `AuthRequest` 结构体用于接口入参的校验，`AppKey `和 `AppSecret` 都设置为了必填项，在 `CheckAuth` 方法中，使用客户端所传入的认证信息作为筛选条件获取数据行，以此根据是否取到认证信息 ID 来进行是否存在的判定。

#### d. 新增路由方法

在 `internal/routers/api` 在新建`auth.go `文件。

```go
package api

import (
   "demo/ch02/global"
   "demo/ch02/internal/service"
   "demo/ch02/pkg/app"
   "demo/ch02/pkg/errcode"
   "github.com/gin-gonic/gin"
)

func GetAuth(c *gin.Context) {
   // 入参绑定与校验
   param := service.AuthRequest{}
   response := app.NewResponse(c)
   valid, errs := app.BindAndValid(c, &param)
   if !valid {
      global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
      response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
      return
   }

   // 判断认证信息
   svc := service.New(c.Request.Context())
   err := svc.CheckAuth(&param)
   if err != nil {
      global.Logger.Errorf(c, "svc.CheckAuth err: %v", err)
      response.ToErrorResponse(errcode.UnauthorizedAuthNotExist)
      return
   }

   // 生成 token
   token, err := app.GenerateToken(param.AppKey, param.AppSecret)
   if err != nil {
      global.Logger.Errorf(c, "app.GenerateToken err: %v", err)
      response.ToErrorResponse(errcode.UnauthorizedTokenGenerate)
      return
   }

   // 返回生成的 token
   response.ToResponse(gin.H{
      "token": token,
   })
}
```

这块的逻辑主要是校验及获取入参后，绑定并获取到的 `app_key` 和 `app_secrect `进行数据库查询，检查认证信息是否存在，若存在则进行 Token 的生成并返回。

接下来修改 `internal/routers` 的 `router.go `文件，新增`auth`路由。至此，就完成了获取 Token 的整套流程。

```go
package routers

import (
   ...
)

func NewRouter() *gin.Engine {
   ...
   // 新增 auth 相关路由
   r.POST("/auth", api.GetAuth)
   ...
}
```

#### e. 接口验证

![image-20220515223328408](https://raw.githubusercontent.com/tonshz/test/master/img/202205152233457.png)

![image-20220515223350469](https://raw.githubusercontent.com/tonshz/test/master/img/202205152233508.png)

### 7. 处理应用中间件

#### a. 编写 JWT 中间件

在完成了获取 Token 的接口后，能获取 Token 了，但是对于其它的业务接口，它还没产生任何作用。涉及特定类别的接口统一处理，那必然是选择应用中间件的方式，接下来在`internal/middleware` 下新建 `jwt.go `文件，写入如下代码：

```go
package middleware

import (
   "demo/ch02/pkg/app"
   "demo/ch02/pkg/errcode"
   "github.com/dgrijalva/jwt-go"
   "github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
   return func(c *gin.Context) {
      var (
         token string
         ecode = errcode.Success
      )
      // 获取 token
      if s, exist := c.GetQuery("token"); exist {
         token = s
      } else {
         token = c.GetHeader("token")
      }
      if token == "" {
         ecode = errcode.InvalidParams
      } else {
         // ParseToken() 解析 token
         _, err := app.ParseToken(token)
         if err != nil {
            switch err.(*jwt.ValidationError).Errors {
            case jwt.ValidationErrorExpired:
               ecode = errcode.UnauthorizedTokenTimeout
            default:
               ecode = errcode.UnauthorizedTokenError
            }
         }
      }
      if ecode != errcode.Success {
         response := app.NewResponse(c)
         response.ToErrorResponse(ecode)
         /*
            Abort() 可防止调用挂起的处理程序
            请注意，这不会停止当前处理程序
            假设您有一个授权中间件来验证当前请求是否已获得授权
            如果授权失败（例如：密码不匹配）
            请调用 Abort 以确保不调用此请求的其余处理程序
         */
         c.Abort()
         return
      }
      c.Next()
   }
}
```

在上述代码中，通过` GetHeader `方法从 Header 中获取 token 参数，并调用` ParseToken` 对其进行解析，再根据返回的错误类型进行断言判定。

#### b. 接入 JWT 中间件

在完成了 JWT 的中间件编写后，需要将其接入到应用流程中，但是需要注意的是，并非所有的接口都需要用到 JWT 中间件，因此需要利用 gin 中的分组路由的概念，只针对 apiv1 的路由分组进行 JWT 中间件的引用，也就是只有 apiv1 路由分组里的路由方法会受此中间件的约束，修改`internal/routers`下的`router.go`。

```go
package routers

import (
   ...
)

func NewRouter() *gin.Engine {
   ...
   // 使用路由组设置访问路由的统一前缀 e.g. /api/v1
   // 此处定义了一个路由组 /api/v1
   apiv1 := r.Group("/api/v1")
   // apiv1 路由分组引入 JWT 中间件
   apiv1.Use(middleware.JWT())
   // 上面花括号是代表中间的语句属于一个空间内，不受外界干扰，可去掉
   {
      ...
   }
   return r
}
```

#### c. 验证接口

##### 没有获取 Token

![image-20220515225527386](https://raw.githubusercontent.com/tonshz/test/master/img/202205152255448.png)

##### Token 错误

![image-20220515230152513](https://raw.githubusercontent.com/tonshz/test/master/img/202205152301590.png)

##### Token 超时

![image-20220515230342431](https://raw.githubusercontent.com/tonshz/test/master/img/202205152303507.png)

### 8. 小结

通过本章节的学习，可以得知 JWT 令牌的内容是非严格加密的，大体上只是进行` base64UrlEncode `的处理，也就是对 JWT 令牌机制有一定了解的人可以进行反向解密，可以编写 base64 的解码代码，也可通过`jwt.io`网站直接进行解码。首先先调用 `/auth` 接口获取一个全新 token，例如：

```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBfa2V5IjoiMjEyMzJmMjk3YTU3YTVhNzQzODk0YTBlNGE4MDFmYzMiLCJhcHBfc2VjcmV0IjoiNjgyYjU1NGRiYmQ5NGE3NDQ0NDU5NDJlOGMyZDk3Y2YiLCJleHAiOjE2NTI2MjcyMTIsImlzcyI6ImJsb2ctc2VydmljZSJ9.siQ-JLv3PZGUtn5OvLGzTOTV69PhkHCrmn1zfwb0dKE"
}
```

接下来针对新获取的 Token 值，只需要手动复制中间那一段（也就是 Payload），然后编写一个测试 Demo 来进行 base64 的解码，Demo 代码如下：

```go
func main() {
    payload, _ := base64.StdEncoding.DecodeString("eyJhcHBfa....DM5MTcsImlzcyI6ImJsb2ctc2VydmljZSJ9")
    fmt.Println(string(payload))
}
```

最终的输出结果如下：

```bash
{"app_key":"21232f297a57a5a743894a0e4a801fc3","app_secret":"682b554dbbd94a744445942e8c2d97cf","exp":1652627212,"iss":"blog-service"}
```

可以看到，假设有人拦截到 Token 后，是可以通过 Token 去解密并获取到 Payload 信息，也就是至少在 Payload 中不应该明文存储重要的信息，若非要存，就必须要进行不可逆加密，这样子才可以确保一定的安全性。

同时也可以发现，过期时间（`ExpiresAt`）是存储在 Payload 中的，也就是 JWT 令牌一旦签发，在没有做特殊逻辑的情况下，过期时间是不可以再度变更的，因此务必根据自己的实际项目情况进行设计和思考。

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

## 十一、应用配置问题

### 1. 配置读取

#### a. 在不同的工作区运行

在前面的部分，运行程序时，默认都是在项目的根目录下运行的，当在项目的其他目录下运行应用程序时会读取不到配置文件。切换到其他目录后，可以发现在其他目录下运行 `go run`时，会提示读取不到配置文件，初始化失败。

```bash
PS C:\Users\zyc\GolandProjects\demo\ch02> cd global
PS C:\Users\zyc\GolandProjects\demo\ch02\global> pwd

Path
----
C:\Users\zyc\GolandProjects\demo\ch02\global


PS C:\Users\zyc\GolandProjects\demo\ch02\global> go run ../main.go
2022/05/19 20:36:11 init.setupSetting err: Config File "config" Not Found i
n "[C:\\Users\\zyc\\GolandProjects\\demo\\ch02\\global\\configs]"
exit status 1
```

使用`go build`后再运行依旧会失败。模拟部署情况后再运行编译后的二进制文件依旧无法启动。**这是因为Go 语言中的编译与其他语言有差别，像配置文件这种非 `.go`文件的文件类型不会被打包进二进制文件中。**

#### b. 路径问题

在配置文件中填写的配置文件路径是相对路径，是相对于执行命令时的目录，因此在前文中应用程序在读取配置文件时读取不到。可以通过拼接可执行文件路径来实现读取配置文件的功能，首先要知道编译后的可执行文件的路径是什么。下面通过一个示例来获取当前可执行文件的路径。

```go
package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	log.Println(path)
}
```

##### go build 

```bash
$ go build . ; ./test
2022/05/19 20:47:38 C:\Users\zyc\GolandProjects\awesomeProject\test\test.ex
e
```

输出的结果与当前目录一致，即当前二进制文件的路径与期望的一致。

##### go run 

```bash
$ go run main.go
2022/05/19 20:44:43 C:\Users\zyc\AppData\Local\Temp\go-build2834558917\b001
\exe\main.exe
```

通过输出结果可以看到，在执行了 `go run`命令后，得到的是一个临时目录地址。如果操作系统是 CentOS，则输出的是`/tmp/go-build`目录，与预期不符，即输出的路径不是当前目录。

这是因为`go run`命令并不像`go build`命令那样可以直接编译输出当前目录，而是将其转换到临时目录下编译并执行，是一个相对临时的运行路径。

另外，通过在示例中打印变量 `os.Args[0]`可知，其传入的就是编译后的可执行文件的绝对路径，即 Go 对 `go run main.go`进行了一定的处理。

#### c. 思考

+ `go run`命令和`go build`命令的不同之处在于，一个是在临时目录下执行，另一个可手动在编译后的目录下执行，路径的处理方式不同。
+ 每次执行 `go run `命令后，生成的新二进制文件不一定在同一个地方。
+ 依赖相对路径读取的文件在没有遵守约定条件时，有可能会出现最终路径出错的问题。

#### d. 解决方案

目前可以确定两点：

+ Go 语言在编译时不会将配置文件这类第三方文件打包进二进制文件中。
+ 即受当前路径的影响，也会相对路径填写的不同而改变，并非时绝对可靠的。

##### 命令行参数

在 Go 语言中，可以通过 `flag`标准库来实现该功能。实现逻辑为：如果存在命令行参数，则优先使用命令行参数，否则使用配置文件中的配置参数。修改`main.go`，针对命令行参数的处理逻辑新增如下代码。

```go
...

var (
   port    string
   runMode string
   config  string
)

// Go 中的执行顺序: 全局变量初始化 =>init() => main()
// 在 main() 之前自动执行，进行初始化操作
func init() {
   // 添加命令行参数处理逻辑
   setupFlag()
   ...
}
...

func setupFlag() error {
   flag.StringVar(&port, "port", "", "启动端口")
   flag.StringVar(&runMode, "mode", "", "启动模式")
   flag.StringVar(&config, "config", "configs", "指定要使用的配置文件路径")
   flag.Parse()

   return nil
}
```

在上述代码中，可以通过标准库`flag`来读取命令行参数，然后根据其默认值判断配置文件是否存在。若存在，则对读取配置的路径进行变更。修改`okg/setting`目录下的`setting.go`。

```go
package setting

import "github.com/spf13/viper"

type Setting struct {
   vp *viper.Viper
}

// 用于初始化项目的基本配置
func NewSetting(configs ...string) (*Setting, error) {
   vp := viper.New()
   vp.SetConfigName("config") // 设置配置文件名称
   // 设置配置文件相对路径， viper 允许多个配置路径，可以不断调用 AddConfigPath()
   //vp.AddConfigPath("configs/")
   
   // 添加可变更配置文件路径
   for _, config := range configs {
      if config != "" {
         vp.AddConfigPath(config)
      }
   }
   vp.SetConfigType("yaml") // 设置配置文件类型

   err := vp.ReadInConfig()
   if err != nil {
      return nil, err
   }
   return &Setting{vp}, nil
}
```

接下来修改`main.go`中的`setupSetting`，对`ServerSetting`配置项进行覆写。如果存在，则覆盖原有的文件配置，使其优先级更高。

```go
...

func setupSetting() error {
	//settings, err := setting.NewSetting()
	// 如果存在则覆盖原有的文件配置
	settings, err := setting.NewSetting(strings.Split(config, ",")...)
	if err != nil {
		return err
	}
	...

	// 如果存在则覆盖原有的文件配置
	if port != "" {
		global.ServerSetting.HttpPort = port
	}
	if runMode != "" {
		global.ServerSetting.RunMode = runMode
	}
	return nil
}
...
```

最后。只需在启动时传入所期望的参数即可。

```bash
$ go run main.go -port 8001 -mode release -config configs/  
```

首先在`ch02`项目根目录下执行以下命令，即可在当前目录下生成`ch02.exe`可执行文件 （win10）。

```bash
$ go build .
```

例如在 `demo/ch02`的上级目录`demo`下运行以下命令即可，通过`config`命令设置`config.yml`配置文件的相**对路径。**

```bash
$ ./ch02/ch02 -port 8001 -config ch02/configs/ 
```

或者在`demo/ch02/global`目录下运行以下命令。

```bash
$ .././ch02 -config ../configs/
```

**也可使用绝对路径：**

```bash
$ .././ch02 -config C:/Users/zyc/GolandProjects/demo/ch02/configs/、
```

##### 系统环境变量

通过设置系统环境变量的方式，由程序去读取配置文件。同样是存在即优先的逻辑处理，`os.GetEnv("ENV")`。也可以将配置文件存放在系统自带的全局变量中，如 `$HOME/conf`或`/etc/conf`中，这样做就不需要重新定义一个新的系统环境变量。一般来说，会在程序内置一些系统环境变量的读取，其优先级低于命令行参数，但是高于文件配置。

##### 打包进二进制文件

可以将配置文件这种第三方文件打包进二进制文件中，这样就不需要过度关注这些第三方文件了，但这样做是有一定代价的，因此要注意使用的应用场景，即并非所有的项目都能这样操作。首先安装`go-bindata`库。

```bash
$ go get -u github.com/go-bindata/go-bindata/...
```

通过`go-bindata`库可以将数据文件转换为 Go 代码。例如，常见 的配置文件、资源文件（如 Swagger UI）等都可以打包进 Go 代码中，这样就可以“摆脱”静态资源文件了。接下来在项目根目录下执行生成命令。

```bash
$ go-bindata -o ./configs/config.go -pkg=configs ./configs/config.yml
```

执行这条命令后，会将`configs/config.yml`文件打包，并输出到`-o`选项指定的路径`configs/config.go`文件中，再通过设置的`-pkg`选项指定生成的`package name`为 `configs`，接下来只需要执行下述代码，就可以读取对应的文件内容了。

```go
b, _ := configs.Asset("configs/config.yml")
```

将第三方文件打包进二进制文件后，二进制文件必然增大，而且在常规方法下无法做到文件的热更新和监听，必须重启和重新打包后才能使用最新的内容。

```bash
$ .././ch02 # ch02 子目录路径执行 ch02.exe
$ ./ch02/ch02 # ch02 父目录路径执行 ch02.exe
```

##### 其他方案

当不使用文件配置时，可以直接使用集中式的配置中心等。

### 配置热更新

#### a. 开源库 fsnotify

既然要做配置热更新，那么首先要知道配置是什么时候修改的，因此需要对配置文件进行监听，以便得知配置文件的修改。开源库`fsnotify`为使用 Go 语言编写的跨平台文件系统监听事件库，常用于文件监听。

##### 安装开源库 fsnotify

```bash
$ go get -u golang.org/x/sys/...
$ go get -u github.com/fsnotify/fsnotify
```

`fsnotify`是基于`golang.org/x/sys`实现的，并非`syscall`标准库，因此在安装时需要更新其版本。

##### 案例

下面通过一个小案例，快速了解和实现文件的监听功能。

```go
package main

import (
   "github.com/fsnotify/fsnotify"
   "log"
)

func main() {
   watcher, _ := fsnotify.NewWatcher()
   defer watcher.Close()

   done := make(chan int)
   go func() {
      for {
         select {
         case event, ok := <-watcher.Events:
            if !ok {
               return
            }
            log.Println("event: ", event)
            if event.Op&fsnotify.Write == fsnotify.Write {
               log.Println("modified file: ", event.Name)
            }
         case err, ok := <-watcher.Errors:
            if !ok {
               return
            }
            log.Println("error: ", err)

         }
      }
   }()

   // 填写你要监听的目录或文件
   path := "C:\\Users\\zyc\\GolandProjects\\awesomeProject\\test\\p\\test.go"
   _ = watcher.Add(path)

   <-done
}
```

```bash
2022/05/19 22:27:12 event:  "C:\\Users\\zyc\\GolandProjects\\awesomeProject\\test\\p\\test.go": WRITE
```

上述代码对项目配置文件进行了监听，因此可以修改配置文件中的值，来查看控制台输出的变更事件。通过监听可以很便捷的知道文件做了哪些变更，可以通过对其进行二次封装，在它的上层实现一些变更动作来完成配置文件的热更新。

##### 如何做

viper 开源库能够很便捷的实现对文件的监听和热更新。修改`pkg/setting/section.go`文件针对重载应用配置项新增处理方法。

```go
...

var sections = make(map[string]interface{})

// 读取相应配置的配置方法
func (s *Setting) ReadSection(k string, v interface{}) error {
	// 将配置文件 按照 父节点读取到相应的struct中
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	// 针对重载应用配置项，新增处理方法
	if _, ok := sections[k]; !ok {
		sections[k] = v
	}
	return nil
}
```

首先修改`ReadSection`方法，增加读取 section 的存储记录，以便在重新加载配置的方法中进行二次处理。接下来新增`ReloadAllSection()`，重新读取配置。

```go
// 重新读取配置
func (s *Setting) ReloadAllSection() error {
	for k, v := range sections {
		err := s.ReadSection(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
```

最后修改`pkg/setting/setting.go`文件，新增文件热更新的监听和变更处理。

```go
// 用于初始化项目的基本配置
func NewSetting(configs ...string) (*Setting, error) {
	vp := viper.New()
	vp.SetConfigName("config") // 设置配置文件名称
	// 设置配置文件相对路径， viper 允许多个配置路径，可以不断调用 AddConfigPath()
	//vp.AddConfigPath("configs/")

	// 添加可变更配置文件路径
	for _, config := range configs {
		if config != "" {
			vp.AddConfigPath(config)
		}
	}
	vp.SetConfigType("yaml") // 设置配置文件类型

	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}
	s := &Setting{vp}
	s.WatchSettingChange()
	return s, nil
}

// 新增热更新的监听和变更处理
func (s *Setting) WatchSettingChange() {
	go func() {
		s.vp.WatchConfig()
        // 如果配置文件发生了改变就重新读取配置项
		s.vp.OnConfigChange(func(in fsnotify.Event) {
			_ = s.ReloadAllSection()
		})
	}()
}
```

在上述代码中，首先在`WatchSettingChange()`中起一个协程，再在里面通过`WatchConfig()`对文件配置进行监听，并在`OnConfigChange()`中调用刚刚编写的重载方法`ReloadAllSection()`来处理热更新文件监听事件回调，这样就可以实现一个文件配置的热更新了。

## 十二、编译程序应用

在编写完应用程序后，下一步就是编译应用程序。Go 语言的应用程序在编译后，只有一个二进制文件。在这个二进制文件中，有许多用法和参数需要深入了解，以便在后续部署时能提供更好的适配性和灵活性。

### 1. 编译简介

#### a. 子命令

Go 语言中有许多的子命令。

```bash
$ go help
Go is a tool for managing Go source code.

Usage:

        go <command> [arguments]

The commands are:

        bug         start a bug report
        build       compile packages and dependencies # 与应用编译有关
        clean       remove object files and cached files
        doc         show documentation for package or symbol
        env         print Go environment information
        fix         update packages to use new APIs
        fmt         gofmt (reformat) package sources
        generate    generate Go files by processing source
        get         add dependencies to current module and install them
        install     compile and install packages and dependencies # 与应用编译有关
        list        list packages or modules
        mod         module maintenance
        run         compile and run Go program # 与应用编译有关
        test        test packages
        tool        run specified go tool
        version     print Go version
        vet         report likely mistakes in packages

Use "go help <command>" for more information about a command.

Additional help topics:

        buildconstraint build constraints
        buildmode       build modes
        c               calling between Go and C
        cache           build and test caching
        environment     environment variables
        filetype        file types
        go.mod          the go.mod file
        gopath          GOPATH environment variable
        gopath-get      legacy GOPATH go get
        goproxy         module proxy protocol
        importpath      import path syntax
        modules         modules, module versions, and more
        module-get      module-aware go get
        module-auth     module authentication using go.sum
        packages        package lists and patterns
        private         configuration for downloading non-public code
        testflag        testing flags
        testfunc        testing functions
        vcs             controlling version control with GOVCS

Use "go help <topic>" for more information about that topic.
```

其中与应用编译相关的是`go run、go install、go build`三个子命令。

#### b. go run 命令

`go run [arguments]`语句的作用是编译并马上运行 Go 程序，它可以接受一个或多个文件参数。当与平常使用的编译命令不同的是，它只接受 main 包下的文件作为参数，如果不是 main 包下的文件，则会出现报错。

```bash
$ go run config.go
package command-line-arguments is not a main package
```

在执行 `go run` 命令后，所编译的二进制文件最终存放在一个临时目录中。可以通过`-n`或`-x`参数进行查看。这两个参数的作用是打印编译过程中的所有执行命令，`-n`参数不会继续执行编译后的二进制文件，而`-x`参数会继续执行编译后的二进制文件。

```go
package main

func main() {
	println("Go 语言编程之旅学习")
}
```

```bash
$ go run -x main.go
# 设置临时环境变量 WORK 创建编译依赖所需的临时目录 可使用 GOTMPDIR 来调整
WORK=C:\Users\zyc\AppData\Local\Temp\go-build3729718146
mkdir -p $WORK\b001\
cat >$WORK\b001\_gomod_.go << 'EOF' # internal
package main
import _ "unsafe"
//go:linkname __debug_modinfo__ runtime.modinfo
var __debug_modinfo__ = "0w\xaf\f\x92t\b\x02A\xe1\xc1\a\xe6\xd6\x18\xe6path\tcommand-line-arguments\nmod\ttest\t(devel)\t\n\xf92C1\x86\x18 r\x00\x82B\x10A\
x16\xd8\xf2"
EOF
cat >$WORK\b001\importcfg << 'EOF' # internal
# import config
packagefile runtime=C:\Users\zyc\sdk\go1.17.2\pkg\windows_amd64\runtime.a
EOF

# 编译和生成编译所需要的依赖
cd C:\Users\zyc\GolandProjects\awesomeProject\test
"C:\\Users\\zyc\\sdk\\go1.17.2\\pkg\\tool\\windows_amd64\\compile.exe" -o "$WORK\\b001\\_pkg_.a" -trimpath "$WORK\\b001=>" -p main -lang=go1.17 -complete -
buildid sXn-0dQ0z65O9G1afwrk/sXn-0dQ0z65O9G1afwrk -dwarf=false -goversion go1.17.2 -D _/C_/Users/zyc/GolandProjects/awesomeProject/test -importcfg "$WORK\\
b001\\importcfg" -pack -c=4 "C:\\Users\\zyc\\GolandProjects\\awesomeProject\\test\\main.go" "$WORK\\b001\\_gomod_.go"
"C:\\Users\\zyc\\sdk\\go1.17.2\\pkg\\tool\\windows_amd64\\buildid.exe" -w "$WORK\\b001\\_pkg_.a" # internal
cp "$WORK\\b001\\_pkg_.a" "C:\\Users\\zyc\\AppData\\Local\\go-build\\77\\7759fd12347d836ad858782e4e83048f0119c09e55da7844d65e4091369164b6-d" # internal
cat >$WORK\b001\importcfg.link << 'EOF' # internal
packagefile command-line-arguments=$WORK\b001\_pkg_.a
packagefile runtime=C:\Users\zyc\sdk\go1.17.2\pkg\windows_amd64\runtime.a
packagefile internal/abi=C:\Users\zyc\sdk\go1.17.2\pkg\windows_amd64\internal\abi.a
packagefile internal/bytealg=C:\Users\zyc\sdk\go1.17.2\pkg\windows_amd64\internal\bytealg.a
packagefile internal/cpu=C:\Users\zyc\sdk\go1.17.2\pkg\windows_amd64\internal\cpu.a
packagefile internal/goexperiment=C:\Users\zyc\sdk\go1.17.2\pkg\windows_amd64\internal\goexperiment.a
packagefile runtime/internal/atomic=C:\Users\zyc\sdk\go1.17.2\pkg\windows_amd64\runtime\internal\atomic.a
packagefile runtime/internal/math=C:\Users\zyc\sdk\go1.17.2\pkg\windows_amd64\runtime\internal\math.a
packagefile runtime/internal/sys=C:\Users\zyc\sdk\go1.17.2\pkg\windows_amd64\runtime\internal\sys.a
EOF

# 创建 exe 目录
mkdir -p $WORK\b001\exe\
cd .

# 生成可执行文件
"C:\\Users\\zyc\\sdk\\go1.17.2\\pkg\\tool\\windows_amd64\\link.exe" -o "$WORK\\b001\\exe\\main.exe" -importcfg "$WORK\\b001\\importcfg.link" -s -w -buildmo
de=pie -buildid=VNr8663EzDjvT1jo-57i/sXn-0dQ0z65O9G1afwrk/FQB3RmqQzJOddr9kzvfg/VNr8663EzDjvT1jo-57i -extld=gcc "$WORK\\b001\\_pkg_.a"
$WORK\b001\exe\main.exe
######################################## 以上部分 go run -n main.go 也会输出同样的
Go 语言编程之旅学习 # 使用 go run -n main.go 命令没有此处的程序输出
```

在上述输出中，编译器执行力绝大部分编译相关的工作。

+ 创建编译依赖所需的临时目录。Go 编译器会设置一个临时环境变量 WORK，用于在此工作去编译应用程序，执行编译后的二进制文件，其默认值为系统的临时文件目录路径。也可以通过设置 GOTMPDIR 来调整其执行目录。
+ 编译和生成编译所需要的依赖，该阶段将会编译和生成标准库中的依赖（如 `flag.a、log.a、net/http`等）、应用程序的外部依赖（如`gin-gonic/gin`等），以及应用程序自身的代码，然后生成、链接对应归档文件（`.a` 文件）和编译配置文件。
+ 创建并进入编译二进制文件所需的临时目录，创建 exe 目录。
+ 生成可执行文件，主要用到的是`link`工具，该工具读取依赖文件的 Go 归档文件或对象及其依赖项，最终将它们组合成可执行的二进制文件，涉及参数如下表所示。
+ 执行可执行文件，到先前指定的目录`$WORK/b001/exe/main`下执行生成的二进制文件。

| 参数名       | 格式              | 含义                                                         |
| ------------ | ----------------- | ------------------------------------------------------------ |
| `-o`         | `-o file`         | 将输出写入文件（在 Windows 上默认为 `a.out` 或 `a.out.exe`） |
| `-importcfg` | `-importcfg file` | 从文件中读取导入配置，文件中通常为`packagefile、packageshlib` |
| `-s`         | `-s`              | 省略符号表并调试信息                                         |
| `-w`         | `-w`              | 省略 DWARF 符号表                                            |
| `-buildmode` | `-buildmode mode` | 设置构建模式（默认为 exe）                                   |
| `-buildid`   | `-buildid id`     | 将 ID 记录为 Go 工具链的构建 ID                              |
| `-extld`     | `-extld linker`   | 设置外部链接器（默认为 clang 或 gcc）                        |

下面将各个步骤与编译的整体行为配套，核心步骤如下图所示。

![image-20220521132057124](https://raw.githubusercontent.com/tonshz/test/master/img/202205211320167.png)

另外，如果要查看对应的生成文件，则需要注意一点，在执行 `go run `命令后，除非设置了`-work`参数，否则会在应用程序结束时自动删除该目录下的相关临时文件（如前面代码中的`b001`）。

#### c. go build 命令

`go builf [-o output] [-i] [build flags] [packages]`语句的作用是编译指定的源文件、软件包及其依赖项，但它不会运行编译后的二进制文件。在`ch02`项目中执行 `go build`命令，会在该目录下生成一个与当前目录名一致的可执行的二进制文件（若为 Windows 系统，则会生成 exe 文件），此时即可直接执行`./ch02`命令，将整个博客后端的应用程序运行起来。

如果要指定所生成二进制文件的名称（win10 下需要添加`.exe`文件后缀名），则可以通过 `-o`参数进行调整。

```bash
$ go build -o blog-service.exe # 注意需要添加文件后缀名 .exe，否则 ./blog-service 无法运行
```

在 `go build`命令中还有许多其他常见的命令行参数（在`go run`命令中ye同样适用）。

`go run` 命令和`go build`命令之间存在区别，首先看一下`go build`命令的编译执行过程。

```bash
$ go build -x
WORK=C:\Users\zyc\AppData\Local\Temp\go-build2306362578
mkdir -p $WORK\b001\
cat >$WORK\b001\importcfg.link << 'EOF' # internal
packagefile demo/ch02=C:\Users\zyc\AppData\Local\go-build\39\394fcdc7023c3b620a9eec9cfd46c681e36cbd0027203c2b8f2009085395201f-d
...
EOF

mkdir -p $WORK\b001\exe\
cd .
"C:\\Users\\zyc\\sdk\\go1.17.2\\pkg\\tool\\windows_amd64\\link.exe" -o "$WORK\\b001\\exe\\a.out.exe" -importcfg "$WORK\\b001\\importcfg.link" -buildmode=pi
e -buildid=QkZMOQd2oFAv3j56SFF4/G3Chj_A4Cz9GWgObuQPV/9pYusRolKC4Q9ufH0Vhd/QkZMOQd2oFAv3j56SFF4 -extld=gcc "C:\\Users\\zyc\\AppData\\Local\\go-build\\39\\39
4fcdc7023c3b620a9eec9cfd46c681e36cbd0027203c2b8f2009085395201f-d"
"C:\\Users\\zyc\\sdk\\go1.17.2\\pkg\\tool\\windows_amd64\\buildid.exe" -w "$WORK\\b001\\exe\\a.out.exe" # internal

cp $WORK\b001\exe\a.out.exe ch02.exe
rm -r $WORK\b001\
```

从本质上讲，`go build`命令和`go run`命令的编译执行过程差不多，唯一不同的是，`go build`命令会生成并执行编译好的二进制文件，将其重命名为`ch02`（当前目录名），并立刻删除编译时生成的临时目录。而在归档文件（.a 文件）上，`go build`命令和`go run`命令的执行结果是一样的，都是对所需的源码文件进行编译。

#### d. go install 命令

`go install [-i] [build flags] [packages]`语句的作用是编译并安装源文件、软件包。实际上`go install、go build、go run`三者在功能上相差不大，最大的区别在于`go install`会见编译后的相关文件（如可执行的二进制文件、归档文件等）安装到指定的目录中。

为了查看`go install`命令的整个执行过程，首先初始化示例项目的 `Go modules`，然后再查看它的编译过程。

```bash
$ go install -x
...
mkdir -p C:\Users\zyc\go\bin\
cp $WORK\b001\exe\a.out.exe C:\Users\zyc\go\bin\awesomeProject.exe
rm -r $WORK\b001\
```

从代码中可以看到，`go install`命令在编译后，会将生成的二进制文件移到 bin 目录下，其文件名称为`Go modules`的项目名，而非目录名。

需要注意的是，但设置了环境变量`$GOBIN`时，会将生成的二进制文件移到`$GOBIN`下，如果禁用了`Go modules`（不建议），那么将安装到`$GOPATH/pkg/$GOOS_$GOARCH`下。

#### e. 常用参数

| 参数名 | 含义                                                         |
| ------ | ------------------------------------------------------------ |
| -x     | 打印编译过程中的所有执行命令，执行生成的二进制文件           |
| -n     | 打印编译过程中的所有执行命令，不执行生成的二进制文件·        |
| -a     | 强制重新编译所有涉及的依赖                                   |
| -o     | 指定生成的二进制文件名称                                     |
| -p     | 指定编译过程中可以并发运行程序的数量，默认值为可用的 CPU 数量（Go 语言默认是支持并发编译的） |
| -work  | 打印临时工作目录的完整路径，再退出时不删除该目录             |
| -race  | 启用数据竞争检测，目前仅支持`Linux/arm64、FreeeBSD/amd64、Darwin/adm64 和 Windows/amd64`平台 |

### 2. 交叉编译

#### a. 什么是交叉编译

交叉编译是指通过编译器在某个系统下编译另一个系统的可执行二进制文件，即目标计算架构的标识与当前运行环境的目标计算架构的标识不同，或者是所构建环境的目标操作系统的标识与当前运行环境的目标操作系统的标识不同。

#### b. 常用参数表

| 参数名      | 含义                                                         |
| ----------- | ------------------------------------------------------------ |
| CGO_ENABLED | 用于标识 CGO 工具是否可用，默认开启，可通过执行 `go env` 进行查看 |
| GOOS        | 用于标识程序构建环境的目标操作系统，如 Linux、Darwin、Windows |
| GOARCH      | 用于标识程序构建环境的目标计算架构，若不设置，则默认值与程序运行环境的目标计算架构一致，如 amd64、386 |
| GOHOSTOS    | 用于标识程序运行环境的目标操作系统                           |
| GOHOSTARCH  | 用于标识程序运行环境的目标计算架构                           |

#### c. 进行交叉编译

Go 编译器是默认支持交叉编译的，只需执行上述日到的一部分参数即可实现。

```bash
$ CGO_ENABLED=0 GOOS=linux go build -a -o blog-service .
```

通过上面的这条命令，关闭了 CGO ，并制定所构建的目标操作系统为 Linux 系统，强制重新编译所有依赖文件，最终输出名为 `	blog-service`的二进制文件到当前目录。

### 3. 编译缓存

在查看编译过程中，会经常看到`C:\Users\zyc\AppData\Local\go-build\...`这类的路径信息。在 Go 语言中，编译设计上是存在缓存机制的（从 Go 1.10 开始引入）,它一般存储在特定的目录下，可以通过以下命令查看。

```bash
$ go env GOCACHE
C:\Users\zyc\AppData\Local\go-build
```

编译缓存能节省大量的时间。

```bash
# 清理编译缓存
$ go clean -cache

# 第一次编译
$ time go build # win10 上无法使用 time 命令
go build 25.62s user 3.89s system 312% cpu 9.447 total

# 第二次编译
$ time go build
go build 0.90s user 0.41s system 215% cpu 0.611 total
```

两者的编译时间分别为 25.62s 和 0.90s，相差巨大，说明编译缓存能大幅提高后续编译的速度，并且目前还支持增量编译。在 Go 语言早期曾使用过时间戳的方式来界定编译器是否需要重新变更，但这存在问题，因为文件的修改时间变更了并不代表它的文件内容与先前的不同，可能存在多次反复修改，故基于时间戳是不正确的。

### 4. 编译文件大小

需要对在编译完应用程序后，编译出来的二进制文件的大小有一定的认知。一个简单的文件输出应用程序大概需要 1 MB以上，`blog-service`项目编译后的二进制文件约 38 M。

#### a. 为什么二进制文件这么大

在默认情况下，gc 工具链中的链接器创建静态链接的二进制文件。因此，所有 Go 二进制文件都包含 Go 运行时的信息、支持动态类型检查，以及在异常抛出时堆栈跟踪所必需的运行时类型信息（文件名/行号）。

在 Linux 上，使用 gcc 静态编译并静态链接的一个简单 C 语言编写的 hello world 程序约 750 KB，而一个等效的 Go 程序使用的 `fmt.Printf()`的大小为 MB 级别，但是它包含更强大的运行时支持，以及类型和调试信息，因此两者实际上并不完全等效。

#### b. 如何缩小二进制文件

最简单的方法是去掉 `DWARF`调试信息和符号表信息。

````bash
$ go build -ldflags="-w -s"
````

此时`blog-service`编译后的大小为 33 M 左右，当这是有代价的。

| 参数名 | 含义                | 副作用                                                      |
| ------ | ------------------- | ----------------------------------------------------------- |
| -w     | 去除 DWARF 调试信息 | 会导致异常（panic）抛出时，调用堆栈信息没有文件名、行号信息 |
| -s     | 去除符号表信息      | 无法使用 `gdb` 调试                                         |

还可以使用 `upx`工具（在 GitHub 上直接搜索 `upx`安装即可）对可执行文件进行压缩。

```bash
$ upx blog-service
```

最后会将目前 33 M 左右的 `blog-service`压缩到 10 M 左右，程序仍可正常运行。

### 5. 编译信息写入

在把应用程序打包成二进制文件后，在多环境部署下容易遇到一个问题，即无法知道这个编译好的二进制文件到底是什么框架版本的，应用版本号又是多少。还需借助其他部署工具进行反查，才能知道这个二进制文件的部署及编译信息。

![image-20220521143535531](https://raw.githubusercontent.com/tonshz/test/master/img/202205211435580.png)

为了解决这类问题，通常会将一些编译信息打包进二进制文件中，这样就可以通过指定命令输出所设置的信息，甚至将编译信息注册到对应的注册中心。

#### 使用 ldflags 设置编译信息

在 Go 语言中，联合使用 `go build`和 `-ldflags`命令，即可在构建时将动态信息设置到二进制文件中。通过 `-ldflags`命令的 `-X`参数可以在链接时将信息写入变量中，其格式如下。

```bash
$ go build -ldflags="-X 'package_path.variable_name=new_value'"
```

```go
package main

import "fmt"

var appName string
func main() {
	fmt.Printf("app_name: %s\n", appName)
}
```

```bash
$ go build -ldflags "-X 'main.appName=Go 编程之旅'"   
$ ./test
app_name: Go 编程之旅
```

至此，基本的使用已经很清楚了，回到项目 `ch02`下的启动文件`main.go`，在其中新增如下代码。

```go
...

var (
   ...
   isVersion    bool
   buildTime    string
   buildVersion string
   gitCommitID  string
)
...
// @title 博客系统
// @version 1.0
// @description Go 语言项目实战学习
// @termsOfService https://github.com/go-programming-tour-book
func main() {
   // 添加版本信息
   if isVersion {
      fmt.Printf("build_time: %s\n", buildTime)
      fmt.Printf("build_version: %s\n", buildVersion)
      fmt.Printf("git_commit_id: %s\n", gitCommitID)
      return
   }
   ...

}
...

func setupFlag() error {
   ...
   // 添加版本信息
   flag.BoolVar(&isVersion, "version", false, "编译信息")
   flag.Parse()

   return nil
}
```

执行下述编译命令，将编译时间、版本号和 Git Hash（前提是安装了 Git，并且这个应用的目录是一个 Git 仓库，否则将无法取到值）设置进去，命令如下（不清楚为什么不能设置日期格式）：

```bash
#  go build -ldflags "-X 'main.buildTime=`date +%Y.%m.%d.%H%M%S`' -X 'main.buildVersion=1.0.0'" 此命令执行无法获取日期及设置日期格式
$ go build -ldflags "-X 'main.buildTime=$(date)' -X 'main.buildVersion=1.0.0'"
# 查看编译后的二进制文件和版本信息
$ ./ch02 -version
build_time: 05/21/2022 15:14:31 +%Y.%m.%d.%H%M%S
build_version: 1.0.0
git_commit_id: # 此应用不为 Git 项目，故未设置
```

至此，就完成了将编译信息打包进二进制文件。在完整流程中，一般会提供程序中的对接，其余的编译、变量设置等工作，都由脚本进行调度和设置，达到一个相对自动化的部署流程。

### 6. 小结

简单介绍了 Go 语言编译的相关命令和知识点，了解了在应用不输钱，应用编译应该做哪些事，以及可能会遇到的问题。

+ 编译速度：了解到 Go 语言的编译器默认支持并发编译和编译缓存，能够明显提升编译效率。
+ 功能使用：Go 提供了多种运行方式，既可以简单快速的使用`go run`命令，也可以在部署时使用`go build`命令。如果仅仅想要库文件，也可以直接执行`go install`命令来获取。
+ 跨平台：Go 的编译器默认支持交叉编译，极大的提高了跨平台的能力。如果需要在一个没有编译环境的操作系统上使用 Go 语言编写程序，则可以在本机对该目标系统进行编译，再将编译好的程序部署过去，就可以使用了。
+ 编译后的二进制文件大小：比普通的 C 程序大，但是 Go 程序支持更强大的功能，并且 Go 程序的编译（不适用 CGO）默认使用了静态编译，也就是说，不需要依赖任何动态链接库。这样一来，就可以将编译好的二进制文件部署到任何适合的运行环境中。不过与动态链接库相比，静态编译出的二进制文件会更大一些，这是 Go 语言的一个权衡。因为多平台的适配性高于存储文件大小的意义，如有特殊需要也可通过`-ldflags="-w -s"`和 `upx`对二进制文件进行压缩，通常不建议。
+ 编译信息：可以检索许哟啊的基本信息打包进二进制文件中，以便后续的使用和排查。

## 十三、优雅重启和停止

在开发完成应用程序后，即可将其部署到测试、预发布或生产环境中。开发人员需要关注的是这个应用程序需要不断的进行更新和发布，即持续集成，在这个应用程序发布时，很可能某个用户正在使用这个应用，直接发布会导致用户的行为被中断。

### 1. 遇到的问题

为了避免这种情况的发生，希望在应用更新或发布时，现有正在处理既有连接的应用不要中断，要先处理完既有连接后再退出。而新发布的应用再部署上去后再开始接受性的请求并进行处理，这样即可避免原来正在处理的链接被中断的问题。

![image-20220521160127002](https://raw.githubusercontent.com/tonshz/test/master/img/202205211601072.png)

### 2. 解决方案

想要解决这个问题，目前最经典的方案就是通过信号量的方式来解决。

#### a. 信号定义（来自维基百科）

信号是 UNIX、类 UNIX，以及其他 POSIX 兼容的操作系统中进程间通信的一种有限制的方法。

它是一种异步的通知机制，用来提醒进程一个事件（硬件异常、程序执行异常、外部发出信号）已经发生。当一个信号发送给一个进程时，操作系统中断了进程正常的控制流程。此时，任何非原子操作都会被中断。如果进程定义了信号的处理函数，那么它将被执行，否则执行默认的处理函数。

#### b. 支持的信号

可以通过 `kill -l`查看系统所支持的所有信号。

```bash
$ kill -l
1) SIGHUP       2) SIGINT       3) SIGQUIT      4) SIGILL       5) SIGTRAP
 6) SIGABRT      7) SIGBUS       8) SIGFPE       9) SIGKILL     10) SIGUSR1
11) SIGSEGV     12) SIGUSR2     13) SIGPIPE     14) SIGALRM     15) SIGTERM
16) SIGSTKFLT   17) SIGCHLD     18) SIGCONT     19) SIGSTOP     20) SIGTSTP
21) SIGTTIN     22) SIGTTOU     23) SIGURG      24) SIGXCPU     25) SIGXFSZ
26) SIGVTALRM   27) SIGPROF     28) SIGWINCH    29) SIGIO       30) SIGPWR
31) SIGSYS      34) SIGRTMIN    35) SIGRTMIN+1  36) SIGRTMIN+2  37) SIGRTMIN+3
38) SIGRTMIN+4  39) SIGRTMIN+5  40) SIGRTMIN+6  41) SIGRTMIN+7  42) SIGRTMIN+8
43) SIGRTMIN+9  44) SIGRTMIN+10 45) SIGRTMIN+11 46) SIGRTMIN+12 47) SIGRTMIN+13
48) SIGRTMIN+14 49) SIGRTMIN+15 50) SIGRTMAX-14 51) SIGRTMAX-13 52) SIGRTMAX-12
53) SIGRTMAX-11 54) SIGRTMAX-10 55) SIGRTMAX-9  56) SIGRTMAX-8  57) SIGRTMAX-7
58) SIGRTMAX-6  59) SIGRTMAX-5  60) SIGRTMAX-4  61) SIGRTMAX-3  62) SIGRTMAX-2
63) SIGRTMAX-1  64) SIGRTMAX
```

### 3. 常用的快捷键

在终端执行特定的组合键可以使系统发送特定的信号给指定进程，并完成一系列动作，常用快捷键如下所示。

| 命令       | 信号    | 含义                   |
| ---------- | ------- | ---------------------- |
| `ctrl + c` | SIGINT  | 希望进程终端，进程结束 |
| `ctrl + z` | SIGTSTP | 任务中断，进程挂起     |
| `ctrl + \` | SIGQUIT | 进程结束和 `dump core` |

因此在使用组合键`ctrl + c`关闭服务端时，会发送希望进程结束的通知（发送 SIGINT 信号），如果没有进行额外处理，该进程会直接退出，最终导致正在访问的用户出现无法访问的情况。

而平时常用的`kill -9 pid`命令，会发送`SIGKILL`信号给进程，作用是强制中断进程。

### 4. 实现优雅重启和停止

#### a. 实现目的

+ 不关闭现有连接（正在运行中的程序）。
+ 新的进程启动并替代旧进程。
+ 新的进程接管新的连接。
+ 连接要随时响应用户的请求，当用户仍在请求就进程时要保持连接，新用户应请求新进程，不可以出现拒绝请求的情况。

#### b. 需要达到的流程

+ 替换可执行文件或修改配置文件。
+ 发送信号量 `SIGHUP`。
+ 拒绝新连接请求旧进程，保证正在处理的连接正常。
+ 启动新的子进程。
+ 新的子进程开始 Accept。
+ 系统将新的请求转交给新的子进程
+ 旧进程处理完所有旧连接后正常退出。

#### c. 实现

在了解实现优雅重启和停止所需的基本概念后，修改`ch02`目录下的`main.go`文件，对项目进行改造，使之能支持优雅重启和停止。

```go
...
// @title 博客系统
// @version 1.0
// @description Go 语言项目实战学习
// @termsOfService https://github.com/go-programming-tour-book
func main() {
   ...
   //// 调用 ListenAndServe() 监听
   //if err := s.ListenAndServe(); err != nil {
   // log.Fatalf("监听失败：%v", err)
   //}
   // 从此处开始修改 使项目支持优雅重启和停止
   go func() {
      err := s.ListenAndServe()
      if err != nil && err != http.ErrServerClosed {
         log.Fatalf("s.ListenAndServe err: %v", err)
      }
   }()

   // 等待中断信号
   quit := make(chan os.Signal)
   // 接受 syscall.SIGINT 和 syscall.SIGTERM 信号 两个都是终止信号
   /*
      signal.Notify()
      通知使包信号将传入信号中继到 quit。
      如果没有提供信号，所有传入的信号将被中继到 quit。
      否则，只有提供的信号会。
   */
   signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
   <-quit
   log.Println("Shut down server...")

   // 最大时间控制，用于通知该服务端它有 5s 的时间来处理原有请求
   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
   // cancel 类: type CancelFunc func()
   defer cancel()
   /*
      Shutdown 优雅地关闭服务器而不中断任何活动连接。
      关闭首先关闭所有打开的侦听器，然后关闭所有空闲连接，
      然后无限期地等待连接返回空闲状态，然后关闭。
      如果提供的上下文在关闭完成之前过期，则 Shutdown 返回上下文的错误，
      否则返回关闭服务器的底层侦听器返回的任何错误。
   */
   if err := s.Shutdown(ctx); err != nil {
      log.Fatalf("Server forced to shutdown: ", err)
   }

   log.Println("Server exiting")
}
...
```

重新启动应用，并制造一个请求比较慢的接口来进行验证，修改`internal/routers/api/v1`目录下的`tag.go`文件，添加请求慢的接口。

```go
// 请求慢的接口
func (t Tag) Get(c *gin.Context) {
	fmt.Println("请求 Google...")
	_, err := ctxhttp.Get(c.Request.Context(), http.DefaultClient, "https://www.google.com/")
	if err != nil {
		log.Fatalf("ctxhttp.Get err: %v", err)
	}
}
```

```bash
$ go run main.go
...
请求 Google...
2022/05/21 17:53:11 Shut down server...
2022/05/21 17:53:28 ctxhttp.Get err: context deadline exceeded
```

![image-20220521180400647](https://raw.githubusercontent.com/tonshz/test/master/img/202205211804813.png)

可以看到，在终端使用组合键`ctrl + c`后，回想应用发送一个`SIGINT`信号，并且会被应用成功捕获到，此时该应用开始停止对外接收新的请求，在原有的请求执行完毕后（可通过输出的 SQL 日志观察到），最终退出旧进程。如果是在一个完整的部署流程中，那么此时就已经完成了交替。

另外需要注意的是，如果没有正在处理的旧请求，那么在按下组合键`ctrl + c`后会直接退出，不需要等待。

### 5. 小结

在 Kubernetes 和 Docker 流行的今天，优雅重启和停止必须要实现的功能，因为 Kubernetes 在发布更新或退出时会向 Pod 发送 `SIGTERM`信号，告诉容器它很快就会被关闭，让应用程序停止接受新的请求，以确保应用在终止时是”干净“的。另外，在 Kubernetes 等待完成的时间，一般称为优雅终止宽限期，在限期到达后（默认 30s），如果仍在运行，那么会发送`SIGKILL`信号将其强制删除。针对这种情况，可以对`SIGKILL`信号进行功能定制。

---------------

## 参考

+ [参考文章](https://golang2.eddycjy.com/)

+ [GitHub代码](https://github.com/go-programming-tour-book)



















