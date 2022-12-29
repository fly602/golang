package svc

import (
	"go-community/docker/mall/product/api/internal/config"
	"go-community/docker/mall/product/rpc/productclient"
	"go-community/docker/mall/vendor/github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	ProductRpc productclient.Product
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		ProductRpc: productclient.NewProduct(zrpc.MustNewClient(c.ProductRpc)),
	}
}
