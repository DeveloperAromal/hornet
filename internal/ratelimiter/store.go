package ratelimiter

import (
	"sync"
	"time"
)

type Store struct {
	mu sync.Mutex

	buckets map[string]*TokenBucket
	maxKeys int
	idleTTL time.Duration

	stop chan struct{}
	done chan struct{}
	once sync.Once
}

func NewStore(maxKeys int, idleTTL, cleanupInterval time.Duration) *Store {
	if maxKeys <= 0 {
		maxKeys = 100_000
	}
	if idleTTL <= 0 {
		idleTTL = 10 * time.Minute
	}
	if cleanupInterval <= 0 {
		cleanupInterval = time.Minute
	}

	s := &Store{
		buckets: make(map[string]*TokenBucket),
		maxKeys: maxKeys,
		idleTTL: idleTTL,
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
	}

	go s.cleanupLoop(cleanupInterval)
	return s
}

func (s *Store) getOrCreate(key string, rate float64, burst int, now time.Time) (*TokenBucket, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if bucket, ok := s.buckets[key]; ok {
		bucket.lastSeen.Store(now.UnixNano())
		return bucket, true
	}

	if len(s.buckets) >= s.maxKeys {
		return nil, false
	}

	bucket := NewTokenBucket(rate, burst, now)
	s.buckets[key] = bucket

	return bucket, true
}

func (s *Store) cleanupLoop(interval time.Duration) {
	defer close(s.done)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.removeExpired(time.Now())
		case <-s.stop:
			return
		}
	}
}

func (s *Store) removeExpired(now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key, bucket := range s.buckets {
		lastSeen := time.Unix(
			0,
			bucket.lastSeen.Load(),
		)

		if now.Sub(lastSeen) >= s.idleTTL {
			delete(s.buckets, key)
		}
	}
}

func (s *Store) Close() {
	s.once.Do(func() {
		close(s.stop)
		<-s.done
	})
}
