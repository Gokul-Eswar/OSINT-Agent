package collector

import (
	"fmt"
	"sync"

	"github.com/spectre/spectre/internal/core"
	"github.com/spectre/spectre/internal/ethics"
	"github.com/spectre/spectre/internal/storage"
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

// RunAndSave executes a collector and automatically persists evidence and ingests it.
func RunAndSave(name string, caseID string, target string, activeAllowed bool) ([]core.Evidence, error) {
	evidenceList, err := Run(name, caseID, target, activeAllowed)
	if err != nil {
		return nil, err
	}

	for i := range evidenceList {
		ev := &evidenceList[i]
		if err := storage.CreateEvidence(ev); err != nil {
			return evidenceList, fmt.Errorf("failed to save evidence: %w", err)
		}
		
		if err := storage.IngestEvidence(ev); err != nil {
			// We log ingestion errors but don't necessarily stop the whole process, 
			// though for "rigorous" we might want to know.
			fmt.Printf("Warning: Ingestion failed for %s: %v\n", ev.Collector, err)
		}
	}

	return evidenceList, nil
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
