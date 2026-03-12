package handler

import (
	"sync"
	"time"
)

type actionDedup struct {
	mu   sync.Mutex
	seen map[string]time.Time
}

func newActionDedup() *actionDedup {
	d := &actionDedup{seen: make(map[string]time.Time)}
	go d.cleanup()
	return d
}

// Allow returns true if the action should be counted (not seen within cooldown).
func (d *actionDedup) Allow(key string, cooldown time.Duration) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.seen[key]; ok && time.Since(t) < cooldown {
		return false
	}
	d.seen[key] = time.Now()
	return true
}

func (d *actionDedup) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		d.mu.Lock()
		now := time.Now()
		for k, t := range d.seen {
			if now.Sub(t) > 24*time.Hour {
				delete(d.seen, k)
			}
		}
		d.mu.Unlock()
	}
}
