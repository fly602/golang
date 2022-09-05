package main

import (
	"go-community/kafka/kafka"
	"log"
)

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	kafka.NewClient()
}
