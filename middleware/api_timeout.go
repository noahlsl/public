package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/noahlsl/public/constants/enums"
	xerror "github.com/zeromicro/x/errors"
	xhttp "github.com/zeromicro/x/http"
)

type TimeoutMiddleware struct {
	timeout time.Duration
}

func NewTimeoutMiddleware(n int) *TimeoutMiddleware {

	return &TimeoutMiddleware{
		timeout: time.Duration(n) * time.Second,
	}
}

func (m *TimeoutMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), m.timeout)
		defer func() {
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				xhttp.JsonBaseResponseCtx(r.Context(), w, xerror.New(enums.ErrTimeout, "Request Timeout"))
				return
			}

			cancel()
		}()
		r = r.WithContext(ctx)
		next(w, r)
	}
}
func (m *TimeoutMiddleware) OriginalHandle(w http.ResponseWriter, r *http.Request) error {

	ctx, cancel := context.WithTimeout(r.Context(), m.timeout)
	defer func() {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			xhttp.JsonBaseResponseCtx(r.Context(), w, xerror.New(enums.ErrTimeout, "Request Timeout"))
			return
		}

		cancel()
	}()
	r = r.WithContext(ctx)

	return nil
}
