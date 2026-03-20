package collector_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/spectre/spectre/internal/collector/active"
	"github.com/spectre/spectre/internal/core"
)

// MockCollector to simulate workload
type MockDelayCollector struct {
	Delay time.Duration
}

func (m *MockDelayCollector) Name() string        { return "mock_delay" }
func (m *MockDelayCollector) Description() string { return "Mock collector for benchmarking" }
func (m *MockDelayCollector) IsActive() bool      { return true }
func (m *MockDelayCollector) Collect(caseID string, target string, options map[string]interface{}) ([]core.Evidence, error) {
	time.Sleep(m.Delay)
	return []core.Evidence{{
		CaseID:      caseID,
		Collector:   "mock_delay",
		FilePath:    "mock/path.json",
		FileHash:    "mock_hash",
		CollectedAt: time.Now(),
		Metadata: map[string]interface{}{
			"target": target,
		},
	}}, nil
}

// BenchmarkCollectorConcurrency benchmarks running multiple collectors concurrently.
func BenchmarkCollectorConcurrency(b *testing.B) {
	concurrencyLevels := []int{1, 5, 10, 50, 100}

	for _, limit := range concurrencyLevels {
		b.Run(fmt.Sprintf("Concurrency-%d", limit), func(b *testing.B) {
			collector := &MockDelayCollector{Delay: 10 * time.Millisecond}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var wg sync.WaitGroup
				sem := make(chan struct{}, limit)

				// Simulate 100 target collections per iteration
				targets := 100

				for j := 0; j < targets; j++ {
					wg.Add(1)
					go func(targetIdx int) {
						defer wg.Done()
						sem <- struct{}{} // Acquire

						_, _ = collector.Collect("test_case", fmt.Sprintf("target-%d", targetIdx), nil)

						<-sem // Release
					}(j)
				}
				wg.Wait()
			}
		})
	}
}

// BenchmarkHTTPCollector tests real-world HTTP concurrency
func BenchmarkHTTPCollector(b *testing.B) {
	collector := &active.HTTPCollector{}
	concurrencyLevels := []int{5, 20} // HTTP limit test

	for _, limit := range concurrencyLevels {
		b.Run(fmt.Sprintf("HTTP-Concurrency-%d", limit), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var wg sync.WaitGroup
				sem := make(chan struct{}, limit)

				targets := 20
				for j := 0; j < targets; j++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						sem <- struct{}{}
						_, _ = collector.Collect("test_case", "http://localhost", nil)
						<-sem
					}()
				}
				wg.Wait()
			}
		})
	}
}
