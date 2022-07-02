## Go 项目案例-即时通信系统

### 架构图

![即时通信系统架构图](https://raw.githubusercontent.com/tonshz/test/master/img/%E5%8D%B3%E6%97%B6%E9%80%9A%E4%BF%A1%E7%B3%BB%E7%BB%9F%E6%9E%B6%E6%9E%84%E5%9B%BE.png?token=AJTN3F2WOQ2ONSBNNY5QUJTCEIIF4 "即时通信系统架构图")

------------------------------

### 构建基础 Server

#### server.go 文件

```go
// server端服务构建
package main

import (
   "fmt"
   "net"
)

type Server struct {
   Ip   string
   Port int
}

// NewServer 创建一个 server 的接口: 类似 Java 中的构造函数
func NewServer(ip string, port int) *Server {
   server := &Server{
      Ip:   ip,
      Port: port,
   }

   return server
}

func (this *Server) Handler(conn net.Conn) {
   //...当前链接的业务
   fmt.Println("链接建立成功")
}

// Start 启动服务器的接口
func (this *Server) Start() {
   // socket listen, fmt.Sprintf() 定义一个格式化的字符串
   listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
   if err != nil {
      fmt.Println("net.Listen err:", err)
      return
   }
   // close listen socket
   defer listener.Close()

   for {
      // accept
      conn, err := listener.Accept()
      if err != nil {
         fmt.Println("listener accept err:", err)
         continue
      }

      // do handler 处理连接业务,为了不阻塞当前 for 循环,使用 go程去处理当前任务
      go this.Handler(conn)
   }
}
```

#### main.go 文件

```go
// 当前进程的主入口
package main

func main() {
   // 在 win10 中安装 netcat64.exe ,执行命令 nc64 127.0.0.1 8888
   server := NewServer("127.0.0.1", 8888)
   server.Start()
}
```

#### 运行程序

```bash
# win10 下执行命令, -o 指定编译输出的名称，代替默认的包名。
go build -o server.exe main.go server.go # win10需为 server.exe ,Linux中为 server 即可
./server # 也可直接run main()
# 另起一个终端执行命令
nc64 127.0.0.1 8888 # win10 nc.exe 会报病毒, nc64.exe 需要放在 C:\Windows\System32 文件夹下
```

-----------------------------

### 用户上线功能

![用户上线及广播功能](https://raw.githubusercontent.com/tonshz/test/master/img/%E7%94%A8%E6%88%B7%E4%B8%8A%E7%BA%BF%E5%8F%8A%E5%B9%BF%E6%92%AD%E5%8A%9F%E8%83%BD.png?token=AJTN3F3ZXRSQR5NVP2HCOCDCEIIH4 "用户上线及广播功能")

#### server.go 文件

```go
// server端服务构建
package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户的列表
	OnlineMap map[string]*User
	// map 是全局需要加锁,Go 中有关同步的机制都在 sync 包中
	mapLock   sync.RWMutex

	// 消息广播的 channel
	Message chan string
}

// NewServer 创建一个 server 的接口: 类似 Java 中的构造函数
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
		// 需要进行初始化,否则为空
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// 监听Message广播消息channel的goroutine，一旦有消息就发送给全部的在线User
func (this *Server) ListenMessager() {
	for {
		// 不断尝试用 Message 中获取数据
		msg := <-this.Message

		//将msg发送给全部的在线User
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 广播消息的方法
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	//...当前链接的业务
	// fmt.Println("链接建立成功")
	
	user := NewUser(conn)

	// 用户上线,将用户加入到onlineMap中,进行广播
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	// 广播当前用户上线消息
	this.BroadCast(user, "login")

	// 当前handler阻塞
	select {}
}

// Start 启动服务器的接口
func (this *Server) Start() {
	// socket listen, fmt.Sprintf() 定义一个格式化的字符串
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// close listen socket
	defer listener.Close()

	// 启动监听 Message 的 goroutine
	go this.ListenMessager()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		// do handler 处理连接业务,为了不阻塞当前 for 循环,使用 go程去处理当前任务
		go this.Handler(conn)
	}
}
```

#### user.go 文件

```go
// 后端服务器表示在线用户的一个结构体封装
package main

import "net"

type User struct {
	// Name 与 Addr 默认都是地址
	Name string
	Addr string
    // 与当前用户绑定的 channel,每个用户都有,C 表示当前是否有数据,回写给对应的客户端
	C    chan string
    // 当前用户与客户端通信的连接(图中 write->client 部分)
	conn net.Conn
}

// 创建一个用户的API
func NewUser(conn net.Conn) *User {
	// 拿到当前客户端连接的地址
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	// 启动监听当前 user channel 消息的 goroutine
	go user.ListenMessage()
	return user
}

// 监听当前User channel的方法,一旦有消息,就直接发送给客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		// 将信息发送至当前用户的客户端
		this.conn.Write([]byte(msg + "\n"))
	}
}
```

#### 运行程序

```bash
# 终端1
nc64 127.0.0.1 8888 # 输入命令,先在当前终端窗口输入 chcp 65001 可解决 win 中文乱码
[127.0.0.1:14053]127.0.0.1:14053:login 
[127.0.0.1:14056]127.0.0.1:14056:login # 终端2输入命令后,终端1输出
[127.0.0.1:14057]127.0.0.1:14057:login
```

```bash
# 终端2
nc64 127.0.0.1 8888
[127.0.0.1:14056]127.0.0.1:14056:login
[127.0.0.1:14057]127.0.0.1:14057:login
```

```bash
# 终端3
nc64 127.0.0.1 8888
[127.0.0.1:14057]127.0.0.1:14057:login
```

---------------------------------

### 用户消息广播机制

#### server.go 

> 完善 handle() 处理业务方法，启动一个针对当前客户端的读 goroutine.

```go
// server端服务构建
package main

import (
   "fmt"
   "io"
   "net"
   "sync"
)

type Server struct {
   Ip   string
   Port int

   // 在线用户的列表
   OnlineMap map[string]*User
   // map 是全局需要加锁,Go 中有关同步的机制都在 sync 包中
   mapLock   sync.RWMutex

   // 消息广播的 channel
   Message chan string
}

// NewServer 创建一个 server 的接口: 类似 Java 中的构造函数
func NewServer(ip string, port int) *Server {
   server := &Server{
      Ip:   ip,
      Port: port,
      // 需要进行初始化,否则为空
      OnlineMap: make(map[string]*User),
      Message:   make(chan string),
   }

   return server
}

// 监听Message广播消息channel的goroutine，一旦有消息就发送给全部的在线User
func (this *Server) ListenMessager() {
   for {
      // 不断尝试用 Message 中获取数据
      msg := <-this.Message

      //将msg发送给全部的在线User
      this.mapLock.Lock()
      for _, cli := range this.OnlineMap {
         cli.C <- msg
      }
      this.mapLock.Unlock()
   }
}

// 广播消息的方法
func (this *Server) BroadCast(user *User, msg string) {
   sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
   this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
   //...当前链接的业务
   // fmt.Println("链接建立成功")

   user := NewUser(conn)

   // 用户上线,将用户加入到onlineMap中,进行广播
   this.mapLock.Lock()
   this.OnlineMap[user.Name] = user
   this.mapLock.Unlock()

   // 广播当前用户上线消息
   this.BroadCast(user, "已上线")

   // 接受客户端发送的消息
   go func() {
      buf := make([]byte, 4096)
      for {
         // Read() 允许从当前连接中读数据到 buf 中,返回已经成功的字节数,失败返回 err
         n, err := conn.Read(buf)

         // 如果读取数据为0,表示客户端合法关闭
         if n == 0 {
            this.BroadCast(user, "下线")
            return
         }

         // 每次读完数据后末尾都会有一个 EOF 标识,若 err != io.EOF 则说明存在一个非法操作
         if err != nil && err != io.EOF {
            fmt.Println("Conn Read err:", err)
            return
         }

         // buf 是一个字节形式的,需要转换成字符串形式,提取用户的消息,去除'\n',最后一个字符为'\n',故为 n-1
         msg := string(buf[:n-1])

         // 将得到的消息进行广播
         this.BroadCast(user, msg)
      }
   }()

   // 当前handler阻塞
   select {}
}

// Start 启动服务器的接口
func (this *Server) Start() {
   // socket listen, fmt.Sprintf() 定义一个格式化的字符串
   listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
   if err != nil {
      fmt.Println("net.Listen err:", err)
      return
   }
   // close listen socket
   defer listener.Close()

   // 启动监听 Message 的 goroutine
   go this.ListenMessager()

   for {
      // accept
      conn, err := listener.Accept()
      if err != nil {
         fmt.Println("listener accept err:", err)
         continue
      }

      // do handler 处理连接业务,为了不阻塞当前 for 循环,使用 go程去处理当前任务
      go this.Handler(conn)
   }
}
```

#### 运行程序

```bash
# 终端3
[127.0.0.1:14665]127.0.0.1:14665:已上线 # 广播上线消息
nihao # 发送消息
[127.0.0.1:14665]127.0.0.1:14665:nihao # 接受到自己发送的消息
```

```bash
# 终端1
[127.0.0.1:14663]127.0.0.1:14663:已上线
[127.0.0.1:14664]127.0.0.1:14664:已上线
[127.0.0.1:14665]127.0.0.1:14665:已上线
[127.0.0.1:14665]127.0.0.1:14665:nihao # 接收到终端3的消息
[127.0.0.1:14665]127.0.0.1:14665:下线 # 此时终端3已停止
```

```bash
# 终端2
[127.0.0.1:14664]127.0.0.1:14664:已上线
[127.0.0.1:14665]127.0.0.1:14665:已上线
[127.0.0.1:14665]127.0.0.1:14665:nihao
[127.0.0.1:14665]127.0.0.1:14665:下线 # 广播下线消息
```

-----------------------

### 用户业务封装

> 将在 server 中有关 user 的业务移至 user.go 中.
>
> server.go: 将之前 user 的业务进行替换.
>
> user.go: 新增与 server 的关联、新增 Online()、新增 Offline()、新增 DoMessage().

#### server.go

```go
// server端服务构建
package main

import (
   "fmt"
   "io"
   "net"
   "sync"
)

type Server struct {
   Ip   string
   Port int

   // 在线用户的列表
   OnlineMap map[string]*User
   // map 是全局需要加锁,Go 中有关同步的机制都在 sync 包中
   mapLock   sync.RWMutex

   // 消息广播的 channel
   Message chan string
}

// NewServer 创建一个 server 的接口: 类似 Java 中的构造函数
func NewServer(ip string, port int) *Server {
   server := &Server{
      Ip:   ip,
      Port: port,
      // 需要进行初始化,否则为空
      OnlineMap: make(map[string]*User),
      Message:   make(chan string),
   }

   return server
}

// 监听Message广播消息channel的goroutine，一旦有消息就发送给全部的在线User
func (this *Server) ListenMessager() {
   for {
      // 不断尝试用 Message 中获取数据
      msg := <-this.Message

      //将msg发送给全部的在线User
      this.mapLock.Lock()
      for _, cli := range this.OnlineMap {
         cli.C <- msg
      }
      this.mapLock.Unlock()
   }
}

// 广播消息的方法
func (this *Server) BroadCast(user *User, msg string) {
   sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
   this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
   //...当前链接的业务
   // fmt.Println("链接建立成功")

   // new user 时将当前的 server 作为参数传值进行初始化,将用户与 server 进行关联
   user := NewUser(conn, this)
   user.Online()

   // 有关用户上线的功能
   /*// 用户上线,将用户加入到onlineMap中,进行广播
   this.mapLock.Lock()
   this.OnlineMap[user.Name] = user
   this.mapLock.Unlock()

   // 广播当前用户上线消息
   this.BroadCast(user, "已上线")*/

   // 接受客户端发送的消息
   go func() {
      buf := make([]byte, 4096)
      for {
         // Read() 允许从当前连接中读数据到 buf 中,返回已经成功的字节数,失败返回 err
         n, err := conn.Read(buf)

         // 如果读取数据为0,表示客户端合法关闭
         if n == 0 {
            // 用户下线的功能
            //this.BroadCast(user, "下线")
            user.Offline()
            return
         }

         // 每次读完数据后末尾都会有一个 EOF 标识,若 err != io.EOF 则说明存在一个非法操作
         if err != nil && err != io.EOF {
            fmt.Println("Conn Read err:", err)
            return
         }

         // buf 是一个字节形式的,需要转换成字符串形式,提取用户的消息,去除'\n',最后一个字符为'\n',故为 n-1
         msg := string(buf[:n-1])

         // 用户正常处理消息的功能
         // 将得到的消息进行广播
         //this.BroadCast(user, msg)

         // 用户针对 msg 进行消息处理
         user.DoMessage(msg)
      }
   }()

   // 当前handler阻塞
   select {}
}

// Start 启动服务器的接口
func (this *Server) Start() {
   // socket listen, fmt.Sprintf() 定义一个格式化的字符串
   listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
   if err != nil {
      fmt.Println("net.Listen err:", err)
      return
   }
   // close listen socket
   defer listener.Close()

   // 启动监听 Message 的 goroutine
   go this.ListenMessager()

   for {
      // accept
      conn, err := listener.Accept()
      if err != nil {
         fmt.Println("listener accept err:", err)
         continue
      }

      // do handler 处理连接业务,为了不阻塞当前 for 循环,使用 go程去处理当前任务
      go this.Handler(conn)
   }
}
```

#### user.go

```go
// 后端服务器表示在线用户的一个结构体封装
package main

import "net"

type User struct {
   // Name 与 Addr 默认都是地址
   Name string
   Addr string
   // 与当前用户绑定的 channel,每个用户都有,C 表示当前是否有数据,回写给对应的客户端
   C    chan string
   // 当前用户与客户端通信的连接
   conn net.Conn

   // 当前用户属于的 server
   server *Server
}

// 创建一个用户的API,添加 Server 传值
func NewUser(conn net.Conn, server *Server) *User {
   // 拿到当前客户端连接的地址
   userAddr := conn.RemoteAddr().String()

   user := &User{
      Name: userAddr,
      Addr: userAddr,
      C:    make(chan string),
      conn: conn,
      server: server,
   }

   // 启动监听当前 user channel 消息的 goroutine
   go user.ListenMessage()
   return user
}

// 用户的上线业务
func (this *User) Online() {

   // 用户上线,将用户加入到 onlineMap 中
   this.server.mapLock.Lock()
   this.server.OnlineMap[this.Name] = this
   this.server.mapLock.Unlock()

   // 广播当前用户上线消息
   this.server.BroadCast(this, "已上线")
}

// 用户的下线业务
func (this *User) Offline() {

   // 用户下线,将用户从 onlineMap 中删除
   this.server.mapLock.Lock()
   delete(this.server.OnlineMap, this.Name)
   this.server.mapLock.Unlock()

   // 广播当前用户上线消息
   this.server.BroadCast(this, "下线")

}

// 用户处理消息的业务
func (this *User) DoMessage(msg string) {
   this.server.BroadCast(this, msg)
}


// 监听当前User channel的方法,一旦有消息,就直接发送给客户端
func (this *User) ListenMessage() {
   for {
      msg := <-this.C
      // 将信息发送至当前用户的客户端
      this.conn.Write([]byte(msg + "\n"))
   }
}
```

---------------------------------

### 在线用户查询

> 假定查询在线用户的消息格式为 who

#### user.go

```go
// 后端服务器表示在线用户的一个结构体封装
package main

import "net"

type User struct {
   // Name 与 Addr 默认都是地址
   Name string
   Addr string
   // 与当前用户绑定的 channel,每个用户都有,C 表示当前是否有数据,回写给对应的客户端
   C    chan string
   // 当前用户与客户端通信的连接
   conn net.Conn

   // 当前用户属于的 server
   server *Server
}

// 创建一个用户的API,添加 Server 传值
func NewUser(conn net.Conn, server *Server) *User {
   // 拿到当前客户端连接的地址
   userAddr := conn.RemoteAddr().String()

   user := &User{
      Name: userAddr,
      Addr: userAddr,
      C:    make(chan string),
      conn: conn,
      server: server,
   }

   // 启动监听当前 user channel 消息的 goroutine
   go user.ListenMessage()
   return user
}

// 用户的上线业务
func (this *User) Online() {

   // 用户上线,将用户加入到 onlineMap 中
   this.server.mapLock.Lock()
   this.server.OnlineMap[this.Name] = this
   this.server.mapLock.Unlock()

   // 广播当前用户上线消息
   this.server.BroadCast(this, "已上线")
}

// 用户的下线业务
func (this *User) Offline() {

   // 用户下线,将用户从 onlineMap 中删除
   this.server.mapLock.Lock()
   delete(this.server.OnlineMap, this.Name)
   this.server.mapLock.Unlock()

   // 广播当前用户上线消息
   this.server.BroadCast(this, "下线")

}

// 给当前 User 对应的客户端发送消息,不是广播
func (this *User) SendMsg(msg string) {
   this.conn.Write([]byte(msg))
}

// 用户处理消息的业务
func (this *User) DoMessage(msg string) {
   // 用户输入 who 就认为客户想查询当前在线用户
   if msg == "who" {
      // 查询当前在线用户
      this.server.mapLock.Lock()
      // 返回全部 OnlineMap
      for _, user := range this.server.OnlineMap {
         onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线...\n"
         this.SendMsg(onlineMsg)
      }
      this.server.mapLock.Unlock()

   } else {
      this.server.BroadCast(this, msg)
   }
}


// 监听当前 User channel 的方法,一旦有消息,就直接发送给客户端
func (this *User) ListenMessage() {
   for {
      msg := <-this.C
      // 将信息发送至当前用户的客户端
      this.conn.Write([]byte(msg + "\n"))
   }
}
```

#### 运行程序

```bash
# 终端1
[127.0.0.1:1321]127.0.0.1:1321:已上线
[127.0.0.1:1322]127.0.0.1:1322:已上线
[127.0.0.1:1323]127.0.0.1:1323:已上线
```

```bash
# 终端2
[127.0.0.1:1322]127.0.0.1:1322:已上线
[127.0.0.1:1323]127.0.0.1:1323:已上线
```

```bash
# 终端3
[127.0.0.1:1323]127.0.0.1:1323:已上线
who # 输入 who 查询当前在线用户
[127.0.0.1:1322]127.0.0.1:1322:在线...
[127.0.0.1:1323]127.0.0.1:1323:在线...
[127.0.0.1:1321]127.0.0.1:1321:在线...
```

-------------------

### 修改用户名

> 假定修改用户名的消息格式为 rename|张三

#### user.go

```go
// 后端服务器表示在线用户的一个结构体封装
package main

import (
   "net"
   "strings"
)

type User struct {
   // Name 与 Addr 默认都是地址
   Name string
   Addr string
   // 与当前用户绑定的 channel,每个用户都有,C 表示当前是否有数据,回写给对应的客户端
   C    chan string
   // 当前用户与客户端通信的连接
   conn net.Conn

   // 当前用户属于的 server
   server *Server
}

// 创建一个用户的API,添加 Server 传值
func NewUser(conn net.Conn, server *Server) *User {
   // 拿到当前客户端连接的地址
   userAddr := conn.RemoteAddr().String()

   user := &User{
      Name: userAddr,
      Addr: userAddr,
      C:    make(chan string),
      conn: conn,
      server: server,
   }

   // 启动监听当前 user channel 消息的 goroutine
   go user.ListenMessage()
   return user
}

// 用户的上线业务
func (this *User) Online() {

   // 用户上线,将用户加入到 onlineMap 中
   this.server.mapLock.Lock()
   this.server.OnlineMap[this.Name] = this
   this.server.mapLock.Unlock()

   // 广播当前用户上线消息
   this.server.BroadCast(this, "已上线")
}

// 用户的下线业务
func (this *User) Offline() {

   // 用户下线,将用户从 onlineMap 中删除
   this.server.mapLock.Lock()
   delete(this.server.OnlineMap, this.Name)
   this.server.mapLock.Unlock()

   // 广播当前用户上线消息
   this.server.BroadCast(this, "下线")

}

// 给当前 User 对应的客户端发送消息,不是广播
func (this *User) SendMsg(msg string) {
   this.conn.Write([]byte(msg))
}

// 用户处理消息的业务
func (this *User) DoMessage(msg string) {
   // 用户输入 who 就认为客户想查询当前在线用户
   if msg == "who" {
      // 查询当前在线用户
      this.server.mapLock.Lock()
      // 返回全部 OnlineMap
      for _, user := range this.server.OnlineMap {
         onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线...\n"
         this.SendMsg(onlineMsg)
      }
      this.server.mapLock.Unlock()

   } else if len(msg) > 7 && msg[:7] == "rename|" {
      // 消息格式: rename|张三
      newName := strings.Split(msg, "|")[1] // 获取张三

      // 判断 name 是否已存在,尝试去 map 中取值,ok 表示是否能取值
      _, ok := this.server.OnlineMap[newName] // _ 部分为 value
      if ok { // 当前 key 已存在
         this.SendMsg("当前用户名已被使用\n")
      } else {
         this.server.mapLock.Lock()
         // 删除旧名称,否则同样的用户会出现两次,原有的 key 和新的 key 指向同一个用户
         delete(this.server.OnlineMap, this.Name)
         this.server.OnlineMap[newName] = this
         this.server.mapLock.Unlock()

         this.Name = newName
         this.SendMsg("您已经更新用户名:" + this.Name + "\n")
      }

   }else {
      this.server.BroadCast(this, msg)
   }
}


// 监听当前 User channel 的方法,一旦有消息,就直接发送给客户端
func (this *User) ListenMessage() {
   for {
      msg := <-this.C
      // 将信息发送至当前用户的客户端
      this.conn.Write([]byte(msg + "\n"))
   }
}
```

#### 运行程序

```bash
# 终端1
[127.0.0.1:1605]127.0.0.1:1605:已上线
[127.0.0.1:1606]127.0.0.1:1606:已上线
[127.0.0.1:1607]127.0.0.1:1607:已上线
who
[127.0.0.1:1605]127.0.0.1:1605:在线...
[127.0.0.1:1606]tomy:在线... # 终端2用户名已被修改
[127.0.0.1:1607]127.0.0.1:1607:在线...
[127.0.0.1:1606]tomy:hello
[127.0.0.1:1606]tomy:下线
```

```bash
# 终端2
[127.0.0.1:1606]127.0.0.1:1606:已上线
[127.0.0.1:1607]127.0.0.1:1607:已上线
rename|tomy # 修改用户名
您已经更新用户名:tomy
hello
[127.0.0.1:1606]tomy:hello
```

```bash
# 终端3
[127.0.0.1:1607]127.0.0.1:1607:已上线
who
[127.0.0.1:1605]127.0.0.1:1605:在线...
[127.0.0.1:1606]127.0.0.1:1606:在线... # 终端2初始用户名
[127.0.0.1:1607]127.0.0.1:1607:在线...
[127.0.0.1:1606]tomy:hello
[127.0.0.1:1606]tomy:下线
```

---------------

### 超时强踢功能

> 用户任意消息表示用户活跃，长时间不发消息认为超时，需要强制关闭用户连接.
>
> 在Hander() goroutine 中，添加用户活跃 channel(isLive)，一旦有消息，就向该channel发送数据(true)，使其不再阻塞.
>
> 在Hander() goroutine 中，添加定时器器功能，超时则进行强踢.

#### server.go

```go
// server端服务构建
package main

import (
   "fmt"
   "io"
   "net"
   "sync"
   "time"
)

type Server struct {
   Ip   string
   Port int

   // 在线用户的列表
   OnlineMap map[string]*User
   // map 是全局需要加锁,Go 中有关同步的机制都在 sync 包中
   mapLock   sync.RWMutex

   // 消息广播的 channel
   Message chan string
}

// NewServer 创建一个 server 的接口: 类似 Java 中的构造函数
func NewServer(ip string, port int) *Server {
   server := &Server{
      Ip:   ip,
      Port: port,
      // 需要进行初始化,否则为空
      OnlineMap: make(map[string]*User),
      Message:   make(chan string),
   }

   return server
}

// 监听Message广播消息channel的goroutine，一旦有消息就发送给全部的在线User
func (this *Server) ListenMessager() {
   for {
      // 不断尝试用 Message 中获取数据
      msg := <-this.Message

      //将msg发送给全部的在线User
      this.mapLock.Lock()
      for _, cli := range this.OnlineMap {
         cli.C <- msg
      }
      this.mapLock.Unlock()
   }
}

// 广播消息的方法
func (this *Server) BroadCast(user *User, msg string) {
   sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
   this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
   //...当前链接的业务
   // fmt.Println("链接建立成功")

   // new user 时将当前的 server 作为参数传值进行初始化,将用户与 server 进行关联
   user := NewUser(conn, this)
   user.Online()

   // 有关用户上线的功能
   /*// 用户上线,将用户加入到onlineMap中,进行广播
   this.mapLock.Lock()
   this.OnlineMap[user.Name] = user
   this.mapLock.Unlock()

   // 广播当前用户上线消息
   this.BroadCast(user, "已上线")*/

   // 监听用户是否活跃的 channel
   isLive := make(chan bool)

   // 接受客户端发送的消息
   go func() {
      buf := make([]byte, 4096)
      for {
         // Read() 允许从当前连接中读数据到 buf 中,返回已经成功的字节数,失败返回 err
         n, err := conn.Read(buf)

         // 如果读取数据为0,表示客户端合法关闭
         if n == 0 {
            // 用户下线的功能
            //this.BroadCast(user, "下线")
            user.Offline()
            return
         }

         // 每次读完数据后末尾都会有一个 EOF 标识,若 err != io.EOF 则说明存在一个非法操作
         if err != nil && err != io.EOF {
            fmt.Println("Conn Read err:", err)
            return
         }

         // buf 是一个字节形式的,需要转换成字符串形式,提取用户的消息,去除'\n',最后一个字符为'\n',故为 n-1
         msg := string(buf[:n-1])

         // 用户正常处理消息的功能
         // 将得到的消息进行广播
         //this.BroadCast(user, msg)

         // 用户针对 msg 进行消息处理
         user.DoMessage(msg)

         // 用户的任意消息,代表当前用户是活跃的,此时 isLive 有可读数据
         isLive <- true
      }
   }()

   // 当前 handler 阻塞
   for {
      select {
      // 只要 isLive 触发,它之下的case都会执行,time.After() 也会执行
      case <-isLive:
         //当前用户是活跃的,应该重置定时器
         //不做任何事情,为了激活 select,更新下面的定时器
         // 下行代码会判断但是不会执行代码块,相当于定时器重置了
      case <-time.After(time.Second * 10): // time.After() 表示一个定时器,是一个 channel,10s之后触发,触发后有数据可读
         //已经超时
         //将当前的User强制的关闭
         user.SendMsg("您被踢了")

         //销毁用的资源
         close(user.C)

         //关闭连接
         conn.Close()

         //退出当前 Handler
         return // runtime.Goexit() 也可以实现类似效果
      }
   }
}

// Start 启动服务器的接口
func (this *Server) Start() {
   // socket listen, fmt.Sprintf() 定义一个格式化的字符串
   listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
   if err != nil {
      fmt.Println("net.Listen err:", err)
      return
   }
   // close listen socket
   defer listener.Close()

   // 启动监听 Message 的 goroutine
   go this.ListenMessager()

   for {
      // accept
      conn, err := listener.Accept()
      if err != nil {
         fmt.Println("listener accept err:", err)
         continue
      }

      // do handler 处理连接业务,为了不阻塞当前 for 循环,使用 go程去处理当前任务
      go this.Handler(conn)
   }
}
```

#### 运行程序

```bash
# 终端1
[127.0.0.1:3473]127.0.0.1:3473:已上线
[127.0.0.1:3474]127.0.0.1:3474:已上线
[127.0.0.1:3474]127.0.0.1:3474:nihao
[127.0.0.1:3474]127.0.0.1:3474:hello
您被踢了 # 10s不发送消息
```

```bash
# 终端2
[127.0.0.1:3474]127.0.0.1:3474:已上线
nihao
[127.0.0.1:3474]127.0.0.1:3474:nihao
hello
[127.0.0.1:3474]127.0.0.1:3474:hello
[127.0.0.1:3473]127.0.0.1:3473:下线 # 终端1被强踢
您被踢了
```

-----------------------

### 私聊功能

> 在DoMessage()方法中，添加对“to|张三|你好啊，我是...”消息格式指的处理，实现私聊功能

#### user.go

```go
// 后端服务器表示在线用户的一个结构体封装
package main

import (
   "net"
   "strings"
)

type User struct {
   // Name 与 Addr 默认都是地址
   Name string
   Addr string
   // 与当前用户绑定的 channel,每个用户都有,C 表示当前是否有数据,回写给对应的客户端
   C    chan string
   // 当前用户与客户端通信的连接
   conn net.Conn

   // 当前用户属于的 server
   server *Server
}

// 创建一个用户的API,添加 Server 传值
func NewUser(conn net.Conn, server *Server) *User {
   // 拿到当前客户端连接的地址
   userAddr := conn.RemoteAddr().String()

   user := &User{
      Name: userAddr,
      Addr: userAddr,
      C:    make(chan string),
      conn: conn,
      server: server,
   }

   // 启动监听当前 user channel 消息的 goroutine
   go user.ListenMessage()
   return user
}

// 用户的上线业务
func (this *User) Online() {

   // 用户上线,将用户加入到 onlineMap 中
   this.server.mapLock.Lock()
   this.server.OnlineMap[this.Name] = this
   this.server.mapLock.Unlock()

   // 广播当前用户上线消息
   this.server.BroadCast(this, "已上线")
}

// 用户的下线业务
func (this *User) Offline() {

   // 用户下线,将用户从 onlineMap 中删除
   this.server.mapLock.Lock()
   delete(this.server.OnlineMap, this.Name)
   this.server.mapLock.Unlock()

   // 广播当前用户上线消息
   this.server.BroadCast(this, "下线")

}

// 给当前 User 对应的客户端发送消息,不是广播
func (this *User) SendMsg(msg string) {
   this.conn.Write([]byte(msg))
}

// 用户处理消息的业务
func (this *User) DoMessage(msg string) {
   // 用户输入 who 就认为客户想查询当前在线用户
   if msg == "who" {
      // 查询当前在线用户
      this.server.mapLock.Lock()
      // 返回全部 OnlineMap
      for _, user := range this.server.OnlineMap {
         onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线...\n"
         this.SendMsg(onlineMsg)
      }
      this.server.mapLock.Unlock()

   } else if len(msg) > 7 && msg[:7] == "rename|" {
      // 消息格式: rename|张三
      newName := strings.Split(msg, "|")[1] // 获取张三

      // 判断 name 是否已存在,尝试去 map 中取值,ok 表示是否能取值
      _, ok := this.server.OnlineMap[newName] // _ 部分为 value
      if ok { // 当前 key 已存在
         this.SendMsg("当前用户名已被使用\n")
      } else {
         this.server.mapLock.Lock()
         // 删除旧名称,否则同样的用户会出现两次,原有的 key 和新的 key 指向同一个用户
         delete(this.server.OnlineMap, this.Name)
         this.server.OnlineMap[newName] = this
         this.server.mapLock.Unlock()

         this.Name = newName
         this.SendMsg("您已经更新用户名:" + this.Name + "\n")
      }

   }else if len(msg) > 4 && msg[:3] == "to|" { // msg[:3]表示0-3个字符
      //消息格式:  to|张三|消息内容

      //1 获取对方的用户名
      remoteName := strings.Split(msg, "|")[1]
      if remoteName == "" {
         this.SendMsg("消息格式不正确，请使用 \"to|张三|你好啊\"格式。\n")
         return
      }

      //2 根据用户名 得到对方User对象
      remoteUser, ok := this.server.OnlineMap[remoteName]
      if !ok {
         this.SendMsg("该用户名不存在\n")
         return
      }

      //3 获取消息内容，通过对方的User对象将消息内容发送过去
      content := strings.Split(msg, "|")[2]
      if content == "" {
         // 给当前用户发送消息
         this.SendMsg("无消息内容，请重发\n")
         return
      }
      // 给私聊对象发送消息
      remoteUser.SendMsg(this.Name + "对您说:" + content)
   } else {
      this.server.BroadCast(this, msg)
   }
}


// 监听当前 User channel 的方法,一旦有消息,就直接发送给客户端
func (this *User) ListenMessage() {
   for {
      msg := <-this.C
      // 将信息发送至当前用户的客户端
      this.conn.Write([]byte(msg + "\n"))
   }
}
```

#### 运行程序

```bash
# 终端1
[127.0.0.1:5370]127.0.0.1:5370:已上线
[127.0.0.1:5373]127.0.0.1:5373:已上线
rename|zhangsan
您已经更新用户名:zhangsan
who
[127.0.0.1:5370]zhangsan:在线...
[127.0.0.1:5373]lisi:在线...
lisi对您说:hello[127.0.0.1:5397]127.0.0.1:5397:已上线

```

```bash
# 终端2
[127.0.0.1:5373]127.0.0.1:5373:已上线
who
[127.0.0.1:5370]127.0.0.1:5370:在线...
[127.0.0.1:5373]127.0.0.1:5373:在线...
rename|lisi
您已经更新用户名:lisi
to|zhangsan|hello
[127.0.0.1:5397]127.0.0.1:5397:已上线
wangwu对您说:nihaoya
```

```bash
# 终端3
rename|wangwu
您已经更新用户名:wangwu
who
[127.0.0.1:5373]lisi:在线...
[127.0.0.1:5397]wangwu:在线...
[127.0.0.1:5370]zhangsan:在线...
to|lisi|nihaoya
```

-------------------------

### 客户端实现

#### 通过解析命令行来获取数据

> init() 初始化命令行参数
>
> main() 解析命令行

##### 命令行解析部分代码实现

```go
... // 以上部分代码省略
// 把变量值绑定到 flag 包中,变量初始化命令行格式: ./client -ip 127.0.0.1 -port 8888
func init() { // go 中每个文件都会有 init(),在 main() 之前执行
   // flag.StringVar(变量, 提示符, 变量默认值, 命令行说明)
   flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认是127.0.0.1)")
   // 根据变量类型选择 xxxVar
   flag.IntVar(&serverPort, "port", 8888, "设置服务器端口(默认是8888)")
}

func main() {
	// 在 main() 处最开始进行命令行解析, go 中可通过 flag 库解析命令行
	flag.Parse() // 判断输入命令行中是否有以上参数,有就会自动进行解析
	... // 以下部分代码省略
}
```

##### 运行程序

```bash
# 获取命令行说明
go build -o client.exe client.go
./client -h
Usage of C:\Users\zyc\go\src\awesomeProject\client\client.exe: # 命令行输出
  -ip string
        设置服务器IP地址(默认是127.0.0.1) (default "127.0.0.1")
  -port int
        设置服务器端口(默认是8888) (default 8888)
```

```bash
# 命令行解析
./client -ip 127.0.0.1 -port 8888
>>>>>链接服务器成功... # 命令行输出
```

#### 菜单显示

>client 新增 flag 属性，默认值为999
>
>新增 menu() 方法，获取用户选择的模式
>
>新增 Run() 主业务循环
>
>在 main() 中调用client.Run()

##### 菜单显示部分代码实现

```go
... 
// 客户端类
type Client struct {
   ServerIp   string
   ServerPort int
   Name       string
   // 连接成功的连接句柄
   conn       net.Conn
   flag       int // 当前 client 的模式
}
...
func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		// 设置 flag 默认值 999
		flag:       999,
	}
	...
}
...
// 获取用户选择的模式
func (client *Client) menu() bool {
	// flag 用来接收用户的输入
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	// 从命令行接收输入
	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>请输入合法范围内的数字<<<<")
		return false
	}

}
...
// 主业务循环
func (client *Client) Run() {
	// 根据 menu() 的结果选择不同的模式,不断循环判断当前 flag 的值
	for client.flag != 0 {
		// 输入合法为止,输入 1-3 之外数据的会循环调用 menu()
		for client.menu() != true {
		}

		// 根据不同的模式处理不同的业务
		switch client.flag {
		case 1:
			// 公聊模式
			client.PublicChat()
			break
		case 2:
			// 私聊模式
			client.PrivateChat()
			break
		case 3:
			// 更新用户名
			client.UpdateName()
			break
		}
	}
}
...
func main() {
	...

	//启动客户端的业务
	client.Run()
}
```

##### 运行程序

```bash
>>>>>链接服务器成功...
1.公聊模式
2.私聊模式
3.更新用户名
0.退出
```

#### 更新用户名

> 新增 UpdateName() 更新用户名，提示用户输入用户名，同时根据服务器的信息格式对数据进行封装
>
> 将方法加入到 Run 业务分支中，提供给用户进行选择
>
> 添加处理 server 回执消息的方法 DealResponse()，与 Run() 并发执行
>
> 在 main() 中开启一个 go 程，去运行 DealResponse() 流程，实现读写并行的效果

##### 更新用户名部分代码实现

```go
...
// 处理 server 回应的消息， 直接显示到标准输出即可
func (client *Client) DealResponse() {
   //一旦 client.conn 有数据，就直接 copy 到 stdout 标准输出上,永久阻塞监听
   io.Copy(os.Stdout, client.conn) // os.Stdout 标准化输出
   /*与上行代码相同的作用
   for{
      buf := make()
      client.conn.Read(buf)
      fmt.Println(buf)
   }*/
}
...
// 更新用户名
func (client *Client) UpdateName() bool {

	fmt.Println(">>>>请输入用户名:")
	// 使用 client.Name 接收用户名
	fmt.Scanln(&client.Name)

	// 根据服务器消息格式对数据进行封装
	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}

	return true
}

// 主业务循环
func (client *Client) Run() {
	// 根据 menu() 的结果选择不同的模式,不断循环判断当前 flag 的值
	for client.flag != 0 {
		...
		switch client.flag {
		...
		case 3:
			// 更新用户名
			client.UpdateName()
			break
		}
	}
}
...
func main() {
	...
	// 单独开启一个 goroutine 去处理 server 的回执消息
	go client.DealResponse()
	...
}
```

#### 公聊模式

>新增 PublicChat() 公聊模式业务
>
>加入Run的分支中

##### 公聊部分代码实现

```go
...
func (client *Client) PublicChat() {
   // 提示用户输入消息
   var chatMsg string // chatMsg 接受当前用户输入内容

   fmt.Println(">>>>请输入聊天内容，exit退出.")
   // 从命令行获取数据
   fmt.Scanln(&chatMsg)

   // 输入 exit 后会回到菜单,可选择其他模式
   for chatMsg != "exit" {
      // 发给服务器
      // 消息不为空则发送
      if len(chatMsg) != 0 {
         sendMsg := chatMsg + "\n"
         _, err := client.conn.Write([]byte(sendMsg))
         if err != nil { // eil 不为空
            fmt.Println("conn Write err:", err)
            break
         }
      }

      chatMsg = "" // 方便下次消息输入,将 chatMsg 置为空
      fmt.Println(">>>>请输入聊天内容，exit退出.")
      fmt.Scanln(&chatMsg)
   }
}
...
// 主业务循环
func (client *Client) Run() {
	// 根据 menu() 的结果选择不同的模式,不断循环判断当前 flag 的值
	for client.flag != 0 {
		...
		// 根据不同的模式处理不同的业务
		switch client.flag {
		case 1:
			// 公聊模式
			client.PublicChat()
			break
		...
		}
	}
}
...
```

#### 私聊模式

>查询当前有哪些用户在线，提示用户选择一个用户进行私聊

```go
...
// 将查询和私聊在服务器设定的消息格式放在私聊模式中
func (client *Client) PrivateChat() {
   var remoteName string
   var chatMsg string

   // 将当前在线用户输出在命令行
   client.SelectUsers()
   fmt.Println(">>>>请输入聊天对象[用户名], exit退出:")
   fmt.Scanln(&remoteName)

   for remoteName != "exit" {
      fmt.Println(">>>>请输入消息内容, exit退出:")
      fmt.Scanln(&chatMsg)

      // 私聊功能
      for chatMsg != "exit" {
         //消息不为空则发送
         if len(chatMsg) != 0 {
            // 根据服务器私聊模式中的消息格式对数据进行封装
            sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
            _, err := client.conn.Write([]byte(sendMsg))
            if err != nil {
               fmt.Println("conn Write err:", err)
               break
            }
         }

         chatMsg = ""
         fmt.Println(">>>>请输入消息内容, exit退出:")
         fmt.Scanln(&chatMsg)
      }

      // 使得用户可以重新选择用户私聊
      client.SelectUsers()
      fmt.Println(">>>>请输入聊天对象[用户名], exit退出:")
      fmt.Scanln(&remoteName)
   }
}
...
// 主业务循环
func (client *Client) Run() {
	// 根据 menu() 的结果选择不同的模式,不断循环判断当前 flag 的值
	for client.flag != 0 {
		...
		// 根据不同的模式处理不同的业务
		switch client.flag {
		...
		case 2:
			// 私聊模式
			client.PrivateChat()
			break
		...
		}
	}
}
...
```

#### client.go

```go
// 实现客户端模拟发送消息
package main

import (
   "flag"
   "fmt"
   "io"
   "net"
   "os"
)

// 客户端类
type Client struct {
   ServerIp   string
   ServerPort int
   Name       string
   // 连接成功的连接句柄
   conn       net.Conn
   flag       int // 当前 client 的模式
}

func NewClient(serverIp string, serverPort int) *Client {
   // 创建客户端对象
   client := &Client{
      ServerIp:   serverIp,
      ServerPort: serverPort,
      // 设置 flag 默认值 999
      flag:       999,
   }

   // 连接 server net.Dial()尝试连接服务器
   conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
   if err != nil {
      fmt.Println("net.Dial error:", err)
      return nil
   }

   client.conn = conn

   //返回对象
   return client
}

// 处理 server 回应的消息， 直接显示到标准输出即可
func (client *Client) DealResponse() {
   //一旦 client.conn 有数据，就直接 copy 到 stdout 标准输出上,永久阻塞监听
   io.Copy(os.Stdout, client.conn) // os.Stdout 标准化输出
   /*与上行代码相同的作用
   for{
      buf := make()
      client.conn.Read(buf)
      fmt.Println(buf)
   }*/
}

// 获取用户选择的模式
func (client *Client) menu() bool {
   // flag 用来接收用户的输入
   var flag int

   fmt.Println("1.公聊模式")
   fmt.Println("2.私聊模式")
   fmt.Println("3.更新用户名")
   fmt.Println("0.退出")

   // 从命令行接收输入
   fmt.Scanln(&flag)

   if flag >= 0 && flag <= 3 {
      client.flag = flag
      return true
   } else {
      fmt.Println(">>>>请输入合法范围内的数字<<<<")
      return false
   }

}

// 查询在线用户
func (client *Client) SelectUsers() {
   // 查询当前用户的消息格式
   sendMsg := "who\n"
   _, err := client.conn.Write([]byte(sendMsg))
   if err != nil {
      fmt.Println("conn Write err:", err)
      return
   }
}

// 将查询和私聊在服务器设定的消息格式放在私聊模式中
func (client *Client) PrivateChat() {
   var remoteName string
   var chatMsg string

   // 将当前在线用户输出在命令行
   client.SelectUsers()
   fmt.Println(">>>>请输入聊天对象[用户名], exit退出:")
   fmt.Scanln(&remoteName)

   for remoteName != "exit" {
      fmt.Println(">>>>请输入消息内容, exit退出:")
      fmt.Scanln(&chatMsg)

      // 私聊功能
      for chatMsg != "exit" {
         //消息不为空则发送
         if len(chatMsg) != 0 {
            // 根据服务器私聊模式中的消息格式对数据进行封装
            sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
            _, err := client.conn.Write([]byte(sendMsg))
            if err != nil {
               fmt.Println("conn Write err:", err)
               break
            }
         }

         chatMsg = ""
         fmt.Println(">>>>请输入消息内容, exit退出:")
         fmt.Scanln(&chatMsg)
      }

      // 使得用户可以重新选择用户私聊
      client.SelectUsers()
      fmt.Println(">>>>请输入聊天对象[用户名], exit退出:")
      fmt.Scanln(&remoteName)
   }
}

func (client *Client) PublicChat() {
   // 提示用户输入消息
   var chatMsg string // chatMsg 接受当前用户输入内容

   fmt.Println(">>>>请输入聊天内容，exit退出.")
   // 从命令行获取数据
   fmt.Scanln(&chatMsg)

   // 输入 exit 后会回到菜单,可选择其他模式
   for chatMsg != "exit" {
      // 发给服务器
      // 消息不为空则发送
      if len(chatMsg) != 0 {
         sendMsg := chatMsg + "\n"
         _, err := client.conn.Write([]byte(sendMsg))
         if err != nil { // eil 不为空
            fmt.Println("conn Write err:", err)
            break
         }
      }

      chatMsg = "" // 方便下次消息输入,将 chatMsg 置为空
      fmt.Println(">>>>请输入聊天内容，exit退出.")
      fmt.Scanln(&chatMsg)
   }
}

// 更新用户名
func (client *Client) UpdateName() bool {

   fmt.Println(">>>>请输入用户名:")
   // 使用 client.Name 接收用户名
   fmt.Scanln(&client.Name)

   // 根据服务器消息格式对数据进行封装
   sendMsg := "rename|" + client.Name + "\n"
   _, err := client.conn.Write([]byte(sendMsg))
   if err != nil {
      fmt.Println("conn.Write err:", err)
      return false
   }

   return true
}

// 主业务循环
func (client *Client) Run() {
   // 根据 menu() 的结果选择不同的模式,不断循环判断当前 flag 的值
   for client.flag != 0 {
      // 输入合法为止,输入 1-3 之外数据的会循环调用 menu()
      for client.menu() != true {
      }

      // 根据不同的模式处理不同的业务
      switch client.flag {
      case 1:
         // 公聊模式
         client.PublicChat()
         break
      case 2:
         // 私聊模式
         client.PrivateChat()
         break
      case 3:
         // 更新用户名
         client.UpdateName()
         break
      }
   }
}

// 全局变量
var serverIp string
var serverPort int

// 把变量值绑定到 flag 包中,变量初始化命令行格式: ./client -ip 127.0.0.1 -port 8888
func init() { // go 中每个文件都会有 init(),在 main() 之前执行
   // flag.StringVar(变量, 提示符, 变量默认值, 命令行说明)
   flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认是127.0.0.1)")
   // 根据变量类型选择 xxxVar
   flag.IntVar(&serverPort, "port", 8888, "设置服务器端口(默认是8888)")
}

func main() {
   // 在 main() 处最开始进行命令行解析, go 中可通过 flag 库解析命令行
   flag.Parse() // 判断输入命令行中是否有以上参数,有就会自动进行解析

   client := NewClient(serverIp, serverPort)
   if client == nil {
      fmt.Println(">>>>> 链接服务器失败...")
      return
   }

   // 单独开启一个 goroutine 去处理 server 的回执消息
   go client.DealResponse()

   fmt.Println(">>>>>链接服务器成功...")

   //启动客户端的业务
   client.Run()
}
```