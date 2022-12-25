package svc

import (
	"go-community/docker/mall/order/model"
	"go-community/docker/mall/order/rpc/internal/config"
	"go-community/docker/mall/product/rpc/productclient"
	"go-community/docker/mall/user/rpc/userclient"
	"go-community/docker/mall/vendor/github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-community/docker/mall/vendor/github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	OrderModel model.OrderModel

	UserRpc    userclient.User
	ProductRpc productclient.Product
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:     c,
		OrderModel: model.NewOrderModel(conn, c.CacheRedis),
		UserRpc:    userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		ProductRpc: productclient.NewProduct(zrpc.MustNewClient(c.ProductRpc)),
	}
}
