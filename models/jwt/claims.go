package jwt

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"time"

	"github.com/pkg/errors"
)

var (
	defaultSalt       = []byte("c400c2a5-7563-40b5-9780-aceAb9e753d8")
	defaultKey        = "sys:jwt:"
	defaultExpiration = 24 * time.Hour
)

type Claims struct {
	Id                 interface{} `json:"id"`
	r                  *redis.Redis
	private            `json:"-"`
	jwt.StandardClaims `json:"-"`
}

type private struct {
	salt []byte // Jwt salt
	key  string // 缓存Key
	exp  time.Duration
}

func NewClaims(id interface{}, inMap ...map[string]interface{}) *Claims {
	claims := make(jwt.MapClaims)
	for _, m := range inMap {
		for k, v := range m {
			claims[k] = v
		}
	}
	claims["exp"] = time.Now().Add(defaultExpiration).Unix() // 默认1天过期
	claims["iat"] = time.Now().Unix()
	return &Claims{
		Id: id,
		private: private{
			salt: defaultSalt,
			key:  defaultKey,
		},
	}
}
func (c *Claims) WithRedis(r *redis.Redis) *Claims {
	c.r = r
	return c
}
func (c *Claims) WithExpiration(in time.Duration) *Claims {
	c.private.exp = in
	return c
}
func (c *Claims) WithSalt(in interface{}) *Claims {
	if v, ok := in.([]byte); ok {
		c.private.salt = v
	} else {
		c.private.salt = []byte(fmt.Sprintf("%v", in))
	}
	return c
}
func (c *Claims) WithKey(in string) *Claims {
	c.private.key = in
	return c
}

// NewToken 只获取Token
func (c *Claims) NewToken(ctx context.Context) (string, error) {

	token, err := generateToken(c, c.private.salt)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if c.r == nil {
		return token, nil
	}

	key := fmt.Sprintf("%v%v", c.private.key, token)
	value := fmt.Sprintf("%v", c.Id)
	err = c.r.SetCtx(ctx, key, value)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return token, nil
}

func generateToken(claims *Claims, salt []byte) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(salt)

	return token, err
}
