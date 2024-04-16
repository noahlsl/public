package nacosx

import (
	"github.com/ghodss/yaml"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// Deprecated: Please Use LoadYaml
func Load[T any](cli config_client.IConfigClient, project string) T {

	content, err := cli.GetConfig(vo.ConfigParam{
		DataId: project,
		Group:  "DEFAULT_GROUP",
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
