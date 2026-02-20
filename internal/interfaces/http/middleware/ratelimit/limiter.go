package ratelimit

import (
	"sync"
	"time"
)

type rateLimiter struct {
	mu       sync.Mutex
	visits   map[string][]time.Time
	limit    int
	interval time.Duration
}

func newRateLimiter(limit int, interval time.Duration) *rateLimiter {
	return &rateLimiter{
		visits:   make(map[string][]time.Time),
		limit:    limit,
		interval: interval,
	}
}

func (rl *rateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	window := now.Add(-rl.interval)
	times := rl.visits[key]
	valid := times[:0]

	for _, t := range times {
		if t.After(window) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= rl.limit {
		rl.visits[key] = valid
		return false
	}

	valid = append(valid, now)
	rl.visits[key] = valid
	return true
}
