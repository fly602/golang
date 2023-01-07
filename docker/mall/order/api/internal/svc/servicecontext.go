package svc

import (
	jwtx "go-community/docker/mall/common/jwt"
	"go-community/docker/mall/order/api/internal/config"
	"go-community/docker/mall/order/rpc/orderclient"
	"go-community/docker/mall/product/rpc/productclient"
	"go-community/docker/mall/vendor/github.com/zeromicro/go-zero/rest"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	OrderRpc   orderclient.Order
	ProductRpc productclient.Product
	JwtHeader  rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		OrderRpc:   orderclient.NewOrder(zrpc.MustNewClient(c.OrderRpc)),
		ProductRpc: productclient.NewProduct(zrpc.MustNewClient(c.ProductRpc)),
		JwtHeader:  jwtx.NewJwtheaderMiddleware(c.Auth).Handle,
	}
}
