package limiter

import (
	"sync"
	"time"
)

var _ RateLimiter = (*TokenBucketLimiter)(nil)

type TokenBucketLimiter struct {
	tokens   uint64
	fillRate float64
	capacity uint64
	lastTime time.Time
	mutex    sync.Mutex
}

func NewTokenBucketLimiter(r float64, b uint64) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		tokens:   b,
		fillRate: r,
		capacity: b,
		lastTime: time.Now(),
	}
}

// Allow implements RateLimiter.
func (t *TokenBucketLimiter) Allow() bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	now := time.Now()

	timePassed := now.Sub(t.lastTime).Seconds()

	toAdd := uint64(timePassed * t.fillRate)

	if toAdd > 0 {
		t.tokens = min(t.capacity, t.tokens+toAdd)
		t.lastTime = now
	}

	if t.tokens > 0 {
		t.tokens--
		return true
	}

	return false
}
