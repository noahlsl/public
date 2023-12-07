package redisx

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gitlab.galaxy123.cloud/base/public/constants/consts"
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
