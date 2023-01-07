package config

import (
	jwtx "go-community/docker/mall/common/jwt"

	"github.com/zeromicro/go-zero/zrpc"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth jwtx.JwtAuth

	OrderRpc   zrpc.RpcClientConf
	ProductRpc zrpc.RpcClientConf
}
