package main

import "go-community/rabbitmq/rabbitmq"

func main() {
	r := rabbitmq.NewRabbitMqRouting("ex-dxl", "dxl")
	r.ReceiveRouting()
}
