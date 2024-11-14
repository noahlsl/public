package nacosx

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	numx "github.com/noahlsl/public/helper/numberx"
	"github.com/noahlsl/public/utils"
	zeroConf "github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
)

var (
	configClient config_client.IConfigClient
	ncOnce       sync.Once
)

func (c *NcConf) InitConfigClient() (err error) {
	ncOnce.Do(func() {
		configClient, err = clients.NewConfigClient(
			vo.NacosClientParam{
				ClientConfig: &constant.ClientConfig{TimeoutMs: 5000, NamespaceId: c.NamespaceID},
				ServerConfigs: []constant.ServerConfig{
					{IpAddr: c.Addr, Port: c.Port},
				},
			},
		)
	})
	return
}

func (c *NcConf) GetConfig() (string, error) {
	var configMap = make(map[interface{}]interface{})
	mainConfig, err := configClient.GetConfig(vo.ConfigParam{DataId: c.DataID, Group: c.Group})
	if err != nil {
		return "", err
	}

	mainMap, err := utils.UnmarshalYamlToMap(mainConfig)
	if err != nil {
		return "", err
	}

	var extMap = make(map[interface{}]interface{})
	for _, dataID := range c.ExtDataIDs {
		extConfig, err := configClient.GetConfig(vo.ConfigParam{DataId: dataID, Group: c.Group})
		if err != nil {
			return "", err
		}

		tmpExtMap, err := utils.UnmarshalYamlToMap(extConfig)
		if err != nil {
			return "", err
		}

		extMap = utils.MergeMap(extMap, tmpExtMap)
	}

	configMap = utils.MergeMap(configMap, extMap)
	configMap = utils.MergeMap(configMap, mainMap)

	yamlString, err := utils.MarshalObjectToYamlString(configMap)
	if err != nil {
		return "", err
	}

	return yamlString, nil
}

func (c *NcConf) Listen(onChange func(string, string, string, string)) error {
	return configClient.ListenConfig(vo.ConfigParam{
		DataId:   c.DataID,
		Group:    c.Group,
		OnChange: onChange,
	})
}

func (c *NcConf) NewZRpcClient(serverName, clientName string) zrpc.Client {
	var target string
	target = fmt.Sprintf("nacos://%s:%d/%s?timeout=%s&namespace_id=%s&group_name=%s&app_name=%s", c.Addr, c.Port, serverName, "3s", c.NamespaceID, c.Group, clientName)
	return zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: target,
	})
}

func MustLoad(path string, v interface{}) *NcConf {
	var (
		err    error
		config string
	)

	var nc NcConf
	zeroConf.MustLoad(path, &nc, zeroConf.UseEnv())
	err = nc.InitConfigClient()
	if err != nil {
		log.Fatalf("init config client error: %v", err)
	}

	config, err = nc.GetConfig()
	if err != nil {
		log.Fatalf("get config error: %v", err)
	}

	err = zeroConf.LoadFromYamlBytes([]byte(config), v)
	if err != nil {
		log.Fatalf("load config error: %v", err)
	}
	return &nc
}

func MustRegister(nc *NcConf, rpcConfig *zrpc.RpcServerConf) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(nc.Addr, nc.Port),
	}

	cc := &constant.ClientConfig{
		NamespaceId:         nc.NamespaceID, // namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogLevel:            "info",
	}
	if nc.LogLevel != "" {
		cc.LogLevel = nc.LogLevel
	}

	// Create Nacos client
	ncClient, err := clients.NewNamingClient(vo.NacosClientParam{
		ServerConfigs: sc,
		ClientConfig:  cc,
	})
	if err != nil {
		log.Fatalf("Failed to create Nacos client: %v", err)
	}

	split := strings.Split(rpcConfig.ListenOn, ":")
	if len(split) != 2 {
		log.Fatalf("invalid rpc listen on: %s", rpcConfig.ListenOn)
	}
	// Prepare service registration options
	registerParam := vo.RegisterInstanceParam{
		Ip:          split[0],                  // IP Address where the service is running
		Port:        numx.Any2Uint64(split[1]), // Port for the service
		ServiceName: rpcConfig.Name,            // The name of the service
		Weight:      1.0,                       // Weight of the instance (default: 1)
	}

	// Register the service to Nacos
	_, err = ncClient.RegisterInstance(registerParam)
	if err != nil {
		log.Fatalf("Failed to register service to Nacos: %v", err)
	}
}
