package rpcx

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/noahlsl/public/constants/consts"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func Load(cli *clientv3.Client, project, env, svc string) zrpc.RpcClientConf {

	c := &Cfg{}
	key := fmt.Sprintf(consts.CfgServer, project, env, svc)
	res, err := cli.Get(context.Background(), key)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(res.Kvs[0].Value, &c)
	if err != nil {
		panic(err)
	}

	return zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: c.Hosts,
			Key:   c.Key,
			User:  c.UserName,
			Pass:  c.PassWord,
		},
		Timeout: c.Timeout * 1000,
	}
}
