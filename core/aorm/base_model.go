package aorm

type BaseModel interface {
	TableName() string
}
type TxModel interface {
	SQL() string
}
