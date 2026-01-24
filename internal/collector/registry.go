package collector

import (
	"fmt"
	"sync"

	"github.com/spectre/spectre/internal/core"
	"github.com/spectre/spectre/internal/ethics"
)

var (
	registry = make(map[string]core.Collector)
	mu       sync.RWMutex
)

// Register adds a collector to the global registry.
func Register(c core.Collector) {
	mu.Lock()
	defer mu.Unlock()
	registry[c.Name()] = c
}

// Run executes a collector by name with ethics enforcement.
func Run(name string, caseID string, target string, activeAllowed bool) ([]core.Evidence, error) {
	c, err := Get(name)
	if err != nil {
		return nil, err
	}

	// 0. Active Consent Check
	if c.IsActive() && !activeAllowed {
		return nil, fmt.Errorf("collector '%s' is an ACTIVE probe. You must provide the --active flag to run it", name)
	}

	// 1. Scope Control
	allowed, err := ethics.IsAllowed(target)
	if !allowed {
		return nil, fmt.Errorf("safety block: %w", err)
	}

	// 2. Rate Limiting
	if err := ethics.Wait(name); err != nil {
		return nil, fmt.Errorf("rate limit error: %w", err)
	}

	return c.Collect(caseID, target)
}

// Get retrieves a collector by name.
func Get(name string) (core.Collector, error) {
	mu.RLock()
	defer mu.RUnlock()
	c, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("collector '%s' not found", name)
	}
	return c, nil
}

// List returns all registered collectors.
func List() []core.Collector {
	mu.RLock()
	defer mu.RUnlock()
	var collectors []core.Collector
	for _, c := range registry {
		collectors = append(collectors, c)
	}
	return collectors
}
