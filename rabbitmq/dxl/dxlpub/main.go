package main

import "go-community/rabbitmq/rabbitmq"

func main() {
	r := rabbitmq.NewRabbitMqRouting("ex-dxl-normal", "normal")
	r.DLXPub("hello dxl")
}
