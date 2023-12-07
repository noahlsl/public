package middleware

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"net/http"

	"github.com/pkg/errors"
	"gitlab.galaxy123.cloud/base/public/constants/consts"
	"gitlab.galaxy123.cloud/base/public/constants/enums"
	"gitlab.galaxy123.cloud/base/public/helper/ipx"
	"gitlab.galaxy123.cloud/base/public/models/res"
)

type IPMiddleware struct {
	r   *redis.Redis
	key string
}

func NewIPMiddleware(r *redis.Redis, key string) *IPMiddleware {
	return &IPMiddleware{
		r:   r,
		key: key,
	}
}

func (l *IPMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		ipStr := ipx.RemoteIp(r)
		result, err := l.r.SismemberCtx(context.Background(), l.key, ipStr)
		if err != nil {
			rs := res.NewRes().WithCode(enums.ErrRequestLimit)
			_, _ = w.Write(rs.ToBytes())
			return
		}

		if !result {
			rs := res.NewRes().WithCode(enums.ErrRequestLimit)
			_, _ = w.Write(rs.ToBytes())
			return
		}

		next(w, r)
	}
}

func (l *IPMiddleware) OriginalHandle(_ http.ResponseWriter, r *http.Request) error {

	ipStr := ipx.RemoteIp(r)
	result, err := l.r.SismemberCtx(context.Background(), l.key, ipStr)
	if err != nil {
		return errors.WithMessage(consts.ErrRequestLimit, err.Error())
	}

	if !result {
		return consts.ErrRequestLimit
	}

	return nil
}
