package main

import (
	"fmt"
	"net/rpc"
)

type GoodsInfo struct {
	ID    uint32
	Name  string
	Price float32
	Rest  uint32
}

type Goods struct {
	ID   uint32
	Name string
}

func main() {
	for i := 0; i < 60000; i++ {
		go func() {
			client, err := rpc.DialHTTP("tcp", ":30001")
			if err != nil {
				fmt.Println("rpc dial err:", err)
				return
			}
			defer client.Close()
			pens := &Goods{
				ID:   1,
				Name: "Pen",
			}
			var gi GoodsInfo
			err = client.Call("Shop.GetGoodsInfo", &pens, &gi)
			if err != nil {
				fmt.Println("rpc call err:", err)
				return
			}
			fmt.Printf("rpc call reply: %+v\n", gi)

			rpccall := client.Go("Shop.GetGoods", &pens, &gi, nil)
			repcall := <-rpccall.Done
			if repcall.Error != nil {
				fmt.Println("rpc Go err,", repcall.Error)
				return
			}
			fmt.Printf("rpc call reply: %+v\n", gi)
		}()
	}
	select {}
}
