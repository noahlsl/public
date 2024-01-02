package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/noahlsl/public/constants/consts"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func BaseMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := map[string]interface{}{}
		lang := r.Header.Get("lang")
		if lang == "" {
			lang = consts.ZH
		}
		m["lang"] = lang
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
		lang := r.Header.Get("lang")
		if lang == "" {
			lang = consts.ZH
		}

		ctx = context.WithValue(ctx, "lang", lang)
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
