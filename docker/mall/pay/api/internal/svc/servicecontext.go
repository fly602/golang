package svc

import (
	jwtx "go-community/docker/mall/common/jwt"
	"go-community/docker/mall/pay/api/internal/config"
	"go-community/docker/mall/pay/rpc/payclient"
	"go-community/docker/mall/vendor/github.com/zeromicro/go-zero/rest"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	PayRpc    payclient.Pay
	JwtHeader rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		PayRpc:    payclient.NewPay(zrpc.MustNewClient(c.PayRpc)),
		JwtHeader: jwtx.NewJwtheaderMiddleware(c.Auth).Handle,
	}
}
