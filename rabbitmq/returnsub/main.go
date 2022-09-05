package main

import "go-community/rabbitmq/rabbitmq"

func main() {
	rabbitmq := rabbitmq.NewRabbitMqPubSub("ex-return")
	rabbitmq.ReceiveReturn()

}
