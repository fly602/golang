package main

import (
	"fmt"
	"go-community/rabbitmq/rabbitmq"
	"log"
	"time"
)

func main() {
	r := rabbitmq.NewRabbitMqRouting("ex-dxl-normal", "normal")
	for i := 1001; i < 1001+1; i++ {
		orderInfo := fmt.Sprintf("订单信息: id=%d,下单时间=%v", i, time.Now())
		r.DLXPub(orderInfo)
		log.Println(orderInfo)
	}
}
