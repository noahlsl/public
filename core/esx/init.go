package esx

import (
	"fmt"
	"github.com/olivere/elastic/v7"
)

func (c *Cfg) NewClient() *elastic.Client {
	var urls []string
	for _, addr := range c.Address {
		url := fmt.Sprintf("http://%s", addr)
		if c.TLS {
			url = fmt.Sprintf("https://%s", addr)
		}
		urls = append(urls, url)
	}

	var opts []elastic.ClientOptionFunc
	// 设置 Elasticsearch 服务器地址
	opts = append(opts, elastic.SetURL(urls...))
	opts = append(opts, elastic.SetSniff(false))       // 禁用 Sniff，防止在集群中查找其他节点
	opts = append(opts, elastic.SetHealthcheck(false)) // 可选，禁用健康检查)
	if c.Username != "" && c.Password != "" {
		// 选填: 设置用户名和密码
		opts = append(opts, elastic.SetBasicAuth(c.Username, c.Password))
	}
	client, err := elastic.NewClient(opts...)
	if err != nil {
		panic(err)
	}

	return client
}
