package server

import (
	"ch04/logic"
	"net/http"
)

func RegisterHandle() {

	// 广播消息处理
	go logic.Broadcaster.Start()

	// 注册两个路由
	// 其中 “/” 代表首页，”/ws” 用来服务 WebSocket 长连接
	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/user_list", userListHandleFunc)
	http.HandleFunc("/ws", WebSocketHandleFunc)
}
