package main

import (
	"fmt"
	"time"
)

func main() {
	name := "test"
	go func() {
		name = "admin"
	}()
	fmt.Println("name is", name)

	time.Sleep(1e9)
}
