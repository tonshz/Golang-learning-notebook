# Go 语言编程之旅(四)：Websocket 应用(一)

## 一、基于 TCP 的聊天室

本节通过命令行来模拟基于 TCP 的简单聊天室。

本程序可以将用户发送的文本消息广播给该聊天室内的所有其他用户。该服务端程序中有四种 goroutine：main goroutine 和 广播消息的 goroutine，以及每一个客户端连接都会有一对读和写的 goroutine。

### 1. 代码实现

