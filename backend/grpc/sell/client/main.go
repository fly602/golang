package main

import (
	"context"
	"fmt"
	"go-community/grpc/sell/proto/sell"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

func main() {
	grpc.EnableTracing = true
	conn, err := grpc.Dial("127.0.0.1:30001", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("dial err,%s", err)
	}
	client := sell.NewDealClient(conn)

	total, err := client.ListGoods(context.Background(), &empty.Empty{})
	if err != nil {
		fmt.Println("listgoods err,", err)
		return
	}
	fmt.Printf("total:%+v\n", total)

	pen := &sell.Goods{
		Id:   1,
		Name: "Pen",
	}
	gi, err := client.Consume(context.Background(), pen)
	if err != nil {
		fmt.Printf("%s consume err,%s\n", gi.G.Name, err)
		return
	}
	fmt.Printf("%s consume success,%+v\n", gi.G.Name, gi)
}
