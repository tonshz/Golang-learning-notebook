package logic

import (
	"container/ring"
	"github.com/spf13/viper"
)

type offlineProcessor struct {
	n int

	// 保存所有用户最近的 n 条消息
	// *ring.Ring 是一个循环链表，将离线消息存储在进程的内存中
	/*
		var r ring.Ring
		fmt.Println(r.Len())    // Output: 1
		fmt.Println(r.Value)  // Output: nil
	*/
	recentRing *ring.Ring

	// 保存某个用户离线消息（一样 n 条）
	userRing map[string]*ring.Ring
}

var OfflineProcessor = newOfflineProcessor()

func newOfflineProcessor() *offlineProcessor {
	n := viper.GetInt("offline-num")

	return &offlineProcessor{
		n:          n,
		recentRing: ring.New(n), // ring.New(n) 创建了 n 个 Ring 实例指针
		userRing:   make(map[string]*ring.Ring),
	}
}

func (o *offlineProcessor) Save(msg *Message) {
	if msg.Type != MsgTypeNormal {
		return
	}
	// 根据 Ring 的使用方式，将用户信息直接存入 recentRing 中，并后移一个位置
	o.recentRing.Value = msg
	o.recentRing = o.recentRing.Next()

	// 判断消息中是否有 @ 谁，需要单独为它保存一个消息列表
	for _, nickname := range msg.Ats {
		// [1:] 第一个为 @ 不需要
		nickname = nickname[1:]
		var (
			r  *ring.Ring
			ok bool
		)
		if r, ok = o.userRing[nickname]; !ok {
			r = ring.New(o.n)
		}
		r.Value = msg
		o.userRing[nickname] = r.Next()
	}
}

func (o *offlineProcessor) Send(user *User) {
	// 遍历最近消息，发送给用户
	o.recentRing.Do(func(value interface{}) {
		if value != nil {
			user.MessageChannel <- value.(*Message)
		}
	})
	// 最近消息是所有用户共有的，不能删除

	if user.isNew {
		return
	}

	// 如果不是新用户，查询是否有 @ 该用户的信息
	if r, ok := o.userRing[user.NickName]; ok {
		// 有则发送给用户
		r.Do(func(value interface{}) {
			if value != nil {
				user.MessageChannel <- value.(*Message)
			}
		})

		// 发送完后将这些消息删除，@ 用户的消息是用户独有的，可以删除
		delete(o.userRing, user.NickName)
	}
}
