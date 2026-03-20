package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spectre/spectre/internal/agent"
	"github.com/spectre/spectre/internal/analyzer"
)

type chatModel struct {
	viewport      viewport.Model
	textinput     textinput.Model
	spinner       spinner.Model
	thinking      bool
	messages      []analyzer.Message
	caseID        string
	engine        *agent.Engine
	width, height int
}

func newChatModel() chatModel {
	ti := textinput.New()
	ti.Placeholder = "Ask the agent anything..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	vp := viewport.New(0, 0)
	vp.SetContent("Welcome to SPECTRE Chat. Select a case to begin investigation.")

	return chatModel{
		textinput: ti,
		viewport:  vp,
		spinner:   s,
		messages:  []analyzer.Message{},
	}
}

func (m *chatModel) setCase(caseID string) {
	if m.caseID == caseID {
		return
	}
	m.caseID = caseID
	m.engine = agent.NewEngine(caseID)
	m.messages = []analyzer.Message{}
	m.viewport.SetContent(fmt.Sprintf("Agent ready for case: %s\nType a message to begin.", caseID))
}

type agentResponseMsg struct {
	content string
	err     error
}

func (m chatModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m chatModel) Update(msg tea.Msg) (chatModel, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		sCmd  tea.Cmd
	)

	m.textinput, tiCmd = m.textinput.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	if m.thinking {
		m.spinner, sCmd = m.spinner.Update(msg)
	}

	switch msg := msg.(type) {
	case agentResponseMsg:
		m.thinking = false
		if msg.err != nil {
			m.appendMessage("system", "Error: "+msg.err.Error())
		} else {
			m.appendMessage("assistant", msg.content)
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.thinking || m.textinput.Value() == "" {
				return m, nil
			}

			input := m.textinput.Value()
			m.appendMessage("user", input)
			m.textinput.SetValue("")
			m.thinking = true

			return m, tea.Batch(
				m.spinner.Tick,
				func() tea.Msg {
					resp, err := m.engine.Execute(input)
					return agentResponseMsg{resp, err}
				},
			)
		}
	}

	return m, tea.Batch(tiCmd, vpCmd, sCmd)
}

func (m *chatModel) appendMessage(role, content string) {
	m.messages = append(m.messages, analyzer.Message{Role: role, Content: content})

	var sb strings.Builder
	for _, msg := range m.messages {
		roleStyle := lipgloss.NewStyle().Bold(true)
		switch msg.Role {
		case "user":
			roleStyle = roleStyle.Foreground(lipgloss.Color("6"))
			sb.WriteString(roleStyle.Render("You: ") + msg.Content + "\n\n")
		case "assistant":
			roleStyle = roleStyle.Foreground(lipgloss.Color("5"))
			sb.WriteString(roleStyle.Render("Agent: ") + msg.Content + "\n\n")
		case "system":
			roleStyle = roleStyle.Foreground(lipgloss.Color("8"))
			sb.WriteString(roleStyle.Render("System: ") + msg.Content + "\n\n")
		}
	}
	m.viewport.SetContent(sb.String())
	m.viewport.GotoBottom()
}

func (m chatModel) View() string {
	spin := ""
	if m.thinking {
		spin = m.spinner.View() + " Agent is thinking..."
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.viewport.View(),
		"\n",
		spin,
		m.textinput.View(),
	)
}

func (m *chatModel) setSize(w, h int) {
	m.width = w
	m.height = h
	m.viewport.Width = w
	m.viewport.Height = h - 6
	m.textinput.Width = w - 2
}
