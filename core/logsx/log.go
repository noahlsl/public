package logsx

import (
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/thinkeridea/go-extend/exbytes"
	"github.com/thinkeridea/go-extend/exstrings"
	"github.com/zeromicro/go-zero/core/logx"
)

func (c *Cfg) NewZeroLogger() zerolog.Logger {
	l := zerolog.New(os.Stderr).
		With().
		Timestamp().
		//Stack().
		Str("service", c.ServiceName).
		//Str("").
		CallerWithSkipFrameCount(2).
		Logger()

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return washPath(file) + ":" + strconv.Itoa(line)
	}

	// 设置日志等级
	//l = l.Level(GetZerologLever(c.Level)).Hook(&ZeroHook{})
	l = l.Level(GetZerologLever(c.Level))
	// 使用官方提供的，输出更友好
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	return l
}

// 路径脱敏
func washPath(s string) string {
	sb := exstrings.Bytes(s)
	path, _ := os.Getwd()
	pathByte := exbytes.Replace([]byte(path+"/"), []byte("\\"), []byte("/"), -1)
	root := os.Getenv("GOROOT")
	rootByte := exbytes.Replace([]byte(root+"/"), []byte("\\"), []byte("/"), -1)
	sb = exbytes.Replace(sb, pathByte, []byte(""), -1)
	sb = exbytes.Replace(sb, rootByte, []byte(""), -1)
	return exbytes.ToString(sb)
}

func GetLogxLever(lever string) uint32 {

	lever = strings.ToLower(lever)
	if lever == "debug" {
		return logx.DebugLevel

	} else if lever == "error" {
		return logx.ErrorLevel
	}

	return logx.InfoLevel
}

func GetZerologLever(lever string) zerolog.Level {

	lever = strings.ToLower(lever)
	if lever == "debug" {
		return zerolog.DebugLevel

	} else if lever == "error" {
		return zerolog.ErrorLevel
	}

	return zerolog.InfoLevel
}

func (c *Cfg) InitLogx() {
	logx.MustSetup(logx.LogConf{
		ServiceName: c.ServiceName,
		TimeFormat:  time.DateTime,
		Path:        c.Path,
		Level:       c.Level,
	})
}

func InitLogX(c logx.LogConf) {
	logx.MustSetup(logx.LogConf{
		ServiceName: c.ServiceName,
		TimeFormat:  time.DateTime,
		Path:        c.Path,
		Level:       c.Level,
		Stat:        c.Stat,
	})
}

func MustSetUp(c logx.LogConf, fields ...zap.Field) *zap.Logger {
	logx.MustSetup(c)
	writer, logger := NewZapWriter(c, fields...)
	logx.SetWriter(writer)
	logger.Level()
	return logger
}
