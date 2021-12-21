package valve

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// clients
type client struct {
	limit     *rate.Limiter
	timestamp time.Time
}

// RateLimiter
type Limiter struct {
	wg      *sync.WaitGroup
	mu      *sync.RWMutex
	rate    rate.Limit // rate
	burst   int        // burst
	ttl     time.Duration
	clients map[string]*client
	close   chan struct{}
}

// NewLimiter returns a new Limiter that allows events up to rate r and permits
// bursts of at most b tokens.
// if ttl equal or less then 1 minute, it will be set to 30 minutes
func NewLimiter(r float64, b int, ttl time.Duration) *Limiter {
	if ttl <= 1*time.Minute {
		ttl = 1 * time.Minute
	}

	limiter := &Limiter{
		wg:      &sync.WaitGroup{},
		mu:      &sync.RWMutex{},
		rate:    rate.Limit(r),
		burst:   b,
		ttl:     ttl,
		clients: map[string]*client{},
		close:   make(chan struct{}, 1),
	}

	go limiter.clean()

	return limiter
}

// clean
func (l *Limiter) clean() {
	tik := time.NewTicker(1 * time.Minute)
	//
	for {
		select {
		case <-l.close:
			l.wg.Done()
			tik.Stop()
			return
		case <-tik.C:
			l.mu.RLock()
			for ip, client := range l.clients {
				if time.Since(client.timestamp) > l.ttl {
					delete(l.clients, ip)
				}
			}
			l.mu.RUnlock()
		}
	}
}

// Allow
func (l *Limiter) Allow(key string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	//
	c, found := l.clients[key]
	// if not found set client ip
	// and return new limiter
	if !found {
		c = &client{
			limit:     rate.NewLimiter(l.rate, l.burst),
			timestamp: time.Now(),
		}
		l.clients[key] = c
	}

	return c.limit.Allow()
}

// Close
func (l *Limiter) Close() {
	l.close <- struct{}{}
	l.wg.Wait()
}
