package middleware

import (
	"fmt"
	"github.com/noahlsl/public/constants/consts"
	"github.com/noahlsl/public/constants/enums"
	"github.com/zeromicro/go-zero/core/stores/redis"
	xerrors "github.com/zeromicro/x/errors"
	xhttp "github.com/zeromicro/x/http"
	"net/http"
)

type UniAuth struct {
	c redis.RedisConf
}

func NewUniAuth(c redis.RedisConf) *UniAuth {
	return &UniAuth{c: c}
}

func (u *UniAuth) Handle(next http.HandlerFunc) http.HandlerFunc {
	rds := redis.MustNewRedis(u.c)
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "" {
			key := fmt.Sprintf(consts.RedisKeyAuth, token)
			exists, err := rds.Exists(key)
			if err != nil || !exists {
				xhttp.JsonBaseResponseCtx(r.Context(), w, xerrors.New(enums.ErrSysTokenExpired, "Login expired"))
				return
			}
		}
		next(w, r)
	}
}

func (u *UniAuth) UniAuthMiddleware(w http.ResponseWriter, r *http.Request) error {
	rds := redis.MustNewRedis(u.c)
	token := r.Header.Get("Authorization")
	if token != "" {
		key := fmt.Sprintf(consts.RedisKeyAuth, token)
		exists, err := rds.Exists(key)
		if err != nil || !exists {
			xhttp.JsonBaseResponseCtx(r.Context(), w, xerrors.New(enums.ErrSysTokenExpired, "Login expired"))
			return consts.ErrSysTokenExpired
		}
	}

	return nil
}
