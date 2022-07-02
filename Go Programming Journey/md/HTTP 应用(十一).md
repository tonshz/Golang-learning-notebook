# Go 语言编程之旅(二)：HTTP 应用(十一) 

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

