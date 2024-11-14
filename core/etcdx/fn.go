package etcdx

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v3"
)

func MustEtcdClient(address string) *clientv3.Client {
	// Initialize etcd client
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{address}, // ETCD 地址
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	return client
}

// MustLoadConfig 加载并返回配置，传入一个空的配置结构体，类型可以动态传入
func MustLoadConfig(addr, key string, config interface{}) {
	client := MustEtcdClient(addr)
	defer client.Close()
	// 从 ETCD 获取配置
	resp, err := client.Get(context.Background(), key)
	if err != nil {
		panic(err)
	}

	if len(resp.Kvs) == 0 {
		log.Fatalf("key %s not found", key)
	}

	// 反序列化配置内容到目标配置结构体
	err = yaml.Unmarshal(resp.Kvs[0].Value, config)
	if err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	go WatchAndUpdate(client, key, config)
}

// WatchAndUpdate 启动监听，自动更新配置
func WatchAndUpdate(client *clientv3.Client, key string, config interface{}) {
	var lock sync.Mutex
	go func() {
		// 监听配置变更
		watchCh := client.Watch(context.Background(), key)

		for wr := range watchCh {
			for _, ev := range wr.Events {
				// 如果 key 发生变化，重新加载配置
				if string(ev.Kv.Key) == key {
					lock.Lock() // 锁住配置更新，避免并发修改
					err := yaml.Unmarshal(ev.Kv.Value, config)
					if err != nil {
						logx.Error(err)
					}
					lock.Unlock() // 解锁
				}
			}
		}
	}()
}
