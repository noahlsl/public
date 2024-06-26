package serverx

import (
	"context"
	"fmt"
	json "github.com/bytedance/sonic"
	"github.com/ghodss/yaml"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"

	"github.com/noahlsl/public/constants/consts"
)

// Deprecated: Please Use AnyLoadYaml
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

func AnyLoadYaml[T any](cli *clientv3.Client, s string) T {

	var c T
	res, err := cli.Get(context.Background(), strings.ToUpper(s[:1])+s[1:])
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(res.Kvs[0].Value, &c)
	if err != nil {
		panic(err)
	}

	return c
}
