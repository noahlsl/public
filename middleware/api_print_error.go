package middleware

import (
	json "github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"io"
	"net/http"
)

type Res struct {
	// Code represents the business code, not the http status code.
	Code int `json:"code" xml:"code"`
	// Msg represents the business message, if Code = BusinessCodeOK,
	// and Msg is empty, then the Msg will be set to BusinessMsgOk.
	Msg string `json:"msg" xml:"msg"`
}

// Custom response writer to capture response data
type responseCaptureWriter struct {
	http.ResponseWriter
	data []byte
}

// PrintErrMiddleware function to log response data
func PrintErrMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captureWriter := &responseCaptureWriter{ResponseWriter: w}
		next.ServeHTTP(captureWriter, r)
		var res Res
		_ = json.Unmarshal(captureWriter.data, &res)
		if res.Code != 0 {
			path := r.RequestURI
			var req string
			// 读取请求体中的JSON数据
			if r.Method == http.MethodGet {
				req = r.URL.Query().Encode()
			} else if r.Method == http.MethodPost {
				body, _ := io.ReadAll(r.Body)
				req = string(body)
			}

			logx.Error("请求错误", logx.Field("path", path),
				logx.Field("ip", httpx.GetRemoteAddr(r)),
				logx.Field("req", req))
		}
	})
}
