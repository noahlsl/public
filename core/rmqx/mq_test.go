package rmqx

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"log"
	"strconv"
	"testing"
	"time"
)

var (
	cfg = Cfg{
		Endpoint:      "159.138.47.60:8200,119.12.165.39:8200",
		AccessKey:     "local02",
		SecretKey:     "eff2wsZD#r83&X%KLbNA3jXX4VIpOiFG",
		ProducerGroup: "test1",
		ConsumerGroup: "test2",
	}
	topic = "test"
)

func TestNewProducer(t *testing.T) {
	conn := cfg.NewProducer()

	for i := 0; i < 10; i++ {
		// new a message
		msg := &primitive.Message{
			Topic: topic,
			Body:  []byte("this is a message : " + strconv.Itoa(i)),
		}
		// send message in sync
		resp, err := conn.SendSync(context.TODO(), msg)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp.Status)
		// wait a moment
		time.Sleep(time.Second * 1)
	}
	_ = conn.Shutdown()
}

func TestNewConsumer(t *testing.T) {

	conn := cfg.NewConsumer("")

	// 必须先在 开始前
	err := conn.Subscribe(topic, consumer.MessageSelector{}, func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range ext {
			fmt.Printf("subscribe callback:%v \n", ext[i])
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		panic(err)
	}
	err = conn.Start()
	if err != nil {
		panic(err)
	}
}

func TestNewPullConsumer(t *testing.T) {
	conn := cfg.NewPullConsumer("")

	result, err := conn.Pull(context.Background(), topic, 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	body := result.GetBody()
	fmt.Println(string(body))

}
