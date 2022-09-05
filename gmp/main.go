package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func test() {
	value := make([]int, 1000)
	for i := range value {
		value[i] = i
	}
	log.Println("test func done...")
	time.Sleep(time.Minute * 10)
}

func main() {
	for i := 0; i < 1000; i++ {
		go test()
	}
	go func() {
		http.ListenAndServe("localhost:28080", nil)
	}()
	select {}
}
