package redlockx

import (
	"context"
	"testing"
	"time"
)

func TestRedLock(t *testing.T) {
	addresses := []string{"localhost:6379", "localhost:6380", "localhost:6381", "localhost:6382", "localhost:6383"}
	InitRedLockClients(addresses)

	lockKey := "myLock"
	redLock := NewRedLock(lockKey, lockTimeout)

	// Create a context with timeout
	lockCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Try to acquire the lock
	ok, err := redLock.AcquireLock(lockCtx)
	if err != nil {
		t.Fatalf("Failed to acquire lock: %v", err)
	}

	if !ok {
		t.Log("Failed to acquire lock")
	}

	t.Log("Lock acquired")
	// Perform critical section operations here

	// Release the lock
	err = redLock.ReleaseLock(lockCtx)
	if err != nil {
		t.Fatalf("Failed to release lock: %v", err)
	}
	t.Log("Lock released")
}
