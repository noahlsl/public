package middleware

import (
	"errors"
	"fmt"
	"github.com/noahlsl/public/constants/consts"
	"github.com/noahlsl/public/constants/enums"
	"github.com/noahlsl/public/helper/strx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	xerrors "github.com/zeromicro/x/errors"
	xhttp "github.com/zeromicro/x/http"
	"net/http"
	"strings"
)

type PremAuth struct {
	r      *redis.Redis
	filter []string
}

func NewPremMiddleware(r *redis.Redis, filters ...string) *PremAuth {
	return &PremAuth{
		r:      r,
		filter: filters,
	}
}

func (m *PremAuth) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Path
		split := strings.Split(path, "?")
		if len(split) > 0 {
			path = split[0]
		}

		var pass bool
		if len(m.filter) > 0 {
			for _, f := range m.filter {
				if strings.Contains(path, f) {
					pass = true
					break
				}
			}
		}

		if !pass {
			id := strx.Any2Str(r.Header.Get("id"))
			val, err := m.r.Hget(consts.RedisKeyAdminRole, id)
			if err != nil {
				xhttp.JsonBaseResponseCtx(r.Context(), w, xerrors.New(enums.ErrSysAuthFailed, "Permission check failed"))
				return
			}

			key := fmt.Sprintf(consts.RedisKeyRolePrem, val)
			ok, err := m.r.Sismember(key, path)
			if err != nil {
				xhttp.JsonBaseResponseCtx(r.Context(), w, xerrors.New(enums.ErrSysAuthFailed, "Permission check failed"))
				return
			}
			if !ok {
				xhttp.JsonBaseResponseCtx(r.Context(), w, xerrors.New(enums.ErrSysAuthFailed, "Permission check failed"))
				return
			}
		}

		next(w, r)
	}
}

func (m *PremAuth) PremAuthMiddleware(_ http.ResponseWriter, r *http.Request) error {
	path := r.URL.Path
	split := strings.Split(path, "?")
	if len(split) > 0 {
		path = split[0]
	}

	if len(m.filter) > 0 {
		for _, f := range m.filter {
			if strings.Contains(path, f) {
				return nil
			}
		}
	}

	id := strx.Any2Str(r.Header.Get("id"))
	val, err := m.r.Hget(consts.RedisKeyAdminRole, id)
	if err != nil {
		return err
	}

	key := fmt.Sprintf(consts.RedisKeyRolePrem, val)
	ok, err := m.r.Sismember(key, path)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("permission check failed")
	}

	return nil
}
