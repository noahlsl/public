package esx

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"gitlab.galaxy123.cloud/base/public/constants/consts"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func Load(cli *clientv3.Client, project, env string) *Cfg {

	c := &Cfg{}
	key := fmt.Sprintf(consts.CfgEs, project, env)
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
