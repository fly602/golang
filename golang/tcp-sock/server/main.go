package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

// 处理：手动接受报文
// 户口端发送n多个报文，server端手动接受报文，查看缓冲区和tcp连接的窗口变化
// netstat -anpt |grep <port>
func processManual(conn net.Conn) {
	defer conn.Close()
	for {
		fmt.Printf("请输入接收次数，长度(0 * 退出)：")
		var times, len int
		fmt.Scanln(&times, &len)
		if times == 0 {
			return
		}
		for i := 0; i < times; i++ {
			var buf []byte = make([]byte, len)
			n, err := conn.Read(buf[:])
			if err != nil {
				if err.Error() == io.EOF.Error() {
					return
				}
				log.Fatalln("Tcp server read err,", err)
			}
			recvStr := string(buf[:n])
			log.Println("接收到客户端数据：", recvStr)
		}
	}
	select {}
}

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
		conn, err := listen.Accept()
		log.Println("New accetp...")
		if err != nil {
			log.Fatalln("Tcp accept err,", err)
		}
		//go process(conn)
		go processManual(conn)
	}
}
