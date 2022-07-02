# Go 语言编程之旅(二)：HTTP 应用(十) 

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

---------------

## 参考

+ [参考文章](https://golang2.eddycjy.com/)

+ [GitHub代码](https://github.com/go-programming-tour-book)

