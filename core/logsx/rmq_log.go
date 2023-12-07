package logsx

import (
	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap"
)

type RMQLogger struct {
	flag string
	lo   *zap.Logger
}

func NewRMQLogger(c logx.LogConf) *RMQLogger {
	c.ServiceName += "-RocketMQ"
	c.Level = "error"
	_, logger, err := newZapWriter(c)
	if err != nil {
		panic(err)
	}
	return &RMQLogger{lo: logger}
}
func (l *RMQLogger) Debug(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}
	l.lo.Debug(msg)
}
func (l *RMQLogger) Info(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}
	l.lo.Info(msg)
}

func (l *RMQLogger) Warning(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}
	l.lo.Warn(msg)
}

func (l *RMQLogger) Error(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}

	if msg == "get consumer list of group from broker error" {
		return
	}

	l.lo.Error(msg)
}

func (l *RMQLogger) Fatal(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}
	l.lo.Fatal(msg)
}
func (l *RMQLogger) Level(level string) {

}

func (l *RMQLogger) OutputPath(path string) error {
	return nil
}
