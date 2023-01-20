package svc

import (
	"go-community/docker/mall/user/model"
	"go-community/docker/mall/user/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(conn, c.CacheRedis),
	}
}

func (c *ServiceContext) Init() {
	// 初始化服务
	// 如果开启Auth 验证，则初始化Token
	key := c.Config.Redis.Key
	cache := c.Config.Redis.NewRedis()
	for _, kv := range c.Config.Tokens {
		value, err := cache.Hget(key, kv.Key)
		if err != nil {
			cache.Hset(key, kv.Key, kv.Value)
		} else {
			if value != kv.Value {
				cache.Hset(key, kv.Key, kv.Value)
			}
		}
	}
}
