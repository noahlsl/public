package middleware

import (
	"context"
	"fmt"
	json "github.com/bytedance/sonic"
	"net/http"

	"github.com/noahlsl/public/helper/bytex"
	"github.com/noahlsl/public/helper/ipx"
	"github.com/noahlsl/public/helper/strx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type BaseDataMiddleware struct {
	key   string
	flags []string
	r     *redis.Redis
}

func NewBaseDataMiddleware(r *redis.Redis, key string, flags ...string) *BaseDataMiddleware {
	return &BaseDataMiddleware{
		key:   key,
		r:     r,
		flags: flags,
	}
}

func (m *BaseDataMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		aid := r.Header.Get("aid")
		data := map[string]interface{}{}
		if aid != "" {
			get, err := m.r.GetCtx(context.Background(),
				fmt.Sprintf(m.key, aid))
			if err != nil {
				return
			}

			data = bytex.ToMap(strx.S2b(get))
		}
		for _, flag := range m.flags {
			data[flag] = r.Header.Get(flag)
		}
		data["ip"] = ipx.RemoteIp(r)
		data["path"] = r.URL.Path
		data["token"] = r.Header.Get("token")
		data["language"] = r.Header.Get("language")
		marshal, _ := json.Marshal(data)
		r.Header.Set("base", string(marshal))
		next(w, r)
	}
}

func (m *BaseDataMiddleware) OriginalHandle(_ http.ResponseWriter, r *http.Request) error {

	aid := r.Header.Get("aid")
	data := map[string]interface{}{}
	if aid != "" {
		get, err := m.r.GetCtx(context.Background(),
			fmt.Sprintf(m.key, aid))
		if err != nil {
			return err
		}

		data = bytex.ToMap(strx.S2b(get))
	}

	data["ip"] = ipx.RemoteIp(r)
	data["path"] = r.URL.Path
	data["token"] = r.Header.Get("token")
	data["language"] = r.Header.Get("language")
	for _, flag := range m.flags {
		data[flag] = r.Header.Get(flag)
	}
	marshal, _ := json.Marshal(data)
	r.Header.Set("base", string(marshal))
	return nil
}
