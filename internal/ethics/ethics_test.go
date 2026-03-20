package ethics

import (
	"testing"
	"time"
)

func TestIsAllowed(t *testing.T) {
	tests := []struct {
		target  string
		allowed bool
	}{
		{"example.com", true},
		{"target.gov", false},
		{"secret.mil", false},
		{"localhost", false},
		{"127.0.0.1", false},
		{"192.168.1.1", true},
	}

	for _, tt := range tests {
		got, _ := IsAllowed(tt.target)
		if got != tt.allowed {
			t.Errorf("IsAllowed(%s) = %v, want %v", tt.target, got, tt.allowed)
		}
	}
}

func TestWhitelist(t *testing.T) {
	SetWhitelist([]string{"trusted.com"})
	defer SetWhitelist([]string{})

	tests := []struct {
		target  string
		allowed bool
	}{
		{"trusted.com", true},
		{"other.com", false},
	}

	for _, tt := range tests {
		got, _ := IsAllowed(tt.target)
		if got != tt.allowed {
			t.Errorf("With whitelist, IsAllowed(%s) = %v, want %v", tt.target, got, tt.allowed)
		}
	}
}

func TestRateLimiter(t *testing.T) {
	collector := "test_limiter"
	SetLimit(collector, 100.0) // 100 req/s

	start := time.Now()
	for i := 0; i < 10; i++ {
		err := Wait(collector)
		if err != nil {
			t.Fatal(err)
		}
	}
	duration := time.Since(start)

	// Should be very fast
	if duration > 500*time.Millisecond {
		t.Errorf("Rate limiter too slow: 10 requests took %v", duration)
	}

	// Strict limit
	SetLimit(collector, 1.0) // 1 req/s
	Wait(collector)          // First one is instant (burst 1)

	start = time.Now()
	Wait(collector) // Should wait ~1s
	duration = time.Since(start)

	if duration < 800*time.Millisecond {
		t.Errorf("Rate limiter too fast: 1 request took %v, want ~1s", duration)
	}
}
