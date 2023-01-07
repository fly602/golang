package svc

import (
	jwtx "go-community/docker/mall/common/jwt"
	"go-community/docker/mall/user/api/internal/config"
	user "go-community/docker/mall/user/rpc/userclient"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	UserRpc   user.User
	JwtHeader rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		UserRpc:   user.NewUser(zrpc.MustNewClient(c.UserRpc)),
		JwtHeader: jwtx.NewJwtheaderMiddleware(c.Auth).Handle,
	}
}
