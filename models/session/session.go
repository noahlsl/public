package session

import (
	"context"
	"fmt"
	"github.com/noahlsl/public/helper/structx"

	"github.com/noahlsl/public/helper/idx"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Session struct {
	r        *redis.Redis
	tokenKey string
	uidKey   string
	dataKey  string
	expired  int
}

func NewSession(r *redis.Redis, tk, uk, dk string, ex ...int) *Session {

	var expired = 3600 * 24 // 默认1天
	if len(ex) != 0 {
		expired = ex[0]
	}

	return &Session{
		r:        r,
		tokenKey: tk,
		uidKey:   uk,
		dataKey:  dk,
		expired:  expired,
	}
}

func (s *Session) Login(ctx context.Context, id interface{}, data interface{}) (string, error) {

	var token string
	// 单点登录
	key := fmt.Sprintf(s.uidKey, id)
	token, err := s.r.Get(key)
	if err != nil && !errors.Is(err, redis.Nil) {
		return token, errors.WithStack(err)
	}

	if token != "" {
		_, _ = s.r.Del(fmt.Sprintf(s.tokenKey, token))
	}

	token = idx.GetNanoId()
	err = s.r.Set(key, token)
	if err != nil {
		return token, errors.WithStack(err)
	}

	err = s.r.ExpireCtx(ctx, key, s.expired)
	if err != nil {
		return token, errors.WithStack(err)
	}

	// 生成Token锁
	key = fmt.Sprintf(s.tokenKey, token)
	err = s.r.SetCtx(ctx, key, fmt.Sprintf("%v", id))
	if err != nil {
		return token, errors.WithStack(err)
	}

	err = s.r.ExpireCtx(ctx, key, s.expired)
	if err != nil {
		return token, errors.WithStack(err)
	}

	key = fmt.Sprintf(s.dataKey, id)
	err = s.r.Set(key, structx.StructToStr(data))
	if err != nil {
		return token, errors.WithStack(err)
	}

	err = s.r.ExpireCtx(ctx, key, s.expired)
	if err != nil {
		return token, errors.WithStack(err)
	}

	return token, nil
}

func (s *Session) Logout(ctx context.Context, id interface{}) error {

	key := fmt.Sprintf(s.uidKey, id)
	result, _ := s.r.GetCtx(ctx, key)
	if result == "" {
		return nil
	}

	key = fmt.Sprintf(s.tokenKey, result)
	_, _ = s.r.DelCtx(ctx, key)

	return nil
}
