package redlockx

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	// RedLock parameters
	lockValue       = "1"
	lockTimeout     = 10 * time.Second
	renewalInterval = 2 * time.Second
	maxRenewalCount = 5 // 最大续约次数
	quorum          = 2 // Number of Redis instances required for quorum
)

// Redis connection pool
var redisClients []*redis.Client

// InitRedLockClients initializes Redis clients as singletons
func InitRedLockClients(addresses []string) {
	for _, addr := range addresses {
		client := redis.NewClient(&redis.Options{
			Addr: addr,
		})
		redisClients = append(redisClients, client)
	}
}

// RedLock structure
type RedLock struct {
	lockKey           string
	lockTimeout       time.Duration
	renewalStopChan   chan struct{}
	renewalStopSignal sync.WaitGroup
}

// NewRedLock creates a new RedLock instance with the provided lock parameters
func NewRedLock(lockKey string, timeout time.Duration) *RedLock {
	return &RedLock{
		lockKey:         lockKey,
		lockTimeout:     timeout,
		renewalStopChan: make(chan struct{}),
	}
}

// AcquireLock tries to acquire a lock with RedLock algorithm
func (r *RedLock) AcquireLock(ctx context.Context) (bool, error) {
	successfulLocks := 0
	for _, client := range redisClients {
		start := time.Now()
		// Try to set the lock
		ok, err := client.SetNX(ctx, r.lockKey, lockValue, r.lockTimeout).Result()
		if err != nil {
			return false, err
		}

		if ok {
			duration := time.Since(start)
			if duration < r.lockTimeout {
				successfulLocks++
			} else {
				client.Del(ctx, r.lockKey)
			}
		}
	}

	// Check if the number of successful locks is enough to form a quorum
	if successfulLocks >= quorum {
		r.renewalStopSignal.Add(1)
		go r.autoRenewLock(ctx)
		return true, nil
	}

	// If not enough quorum, release acquired locks
	for _, client := range redisClients {
		client.Del(ctx, r.lockKey)
	}

	return false, nil
}

// autoRenewLock handles automatic lock renewal
func (r *RedLock) autoRenewLock(ctx context.Context) {
	defer r.renewalStopSignal.Done()

	renewalCount := 0
	for {
		select {
		case <-r.renewalStopChan:
			// Stop renewal process if lock is released
			return
		case <-ctx.Done():
			// Stop renewal process if context is canceled
			return
		case <-time.After(renewalInterval):
			renewalCount++
			if renewalCount > maxRenewalCount {
				return
			}

			_, err := r.renewLock(ctx)
			if err != nil {
				log.Printf("Failed to renew lock: %v", err)
				return
			}
		}
	}
}

// renewLock attempts to renew the lock
func (r *RedLock) renewLock(ctx context.Context) (bool, error) {
	successfulLocks := 0
	for _, client := range redisClients {
		ok, err := client.Expire(ctx, r.lockKey, r.lockTimeout).Result()
		if err != nil {
			return false, err
		}

		if ok {
			successfulLocks++
		}
	}

	return successfulLocks >= quorum, nil
}

// ReleaseLock releases the lock
func (r *RedLock) ReleaseLock(ctx context.Context) error {
	close(r.renewalStopChan)   // Signal to stop renewal
	r.renewalStopSignal.Wait() // Wait for renewal goroutine to finish

	for _, client := range redisClients {
		_, err := client.Del(ctx, r.lockKey).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return err
		}
	}
	return nil
}
