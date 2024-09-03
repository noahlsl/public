package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/noahlsl/public/constants/consts"
	"github.com/noahlsl/public/helper/jwtx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Ses struct {
	r      *redis.Redis
	expire int64 // 过期时间秒
}

func NewSes(r *redis.Redis, expire int64) *Ses {
	return &Ses{
		r:      r,
		expire: expire,
	}
}

// Login 多点登录
// values 选填参数填入登录信息,设备信息等的JSON字符串
func (s *Ses) Login(ctx context.Context, secret string, id string, values ...any) (string, error) {

	// 生成 jwt 响应
	now := time.Now().Unix()
	token, err := jwtx.GenJwtToken(secret, now, s.expire, id)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf(consts.RedisKeyAuth, token)
	err = s.r.SetexCtx(ctx, key, id, int(s.expire))
	if err != nil {
		return "", err
	}

	key = fmt.Sprintf(consts.RedisKeyUidMap, id)
	var val = id
	if len(values) > 0 {
		if v, ok := values[0].(string); ok {
			val = v
		}
	}

	err = s.r.HsetCtx(ctx, key, token, val)
	if err != nil {
		return "", err
	}

	_ = s.r.ExpireCtx(ctx, key, int(s.expire))
	return token, nil
}

// LoginOnce 单点登录
// values 选填参数填入登录信息,设备信息等的JSON字符串
func (s *Ses) LoginOnce(ctx context.Context, secret string, id string, values ...any) (string, error) {
	// 单点登录
	key := fmt.Sprintf(consts.RedisKeyUidMap, id)
	_, _ = s.r.DelCtx(ctx, key)

	// 生成 jwt 响应
	now := time.Now().Unix()
	token, err := jwtx.GenJwtToken(secret, now, s.expire, id)
	if err != nil {
		return "", err
	}

	tokenKey := fmt.Sprintf(consts.RedisKeyAuth, token)
	err = s.r.SetexCtx(ctx, tokenKey, id, int(s.expire))
	if err != nil {
		return "", err
	}

	var val = id
	if len(values) > 0 {
		if v, ok := values[0].(string); ok {
			val = v
		}
	}

	err = s.r.HsetCtx(ctx, key, token, val)
	if err != nil {
		return "", err
	}

	_ = s.r.ExpireCtx(ctx, key, int(s.expire))
	return token, nil
}

func (s *Ses) Logout(ctx context.Context, token string) error {
	key := fmt.Sprintf(consts.RedisKeyAuth, token)
	id, err := s.r.GetCtx(ctx, key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}

		return err
	}

	_, err = s.r.DelCtx(ctx, key)
	if err != nil {
		return err
	}

	key = fmt.Sprintf(consts.RedisKeyUidMap, id)
	_, _ = s.r.HdelCtx(ctx, key, token)
	return nil
}
