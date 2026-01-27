package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors - Stealth/Retro (Gruvbox Inspired)
	ColorPrimary = lipgloss.Color("#689D6A") // Retro Green/Aqua
	ColorAccent  = lipgloss.Color("#FABD2F") // Retro Yellow/Amber
	ColorMuted   = lipgloss.Color("#7C6F64") // Warm Gray
	ColorError   = lipgloss.Color("#CC241D") // Dim Red
	ColorSuccess = lipgloss.Color("#98971A") // Olive Green
	ColorBG      = lipgloss.Color("#1D2021") // Dark Hard Background

	// Styles
	StyleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Padding(0, 1)

	StyleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Border(lipgloss.DoubleBorder(), false, false, true, false).
			BorderForeground(ColorMuted).
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
