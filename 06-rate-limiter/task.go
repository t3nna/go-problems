package main

import (
	"sync"
	"time"
)

type client struct {
	count    int
	lastSeen time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*client
	limit   int
	window  time.Duration
}

func NewRateLimiter(n int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*client),
		limit:   n,
		window:  window,
	}

	go rl.cleanup()
	return rl
}

func (r *RateLimiter) CanTake(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	c, ok := r.clients[ip]
	if !ok || now.Sub(c.lastSeen) > r.window {
		r.clients[ip] = &client{count: 1, lastSeen: now}
		return true
	}

	// less then window
	if c.count >= r.limit {
		return false
	}
	c.count++
	return true
}

func (r *RateLimiter) cleanup() {
	for {
		time.Sleep(r.window)
		r.mu.Lock()
		for k, v := range r.clients {
			if time.Since(v.lastSeen) > r.window {
				delete(r.clients, k)
			}
		}
		r.mu.Unlock()
	}
}

func (r *RateLimiter) Take() {
}
