package store

import (
	"errors"
	"strings"
	"testing"

	"github.com/spectre/spectre/internal/extensions"
)

func TestItemMetadata(t *testing.T) {
	e := extensions.Extension{
		Name:        "dns-probe",
		Version:     "1.2.3",
		Description: "Probe DNS metadata",
	}
	i := item{ext: e}

	if got := i.Title(); got != "dns-probe" {
		t.Fatalf("unexpected title: %q", got)
	}

	if got := i.FilterValue(); got != "dns-probe" {
		t.Fatalf("unexpected filter value: %q", got)
	}

	if got := i.Description(); got != "[1.2.3] Probe DNS metadata" {
		t.Fatalf("unexpected description: %q", got)
	}
}

func TestView_WhenInstalling_ShowsSpinnerMessage(t *testing.T) {
	m := NewModel(nil)
	m.installing = "dns-probe"

	view := m.View()
	if !strings.Contains(view, "Installing dns-probe") {
		t.Fatalf("expected installing message in view, got: %q", view)
	}
}

func TestView_WhenError_ShowsErrorMessage(t *testing.T) {
	m := NewModel(nil)
	m.err = errors.New("install failed")

	view := m.View()
	if !strings.Contains(view, "Error: install failed") {
		t.Fatalf("expected error message in view, got: %q", view)
	}
}
