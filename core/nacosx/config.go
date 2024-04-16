package nacosx

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type Conf struct {
	Host     string
	Port     uint64
	Username string
	Password string
}

func NewConf(address string) *Conf {
	return &Conf{
		Host: address,
		Port: 8848,
	}
}

func (c *Conf) WithPort(port uint64) *Conf {
	c.Port = port
	return c
}

func (c *Conf) WithUsername(username string) *Conf {
	c.Username = username
	return c
}

func (c *Conf) WithPassword(password string) *Conf {
	c.Password = password
	return c
}

func (c *Conf) NewConfigClient() config_client.IConfigClient {
	sc := []constant.ServerConfig{{
		IpAddr: c.Host,
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
