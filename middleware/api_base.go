package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/noahlsl/public/constants/consts"
	"github.com/ua-parser/uap-go/uaparser"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func BaseMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := map[string]interface{}{}
		lang := r.Header.Get("language")
		if lang == "" {
			lang = consts.ZH
		}
		m["language"] = lang
		m["ip"] = httpx.GetRemoteAddr(r)
		m["path"] = r.URL.Path
		driver := r.Header.Get("driver")
		if driver == "" {
			driver = "1"
		}
		m["driver"] = driver
		siteCode := r.Header.Get("website")
		if siteCode == "" {
			siteCode = "mall"
		}
		m["website"] = siteCode
		siteLang := r.Header.Get("lang")
		if siteLang == "" {
			siteLang = consts.JA
		}
		m["lang"] = siteLang
		m["debug"] = r.Header.Get("debug")
		marshal, _ := json.Marshal(m)
		r.Header.Set("base", string(marshal))
		next(w, r)
	}
}

func BaseCtxMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		lang := r.Header.Get("language")
		if lang == "" {
			lang = consts.ZH
		}

		ctx = context.WithValue(ctx, "language", lang)
		ctx = context.WithValue(ctx, "ip", httpx.GetRemoteAddr(r))
		ctx = context.WithValue(ctx, "driver", getDriver(r))
		website := r.Header.Get("website")
		if website == "" {
			website = "mall"
		}
		ctx = context.WithValue(ctx, "website", website)
		siteLang := r.Header.Get("lang")
		if siteLang == "" {
			siteLang = consts.JA
		}
		ctx = context.WithValue(ctx, "lang", siteLang)
		ctx = context.WithValue(ctx, "debug", r.Header.Get("debug"))
		ctx = context.WithValue(ctx, "path", r.URL.Path)
		r = r.WithContext(ctx)
		next(w, r)
	}
}

func getDriver(r *http.Request) string {
	ua := r.UserAgent()
	parser := uaparser.NewFromSaved()
	client := parser.Parse(ua)
	deviceFamily := client.Device.Family
	if deviceFamily == "Other" {
		return "1"
	} else {
		return "2"
	}
}
