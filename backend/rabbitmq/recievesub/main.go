package main

import "go-community/rabbitmq/rabbitmq"

func main() {
	rabbitmq := rabbitmq.NewRabbitMqPubSub("" + "newProduct")
	rabbitmq.ReceiveSub()
}
