package nacosx

import (
	"github.com/ghodss/yaml"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// Deprecated: Please Use LoadYaml
func Load[T any](cli config_client.IConfigClient, project string, groups ...string) T {

	group := "DEFAULT_GROUP"
	if len(groups) > 0 {
		group = groups[0]
	}
	content, err := cli.GetConfig(vo.ConfigParam{
		DataId: project,
		Group:  group,
	})

	if err != nil {
		panic(err)
	}
	var out T
	err = yaml.Unmarshal([]byte(content), &out)
	if err != nil {
		panic(err)
	}

	return out
}
