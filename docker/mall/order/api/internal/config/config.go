package config

import (
	"github.com/zeromicro/go-zero/zrpc"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	OrderRpc   zrpc.RpcClientConf
	ProductRpc zrpc.RpcClientConf
}
