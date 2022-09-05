package main

import "go-community/rabbitmq/rabbitmq"

func main() {
	rabbitmq := rabbitmq.NewRabbitMqPubSub("ex-confirm")
	rabbitmq.PublishConfirm("hello confirm")
}
