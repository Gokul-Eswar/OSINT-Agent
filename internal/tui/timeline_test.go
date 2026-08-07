package tui

import (
	"testing"
)

func TestRenderTimeline_EmptyCaseID(t *testing.T) {
	result := RenderTimeline("")
	if result != "No case selected." {
		t.Errorf("expected 'No case selected.', got %q", result)
	}
}
