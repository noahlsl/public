package middleware

import (
	"net/http"

	"github.com/noahlsl/public/core/logsx"
	"github.com/zeromicro/go-zero/core/logx"
)

type LogMiddleware struct {
}

func NewLogMiddleware() *LogMiddleware {
	return &LogMiddleware{}
}

func (m *LogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		logx.Debugw("request", logsx.GetFields(r)...)
		next(w, r)
	}
}
func (m *LogMiddleware) OriginalHandle(_ http.ResponseWriter, r *http.Request) error {

	logx.Debugw("request", logsx.GetFields(r)...)
	return nil
}
