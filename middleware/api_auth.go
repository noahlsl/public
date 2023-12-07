package middleware

import (
	"gitlab.galaxy123.cloud/base/public/constants/enums"
	"gitlab.galaxy123.cloud/base/public/models/res"
	"net/http"

	"github.com/zeromicro/go-zero/rest/handler"
)

type AuthMiddleware struct {
	secret string
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{
		secret: secret,
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			rs := res.NewRes().WithCode(enums.ErrSysTokenExpired)
			_, _ = w.Write(rs.ToBytes())
			return
		}

		authHandler := handler.Authorize(m.secret)
		authHandler(next).ServeHTTP(w, r)
		return
	}
}
