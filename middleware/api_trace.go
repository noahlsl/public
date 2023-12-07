package middleware

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/trace"
	"gitlab.galaxy123.cloud/base/public/helper/idx"
)

func TraceMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		traceId := r.Header.Get(trace.TraceIdKey)
		if traceId == "" {
			traceId = idx.GetNanoId()
		}

		//traceId = trace.TraceIDFromContext(ctx)
		r.Header.Set(trace.TraceIdKey, traceId)
		ctx = context.WithValue(ctx, trace.TraceIdKey, traceId)
		r = r.WithContext(ctx)
		next(w, r)
	}
}

func OriginalTraceMiddleware(_ http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	traceId := r.Header.Get(trace.TraceIdKey)
	if traceId == "" {
		traceId = idx.GetNanoId()
	}

	//traceId = trace.TraceIDFromContext(ctx)
	r.Header.Set(trace.TraceIdKey, traceId)
	ctx = context.WithValue(ctx, trace.TraceIdKey, traceId)
	r = r.WithContext(ctx)

	return nil
}
