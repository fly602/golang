package config

import (
	jwtx "go-community/docker/mall/common/jwt"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf

	Auth jwtx.JwtAuth

	UserRpc zrpc.RpcClientConf
}
