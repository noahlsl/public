package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/noahlsl/public/constants/consts"
	"github.com/noahlsl/public/constants/enums"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest/httpx"
	xerrors "github.com/zeromicro/x/errors"
	xhttp "github.com/zeromicro/x/http"
)

func BaseMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := map[string]interface{}{}
		language := r.Header.Get("language")
		if language == "" {
			language = consts.ZH
		}
		m["language"] = language
		m["ip"] = httpx.GetRemoteAddr(r)

		driver := r.Header.Get("driver")
		if driver == "" {
			driver = "1"
		}
		m["driver"] = driver
		siteCode := r.Header.Get("site_code")
		if siteCode == "" {
			siteCode = "mall"
		}
		m["site_code"] = siteCode
		m["debug"] = r.Header.Get("debug")
		marshal, _ := json.Marshal(m)
		r.Header.Set("base", string(marshal))
		next(w, r)
	}
}

func BaseCtxMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		language := r.Header.Get("language")
		if language == "" {
			language = consts.ZH
		}

		ctx = context.WithValue(ctx, "language", language)
		ctx = context.WithValue(ctx, "ip", httpx.GetRemoteAddr(r))

		driver := r.Header.Get("driver")
		if driver == "" {
			driver = "1"
		}
		ctx = context.WithValue(ctx, "driver", driver)
		website := r.Header.Get("website")
		if website == "" {
			website = "mall"
		}
		ctx = context.WithValue(ctx, "website", website)
		ctx = context.WithValue(ctx, "debug", r.Header.Get("debug"))
		r = r.WithContext(ctx)
		next(w, r)
	}
}

type UniAuth struct {
	c redis.RedisConf
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
