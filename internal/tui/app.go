package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spectre/spectre/internal/analysis"
	"github.com/spectre/spectre/internal/config"
	"github.com/spectre/spectre/internal/core"
	"github.com/spectre/spectre/internal/report"
	"github.com/spf13/viper"
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
	ViewDashboard
)

type analysisStatus int

const (
	AnalysisIdle analysisStatus = iota
	AnalysisRunning
	AnalysisComplete
	AnalysisError
)

type model struct {
	state          sessionState
	cursor         int
	quitting       bool
	width          int
	height         int
	selectedCaseID string
	modelName      string // For status bar
	focusNav       bool   // Focus state: true=Sidebar, false=MainContent
	navCursor      int    // Cursor for sidebar selection

	// Settings State
	settingsCursor int

	availableModels []string

	// Analysis State
	analysisStatus analysisStatus
	analysisStep   int
	analysisResult string
	analysisError  string

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
		state:          ViewCases,
		navCursor:      0,
		settingsCursor: 0,
		caseList:       l,
		entityTable:    NewEntityTable(),
		relTable:       NewRelationshipTable(),
		runner:         NewRunnerModel(),
		modelName:      "llama3:8b",
		availableModels: []string{"llama3:8b", "mistral"},
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

	case TickMsg:
		if m.state == ViewAnalysis && m.analysisStatus == AnalysisRunning {
			m.analysisStep++
			if m.analysisStep >= 4 {
				return m, PerformActualAnalysis(m.selectedCaseID, m.modelName)
			}
			return m, tickCmd()
		}

	case AnalysisFinishedMsg:
		if msg.Result != nil {
			m.analysisStatus = AnalysisComplete
			m.analysisResult = FormatAnalysis(msg.Result)
		} else {
			// This was a report generation success
			m.analysisStatus = AnalysisComplete
			m.analysisResult = "Report generated successfully! Check report_<case_id>.md"
		}

	case ModelsFoundMsg:
		m.availableModels = msg
		// If current model is not in list, switch to first
		if len(m.availableModels) > 0 {
			found := false
			for _, mod := range m.availableModels {
				if mod == m.modelName {
					found = true
					break
				}
			}
			if !found {
				m.modelName = m.availableModels[0]
				viper.Set("llm.model", m.modelName)
			}
		}

	case AnalysisErrorMsg:
		m.analysisStatus = AnalysisError
		m.analysisError = string(msg)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "tab":
			m.focusNav = !m.focusNav
			return m, nil
		case "right":
			m.focusNav = false
			return m, nil
		case "left":
			m.focusNav = true
			return m, nil
		case "?":
			// Show help overlay?
		}

		if m.focusNav {
			switch msg.String() {
			case "j", "down":
				if m.navCursor < int(ViewDashboard) {
					m.navCursor++
				}
			case "k", "up":
				if m.navCursor > 0 {
					m.navCursor--
				}
			case "enter":
				m.state = sessionState(m.navCursor)
				// Trigger specific view logic if needed
				if m.state == ViewAnalysis && m.selectedCaseID != "" && m.analysisStatus == AnalysisIdle {
					m.analysisStatus = AnalysisRunning
					m.analysisStep = 0
					return m, StartAnalysis(m.selectedCaseID, m.modelName)
				}
				if m.state == ViewDashboard {
					openBrowser("http://localhost:8080")
				}
			}
			return m, nil
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

		if m.state == ViewEvidence {
			var cmd2 tea.Cmd
			m.entityTable, cmd2 = m.entityTable.Update(msg)
			return m, cmd2
		}

		// Global navigation if not in a list
		switch msg.String() {
		case "1": m.state = ViewCases
		case "2":
			m.state = ViewAnalysis
			if m.selectedCaseID != "" && m.analysisStatus == AnalysisIdle {
				m.analysisStatus = AnalysisRunning
				m.analysisStep = 0
				return m, StartAnalysis(m.selectedCaseID, m.modelName)
			}
		case "3": m.state = ViewEvidence
		case "4": m.state = ViewGraph
		case "5": m.state = ViewTimeline
		case "6": m.state = ViewReports
		case "7": m.state = ViewSettings
		}

		// Specific View Keybindings
		if m.state == ViewSettings {
			if msg.String() == "r" && m.settingsCursor == 1 {
				return m, fetchModelsCmd
			}

			switch msg.String() {
			case "j", "down":
				if m.settingsCursor < 6 {
					m.settingsCursor++
				}
			case "k", "up":
				if m.settingsCursor > 0 {
					m.settingsCursor--
				}
			case "enter", " ":
				switch m.settingsCursor {
				case 0: // Ghost Mode
					viper.Set("ghost_mode", !viper.GetBool("ghost_mode"))
				case 1: // Model
					idx := 0
					for i, mod := range m.availableModels {
						if mod == m.modelName {
							idx = i
							break
						}
					}
					if len(m.availableModels) > 0 {
						idx = (idx + 1) % len(m.availableModels)
						m.modelName = m.availableModels[idx]
					}
					viper.Set("llm.model", m.modelName)
				case 2: // DNS
					viper.Set("collectors.dns.enabled", !viper.GetBool("collectors.dns.enabled"))
				case 3: // Whois
					viper.Set("collectors.whois.enabled", !viper.GetBool("collectors.whois.enabled"))
				case 4: // GitHub
					viper.Set("collectors.github.enabled", !viper.GetBool("collectors.github.enabled"))
				case 5: // Geo
					viper.Set("collectors.geo.enabled", !viper.GetBool("collectors.geo.enabled"))
				case 6: // Ports
					viper.Set("collectors.ports.enabled", !viper.GetBool("collectors.ports.enabled"))
				}
				// Save config
				config.ApplyEthicsConfig() // Re-apply ethics if needed
				viper.WriteConfig()
			}
			return m, nil
		}

		if m.state == ViewReports {
			if msg.String() == "1" && m.selectedCaseID != "" {
				return m, func() tea.Msg {
					content, err := report.GenerateMarkdownReport(m.selectedCaseID)
					if err != nil {
						return AnalysisErrorMsg(err.Error())
					}
					// Save to file
					filename := fmt.Sprintf("report_%s.md", m.selectedCaseID)
					os.WriteFile(filename, []byte(content), 0644)
					return AnalysisFinishedMsg{nil} // Signal success (nil result means just update status)
				}
			}
			if msg.String() == "p" && m.selectedCaseID != "" {
				return m, func() tea.Msg {
					filename, err := report.GeneratePDFReport(m.selectedCaseID)
					if err != nil {
						return AnalysisErrorMsg(err.Error())
					}
					return AnalysisFinishedMsg{&core.Analysis{
						Findings: []string{fmt.Sprintf("PDF Report generated: %s", filename)},
					}}
				}
			}
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
	
	gapSize := m.width - lipgloss.Width(title) - lipgloss.Width(info) - 4
	if gapSize < 0 {
		gapSize = 0
	}
	
	header := lipgloss.JoinHorizontal(lipgloss.Center, title, strings.Repeat(" ", gapSize), info)
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
		{ViewDashboard, "Web Dashboard"},
	}

	var s strings.Builder
	s.WriteString("\n")
	for i, v := range views {
		line := v.label
		
		// Prefix: Active indicator
		prefix := "  "
		if m.state == v.state {
			prefix = "▶ "
		}

		// Selection style
		content := line
		if m.navCursor == i {
			content = StyleSelectedNav.Render(line)
		} else if m.state == v.state {
			content = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true).Render(line)
		}

		s.WriteString(prefix + content + "\n")
	}

	style := StyleNav.Height(m.height - 4).Width(20)
	if m.focusNav {
		style = style.BorderForeground(ColorAccent)
	}
	return style.Render(s.String())
}

func (m model) renderContent() string {
	var content string
	switch m.state {
	case ViewCases:
		content = m.caseList.View()
	case ViewEvidence:
		content = fmt.Sprintf("EVIDENCE — %s\n\n", m.selectedCaseID) + m.entityTable.View()
	case ViewAnalysis:
		if m.selectedCaseID == "" {
			content = "No case selected. Please select a case first (Press 1)."
		} else {
			switch m.analysisStatus {
			case AnalysisIdle:
				content = "Press '2' to start analysis."
			case AnalysisRunning:
				steps := []string{
					"Collecting evidence...",
					"Cross-checking sources...",
					"Reasoning over timeline...",
					"Synthesizing conclusion...",
				}
				var s strings.Builder
				s.WriteString(fmt.Sprintf("ANALYSIS — Case %s\n", m.selectedCaseID))
				s.WriteString("────────────────────────────────────\n\n")
				s.WriteString("⠋ Running analysis pipeline...\n\n")

				for i, step := range steps {
					if m.analysisStep > i {
						s.WriteString(fmt.Sprintf(" [✓] %s\n", step))
					} else if m.analysisStep == i {
						s.WriteString(fmt.Sprintf(" [▶] %s\n", step))
					} else {
						s.WriteString(fmt.Sprintf(" [ ] %s\n", step))
					}
				}
				content = s.String()
			case AnalysisComplete:
				content = m.analysisResult
			case AnalysisError:
				content = StyleMuted.Foreground(ColorError).Render(fmt.Sprintf("Error: %s", m.analysisError))
			}
		}
	case ViewGraph:
		content = RenderASCIIGraph(m.selectedCaseID)
	case ViewTimeline:
		content = RenderTimeline(m.selectedCaseID)
	case ViewReports:
		status := ""
		if m.analysisStatus == AnalysisComplete {
			status = "\n\n" + StyleMuted.Foreground(ColorSuccess).Render(m.analysisResult)
		}
		content = "REPORTS\n───────\n\n[1] Generate Markdown Report\n[p] Generate Professional PDF\n[2] View Latest Analysis" + status
	case ViewSettings:
		// Settings Menu
		modelLabel := m.modelName
		if m.settingsCursor == 1 {
			modelLabel += " (Press 'r' to scan)"
		}

		opts := []struct {
			label string
			val   string
		}{
			{"Ghost Mode", formatBool(viper.GetBool("ghost_mode"))},
			{"Model", modelLabel},
			{"DNS Collector", formatBool(viper.GetBool("collectors.dns.enabled"))},
			{"Whois Collector", formatBool(viper.GetBool("collectors.whois.enabled"))},
			{"GitHub Collector", formatBool(viper.GetBool("collectors.github.enabled"))},
			{"Geo Collector", formatBool(viper.GetBool("collectors.geo.enabled"))},
			{"Ports Collector", formatBool(viper.GetBool("collectors.ports.enabled"))},
		}

		var s strings.Builder
		s.WriteString("SETTINGS\n────────\n\n")

		for i, opt := range opts {
			cursor := "  "
			if m.settingsCursor == i {
				cursor = "> "
			}

			style := lipgloss.NewStyle()
			if m.settingsCursor == i {
				style = StyleSelectedNav
			}

			s.WriteString(fmt.Sprintf("%s%s: %s\n", cursor, style.Render(opt.label), opt.val))
		}
		
		s.WriteString("\n(Press 'Enter' or 'Space' to toggle/change)\n")
		content = s.String()

	case ViewDashboard:
		content = "Opening Web Dashboard in your default browser...\n\nURL: http://localhost:8080"
	
	default:
		content = "View not implemented yet."
	}

	return StyleMain.Width(m.width - 25).Height(m.height - 4).Render(content)
}

func (m model) renderFooter() string {
	status := "● Connected"
	if viper.GetBool("ghost_mode") {
		status = "👻 GHOST MODE ACTIVE"
	}

	info := "ollama:localhost:11434  |  Press ? for help"

	

	gapSize := m.width - lipgloss.Width(status) - lipgloss.Width(info) - 4

	if gapSize < 0 {

		gapSize = 0

	}



	footer := lipgloss.JoinHorizontal(lipgloss.Center, 

		lipgloss.NewStyle().Foreground(ColorSuccess).Render(status),

		strings.Repeat(" ", gapSize),

		info,

	)

	

	return StyleStatus.Width(m.width).Render(footer)

}

func formatBool(b bool) string {
	if b {
		return "[ON]"
	}
	return "[OFF]"
}

type ModelsFoundMsg []string

func fetchModelsCmd() tea.Msg {
	models, err := analysis.FetchAvailableModels()
	if err != nil {
		return AnalysisErrorMsg("Scan failed: " + err.Error())
	}
	return ModelsFoundMsg(models)
}
