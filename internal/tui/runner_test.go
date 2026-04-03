package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func TestRunner_SelectCaseToCollectorTransition(t *testing.T) {
	m := NewRunnerModel()
	m.caseList.SetItems([]list.Item{item{id: "case-1", title: "Case 1", desc: "active"}})

	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if next.state != selectCollector {
		t.Fatalf("expected state %v, got %v", selectCollector, next.state)
	}
	if next.selectedCaseID != "case-1" {
		t.Fatalf("expected selectedCaseID case-1, got %q", next.selectedCaseID)
	}
}

func TestRunner_ToggleActiveMode(t *testing.T) {
	m := NewRunnerModel()
	m.state = selectCollector

	next, _ := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	if !next.activeAllowed {
		t.Fatal("expected activeAllowed to be true after space toggle")
	}
	if !strings.Contains(next.collList.Title, "ACTIVE/DANGEROUS") {
		t.Fatalf("expected collector title to reflect active mode, got %q", next.collList.Title)
	}
}

func TestRunner_PortsFlowTransitionsToOptions(t *testing.T) {
	m := NewRunnerModel()
	m.state = selectCollector
	m.collList.SetItems([]list.Item{item{title: "ports", desc: "port scan"}})

	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if m2.state != inputTarget {
		t.Fatalf("expected state %v, got %v", inputTarget, m2.state)
	}
	if m2.selectedColl != "ports" {
		t.Fatalf("expected selected collector ports, got %q", m2.selectedColl)
	}

	m2.textInput.SetValue("example.com")
	m3, _ := m2.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if m3.state != inputOptions {
		t.Fatalf("expected state %v, got %v", inputOptions, m3.state)
	}
}
