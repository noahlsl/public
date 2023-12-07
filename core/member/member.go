package member

import (
	"github.com/noahlsl/public/core/aorm"
)

type ModelMember struct {
	MemberId int64 `db:"member_id" json:"member_id"` //  MemberId
}

const (
	MemberId = "member_id" //MemberId
)

func (ModelMember) TableName() string {
	return "member"
}

type Member struct {
	Model ModelMember
	*aorm.Dataset[ModelMember]
}

func NewMember(client *aorm.Client) *Member {
	return &Member{
		Dataset: aorm.NewDataset[ModelMember](client),
	}
}
