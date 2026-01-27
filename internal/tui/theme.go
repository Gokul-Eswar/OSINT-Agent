package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors - "Sleeper Sport" (Matte Slate & Caliper Red)
	// Concept: Looks like a standard high-end dashboard (Normal Car)
	// but has the redline performance accents (Sports Car).
	ColorPrimary = lipgloss.Color("#E2E8F0") // Platinum / Silver
	ColorAccent  = lipgloss.Color("#EF4444") // Sport Red
	ColorMuted   = lipgloss.Color("#475569") // Steel Gray
	ColorError   = lipgloss.Color("#DC2626") // Deep Red
	ColorSuccess = lipgloss.Color("#059669") // Emerald
	ColorBG      = lipgloss.Color("#0F172A") // Midnight Slate

	// Styles
	StyleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorError). // The "Badge" look
			Padding(0, 1)

	StyleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Background(lipgloss.Color("#1E293B")). // Subtle lift from BG
			Padding(0, 1)

	StyleNav = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(ColorMuted).
			Padding(1, 2).
			Background(ColorBG)

	StyleMain = lipgloss.NewStyle().
			Padding(1, 2).
			Background(ColorBG)

	StyleStatus = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(ColorMuted).
			Padding(0, 1).
			Background(ColorBG)

	StyleSelectedNav = lipgloss.NewStyle().
				Foreground(ColorAccent).
				Bold(true)

	StyleMuted = lipgloss.NewStyle().
			Foreground(ColorMuted)
)
