package mongox

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func (c *Cfg) NewDatabase() *mongo.Database {
	//ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancelFunc()
	ctx := context.Background()
	var addr string
	if c.Username != "" && c.Password != "" {
		addr = fmt.Sprintf("mongodb://%s:%s@%s:%d", c.Username, c.Password, c.Host, c.Port)
	} else {
		addr = fmt.Sprintf("mongodb://%s:%d", c.Host, c.Port)
	}
	if c.Database != "" {
		addr += "/" + c.Database
	}
	clientOptions := options.Client().ApplyURI(addr)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database(c.Database)
}
