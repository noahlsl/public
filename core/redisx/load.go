package redisx

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/noahlsl/public/constants/consts"
	"github.com/zeromicro/go-zero/core/stores/redis"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func Load(cli *clientv3.Client, project, env string) redis.RedisConf {

	c := redis.RedisConf{}
	key := fmt.Sprintf(consts.CfgRedis, project, env)
	res, err := cli.Get(context.Background(), key)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(res.Kvs[0].Value, &c)
	if err != nil {
		panic(err)
	}

	return c
}
