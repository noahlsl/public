package orm

import (
	"gitlab.galaxy123.cloud/base/public/core/aorm"
	"gitlab.galaxy123.cloud/base/public/core/member"
)

type Client struct {
	*aorm.Client
	Member *member.Member
}

func NewClient(dsn string, poolSize ...int) *Client {
	client := aorm.NewClient(dsn, poolSize...)
	return &Client{
		Client: client,
		Member: member.NewMember(client),
	}
}
