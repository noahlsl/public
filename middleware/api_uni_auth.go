package middleware

import (
	"context"
	"fmt"
	"github.com/noahlsl/public/constants/consts"
	"github.com/noahlsl/public/constants/enums"
	"github.com/zeromicro/go-zero/core/stores/redis"
	xerrors "github.com/zeromicro/x/errors"
	xhttp "github.com/zeromicro/x/http"
	"net/http"
)

type UniAuth struct {
	r *redis.Redis
}

func NewUniAuth(r *redis.Redis) *UniAuth {
	return &UniAuth{r: r}
}

func (u *UniAuth) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "" {
			key := fmt.Sprintf(consts.RedisKeyAuth, token)
			v, err := u.r.Get(key)
			if err != nil || v == "" {
				xhttp.JsonBaseResponseCtx(r.Context(), w, xerrors.New(enums.ErrSysTokenExpired, "Login expired"))
				return
			}
			ctx := context.WithValue(r.Context(), "id", v)
			r = r.WithContext(ctx)
		} else {
			next(w, r)
		}
	}
}

func (u *UniAuth) UniAuthMiddleware(w http.ResponseWriter, r *http.Request) error {
	token := r.Header.Get("Authorization")
	if token != "" {
		key := fmt.Sprintf(consts.RedisKeyAuth, token)
		v, err := u.r.Get(key)
		if err != nil || v == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return consts.ErrSysTokenExpired
		}

		ctx := context.WithValue(r.Context(), "id", v)
		r = r.WithContext(ctx)
	}

	return nil
}
