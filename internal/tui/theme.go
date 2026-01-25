package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	ColorPrimary = lipgloss.Color("#7C3AED") // Violet
	ColorAccent  = lipgloss.Color("#22D3EE") // Cyan
	ColorMuted   = lipgloss.Color("#6B7280") // Gray
	ColorError   = lipgloss.Color("#EF4444") // Red
	ColorSuccess = lipgloss.Color("#10B981") // Green
	ColorBG      = lipgloss.Color("#111827") // Dark background

	// Styles
	StyleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Padding(0, 1)

	StyleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorPrimary).
			Padding(0, 1)

	StyleNav = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(ColorMuted).
			Padding(1, 2)

	StyleMain = lipgloss.NewStyle().
			Padding(1, 2)

	StyleStatus = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(ColorMuted).
			Padding(0, 1)

	StyleSelectedNav = lipgloss.NewStyle().
				Foreground(ColorAccent).
				Bold(true)

	StyleMuted = lipgloss.NewStyle().
			Foreground(ColorMuted)
)
