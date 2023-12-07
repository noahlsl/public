package lockx

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

var (
	key        = "public:lock:%v"
	expiration = 10 * time.Second
)

func SetLock(ctx context.Context, r *redis.Redis, val interface{}) bool {

	k := fmt.Sprintf(key, val)
	ok, err := r.SetnxCtx(ctx, key, "1")
	if err != nil || !ok {
		return false
	}

	go func(k string) {
		for {
			select {
			case <-ctx.Done():
				_, _ = r.Del(k)
				return

			case <-time.After(expiration):
				_, _ = r.Del(k)
				return
			}
		}
	}(k)

	return true
}

// TryLock 在指定的时间范围内尝试获取锁
// maxSecond 指定的时间秒数,默认3秒
func TryLock(ctx context.Context, r *redis.Redis, key interface{}, maxSecond ...int) bool {

	max := 3
	if len(maxSecond) == 1 {
		if maxSecond[0] > 1 {
			max = maxSecond[0]
		}
	}
	keyStr := fmt.Sprintf("lock:%v", key)
	for i := 0; i < max*2; i++ {
		ok, err := r.SetnxCtx(ctx, keyStr, "1")
		if err != nil {
			// 报错直接返回
			return false
		}

		if !ok {
			// 睡眠半秒再次尝试获取锁
			time.Sleep(500 * time.Millisecond)
			continue
		}

		// 设置过期时间2秒
		_ = r.Expire(keyStr, 2)
		go func() {
			// 2秒内主动监听退出信号主动释放锁
			for i := 0; i < 2000; i++ {
				select {
				case <-ctx.Done():
					return
				default:
					// 停顿1毫秒，出让CPU时间片
					time.Sleep(time.Millisecond)
				}
			}
			_, _ = r.Del(keyStr)
		}()
		return true
	}

	return false
}
