package sqx

import (
	"fmt"
	"testing"
)

type Person struct {
	Name string `json:"name" db:"name"`
	Age  int    `json:"age" db:"age"`
	VV   int    `json:"vv" db:"-"`
	V2   string `json:"v2" db:"v2"`
}

func TestGetFields(t *testing.T) {
	f := GetFields(Person{})
	fmt.Println(f)
}
