package main

import (
	"context"
	pb "go-community/grpc/helloworld/proto/hello"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	Address = "127.0.0.1:50052"
)

func main() {
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Dail failed,err:", err)
	}
	// 初始化客户端
	c := pb.NewHelloClient(conn)
	// 调用方法
	req := &pb.HelloRequest{Name: "gRPC"}
	res, err := c.Sayhello(context.Background(), req)
	if err != nil {
		log.Fatalln("Sayhello err,", err)
	}
	log.Println("res=", res.Message)

	_, err = c.SayHelloNoResp(context.Background(), wrapperspb.Int32(100))
	if err != nil {
		log.Fatalln("SayHelloNoResp err,", err)
	}

}
