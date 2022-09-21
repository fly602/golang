package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
)

const (
	TCP_OPTION_NODELAY = true
	TCP_OPTION_LINGER  = 0
)

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return b
}

func tcpSetOPtion(conn net.Conn) {
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return
	}
	tcpConn.SetNoDelay(TCP_OPTION_NODELAY)
	tcpConn.SetLinger(TCP_OPTION_LINGER)
}

func sendManual(conn net.Conn) {
	defer conn.Close()
	for {
		fmt.Printf("请输入发送[次数，长度, 0-*退出]：")
		var times, len int
		fmt.Scanln(&times, &len)
		if times == 0 {
			return
		}
		for i := 0; i < times; i++ {
			buff := randBytes(len)
			_, err := conn.Write(buff)
			if err != nil {
				log.Println("send err:", err)
				return
			}
			log.Println("send success:", string(buff))
		}
	}
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:29999")
	if err != nil {
		log.Fatalln("Client dail err,", err)
	}
	tcpSetOPtion(conn)
	sendManual(conn)
}
