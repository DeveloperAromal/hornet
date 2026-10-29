package fraud

import (
	"sync"
	"time"
)

type FailureTracker struct {
	mu       sync.Mutex
	failures map[string][]time.Time
}

func NewFailureTracker() *FailureTracker {
	return &FailureTracker{
		failures: make(
			map[string][]time.Time,
		),
	}
}

func (track *FailureTracker) RecordFailure(identity string, now time.Time) int {
	track.mu.Lock()
	defer track.mu.Unlock()

	cutoff := now.Add(-time.Minute)

	previous := track.failures[identity]
	recent := make([]time.Time, 0, len(previous))

	for _, timestamp := range previous {
		if timestamp.After(cutoff) {
			recent = append(recent, timestamp)
		}
	}

	recent = append(recent, now)
	track.failures[identity] = recent

	return len(recent)
}
