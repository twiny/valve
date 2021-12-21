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
	r       rate.Limit // rate
	b       int        // burst
	clients map[string]*client
	close   chan struct{}
}

// NewLimiter returns a new Limiter that allows events up to rate r and permits
// bursts of at most b tokens.
func NewLimiter(r float64, b int) *Limiter {
	limiter := &Limiter{
		wg:      &sync.WaitGroup{},
		mu:      &sync.RWMutex{},
		r:       rate.Limit(r),
		b:       b,
		clients: map[string]*client{},
		close:   make(chan struct{}, 1),
	}

	go limiter.clean()

	return limiter
}

// clean
func (l *Limiter) clean() {
	ticker := time.NewTicker(1 * time.Minute)

	//
	for {
		select {
		case <-l.close:
			l.wg.Done()
			ticker.Stop()
			return
		case now := <-ticker.C:
			l.mu.RLock()
			for ip, client := range l.clients {
				if now.Sub(client.timestamp) > 10*time.Minute {
					delete(l.clients, ip)
				}
			}
			l.mu.RUnlock()
		}
	}
}

// Get
func (l *Limiter) Get(key string) *rate.Limiter {
	l.mu.RLock()
	defer l.mu.RUnlock()
	//
	c, found := l.clients[key]
	// if not found set client ip
	// and return new limiter
	if !found {
		c = &client{
			limit:     rate.NewLimiter(l.r, l.b),
			timestamp: time.Now(),
		}
		l.clients[key] = c
	}

	return c.limit
}

// Close
func (l *Limiter) Close() {
	l.close <- struct{}{}
	l.wg.Wait()
}
