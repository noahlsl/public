package mongox

import (
	"context"
	"fmt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
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
	err := mgm.SetDefaultConfig(nil, c.Database, options.Client().ApplyURI("mongodb://root:12345@localhost:27017"))
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

func (c *Cfg) InitMgm() {
	var addr string
	if c.Username != "" && c.Password != "" {
		addr = fmt.Sprintf("mongodb://%s:%s@%s:%d", c.Username, c.Password, c.Host, c.Port)
	} else {
		addr = fmt.Sprintf("mongodb://%s:%d", c.Host, c.Port)
	}
	cfg := &mgm.Config{CtxTimeout: 10 * time.Second}
	if c.Timeout != 0 {
		cfg = &mgm.Config{CtxTimeout: time.Duration(c.Timeout) * time.Second}
	}
	err := mgm.SetDefaultConfig(cfg, c.Database, options.Client().ApplyURI(addr))
	if err != nil {
		log.Fatal(err)
	}
}
