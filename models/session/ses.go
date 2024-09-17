package session

import (
	"context"
	"errors"
	"fmt"
	json "github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"time"

	"github.com/noahlsl/public/constants/consts"
	"github.com/noahlsl/public/helper/jwtx"
	"github.com/noahlsl/public/helper/structx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Ses struct {
	r         *redis.Redis
	expire    int64 // 过期时间秒
	maxLogins int   // 最大登录数量
}

func NewSes(r *redis.Redis, expire int64) *Ses {
	return &Ses{
		maxLogins: 5,
		r:         r,
		expire:    expire,
	}
}

func (s *Ses) WithMaxLogins(max int) *Ses {
	s.maxLogins = max
	return s
}

type LoginData struct {
	Token string `json:"token"`
	Data  any    `json:"data"`
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

	var val any
	if len(values) > 0 {
		val = values[0]
	} else {
		val = id
	}

	key = fmt.Sprintf(consts.RedisKeyUserLogin, id)
	score := time.Now().UnixMilli()
	value := LoginData{
		Token: token,
		Data:  val,
	}
	count, err := s.r.ZcardCtx(ctx, key)
	if err != nil {
		return "", err
	}

	if count >= s.maxLogins {
		var stop = count - s.maxLogins
		if stop > 0 {
			stop -= 1
		}
		result, err := s.r.ZrangeCtx(ctx, key, 0, int64(stop))
		if err != nil {
			return "", err
		}

		for _, v := range result {
			_, err = s.r.ZremCtx(ctx, key, v)
			if err != nil {
				return "", err
			}
		}
	}

	_, err = s.r.ZaddCtx(ctx, key, score, structx.StructToStr(value))
	if err != nil {
		return "", err
	}

	end := score - (s.expire * 1000)
	_, err = s.r.ZremrangebyscoreCtx(ctx, key, 0, end)
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
	key := fmt.Sprintf(consts.RedisKeyUserLogin, id)
	result, err := s.r.ZrangeCtx(ctx, key, 0, -1)
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}

	var item LoginData
	for _, value := range result {
		err = json.Unmarshal([]byte(value), &item)
		if err != nil {
			logx.Error(err)
			continue
		}

		tokenKey := fmt.Sprintf(consts.RedisKeyAuth, item.Token)
		_, err = s.r.DelCtx(ctx, tokenKey)
	}

	_, _ = s.r.DelCtx(ctx, key)

	return s.Login(ctx, secret, id, values...)
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

	key = fmt.Sprintf(consts.RedisKeyUserLogin, id)
	result, err := s.r.ZrangeCtx(ctx, key, 0, -1)
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	for _, v := range result {
		if strings.Contains(v, token) {
			_, _ = s.r.ZremCtx(ctx, key, v)
			break
		}
	}

	return nil
}

// KickOff 用户踢线,全部登录的地方都踢下线
func (s *Ses) KickOff(ctx context.Context, id string) error {
	key := fmt.Sprintf(consts.RedisKeyUserLogin, id)
	result, err := s.r.ZrangeCtx(ctx, key, 0, -1)
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	for _, v := range result {
		authKey := fmt.Sprintf(consts.RedisKeyAuth, v)
		_, err = s.r.DelCtx(ctx, authKey)
		if err != nil {
			logx.Error(err)
		}
	}

	_, err = s.r.DelCtx(ctx, key)

	return err
}
