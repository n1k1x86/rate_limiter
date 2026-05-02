package main

import (
	"fmt"
	"time"
)

// попробуем реализацию Token Bucket
// пусть есть размер бакета - size, на сколько в секунду пополнять бакет - burst, сколько токенов свободно - tokens
// при обращении клиента забираем токен + с учетом прошедшего времени пополняем перед работой

type limiter struct {
	size     float64
	burst    float64
	lastTime time.Time
	tokens   float64
}

func newLimiter(size int, burst int) *limiter {
	return &limiter{
		size:     float64(size),
		burst:    float64(burst),
		lastTime: time.Now(),
		tokens:   float64(size),
	}
}

func (l *limiter) allow() bool {
	now := time.Now()
	duration := time.Since(l.lastTime).Seconds()

	l.lastTime = now

	newTokens := duration * l.burst
	if newTokens+float64(l.tokens) > l.size {
		l.tokens = l.size
	} else {
		l.tokens = newTokens + l.tokens
	}

	if l.tokens-1 < 0 {
		return false
	}
	l.tokens--
	return true
}

func main() {
	l := newLimiter(10, 10)

	for range 300 {
		fmt.Println(l.allow())
	}
}
