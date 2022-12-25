package svc

import (
	"go-community/docker/mall/order/api/internal/config"
	"go-community/docker/mall/order/rpc/orderclient"
	"go-community/docker/mall/vendor/github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	OrderRpc orderclient.Order
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		OrderRpc: orderclient.NewOrder(zrpc.MustNewClient(c.OrderRpc)),
	}
}
