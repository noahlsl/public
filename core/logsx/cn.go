package logsx

import (
	"net/http"

	"github.com/rs/zerolog"
	"gitlab.galaxy123.cloud/base/public/helper/ipx"
)

type GinCn struct {
	r *http.Request
}

func NewGinCn(r *http.Request) *GinCn {
	return &GinCn{
		r: r,
	}
}

func (c *GinCn) Fn(e *zerolog.Event) {
	e.Str("trace", c.r.Header.Get("trace")).
		Str("path", c.r.URL.Path).
		Str("ip", ipx.RemoteIp(c.r)).
		Str("param", c.r.Header.Get("param")).
		Str("method", c.r.Method).
		Str("content_type", c.r.Header.Get("Content-Type"))
	return
}
