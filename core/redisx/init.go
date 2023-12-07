package redisx

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"runtime"
	"time"
)

func (c *Cfg) NewRedis() *redis.Redis {

	if c.PoolSize == 0 {
		c.PoolSize = 4 * runtime.NumCPU()
	}
	cli, err := redis.NewRedis(redis.RedisConf{
		Host:        c.Address[0],
		Pass:        c.Password,
		Type:        c.Type,
		PingTimeout: 10 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	return cli
}

func (c *Cfg) NewRedisCacheConf() cache.CacheConf {

	return cache.CacheConf{
		cache.NodeConf{
			RedisConf: redis.RedisConf{
				Host: c.Address[0],
				Pass: c.Password,
				Type: redis.NodeType,
			},
			Weight: c.PoolSize,
		},
	}
}
func (c *Cfg) NewRedisConf() redis.RedisConf {

	return redis.RedisConf{
		Host:        c.Address[0],
		Pass:        c.Password,
		Type:        c.Type,
		PingTimeout: 10 * time.Second,
	}
}
