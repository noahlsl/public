package middleware

import (
	"context"
	"net/http"

	"github.com/noahlsl/public/constants/consts"
	"github.com/noahlsl/public/constants/enums"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest/httpx"
	xerror "github.com/zeromicro/x/errors"
	xhttp "github.com/zeromicro/x/http"
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
		ipStr := httpx.GetRemoteAddr(r)
		result, err := l.r.SismemberCtx(context.Background(), l.key, ipStr)
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, xerror.New(enums.ErrRequestLimit, "IP limit"))
			return
		}

		if !result {
			xhttp.JsonBaseResponseCtx(r.Context(), w, xerror.New(enums.ErrRequestLimit, "IP limit"))
			return
		}

		next(w, r)
	}
}

func (l *IPMiddleware) OriginalHandle(_ http.ResponseWriter, r *http.Request) error {

	ipStr := httpx.GetRemoteAddr(r)
	result, err := l.r.SismemberCtx(context.Background(), l.key, ipStr)
	if err != nil {
		return errors.WithMessage(consts.ErrRequestLimit, err.Error())
	}

	if !result {
		return consts.ErrRequestLimit
	}

	return nil
}
