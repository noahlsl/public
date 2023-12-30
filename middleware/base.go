package middleware

import (
	"net/http"
	"strings"

	"github.com/noahlsl/public/constants/consts"
	"github.com/noahlsl/public/helper/ipx"
	"github.com/goccy/go-json"
)

func BaseMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := map[string]interface{}{}
		language := r.Header.Get("language")
		if language == "" {
			language = consts.ZH
		}
		m["language"] = language

		var ip string
		ips := strings.Split(ipx.RemoteIp(r), ",")
		if len(ips) > 0 {
			ip = ips[0]
		}
		m["ip"] = ip

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
