package models

import (
	"context"
	"fmt"
	"testing"

	"gitlab.galaxy123.cloud/base/public/models/jwt"
)

func TestNewClaims(t *testing.T) {
	j := jwt.NewClaims("123")
	token, err := j.NewToken(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(token)
}
