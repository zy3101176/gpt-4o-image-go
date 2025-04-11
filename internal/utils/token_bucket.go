package utils

import (
	"context"
	"sync"
	"time"
)

// TokenBucket 实现令牌桶算法用于限流
type TokenBucket struct {
	rateLimit   time.Duration
	lastRequest time.Time
	mu          sync.Mutex
}

// NewTokenBucket 创建新的令牌桶
func NewTokenBucket(rateLimit time.Duration) *TokenBucket {
	return &TokenBucket{
		rateLimit:   rateLimit,
		lastRequest: time.Now(),
	}
}

// Wait 等待直到可以发送请求
func (tb *TokenBucket) Wait(ctx context.Context) error {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRequest)

	if elapsed < tb.rateLimit {
		waitTime := tb.rateLimit - elapsed
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitTime):
			tb.lastRequest = time.Now()
			return nil
		}
	}

	tb.lastRequest = now
	return nil
}
