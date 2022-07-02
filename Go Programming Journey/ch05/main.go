package main

import (
	"encoding/json"
	"github.com/allegro/bigcache"
	"strconv"
	"time"
)

//func main() {
//	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	entry, err := cache.Get("my-uinique-key")
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	if entry == nil {
//		// 从缓存中没有获取到，则从数据源中获取（一般是数据库），然后设置到缓存
//		entry = []byte("value") // 实际从数据库中获取
//		cache.Set("my-uinique-key", entry)
//	}
//	log.Println(string(entry))
//}

type Value struct {
	A string
	B int
	C time.Time
	D []byte
	E float32
	F *string
	T T
}

type T struct {
	G int
	I int
	K int
	L int
	M int
	N int
}

//func main() {
//	num := 10000000
//	config := bigcache.DefaultConfig(10 * time.Minute)
//	cache, _ := bigcache.NewBigCache(config)
//
//	for i := 0; i < num; i++ {
//		cache.Set(strconv.Itoa(i), []byte(cast.ToString(&Value{})))
//	}
//
//	for i := 0; ; i++ {
//		cache.Delete(strconv.Itoa(i))
//		cache.Set(strconv.Itoa(num+i), []byte(cast.ToString(&Value{})))
//		time.Sleep(5 * time.Millisecond)
//	}
//}

func main() {
	num := 10000000
	config := bigcache.DefaultConfig(10 * time.Minute)
	cache, _ := bigcache.NewBigCache(config)

	for i := 0; i < num; i++ {
		j, _ := json.Marshal(&Value{})
		cache.Set(strconv.Itoa(i), j)
	}

	for i := 0; ; i++ {
		cache.Delete(strconv.Itoa(i))
		j, _ := json.Marshal(&Value{})
		cache.Set(strconv.Itoa(num+i), j)
		time.Sleep(5 * time.Millisecond)
	}
}
