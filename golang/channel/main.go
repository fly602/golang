package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	taskChan := make(chan int, 100)

	go func() {
		for i := 0; i < 100; i++ {
			taskChan <- i
		}
	}()
	go func() {
		time.Sleep(time.Second * 10)
		close(taskChan)
	}()
	for task := range taskChan {
		fmt.Println(task)
		log.Println("")
	}
}
