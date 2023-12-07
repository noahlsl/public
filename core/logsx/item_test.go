package logsx

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	err := A()
	stack := GetStack(err)
	fmt.Println(stack)
}
