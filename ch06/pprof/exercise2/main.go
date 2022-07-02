package main

import (
	"net/http"
	// "net/http/pprof" 中对调用的 pprof 接口进行了初始化
	/*
		func init() {
			http.HandleFunc("/debug/pprof/", Index)
			http.HandleFunc("/debug/pprof/cmdline", Cmdline)
			http.HandleFunc("/debug/pprof/profile", Profile)
			http.HandleFunc("/debug/pprof/symbol", Symbol)
			http.HandleFunc("/debug/pprof/trace", Trace)
		}
		实际上 net/http/pprof 会在初始化函数中对标准库中 net/http 所默认提供的 DefaultServeMux 进行路由注册，源码如下
		var DefaultServeMux = &defaultServeMux
		var defaultServeMux ServeMux
		func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
		    DefaultServeMux.HandleFunc(pattern, handler)
		}
	*/
	/*
		在实际项目中，都是有相对独立的 ServeMux 的，只要仿照着将 pprof 对应的路由注册进去就可以了
		mux := http.NewServeMux()
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	*/
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
