package redisx

import (
	"context"
	"fmt"

	json "github.com/bytedance/sonic"
	"github.com/ghodss/yaml"
	"github.com/noahlsl/public/constants/consts"
	"github.com/zeromicro/go-zero/core/stores/redis"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// Deprecated: Please Use LoadYaml
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
func LoadYaml(cli *clientv3.Client) redis.RedisConf {

	c := redis.RedisConf{}
	res, err := cli.Get(context.Background(), consts.ConfYamlRedis)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(res.Kvs[0].Value, &c)
	if err != nil {
		panic(err)
	}

	return c
}
