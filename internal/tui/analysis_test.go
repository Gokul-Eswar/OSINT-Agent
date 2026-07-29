package tui

import (
	"strings"
	"testing"

	"github.com/Gokul-Eswar/Spectre/internal/core"
)

func TestFormatAnalysis_WithAllSections(t *testing.T) {
	res := &core.Analysis{
		Findings:            []string{"open admin panel"},
		Risks:               []string{"sensitive data exposure"},
		NextSteps:           []string{"validate auth boundary"},
		MissingData:         []string{"WAF headers"},
		SuggestedCollectors: []string{"http"},
		Confidence:          0.87,
	}

	out := FormatAnalysis(res)
	checks := []string{
		"PRELIMINARY FINDINGS",
		"IDENTIFIED RISKS",
		"RECOMMENDED NEXT STEPS",
		"MISSING DATA",
		"SUGGESTED COLLECTORS",
		"Confidence Level: 0.87",
	}

	for _, check := range checks {
		if !strings.Contains(out, check) {
			t.Fatalf("expected output to contain %q, got %q", check, out)
		}
	}
}

func TestFormatAnalysis_NilResult(t *testing.T) {
	if out := FormatAnalysis(nil); out != "" {
		t.Fatalf("expected empty string for nil result, got %q", out)
	}
}
