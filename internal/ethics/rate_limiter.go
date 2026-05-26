package ethics

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

var (
	limiters = make(map[string]*rate.Limiter)
	mu       sync.RWMutex
)

// Default limits (requests per second)
var defaultLimits = map[string]rate.Limit{
	"dns":    10.0,
	"whois":  1.0,
	"github": 2.0,
}

// Wait blocks until the collector is allowed to proceed based on its rate limit.
func Wait(collectorName string) error {
	l := getLimiter(collectorName)
	return l.Wait(context.Background())
}

func getLimiter(name string) *rate.Limiter {
	mu.RLock()
	l, ok := limiters[name]
	mu.RUnlock()

	if ok {
		return l
	}

	mu.Lock()
	defer mu.Unlock()

	// Double check
	if l, ok = limiters[name]; ok {
		return l
	}

	limit, ok := defaultLimits[name]
	if !ok {
		limit = 5.0 // Default for unknown collectors
	}

	// Burst size of 1 for strict enforcement
	l = rate.NewLimiter(limit, 1)
	limiters[name] = l
	return l
}

// SetLimit allows overriding limits at runtime (e.g. from config)
func SetLimit(name string, r float64) {
	mu.Lock()
	defer mu.Unlock()
	limiters[name] = rate.NewLimiter(rate.Limit(r), 1)
}
