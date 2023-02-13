package main

import (
	"go-community/backend/kafka/kafka"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	g := kafka.NewGroup()
	c := g.Connect()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigterm:
		log.Println("terminating: via signal")
	}
	c()
}
