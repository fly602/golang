package main

import (
	"log"
	"time"
)

func main() {
	var taskChan = make(chan int, 100)

	go func() {
		for i := 0; i < 10; i++ {
			taskChan <- i
		}
	}()
	go func() {
		time.Sleep(time.Second * 3)
		close(taskChan)
	}()

	f1 := func(i int) {
		for task := range taskChan {
			log.Println("go", i, " =", task)
		}
	}

	go f1(1)
	f1(2)

	// 给已关闭的管道发送数据，会崩溃
	// taskChan <- 10

	task, ok := <-taskChan
	log.Println("Last get task=", task, "ok =", ok)
}
