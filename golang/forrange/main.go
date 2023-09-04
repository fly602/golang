package main

import (
	"fmt"
	"time"
)

func main() {
	var m = [...]int{1, 2, 3, 4, 5}

	for i, v := range m {
		go func() {
			time.Sleep(time.Second * 3)
			fmt.Println(i, v)
		}()
	}

	var a = 1

	f1 := func() {
		time.Sleep(time.Second * 3)
		fmt.Println(a)
	}

	a = 2
	f1()
	time.Sleep(time.Second * 10)
}
