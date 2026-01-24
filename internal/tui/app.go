package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState int

const (
	homeView sessionState = iota
	caseView
	collectView
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
	caseList list.Model
}

func InitialModel() model {
	// Initialize list with dummy data, will be updated in Init or Update
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Investigation Cases"

	return model{
		state:   homeView,
		choices: []string{"Investigation Cases", "Run Collector", "System Stats", "Quit"},
		caseList: l,
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
		m.caseList.SetSize(msg.Width-4, msg.Height-4)

	case []list.Item:
		m.caseList.SetItems(msg)

	case tea.KeyMsg:
		if m.state == caseView {
			if msg.String() == "esc" {
				m.state = homeView
				return m, nil
			}
			m.caseList, cmd = m.caseList.Update(msg)
			return m, cmd
		}

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
			case 3: // Quit
				m.quitting = true
				return m, tea.Quit
			}
		}
	}
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return "Goodbye, Agent.\n"
	}

	if m.state == caseView {
		return docStyle.Render(m.caseList.View())
	}

	var s strings.Builder

	// Header
	s.WriteString("\n")
	s.WriteString(titleStyle.Render("SPECTRE COMMAND CENTER"))
	s.WriteString("\n\n")

	// ASCII Art (Fixed escapes)
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