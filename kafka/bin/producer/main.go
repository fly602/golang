package main

import (
	"fmt"
	"go-community/kafka/kafka"
	"log"

	"github.com/Shopify/sarama"
)

var auto = false
var topic = "zaplog"

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)

	producer := &kafka.Producer{}
	producer.NewProducer(sarama.NewRandomPartitioner)
	//producer.InitMsgPool("second", "", "", "This is a test log")
	producer.NewMsgPool(topic, "", "")
	if auto {
		automsg(producer)
	} else {
		manualmsg(producer)
	}
	defer producer.MsgPool.CloseMsgPool()
}

func automsg(producer *kafka.Producer) {
	var msg string
	for i := 0; i < 2; i++ {

		msg = fmt.Sprintf("This is a test log,id = %d", i)
		producer.MsgPool.PushMsg(msg)
	}
}

func manualmsg(producer *kafka.Producer) {
	var msg string
	for {
		fmt.Scanln(&msg)
		if msg == "exit" {
			break
		}
		producer.MsgPool.PushMsg(msg)
	}
}
