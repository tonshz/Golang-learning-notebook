# Go 语言编程之旅(二)：HTTP 应用(二) 

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

----------------------

## 参考

+ [参考文章](https://golang2.eddycjy.com/)

+ [GitHub代码](https://github.com/go-programming-tour-book)


