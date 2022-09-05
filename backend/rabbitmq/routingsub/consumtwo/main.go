package main

import (
	"go-community/rabbitmq/rabbitmq"
)

func main() {
	kutengone := rabbitmq.NewRabbitMqRouting("kuteng", "kuteng_two")
	kutengone.ReceiveRouting()
}
