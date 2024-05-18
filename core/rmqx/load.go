package rmqx

import (
	"context"
	"fmt"

	json "github.com/bytedance/sonic"
	"github.com/ghodss/yaml"
	"github.com/noahlsl/public/constants/consts"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// Deprecated: Please Use LoadYaml
func Load(cli *clientv3.Client, project, env string) *Cfg {

	c := &Cfg{}
	key := fmt.Sprintf(consts.CfgRocketMQ, project, env)
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

func LoadYaml(cli *clientv3.Client) *Cfg {

	c := &Cfg{}
	res, err := cli.Get(context.Background(), consts.ConfYamlRocketMQ)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(res.Kvs[0].Value, &c)
	if err != nil {
		panic(err)
	}

	return c
}
