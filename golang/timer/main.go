package main

import (
	"fmt"
	"sync"
	"time"
)

func aaa() {

}

func calcTime() {
	time.Now().Second()
	t1 := time.Now()
	time.Sleep(time.Second)
	fmt.Println("===>>>", time.Until(t1).Seconds())
}

func main() {
	// ticker := time.NewTicker(time.Second * 1)
	// for {
	// 	<-ticker.C
	// 	log.Println("hello world...")
	// }
	calcTime()
	wg := sync.WaitGroup{}
	// 处理定时器任务
	counter := func() func() int {
		count := 5 + 1
		return func() int {
			count--
			return count
		}
	}()
	timer := time.NewTimer(time.Second * 1)
	time.AfterFunc(time.Second, aaa)
	timer.Reset(0)
	wg.Add(1)
	go func() {
		for range timer.C {
			count := counter()
			if count == 0 {
				fmt.Println("counting:", count)
				wg.Done()
			} else {
				fmt.Println("counting:", count)
				timer.Reset(0)
			}
		}
	}()

	fmt.Println("time done")
	wg.Wait()
	timer.Stop()
	fmt.Printf("%p\n", timer)
}
