package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spectre/spectre/internal/server"
)

type sessionState int

const (
	homeView sessionState = iota
	caseView
	detailView
	relView
	collectView
	webView
	statsView
)

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		MarginLeft(2)

	itemStyle = lipgloss.NewStyle().
		PaddingLeft(4)

	selectedItemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.Color("#10b981")).
			Bold(true)

	docStyle = lipgloss.NewStyle().Margin(1, 2)
)

type model struct {
	state    sessionState
	cursor   int
	choices  []string
	quitting bool
	width    int
	height   int

	// Sub-models
	caseList    list.Model
	entityTable table.Model
	relTable    table.Model
	runner      runnerModel
	selectedCaseID string

	// Server State
	serverRunning bool
	serverPort    int
	serverMsg     string
}

func InitialModel() model {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Investigation Cases"

	return model{
		state:       homeView,
		choices:     []string{"Investigation Cases", "Run Collector", "Web Dashboard", "System Stats", "Quit"},
		caseList:    l,
		entityTable: NewEntityTable(),
		relTable:    NewRelationshipTable(),
		runner:      NewRunnerModel(),
		serverPort:  8080,
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
		h := msg.Height - 4
		if h < 0 {
			h = 0
		}
		w := msg.Width - 4
		if w < 0 {
			w = 0
		}
		m.caseList.SetSize(w, h)
		m.runner.caseList.SetSize(w, h)
		m.runner.collList.SetSize(w, h)

	case []list.Item:
		m.caseList.SetItems(msg)
		m.runner.caseList.SetItems(msg)

	case []table.Row:
		if m.state == detailView {
			m.entityTable.SetRows(msg)
		} else if m.state == relView {
			m.relTable.SetRows(msg)
		}

	case tea.KeyMsg:
		// Global back to home
		if msg.String() == "esc" && m.state != homeView {
			if m.state == relView {
				m.state = detailView
			} else {
				m.state = homeView
			}
			return m, nil
		}

		switch m.state {
		case detailView:
			if msg.String() == "r" {
				m.state = relView
				return m, func() tea.Msg {
					rows, err := FetchRelationships(m.selectedCaseID)
					if err != nil {
						return err
					}
					return rows
				}
			}
			m.entityTable, cmd = m.entityTable.Update(msg)
			return m, cmd

		case relView:
			m.relTable, cmd = m.relTable.Update(msg)
			return m, cmd

		case caseView:
			if msg.String() == "enter" {
				selected := m.caseList.SelectedItem().(item)
				m.selectedCaseID = selected.id
				m.state = detailView
				return m, func() tea.Msg {
					rows, err := FetchEntities(selected.id)
					if err != nil {
						return err
					}
					return rows
				}
			}
			m.caseList, cmd = m.caseList.Update(msg)
			return m, cmd

		case collectView:
			m.runner, cmd = m.runner.Update(msg)
			return m, cmd

		case webView:
			switch msg.String() {
			case "enter", "o":
				// Open browser
				openBrowser(fmt.Sprintf("http://localhost:%d", m.serverPort))
				m.serverMsg = "Opened in browser"
			case "s":
				if !m.serverRunning {
					go func() {
						if err := server.Start(m.serverPort); err != nil {
							// Log error? For now just ignore as it might be blocking
						}
					}()
					// Give it a moment to start
					time.Sleep(500 * time.Millisecond)
					m.serverRunning = true
					m.serverMsg = "Server started"
				}
			}
			return m, nil

		case homeView:
			switch msg.String() {
			case "ctrl+c", "q":
				m.quitting = true
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			case "enter", " ":
				switch m.cursor {
				case 0: // Cases
					m.state = caseView
				case 1: // Run Collector
					m.state = collectView
				case 2: // Web Dashboard
					m.state = webView
					if !m.serverRunning {
						// Auto-start server if not running
						go server.Start(m.serverPort)
						m.serverRunning = true
						m.serverMsg = "Server auto-started..."
						// Auto-open browser
						go func() {
							time.Sleep(1 * time.Second)
							openBrowser(fmt.Sprintf("http://localhost:%d", m.serverPort))
						}()
					}
				case 3: // Stats
					m.state = statsView
				case 4: // Quit
					m.quitting = true
					return m, tea.Quit
				}
			}
		}
	}
	
	if m.state == collectView {
		m.runner, cmd = m.runner.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return "Goodbye, Agent.\n"
	}

	switch m.state {
	case caseView:
		return docStyle.Render(m.caseList.View())
	case detailView:
		return docStyle.Render(
			fmt.Sprintf("Case Entities: %s\n\n", m.selectedCaseID) +
			m.entityTable.View() +
			"\n\n(r: relationships • esc: back)",
		)
	case relView:
		return docStyle.Render(
			fmt.Sprintf("Case Relationships: %s\n\n", m.selectedCaseID) +
			m.relTable.View() +
			"\n\n(esc: back)",
		)
	case collectView:
		return docStyle.Render(m.runner.View())
	case webView:
		status := "STOPPED"
		color := "#ef4444" // red
		if m.serverRunning {
			status = "RUNNING"
			color = "#10b981" // green
		}
		
		s := fmt.Sprintf("\n  WEB COMMAND CENTER\n\n  Status: %s\n  URL:    http://localhost:%d\n\n", 
			lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Bold(true).Render(status),
			m.serverPort,
		)

		if m.serverMsg != "" {
			s += fmt.Sprintf("  > %s\n\n", m.serverMsg)
		}

		s += "  (o: open browser • esc: back)"
		return docStyle.Render(s)
	case statsView:
		return docStyle.Render(GetSystemStats())
	}

	var s strings.Builder

	// Header
	s.WriteString("\n")
	s.WriteString(titleStyle.Render("SPECTRE COMMAND CENTER"))
	s.WriteString("\n\n")

	// ASCII Art
	ascii := `    .---.      .---.    
   /  _  \    /  _  \   
   | (_) |    | (_) |   
   \  ^  /    \  ^  /   
    '---'      '---'    `
	s.WriteString(ascii)
	s.WriteString("\n\n")

	// Menu
	for i, choice := range m.choices {
		if m.cursor == i {
			s.WriteString(selectedItemStyle.Render(fmt.Sprintf("> %s", choice)))
		} else {
			s.WriteString(itemStyle.Render(choice))
		}
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render(" (up/down: navigate • enter: select • q: quit)"))

	return docStyle.Render(s.String())
}
