package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState int

const (
	ViewCases sessionState = iota
	ViewAnalysis
	ViewEvidence
	ViewGraph
	ViewTimeline
	ViewReports
	ViewSettings
)

type model struct {
	state          sessionState
	cursor         int
	quitting       bool
	width          int
	height         int
	selectedCaseID string
	modelName      string // For status bar

	// Sub-models
	caseList    list.Model
	entityTable table.Model
	relTable    table.Model
	runner      runnerModel
}

func InitialModel() model {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "CASES"

	return model{
		state:       ViewCases,
		caseList:    l,
		entityTable: NewEntityTable(),
		relTable:    NewRelationshipTable(),
		runner:      NewRunnerModel(),
		modelName:   "llama3:8b",
	}
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		items, err := FetchCases()
		if err != nil {
			return err
		}
		return items
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		
		// Update sub-models sizes
		h := msg.Height - 6 // Space for header and footer
		w := msg.Width - 25 // Space for nav
		
		m.caseList.SetSize(w, h)
		m.runner.caseList.SetSize(w, h)
		m.runner.collList.SetSize(w, h)

	case []list.Item:
		m.caseList.SetItems(msg)
		m.runner.caseList.SetItems(msg)

	case []table.Row:
		if m.state == ViewEvidence {
			m.entityTable.SetRows(msg)
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "tab":
			// Toggle between nav and content? For now just simple
		case "?":
			// Show help overlay?
		}

		// Navigation logic (j/k or arrows)
		if m.state == ViewCases {
			// In Cases view, we might want to navigate the list
			var cmd2 tea.Cmd
			m.caseList, cmd2 = m.caseList.Update(msg)
			
			if msg.String() == "enter" {
				if selected, ok := m.caseList.SelectedItem().(item); ok {
					m.selectedCaseID = selected.id
					// Transition to Analysis or Evidence?
					m.state = ViewEvidence
					return m, func() tea.Msg {
						rows, err := FetchEntities(selected.id)
						if err != nil {
							return err
						}
						return rows
					}
				}
			}
			return m, cmd2
		}

		// Global navigation if not in a list
		switch msg.String() {
		case "1": m.state = ViewCases
		case "2": m.state = ViewAnalysis
		case "3": m.state = ViewEvidence
		case "4": m.state = ViewGraph
		case "5": m.state = ViewTimeline
		case "6": m.state = ViewReports
		case "7": m.state = ViewSettings
		}
	}

	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return "Goodbye, Agent.\n"
	}

	header := m.renderHeader()
	nav := m.renderNav()
	content := m.renderContent()
	footer := m.renderFooter()

	// Layout:
	// Header
	// Nav | Content
	// Footer

	mainArea := lipgloss.JoinHorizontal(lipgloss.Top, nav, content)
	
	return lipgloss.JoinVertical(lipgloss.Left, header, mainArea, footer)
}

func (m model) renderHeader() string {
	caseInfo := "No Case Selected"
	if m.selectedCaseID != "" {
		caseInfo = fmt.Sprintf("Case: %s", m.selectedCaseID)
	}

	title := StyleTitle.Render("SPECTRE v1.0")
	info := StyleMuted.Render(fmt.Sprintf("%s  |  Model: %s", caseInfo, m.modelName))
	
	header := lipgloss.JoinHorizontal(lipgloss.Center, title, strings.Repeat(" ", m.width-lipgloss.Width(title)-lipgloss.Width(info)-4), info)
	return StyleHeader.Width(m.width).Render(header)
}

func (m model) renderNav() string {
	views := []struct {
		state sessionState
		label string
	}{
		{ViewCases, "Cases"},
		{ViewAnalysis, "Analysis"},
		{ViewEvidence, "Evidence"},
		{ViewGraph, "Graph"},
		{ViewTimeline, "Timeline"},
		{ViewReports, "Reports"},
		{ViewSettings, "Settings"},
	}

	var s strings.Builder
	s.WriteString("\n")
	for _, v := range views {
		label := v.label
		if m.state == v.state {
			s.WriteString(StyleSelectedNav.Render("▶ " + label))
		} else {
			s.WriteString("  " + label)
		}
		s.WriteString("\n")
	}

	return StyleNav.Height(m.height - 4).Width(20).Render(s.String())
}

func (m model) renderContent() string {
	var content string
	switch m.state {
	case ViewCases:
		content = m.caseList.View()
	case ViewEvidence:
		content = fmt.Sprintf("EVIDENCE — %s\n\n", m.selectedCaseID) + m.entityTable.View()
	case ViewAnalysis:
		content = "ANALYSIS BRAIN\n\n⠋ Thinking..."
	default:
		content = "View not implemented yet."
	}

	return StyleMain.Width(m.width - 25).Height(m.height - 4).Render(content)
}

func (m model) renderFooter() string {
	status := "● Connected"
	info := "ollama:localhost:11434  |  Press ? for help"
	
	footer := lipgloss.JoinHorizontal(lipgloss.Center, 
		lipgloss.NewStyle().Foreground(ColorSuccess).Render(status),
		strings.Repeat(" ", m.width-lipgloss.Width(status)-lipgloss.Width(info)-4),
		info,
	)
	
	return StyleStatus.Width(m.width).Render(footer)
}