package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

type User struct {
	// ID 是用户唯一标识，通过 GenUserID 函数生成
	ID int
	// Addr 是用户的 IP 地址和端口
	Addr string
	// EnterAt 是用户进入时间
	EnterAt time.Time
	// MessageChannel 是当前用户发送消息的通道
	MessageChannel chan string
}

// 给用户发送的消息
type Message struct {
	OwnerID int
	Content string
}

func (u *User) String() string {
	return u.Addr + ", UID:" + strconv.Itoa(u.ID) + ", Enter At:" +
		u.EnterAt.Format("2006-01-02 15:04:05+8000")
}

var (
	// 新用户到来，通过该 channel 进行登记
	enteringChannel = make(chan *User)
	// 用户离开，通过该 channel 进行登记
	leavingChannel = make(chan *User)
	// 广播专用的用户普通消息 channel，缓冲是尽可能避免出现异常情况堵塞，这里简单给了 8，具体值根据情况调整
	messageChannel = make(chan Message, 8)
)

func main() {
	/*
		在 listen 时没有指定 IP，表示绑定到当前机器的所有 IP 上。
		根据具体情况可以限制绑定具体的 IP，比如只绑定在 127.0.0.1 上
		net.Listen(“tcp”, “127.0.0.1:2020”)
	*/
	listener, err := net.Listen("tcp", ":9099")
	if err != nil {
		panic(err)
	}
	// 用于广播消息
	go broadcaster()
	for {
		// 监听连接
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

// broadcaster 用于记录聊天室用户，并进行消息广播：
// 1. 新用户进来；2. 用户普通消息；3. 用户离开
func broadcaster() {
	// 负责登记/注销用户，通过 map 存储在线用户
	users := make(map[*User]struct{})
	for {
		select {
		case user := <-enteringChannel:
			// 新用户进入
			users[user] = struct{}{}
		case user := <-leavingChannel:
			// 用户离开
			delete(users, user)
			// 避免 goroutine 泄露
			close(user.MessageChannel)
		case msg := <-messageChannel:
			// 给所有在线用户发送消息
			for user := range users {
				if user.ID == msg.OwnerID {
					continue
				}
				user.MessageChannel <- msg.Content
			}
		}
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	// 1. 新用户进来，构建该用户的实例
	user := &User{
		ID:             GenUserID(),
		Addr:           conn.RemoteAddr().String(),
		EnterAt:        time.Now(),
		MessageChannel: make(chan string, 8),
	}
	// 2. 当前在一个新的 goroutine 中，用来进行读操作，因此需要开一个 goroutine 用于写操作
	// 读写 goroutine 之间可以通过 channel 进行通信
	go sendMessage(conn, user.MessageChannel)
	// 3. 给当前用户发送欢迎信息；给所有用户告知新用户到来
	user.MessageChannel <- "Welcome, " + user.String()
	msg := Message{
		OwnerID: user.ID,
		Content: "user:`" + strconv.Itoa(user.ID) + "` has enter",
	}
	messageChannel <- msg

	// 4. 将该记录到全局的用户列表中，通过 channel 来写入，避免用锁
	enteringChannel <- user
	// 控制超时用户踢出
	/*
		struct {}是一个无元素的结构体类型，通常在没有信息存储时使用。
		优点是大小为0，不需要内存来存储struct {}类型的值。
		可以用map[string]struct{}来当作成一个set来用。

		var set map[string]struct{}
		set = make(map[string]struct{})

		set["red"] = struct{}{} // struct{}{}  构造了一个struct {}类型的值
		set["blue"] = struct{}{}

		_, ok := set["red"]
		fmt.Println("Is red in the map?", ok) // true
		_, ok = set["green"]
		fmt.Println("Is green in the map?", ok) // false
	*/
	var userActive = make(chan struct{})
	go func() {
		// 设置超时时间为 5封装
		d := 5 * time.Minute
		/*
			NewTimer 创建一个新的 Timer，
			它将在至少持续时间 d 之后在其通道上发送当前时间。
		*/
		timer := time.NewTimer(d)
		for {
			select {
			// 超时关闭连接
			case <-timer.C:
				conn.Close()
			case <-userActive:
				timer.Reset(d)
			}
		}
	}()

	// 5. 循环读取用户的输入
	// bufio.NewScanner() 基于行的输入，每行都剔除分隔标识
	input := bufio.NewScanner(conn)
	for input.Scan() {
		msg.Content = strconv.Itoa(user.ID) + ":" + input.Text()
		messageChannel <- msg

		// 每次接收到用户的消息后，往 userActive 中写入消息，表示用户活跃
		/*
			struct{}{}是一个复合字面量，两个 struct{}{}地址相等
			它构造了一个struct {}类型的值，该值也是空。
		*/
		userActive <- struct{}{}
	}
	if err := input.Err(); err != nil {
		log.Println("读取错误：", err)
	}
	// 6. 用户离开
	leavingChannel <- user
	msg.Content = "user:`" + strconv.Itoa(user.ID) + "` has left"
	messageChannel <- msg
}

// 用于给用户发送信息，ch <-chan string 只允许从 channel 中读取数据
func sendMessage(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

// 生成用户 ID
var (
	globalID int
	idLocker sync.Mutex
)

func GenUserID() int {
	idLocker.Lock()
	defer idLocker.Unlock()

	globalID++
	return globalID
}
