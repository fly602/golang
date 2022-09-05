package main

import (
	"log"
	"time"
)

type pool struct {
	ch chan bool
}

var p = pool{}

func runTask(i int) {
	for value := range p.ch {
		log.Println("Goroutine", i, "running...", value)
	}
}

func custom() {
	for {
		p.ch <- true
		time.Sleep(time.Microsecond)
	}
}

func main() {
	p.ch = make(chan bool, 100)
	// for i := 0; i < 100; i++ {
	// 	go custom()
	// }
	for i := 0; i < 100; i++ {
		go runTask(i)
	}
	for i := 0; i < 10; i++ {
		p.ch <- true
	}
	time.Sleep(time.Second)
	log.Println("task len=", len(p.ch))
	close(p.ch)
}
