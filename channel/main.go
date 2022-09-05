package main

import (
	"fmt"
	"log"
)

type NUM uint32
type FLogger struct {
	log.Logger
}

func (l *FLogger) Printf(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf(format, v...))
}

func (n NUM) Add(a, b uint32) {
	n = NUM(a + b)
	fmt.Printf("n: %+v\n", n)
}

func main() {
	n := new(NUM)
	n.Add(1, 2)
	flogger.New
	flogger.Printf("======")
	// taskChan := make(chan int, 100)

	// go func() {
	// 	for i := 0; i < 100; i++ {
	// 		taskChan <- i
	// 	}
	// }()
	// go func() {
	// 	time.Sleep(time.Second * 10)
	// 	close(taskChan)
	// }()
	// for task := range taskChan {
	// 	fmt.Println(task)
	// 	log.Println("")
	// }
}
