package ratelimiter

import (
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type TokenBucket struct {
	mu         sync.Mutex
	rate       float64
	burst      int
	tokens     float64
	lastRefill time.Time
	lastSeen   atomic.Int64
}

func NewTokenBucket(rate float64, burst int, now time.Time) *TokenBucket {
	bkt := &TokenBucket{
		rate:       rate,
		burst:      burst,
		tokens:     float64(burst),
		lastRefill: now,
	}
	bkt.lastSeen.Store((now.UnixNano()))

	return bkt
}

func (bkt *TokenBucket) Allow(now time.Time) (bool, time.Duration, int) {

	bkt.mu.Lock()
	defer bkt.mu.Unlock()

	bkt.lastSeen.Store(now.UnixNano())

	elapsed := now.Sub(bkt.lastRefill).Seconds()
	if elapsed > 0 {
		bkt.tokens += elapsed * bkt.rate

		if bkt.tokens > float64(bkt.burst) {
			bkt.tokens = float64(bkt.burst)
		}

		bkt.lastRefill = now
	}

	if bkt.tokens >= 1 {
		bkt.tokens--

		return true, 0, int(math.Floor(bkt.tokens))
	}

	wait := time.Duration(
		math.Ceil((1 - bkt.tokens) / bkt.rate * float64(time.Second)),
	)

	if wait < time.Millisecond {
		wait = time.Millisecond
	}

	return false, wait, 0
}
