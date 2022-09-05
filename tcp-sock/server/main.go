package main

import (
	"bufio"
	"io"
	"log"
	"net"
)

func process(conn net.Conn) {
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		var buf [1024]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			if err.Error() == io.EOF.Error() {
				return
			}
			log.Fatalln("Tcp server read err,", err)
		}
		recvStr := string(buf[:n])
		log.Println("接收到客户端数据：", recvStr)
		conn.Write([]byte(recvStr))
	}
}

func main() {
	log.Println("Tcp Server Starting...")
	listen, err := net.Listen("tcp", "127.0.0.1:29999")
	if err != nil {
		log.Fatalln("Tcp listen err,", err)
	}
	for {
		_, err := listen.Accept()
		log.Println("New accetp...")
		if err != nil {
			log.Fatalln("Tcp accept err,", err)
		}
		//go process(conn)
	}
}
