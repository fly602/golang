package main

import (
	"go-community/rabbitmq/rabbitmq"
	"log"
	"strconv"
	"time"
)

func main() {
	kutengone := rabbitmq.NewRabbitMqRouting("kuteng", "kuteng_one")
	kutengtwo := rabbitmq.NewRabbitMqRouting("kuteng", "kuteng_two")
	for i := 0; i < 100; i++ {
		kutengone.PublishRouting("Hello kuteng one!" + strconv.Itoa(i))
		log.Println("Hello kuteng one!" + strconv.Itoa(i))
		kutengtwo.PublishRouting("Hello kuteng two!" + strconv.Itoa(i))
		log.Println("Hello kuteng two!" + strconv.Itoa(i))
		time.Sleep(time.Millisecond)
	}
}
