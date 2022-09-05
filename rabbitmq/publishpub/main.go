package main

import (
	"go-community/rabbitmq/rabbitmq"
	"log"
	"strconv"
)

func main() {
	rabbitmq := rabbitmq.NewRabbitMqPubSub("" + "newProduct")
	for i := 0; i < 100; i++ {
		rabbitmq.PublishPub("订阅模式生产第" + strconv.Itoa(i) + "条数据")
		log.Println("订阅模式生产第" + strconv.Itoa(i) + "条数据")
	}
}
