package rpcx

import (
	"context"
	"fmt"
	json "github.com/bytedance/sonic"
	"github.com/ghodss/yaml"
	"github.com/noahlsl/public/constants/consts"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
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
func LoadYaml(cli *clientv3.Client, s string) zrpc.RpcClientConf {

	c := &Cfg{}
	res, err := cli.Get(context.Background(), strings.ToUpper(s[:1])+s[1:])
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(res.Kvs[0].Value, &c)
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
