package rmqx

import (
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"strings"
	"time"
)

func (c *Cfg) NewConsumer(flag string) rocketmq.PushConsumer {

	// 设置推送消费者
	var opts = []consumer.Option{
		//消费组
		consumer.WithGroupName(flag + "-" + c.ConsumerGroup),
		consumer.WithInstance(flag + "-" + time.Now().Format(time.DateTime)),
	}
	address := strings.Split(c.Endpoint, ",")
	opts = append(opts, consumer.WithNameServer(address))
	if c.SecretKey != "" && c.AccessKey != "" {
		opts = append(opts, consumer.WithCredentials(primitive.Credentials{
			AccessKey: c.AccessKey,
			SecretKey: c.SecretKey,
		}))
	}

	// 设置推送消费者
	conn, err := rocketmq.NewPushConsumer(opts...)
	if err != nil {
		panic(err)
	}

	return conn
}

func (c *Cfg) NewPushConsumer(flag string) rocketmq.PushConsumer {
	return c.NewConsumer(flag)
}

func (c *Cfg) NewPullConsumer(flag string) rocketmq.PullConsumer {

	// 设置推送消费者
	var opts = []consumer.Option{
		//消费组
		consumer.WithGroupName(flag + "-" + c.ConsumerGroup),
		consumer.WithInstance(flag + "-" + time.Now().Format(time.DateTime)),
	}
	address := strings.Split(c.Endpoint, ",")
	opts = append(opts, consumer.WithNameServer(address))
	if c.SecretKey != "" && c.AccessKey != "" {
		opts = append(opts, consumer.WithCredentials(primitive.Credentials{
			AccessKey: c.AccessKey,
			SecretKey: c.SecretKey,
		}))
	}

	conn, err := rocketmq.NewPullConsumer(opts...)
	if err != nil {
		panic(err)
	}
	return conn
}

func (c *Cfg) NewProducer() rocketmq.Producer {

	var opts = []producer.Option{
		// 指定发送失败时的重试时间
		producer.WithRetry(2),
		// 设置 Group
		producer.WithGroupName(c.ProducerGroup),
	}
	address := strings.Split(c.Endpoint, ",")
	opts = append(opts, producer.WithNameServer(address))
	if c.AccessKey != "" && c.SecretKey != "" {
		opts = append(opts, producer.WithCredentials(primitive.Credentials{
			AccessKey: c.AccessKey,
			SecretKey: c.SecretKey,
		}))
	}

	p, err := rocketmq.NewProducer(opts...)
	if err != nil {
		panic(fmt.Sprintf("start producer error: %s", err.Error()))
	}

	// 开始连接
	err = p.Start()
	if err != nil {
		panic(fmt.Sprintf("start producer error: %s", err.Error()))
	}

	return p
}
