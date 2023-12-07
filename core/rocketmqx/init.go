package rocketmqx

import (
	"time"

	client "github.com/noahlsl/rocketmq"
	"github.com/noahlsl/rocketmq/credentials"
)

func (c *Cfg) SetConsumerGroup(s string) *Cfg {
	c.ConsumerGroup = s
	return c
}

func (c *Cfg) SetProducerGroup(s string) *Cfg {
	c.ProducerGroup = s
	return c
}

func (c *Cfg) NewConsumer(topic string) client.SimpleConsumer {

	client.ResetLogger()
	if c.AwaitDuration == 0 {
		c.AwaitDuration = 5
	}
	// new simpleConsumer instance
	consumer, err := client.NewSimpleConsumer(&client.Config{
		Endpoint:      c.Endpoint,
		ConsumerGroup: c.ConsumerGroup,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    c.AccessKey,
			AccessSecret: c.SecretKey,
		},
	},
		client.WithAwaitDuration(time.Duration(c.AwaitDuration)*time.Second),
		client.WithSubscriptionExpressions(map[string]*client.FilterExpression{
			topic: client.SUB_ALL,
		}),
	)

	if err != nil {
		panic(err)
	}

	return consumer
}

func (c *Cfg) NewProducer() client.Producer {

	producer, err := client.NewProducer(&client.Config{
		Endpoint: c.Endpoint,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    c.AccessKey,
			AccessSecret: c.SecretKey,
		},
	})
	if err != nil {
		panic(err)
	}

	// start producer
	err = producer.Start()
	if err != nil {
		panic(err)
	}

	return producer
}
