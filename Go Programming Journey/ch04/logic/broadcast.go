package logic

import (
	"ch04/global"
	"log"
)

// broadcaster 广播器
type broadcaster struct {
	// 所有聊天室用户
	users map[string]*User
	// 所有 channel 统一管理，可以避免外部乱用
	enteringChannel chan *User
	leavingChannel  chan *User
	messageChannel  chan *Message

	// 判断该昵称用户是否可以进入聊天室（重复与否）：true 能, false 不能
	checkUserChannel      chan string
	checkUserCanInChannel chan bool

	// 获取用户列表
	requestUsersChannel chan struct{}
	usersChannel        chan []*User
}

// 判断用户是否存在
func (b *broadcaster) CanEnterRoom(nickname string) bool {
	b.checkUserChannel <- nickname
	return <-b.checkUserCanInChannel
}

var Broadcaster = &broadcaster{
	users: make(map[string]*User),

	enteringChannel: make(chan *User),
	leavingChannel:  make(chan *User),
	messageChannel:  make(chan *Message, global.MessageQueueLen),

	checkUserChannel:      make(chan string),
	checkUserCanInChannel: make(chan bool),

	requestUsersChannel: make(chan struct{}),
	usersChannel:        make(chan []*User),
}

// Start 启动器
// 需要在一个新的 goroutine 中运行，因为它不会返回
func (b *broadcaster) Start() {
	for {
		// select case 语句中不能使用 fallthrough
		select {
		/*
			每个 case 后都必须有一个 channel 的接受或者发送操作
			所有非阻塞 case 操作将会有一个被随机选择执行（不是按照从上至下的顺序）
			在所有 case 操作均阻塞的情况下，如果存在 default 分支，则执行该分支
			否则当前 goroutine 进入阻塞状态
			***********
			并且，根据以上规则，可以得出
			一个不含任何分支的 select-case 代码块 select{}
			将使当前 goroutine 处于永久阻塞状态
			e.g.
			func main() {
				go func() {
					// 该函数不会退出
					for {
						// 省略代码
					}
				} ()
				select {}
			}
			***********
			这样可以让 main goroutine 永远阻塞，让其他 goroutine 运行
			如果没有其他可运行的 goroutine 将会导致死锁
			fatal error: all goroutines are asleep - deadlock!
			e.g.
			func main(){
				select()
			}
		*/
		case user := <-b.enteringChannel:
			// 新用户进入
			b.users[user.NickName] = user

			OfflineProcessor.Send(user)
		case user := <-b.leavingChannel:
			// 用户离开
			delete(b.users, user.NickName)
			// 避免 goroutine 泄露
			user.CloseMessageChannel()
		case msg := <-b.messageChannel:
			// 给所有在线用户发送消息
			for _, user := range b.users {
				if user.UID == msg.User.UID {
					continue
				}
				user.MessageChannel <- msg
			}
			OfflineProcessor.Save(msg)
		case nickname := <-b.checkUserChannel:
			if _, ok := b.users[nickname]; ok {
				b.checkUserCanInChannel <- false
			} else {
				b.checkUserCanInChannel <- true
			}
		case <-b.requestUsersChannel:
			userList := make([]*User, 0, len(b.users))
			for _, user := range b.users {
				userList = append(userList, user)
			}

			b.usersChannel <- userList
		}
	}
}

func (b *broadcaster) UserEntering(u *User) {
	b.enteringChannel <- u
}

func (b *broadcaster) UserLeaving(u *User) {
	b.leavingChannel <- u
}

func (b *broadcaster) Broadcast(msg *Message) {
	if len(b.messageChannel) >= global.MessageQueueLen {
		log.Println("broadcast queue 满了")
	}
	b.messageChannel <- msg
}

func (b *broadcaster) GetUserList() []*User {
	b.requestUsersChannel <- struct{}{}
	return <-b.usersChannel
}
