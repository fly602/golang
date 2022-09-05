package main

import (
	"fmt"
	"go-community/rabbitmq/rabbitmq"
	"log"
	"time"
)

func main() {
	rabbitmq := rabbitmq.NewRabbitMqSimple("" + "guest")
	for i := 0; i < 100; i++ {
		rabbitmq.PublishSimple(fmt.Sprintf("Hello guest!!!,%d", i))
		time.Sleep(time.Millisecond * 10)
		log.Println("Send success!!!", i)
	}
	select {}
}
