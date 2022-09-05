package main

import (
	"log"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second * 1)
	for {
		<-ticker.C
		log.Println("hello world...")
	}
}
