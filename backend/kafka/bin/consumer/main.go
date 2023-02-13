package main

import (
	"go-community/backend/kafka/kafka"
	"log"
)

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	kafka.NewClient()
}
