package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var datas []string

func main() {
	go func() {
		log.Printf("len: %d", Add("pprof-test"))
		time.Sleep(time.Millisecond * 10)
	}()

	_ = http.ListenAndServe(":8080", nil)
}

func Add(str string) int {
	data := []byte(str)
	datas = append(datas, string(data))
	return len(data)
}
