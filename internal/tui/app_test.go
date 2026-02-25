package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestInitialModel(t *testing.T) {
	m := InitialModel()
	assert.Equal(t, ViewCases, m.state)
	assert.False(t, m.quitting)
	assert.Equal(t, 0, m.navCursor)
}

func TestModelTransitions(t *testing.T) {
	m := InitialModel()
	
	// Switch to sidebar focus
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	assert.True(t, m2.(model).focusNav)

	// Move down to 'Analysis' (navCursor 1)
	m3, _ := m2.(model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	assert.Equal(t, 1, m3.(model).navCursor)

	// Press enter to switch state
	m4, _ := m3.(model).Update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.Equal(t, ViewAnalysis, m4.(model).state)
}

func TestQuit(t *testing.T) {
	m := InitialModel()
	m2, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	
	assert.True(t, m2.(model).quitting)
	assert.NotNil(t, cmd)
}

func TestNavFocus(t *testing.T) {
	m := InitialModel()
	assert.False(t, m.focusNav) // Default: sidebar not focused for arrows
	
	// Tab toggles focus
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	assert.True(t, m2.(model).focusNav)
	
	m3, _ := m2.(model).Update(tea.KeyMsg{Type: tea.KeyTab})
	assert.False(t, m3.(model).focusNav)
}
