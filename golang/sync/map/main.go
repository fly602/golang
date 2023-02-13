package main

import (
	"fmt"
	"sync"
)

func main() {
	sm := new(sync.Map)
	sm.Store("aaa", "aaa")
	val, ok := sm.Load("aaa")
	if ok {
		fmt.Printf("sync.load[aaa]=%v\n", val)
	}
}
