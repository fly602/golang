package main

import (
	"go-community/rabbitmq/rabbitmq"
)

func main() {
	kutengone := rabbitmq.NewRabbitMqRouting("kuteng", "kuteng_one")
	kutengone.ReceiveRouting()
}
