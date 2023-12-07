package serverx

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	clientv3 "go.etcd.io/etcd/client/v3"

	"gitlab.galaxy123.cloud/base/public/constants/consts"
)

func AnyLoad[T any](cli *clientv3.Client, project, env, svc string) T {

	var c T
	key := fmt.Sprintf(consts.CfgServer, project, env, svc)
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
