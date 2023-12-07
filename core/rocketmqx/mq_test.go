package rocketmqx

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/noahlsl/rocketmq"
)

var (
	cfg = Cfg{
		Endpoint:      "127.0.0.1:10911",
		AccessKey:     "",
		SecretKey:     "",
		ProducerGroup: "test1",
		ConsumerGroup: "test2",
	}
	topic = "test01"
)

func TestNewProducer(t *testing.T) {
	producer := cfg.NewProducer()
	// start producer
	err := producer.Start()
	if err != nil {
		log.Fatal(err)
	}
	// gracefule stop producer

	for i := 0; i < 10; i++ {
		// new a message
		msg := &golang.Message{
			Topic: topic,
			Body:  []byte("this is a message : " + strconv.Itoa(i)),
		}
		// set keys and tag
		msg.SetKeys("a", "b")
		msg.SetTag("ab")
		// send message in sync
		resp, err := producer.Send(context.TODO(), msg)
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < len(resp); i++ {
			fmt.Printf("%#v\n", resp[i])
		}
		// wait a moment
		time.Sleep(time.Second * 1)
	}
	_ = producer.GracefulStop()
}

func TestNewConsumer(t *testing.T) {

}
