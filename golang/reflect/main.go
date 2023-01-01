package main

import (
	"fmt"
	"reflect"
	"time"
)

type Order struct {
	Id         int64     `db:"id"`
	Uid        int64     `db:"uid"`    // 用户ID
	Pid        int64     `db:"pid"`    // 产品ID
	Amount     int64     `db:"amount"` // 订单金额
	Status     int64     `db:"status"` // 订单状态
	CreateTime time.Time `db:"create_time"`
	UpdateTime time.Time `db:"update_time"`
}

func main() {
	var resp []Order

	var v interface{} = resp
	fmt.Println("typeof(v)=", reflect.TypeOf(v))
	fmt.Println("typeof(v).Elem=", reflect.TypeOf(v).Elem())
	fmt.Println("typeof(v).Elem.Kind=", reflect.TypeOf(v).Elem().Kind())

}
