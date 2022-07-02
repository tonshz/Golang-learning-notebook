package main

import (
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
)

func init() {
	// 值小于等于 0 表示关闭该数据的采集
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)
}

func main() {

	var m sync.Mutex
	var datas = make(map[int]struct{})
	for i := 0; i < 999; i++ {
		go func(i int) {
			m.Lock()
			defer m.Unlock()
			datas[i] = struct{}{}
		}(i)
	}

	_ = http.ListenAndServe(":8080", nil)
}
