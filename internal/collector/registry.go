package collector

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/spectre/spectre/internal/core"
	"github.com/spectre/spectre/internal/ethics"
	"github.com/spectre/spectre/internal/storage"
)

// Registry defines the contract for managing collectors.
type Registry interface {
	Register(c core.Collector) error
	Get(name string) (core.Collector, error)
	List() []core.Collector
}

// DefaultRegistry is the concrete implementation of the Registry interface.
type DefaultRegistry struct {
	collectors map[string]core.Collector
	mu         sync.RWMutex
}

// NewRegistry creates a new instance of DefaultRegistry.
func NewRegistry() *DefaultRegistry {
	return &DefaultRegistry{
		collectors: make(map[string]core.Collector),
	}
}

// Register adds a collector to the registry.
func (r *DefaultRegistry) Register(c core.Collector) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.collectors[c.Name()]; exists {
		return fmt.Errorf("collector '%s' already registered", c.Name())
	}
	r.collectors[c.Name()] = c
	log.Debug().Str("collector", c.Name()).Msg("collector registered")
	return nil
}

// Get retrieves a collector by name.
func (r *DefaultRegistry) Get(name string) (core.Collector, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.collectors[name]
	if !ok {
		return nil, fmt.Errorf("collector '%s' not found", name)
	}
	return c, nil
}

// List returns all registered collectors.
func (r *DefaultRegistry) List() []core.Collector {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var collectors []core.Collector
	for _, c := range r.collectors {
		collectors = append(collectors, c)
	}
	return collectors
}

// Global instance for convenience
var globalRegistry = NewRegistry()

// Register adds a collector to the global registry.
func Register(c core.Collector) {
	if err := globalRegistry.Register(c); err != nil {
		log.Warn().Err(err).Msg("failed to register collector")
	}
}

// Get retrieves a collector from the global registry.
func Get(name string) (core.Collector, error) {
	return globalRegistry.Get(name)
}

// List returns all collectors from the global registry.
func List() []core.Collector {
	return globalRegistry.List()
}

// Run executes a collector by name with ethics enforcement.
func Run(name string, caseID string, target string, activeAllowed bool, options map[string]interface{}) ([]core.Evidence, error) {
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

	return c.Collect(caseID, target, options)
}

// RunAndSave executes a collector and automatically persists evidence and ingests it.
func RunAndSave(name string, caseID string, target string, activeAllowed bool, options map[string]interface{}) ([]core.Evidence, error) {
	evidenceList, err := Run(name, caseID, target, activeAllowed, options)
	if err != nil {
		return nil, err
	}

	for i := range evidenceList {
		ev := &evidenceList[i]
		if err := storage.CreateEvidence(ev); err != nil {
			return evidenceList, fmt.Errorf("failed to save evidence: %w", err)
		}
		
		if err := storage.IngestEvidence(ev); err != nil {
			log.Warn().
				Err(err).
				Str("collector", ev.Collector).
				Msg("ingestion failed for evidence")
		}
	}

	return evidenceList, nil
}
