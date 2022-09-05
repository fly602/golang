package main

import (
	"fmt"
	"go-community/rpc/tcp/server/goods"
	"net/http"
	"net/rpc"

	"net"
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
	shop := goods.InitShop("fuleyi")
	shop.PutGoods(pens)
	shop.PutGoods(eggs)
	shop.PutGoods(clothes)

	rpc.Register(shop)
	rpc.HandleHTTP()
	lis, err := net.Listen("tcp", ":30001")
	if err != nil {
		fmt.Println("listen err:", err)
		return
	}
	defer lis.Close()
	http.Serve(lis, nil)
}
