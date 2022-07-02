package main

import (
	"sync"
)

// GODEBUG=gctrace=1,scheddetail=1,schedtrace=10000
// scheddetail: 显示调度器完整信息
// schedtrace: 运行时每 n 毫秒输出到标准err输出关于调度器的摘要信息
// gctrace: 向标准错误流(os.Stderr)输出 GC 运行信息
func main() {
	// 程序运行时，通过debug.SetGCPercent(-1)来动态调整GOGC的值，默认为100
	/*
		SetGCPercent 设置垃圾回收目标百分比：当新分配的数据与上一次回收后剩余的活数据的比率达到该百分比时触发回收。
		SetGCPercent 返回之前的设置。初始设置是启动时 GOGC 环境变量的值，如果未设置该变量，则为 100。
		负百分比禁用垃圾收集。
	*/
	//debug.SetGCPercent(-1)
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(wg *sync.WaitGroup) {
			var couter int
			for i := 0; i < 1e10; i++ {
				couter++
			}
			wg.Done()
		}(&wg)
	}

	wg.Wait()
	// 主动触发 GC 行为
	//runtime.GC()
	// 与之配套的内存相关的归还，通过调用 debug.FreeOSMemory() 实现
	//debug.FreeOSMemory()
}
