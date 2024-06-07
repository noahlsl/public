package mongox

import (
	"context"
	"github.com/noahlsl/public/helper/idx"
	"testing"
	"time"
)

func TestMgo(t *testing.T) {
	cfg := &Cfg{
		Host:     "127.0.0.1",
		Port:     27017,
		Username: "msg",
		Password: "3QPEtoht7yXYRavp",
		Database: "msg",
	}
	c := cfg.NewDatabase()
	v := NewPrivateMessage(1, 3206523318469632, 1, "hello")
	one, err := c.Collection(cfg.Database).InsertOne(
		context.TODO(), v)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(one)
}

type PrivateMessage struct {
	ID         int64  `bson:"_id" json:"id"`
	SenderID   int64  `bson:"sender_id" json:"sender_id"`
	ReceiverID int64  `bson:"receiver_id" json:"receiver_id"`
	Mode       int64  `bson:"mode" json:"mode"`
	Content    string `bson:"content" json:"content"`
	CreatedAt  int64  `bson:"created_at" json:"created_at"`
}

func NewPrivateMessage(senderID, receiverID, mode int64, content string) *PrivateMessage {
	return &PrivateMessage{
		ID:         idx.GenSnId(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Mode:       mode,
		Content:    content,
		CreatedAt:  time.Now().UnixMilli(),
	}
}
