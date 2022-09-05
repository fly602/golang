package main

import (
	"go-community/rabbitmq/rabbitmq"
)

func main() {
	rabbitymq := rabbitmq.NewRabbitMqSimple("" + "guest")
	rabbitymq.ComsumeSimple()
}
