package main

import (
	"log"
	"net"
	"strconv"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:29999")
	if err != nil {
		log.Fatalln("Client dail err,", err)
	}
	defer conn.Close()
	for i := 0; i < 1; i++ {
		var buff string = "hello world [" + strconv.Itoa(i) + "]"
		_, err := conn.Write([]byte(buff))
		if err != nil {
			log.Println("send err:", err)
			return
		}
		log.Println("send success:", buff)
	}
	select {}
}
