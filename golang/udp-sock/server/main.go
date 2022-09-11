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
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("udp listen err:", err)
		return
	}
	defer conn.Close()

	for {
		buff := make([]byte, 1024)
		_, remote, err := conn.ReadFromUDP(buff)
		if err != nil {
			fmt.Println("udp recv err:", err)
			break
		}
		fmt.Printf("[%v]: %s\n", remote.String(), string(buff))
		// 一旦服务器接收效率底下，查看udp缓冲器大小
		//time.Sleep(time.Second)
	}
}
