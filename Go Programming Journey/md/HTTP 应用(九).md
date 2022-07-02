# Go 语言编程之旅(二)：HTTP 应用(九) 

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

-----------

## 参考

+ [参考文章](https://golang2.eddycjy.com/)

+ [GitHub代码](https://github.com/go-programming-tour-book)

