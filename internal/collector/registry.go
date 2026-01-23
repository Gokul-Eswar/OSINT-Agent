package collector

import (
	"fmt"
	"sync"

	"github.com/spectre/spectre/internal/core"
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
