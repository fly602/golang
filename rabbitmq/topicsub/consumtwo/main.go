package main

import "go-community/rabbitmq/rabbitmq"

func main() {
	kutengone := rabbitmq.NewRabbitMqRouting("exkutengTopic", "kuteng.*.one")
	kutengone.ReceivedTopic()
}
