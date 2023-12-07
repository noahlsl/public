package logsx

import (
	"github.com/goccy/go-json"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/zero-contrib/logx/zerologx"
)

func Init(cfg interface{}) {

	c := &Cfg{}
	marshal, _ := json.Marshal(cfg)
	_ = json.Unmarshal(marshal, c)
	// 初始化日志依赖
	logger := c.NewZeroLogger()
	writer := zerologx.NewZeroLogWriter(logger)
	logx.DisableStat()
	logx.SetWriter(writer)
	logx.SetLevel(GetLogxLever(c.Level))
}
