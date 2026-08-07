package tui

import (
	"testing"
)

func TestRenderASCIIGraph_EmptyCaseID(t *testing.T) {
	result := RenderASCIIGraph("")
	if result != "No case selected." {
		t.Errorf("expected 'No case selected.', got %q", result)
	}
}
