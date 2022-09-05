package main

import "fmt"

func addN(num int) func(p *int) int {
	n := num
	return func(p *int) int {
		*p += n
		return *p
	}
}

// test1 内存逃逸1
func test1() {
	f := addN(10)
	i := 11
	fmt.Println("i=", f(&i))
}

// test2 内存逃逸2
func test2() {
	i := 10
	var m []*int
	m = append(m, &i)
}

// test3 内存逃逸3
func test3() {
	s := make([]int, 10000, 10000)
	for idx, _ := range s {
		s[idx] = idx
	}
}

// test4 内存逃逸4：使用接口调用方法
type itf interface{ Get() }
type me struct{}

func (me me) Get() {}
func test4() {
	var i itf = me{}
	i.Get()
}

func main() {
	var a [1024 * 1024 * 1024]uint64
	a[100] = 100
	test2()
}
