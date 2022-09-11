package main

import (
	"fmt"
	"net"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", ":30001")
	if err != nil {
		fmt.Println("udp addr err:", err)
		return
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("udp dial err:", err)
		return
	}
	defer conn.Close()
	buff := "hello world"
	for i := 0; i < 1000000; i++ {
		_, err = conn.Write([]byte(buff))
		if err != nil {
			fmt.Println("udp write err:", err)
			return
		}
	}
}
