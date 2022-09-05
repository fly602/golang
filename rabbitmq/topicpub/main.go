package main

import (
	"go-community/rabbitmq/rabbitmq"
	"log"
	"strconv"
)

func main() {
	kutengone := rabbitmq.NewRabbitMqRouting("exkutengTopic", "kuteng.topic.one")
	kutengtwo := rabbitmq.NewRabbitMqRouting("exkutengTopic", "kuteng.topic.two")
	var i int
	for {
		i++
		kutengone.PublishTopic("Hello kuteng topic one!" + strconv.Itoa(i))
		log.Println("Hello kuteng topic one!" + strconv.Itoa(i))
		kutengtwo.PublishTopic("Hello kuteng topic two!" + strconv.Itoa(i))
		log.Println("Hello kuteng topic two!" + strconv.Itoa(i))
	}
}
