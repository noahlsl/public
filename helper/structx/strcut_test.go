package structx

import (
	"fmt"
	"testing"
)

type Person struct {
	Name string
}

func TestStructCause(t *testing.T) {
	p := &Person{"John"}
	cause := StructCause(p)
	fmt.Println(cause)

	p1 := Person{"John"}
	cause = StructCause(p1)
	fmt.Println(cause)

	cause = StructCause(1)
	fmt.Println(cause)

	cause = StructCause("p")
	fmt.Println(cause)

	cause = StructCause([]string{"1", "2"})
	fmt.Println(cause)

	cause = StructCause(struct {
		Data map[string]interface{}
	}{Data: map[string]interface{}{"1": "2"}})
	fmt.Println(cause)
}

func TestMapToStruct(t *testing.T) {
	m := map[string]string{
		"name": "123",
	}
	toStruct, err := MapToStruct[Person](m)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(toStruct)
}
