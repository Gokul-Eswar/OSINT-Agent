package main

import (
	"testing"
)

func TestRequiredMean(t *testing.T) {
	tests := []struct {
		baseline          float64
		minImprovementPct float64
		want              float64
	}{
		{100.0, 10.0, 90.0},
		{100.0, 0.0, 100.0},
		{200.0, 50.0, 100.0},
		{50.0, 20.0, 40.0},
	}

	for _, tt := range tests {
		got := requiredMean(tt.baseline, tt.minImprovementPct)
		if got != tt.want {
			t.Errorf("requiredMean(%v, %v) = %v, want %v", tt.baseline, tt.minImprovementPct, got, tt.want)
		}
	}
}
