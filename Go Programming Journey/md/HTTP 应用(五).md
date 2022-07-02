# Go 语言编程之旅(二)：HTTP 应用(五) 

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

--------------

## 参考

+ [参考文章](https://golang2.eddycjy.com/)

+ [GitHub代码](https://github.com/go-programming-tour-book)



