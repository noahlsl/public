package middleware

import (
	"context"
	"fmt"
	"github.com/noahlsl/public/constants/consts"
	"github.com/noahlsl/public/constants/enums"
	"github.com/noahlsl/public/helper/md5x"
	"github.com/noahlsl/public/models/res"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"net/http"
	"strings"
)

// IdempotenceMiddleware 幂等性中间件
type IdempotenceMiddleware struct {
	num    int
	r      *redis.Redis
	key    string
	filter []string
}

func NewIdempotenceMiddleware(r *redis.Redis, num int, filter ...string) *IdempotenceMiddleware {
	return &IdempotenceMiddleware{
		r:      r,
		num:    num,
		key:    "base:limit:idempotence:%v",
		filter: filter,
	}
}

func (m *IdempotenceMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		param := r.Header.Get("param")
		if param != "" && m.num != 0 {
			key := fmt.Sprintf(m.key, md5x.ByString(param, r.RequestURI, token))
			ok, err := m.r.SetnxExCtx(context.Background(), key, "0", m.num)
			if err != nil || !ok {
				rs := res.NewRes().WithCode(enums.ErrRequestLimit)
				_, _ = w.Write(rs.ToBytes())
				return
			}
		}

		next(w, r)
	}
}

func (m *IdempotenceMiddleware) OriginalHandle(_ http.ResponseWriter, r *http.Request) error {

	path := r.RequestURI
	for _, i := range m.filter {
		if strings.Contains(path, i) {
			return nil
		}
	}

	param := r.Header.Get("param")
	token := r.Header.Get("token")
	if param != "" && m.num != 0 {
		key := fmt.Sprintf(m.key, md5x.ByString(param, r.RequestURI, token))
		ok, err := m.r.SetnxExCtx(context.Background(), key, "0", m.num)
		if err != nil || !ok {
			return consts.ErrRequestLimit
		}
	}

	return nil
}
