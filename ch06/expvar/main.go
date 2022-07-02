package main

import (
	"expvar"
	"fmt"
	"net/http"
	"runtime"
)

func main() {
	expvarFunc := expvar.Get("memstats").(expvar.Func)
	memstats := expvarFunc().(runtime.MemStats)
	fmt.Printf("memstats.GCSys: %d", memstats.GCSys)
	_ = http.ListenAndServe(":6060", http.DefaultServeMux)
}
