package main

import (
	"fmt"
	"go-community/grpc/sell/proto/sell"
	"go-community/grpc/sell/server/goods"
	"net"
	"net/http"

	"google.golang.org/grpc"
)

var pens = &goods.GoodsInfo{
	Name:  "Pen",
	Price: 10.50,
	Rest:  99,
}

var eggs = &goods.GoodsInfo{
	Name:  "Egg",
	Price: 0.50,
	Rest:  10000,
}

var clothes = &goods.GoodsInfo{
	Name:  "Clothes",
	Price: 1000,
	Rest:  1,
}

func main() {
	grpc.EnableTracing = true
	listen, err := net.Listen("tcp", "127.0.0.1:30001")
	if err != nil {
		fmt.Printf("net listen err: %s", err)
		return
	}

	service := grpc.NewServer()
	go http.ListenAndServe("127.0.0.1:30002", http.DefaultServeMux)
	shop := goods.InitShop("fuleyi")
	shop.PutGoods(pens)
	shop.PutGoods(eggs)
	shop.PutGoods(clothes)

	sell.RegisterDealServer(service, shop)
	service.Serve(listen)
}
