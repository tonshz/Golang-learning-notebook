package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":9099")
	if err != nil {
		panic(err)
	}
	done := make(chan struct{})
	go func() {
		// 包括从标准输入读取数据写入 TCP 连接中
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		log.Println("done")
		// 新开的 goroutine 通过一个 channel 来和 main goroutine
		done <- struct{}{} // signal the main goroutine
	}()
	mustCopy(conn, os.Stdin)
	defer conn.Close()
	<-done
}

func mustCopy(dst io.Writer, src io.Reader) {
	// 从 TCP 连接中读取数据写入标准输出
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
