package ratelimiter

import (
	"errors"
	"time"
)

var (
	ErrInvalidRate  = errors.New("rate must be greater than zero")
	ErrInvalidBurst = errors.New("burst must be at least one")
)

type Config struct {
	// Number of tokens refilled per second.
	RatePerSecond float64

	// Maximum number of tokens in each bucket.
	Burst int

	// Maximum number of distinct client keys.
	MaxKeys int

	// How long an inactive client bucket is retained.
	IdleTTL time.Duration

	// Frequency of inactive bucket cleanup.
	CleanupInterval time.Duration
}

type Decision struct {
	Allowed    bool
	RetryAfter time.Duration
	Remaining  int
}

type Limiter struct {
	rate  float64
	burst int
	store *Store

	now func() time.Time
}

func NewLimiter(cfg Config) (*Limiter, error) {
	if cfg.RatePerSecond <= 0 {
		return nil, ErrInvalidRate
	}
	if cfg.Burst < 1 {
		return nil, ErrInvalidBurst
	}

	return &Limiter{
		rate:  cfg.RatePerSecond,
		burst: cfg.Burst,
		store: NewStore(
			cfg.MaxKeys,
			cfg.IdleTTL,
			cfg.CleanupInterval,
		),
		now: time.Now,
	}, nil
}

func (l *Limiter) Allow(key string) Decision {
	if key == "" {
		return Decision{Allowed: false}
	}

	now := l.now()

	bucket, ok := l.store.getOrCreate(
		key,
		l.rate,
		l.burst,
		now,
	)

	if !ok {
		// Store capacity reached. Fail closed for new keys.
		return Decision{
			Allowed:    false,
			RetryAfter: time.Second,
		}
	}

	allowed, retryAfter, remaining := bucket.Allow(now)

	return Decision{
		Allowed:    allowed,
		RetryAfter: retryAfter,
		Remaining:  remaining,
	}
}

func (l *Limiter) Close() {
	l.store.Close()
}
