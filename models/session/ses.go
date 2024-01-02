package session

import (
	"context"
	"fmt"
	"github.com/noahlsl/public/helper/strx"
	"time"

	"github.com/noahlsl/public/constants/consts"
	"github.com/noahlsl/public/helper/jwtx"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	AccessExpire = int64(3600 * 24)
)

type Ses struct {
	r *redis.Redis
}

func NewSes(r *redis.Redis) *Ses {
	return &Ses{
		r: r,
	}
}

func (s *Ses) Login(ctx context.Context, secret string, id interface{}, ex ...int64) (string, error) {

	accessExpire := AccessExpire
	if len(ex) > 0 {
		accessExpire = ex[0]
	}
	var token string
	// 单点登录
	key := fmt.Sprintf(consts.RedisKeyUid, id)
	token, err := s.r.GetCtx(ctx, fmt.Sprintf(consts.RedisKeyUid, id))
	if err != nil && !errors.Is(err, redis.Nil) {
		return token, errors.WithStack(err)
	}

	if token != "" {
		_, _ = s.r.DelCtx(ctx, fmt.Sprintf(consts.RedisKeyAuth, token))
	}

	// 生成 jwt 响应
	now := time.Now().Unix()
	token, err = jwtx.GenJwtToken(secret, now, accessExpire, id)
	if err != nil {
		return "", err
	}

	key = fmt.Sprintf(consts.RedisKeyAuth, token)
	err = s.r.SetexCtx(ctx, key, strx.Any2Str(id), int(accessExpire))
	if err != nil {
		return "", err
	}

	key = fmt.Sprintf(consts.RedisKeyUid, id)
	err = s.r.SetexCtx(ctx, key, token, int(accessExpire))
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Ses) Logout(ctx context.Context, id interface{}) error {

	key := fmt.Sprintf(consts.RedisKeyUid, id)
	token, err := s.r.GetCtx(ctx, key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}

		return err
	}

	err = s.r.PipelinedCtx(ctx, func(pipe redis.Pipeliner) error {
		key = fmt.Sprintf(consts.RedisKeyAuth, token)
		err = pipe.Del(ctx, key).Err()
		if err != nil {
			return err
		}

		key = fmt.Sprintf(consts.RedisKeyUid, id)
		err = pipe.Del(ctx, key).Err()
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
