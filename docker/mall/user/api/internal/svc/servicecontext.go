package svc

import (
	"go-community/docker/mall/user/api/internal/config"
	"go-community/docker/mall/user/api/internal/middleware"
	user "go-community/docker/mall/user/rpc/userclient"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	UserRpc user.User
	Example rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: user.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Example: middleware.NewExampleMiddleware(c).Handle,
	}
}
