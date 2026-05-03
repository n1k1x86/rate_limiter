package limiter

import (
	"sync"
	"time"
)

type RateLimiterRule struct {
	capacity     float64
	refillPerSec float64
}

func NewRateLimiterRule(capacity, refillPerSec int) (*RateLimiterRule, error) {
	if capacity <= 0 || refillPerSec <= 0 {
		return nil, ErrInvalidRuleParams
	}

	return &RateLimiterRule{
		capacity:     float64(capacity),
		refillPerSec: float64(refillPerSec),
	}, nil
}

type RateLimiter struct {
	mu           sync.Mutex
	capacity     float64
	refillPerSec float64
	lastTime     time.Time
	tokens       float64
}

func NewLimiter(rule *RateLimiterRule) (*RateLimiter, error) {
	if rule == nil {
		return nil, ErrRuleIsNil
	}

	return &RateLimiter{
		capacity:     rule.capacity,
		refillPerSec: rule.refillPerSec,
		lastTime:     time.Now(),
		tokens:       rule.capacity,
	}, nil
}

func (l *RateLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	duration := now.Sub(l.lastTime).Seconds()

	l.lastTime = now

	newTokens := duration * l.refillPerSec
	if newTokens+l.tokens > l.capacity {
		l.tokens = l.capacity
	} else {
		l.tokens = newTokens + l.tokens
	}

	if l.tokens < 1 {
		return false
	}

	l.tokens--
	return true
}
