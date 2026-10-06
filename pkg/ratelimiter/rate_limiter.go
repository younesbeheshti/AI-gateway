package ratelimiter

import (
	"sync"
	"time"
)

type Config struct {
	Limit  int           `koanf:"limit"`
	Window time.Duration `koanf:"window"`
}

type RateLimit struct {
	Count       int
	WindowStart time.Time
}
type RateLimiter struct {
	mu      sync.Mutex
	clients map[int64]*RateLimit

	limit  int
	window time.Duration
}

func New(cfg Config) *RateLimiter {
	return &RateLimiter{
		clients: make(map[int64]*RateLimit),
		limit:   cfg.Limit,
		window:  cfg.Window,
	}
}

func (r *RateLimiter) Allow(userId int64) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	rate, exist := r.clients[userId]

	if !exist {
		r.clients[userId] = &RateLimit{
			Count:       1,
			WindowStart: now,
		}
		return true
	}

	if now.Sub(rate.WindowStart) >= r.window {
		rate.Count = 1
		rate.WindowStart = now
		return true
	}

	if rate.Count >= r.limit {
		return false
	}

	rate.Count++
	return true

}
