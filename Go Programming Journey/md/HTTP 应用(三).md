# Go 语言编程之旅(二)：HTTP 应用(三) 

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

-----------------------

## 参考

+ [参考文章](https://golang2.eddycjy.com/)

+ [GitHub代码](https://github.com/go-programming-tour-book)









