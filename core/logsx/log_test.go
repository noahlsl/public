package logsx

import (
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestNewLogLogger(t *testing.T) {

	cfg := logx.LogConf{
		ServiceName: "test",
		Level:       "info",
	}
	err := A()
	MustSetUp(cfg)
	logx.Errorw(err.Error(), logx.Field("stack", GetStack(err)))
	logx.Info("test")
}

func C() error {
	_, err := os.Open("err")
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
func B() error {
	return C()
}
func A() error {
	return B()
}
