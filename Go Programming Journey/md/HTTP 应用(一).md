# Go 语言编程之旅(二)：HTTP 应用(一) 

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

-----------------------

## 参考

+ [参考文章](https://golang2.eddycjy.com/)

+ [GitHub代码](https://github.com/go-programming-tour-book)



