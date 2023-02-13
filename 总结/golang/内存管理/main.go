package main

import (
	"fmt"
	"runtime"
	"time"
)

var stat runtime.MemStats

func zerotest() {
	var (
		a struct{}
		b [0]int
		c [100]struct{}
		d = make([]struct{}, 1024)
	)
	fmt.Printf("%p\n", &a)
	fmt.Printf("%p\n", &b)
	fmt.Printf("%p\n", &c)
	fmt.Printf("%p\n", &(d[0]))
	fmt.Printf("%p\n", &(d[1]))
	fmt.Printf("%p\n", &(d[1000]))
}

// 变量a内存未逃逸:go build -gcflags="-m -l"
func stack1() {
	var a int = 100
	go func(i int) {
		time.Sleep(time.Second)

	}(a)
	a = 0
}

// 变量a内存逃逸:go build -gcflags="-m -l"
// moved to heap: a
func stack2() {
	var a int = 100
	go func() {
		a = a + 1
	}()
	a = 0
}

func main() {
	runtime.ReadMemStats(&stat)
	println(stat.HeapSys)
	zerotest()
	stack1()
	stack2()
}
