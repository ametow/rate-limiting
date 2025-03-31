package limiter

import (
	"container/list"
	"sync"
	"time"
)

type RateLimiter interface {
	Allow() bool
}

var _ RateLimiter = (*SlidingWindowLimiter)(nil)

type SlidingWindowLimiter struct {
	window int64
	limit  int
	logs   *list.List // deque // push_back, push_front -> in O(1)
	mutex  sync.Mutex
}

func NewSlidingWindowLimiter(window int64, limit int) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		window: window,
		limit:  limit,
		logs:   list.New(),
	}
}

func (this *SlidingWindowLimiter) Allow() bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	now := time.Now()
	delta := now.Unix() - this.window
	edgeTime := time.Unix(delta, 0)

	// Remove outdated logs
	for this.logs.Len() > 0 {
		front := this.logs.Front()
		if front.Value.(time.Time).Before(edgeTime) {
			this.logs.Remove(front)
		} else {
			break
		}
	}

	// Check if we can accept the request
	if this.logs.Len() < this.limit {
		this.logs.PushBack(now)
		return true
	}

	return false
}
