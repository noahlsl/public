package nacosx

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type NcConf struct {
	Addr        string
	Port        uint64
	Group       string
	DataID      string
	ExtDataIDs  []string `json:",optional"`
	NamespaceID string
	Username    string
	Password    string
	LogLevel    string
}

func NewConf(address string) *NcConf {
	return &NcConf{
		Addr: address,
		Port: 8848,
	}
}

func (c *NcConf) WithPort(port uint64) *NcConf {
	c.Port = port
	return c
}

func (c *NcConf) WithUsername(username string) *NcConf {
	c.Username = username
	return c
}

func (c *NcConf) WithPassword(password string) *NcConf {
	c.Password = password
	return c
}

func (c *NcConf) NewConfigClient() config_client.IConfigClient {
	sc := []constant.ServerConfig{{
		IpAddr: c.Addr,
		Port:   c.Port,
	}}

	cc := &constant.ClientConfig{
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "nacos_log",
		CacheDir:            "nacos_cache",
		LogLevel:            "debug",
	}
	if c.Username != "" && c.Password != "" {
		cc.Username = c.Username
		cc.Password = c.Password
	}
	client, err := clients.NewConfigClient(vo.NacosClientParam{
		ServerConfigs: sc,
		ClientConfig:  cc,
	})
	if err != nil {
		panic(err)
	}

	return client
}
