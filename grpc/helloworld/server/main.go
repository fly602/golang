package main

import (
	"context"
	"fmt"
	pb "go-community/grpc/helloworld/proto/hello"

	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
)

const (
	Address = "127.0.0.1:50052"
)

type HelloService struct{}

var hs = HelloService{}

func (h HelloService) Sayhello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	resp := new(pb.HelloResponse)
	resp.Message = fmt.Sprintf("Hello %s", in.Name)
	log.Println("Sayhello", resp.Message)
	return resp, nil
}

func (h HelloService) SayHelloNoResp(ctx context.Context, in *wrappers.Int32Value) (*empty.Empty, error) {
	resp := new(empty.Empty)
	log.Println("SayHelloNoResp value=", in.Value)
	return resp, nil
}

func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		log.Fatalln("Failed to listen，", err)
	}
	// 实例化grpc server
	s := grpc.NewServer()

	// 注册HelloService
	pb.RegisterHelloServer(s, hs)

	log.Println("Listen on:", Address)
	s.Serve(listen)

}
