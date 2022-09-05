package main

import "go-community/rabbitmq/rabbitmq"

func main() {
	kutengone := rabbitmq.NewRabbitMqRouting("exkutengTopic", "#")
	kutengone.ReceivedTopic()
}
