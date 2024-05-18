package loggerx

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm/logger"
	"time"
)

type Logger struct {
	LogLevel logger.LogLevel
}

func New(level string) *Logger {
	l := new(Logger)
	switch level {
	case "debug": // 打印SQL
		l.LogMode(logger.Info)
	default:
		l.LogMode(logger.Error)
	}

	return l
}

func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l
}

func (l *Logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel < logger.Info {
		return
	}
	logx.WithContext(ctx).Debugf(msg, data)
}
func (l *Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel < logger.Warn {
		return
	}
	logx.WithContext(ctx).Infof(msg, data)
}

func (l *Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel < logger.Error {
		return
	}
	logx.WithContext(ctx).Errorf(msg, data)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	logx.WithContext(ctx).WithDuration(elapsed).Slowf("Trace sql: %v  row： %v  err: %v", sql, rows, err)
}
