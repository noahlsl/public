package middleware

import (
	"net/http"
)

func CorsMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		origin := r.Header.Get("Origin") //请求头部
		//接收客户端发送的origin （重要！）
		r.Header.Set("Access-Control-Allow-Origin", origin)
		//服务器支持的所有跨域请求的方法
		r.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE,Debug ,debug,*")
		//允许跨域设置可以返回其他子段，可以自定义字段
		r.Header.Set("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,Debug ,debug,*")
		// 允许浏览器（客户端）可以解析的头部 （重要）
		r.Header.Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Debug ,debug,*")
		//设置缓存时间
		r.Header.Set("Access-Control-Max-Age", "172800")
		//允许客户端传递校验信息比如 cookie (重要)
		r.Header.Set("Access-Control-Allow-Credentials", "true")
		// 允许放行OPTIONS请求
		if method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

func OriginalCorsMiddleware(w http.ResponseWriter, r *http.Request) error {

	method := r.Method
	r.Header.Set("Access-Control-Allow-Origin", "*")
	r.Header.Set("Access-Control-Allow-Headers", "*")
	r.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	r.Header.Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	r.Header.Set("Access-Control-Allow-Credentials", "true")
	// 允许放行OPTIONS请求
	if method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	return nil
}
