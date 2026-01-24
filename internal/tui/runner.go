package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spectre/spectre/internal/collector"
)

type runnerState int

const (
	selectCase runnerState = iota
	selectCollector
	inputTarget
	executing
	result
)

type runnerModel struct {
	state          runnerState
	caseList       list.Model
	collList       list.Model
	textInput      textinput.Model
	selectedCaseID string
	selectedColl   string
	activeAllowed  bool
	err            error
	message        string
}

func NewRunnerModel() runnerModel {
	ti := textinput.New()
	ti.Placeholder = "example.com"
	ti.Focus()

	// Dynamic Collector list
	var collectors []list.Item
	for _, c := range collector.List() {
		collectors = append(collectors, item{
			title: c.Name(),
			desc:  c.Description(),
		})
	}

	cl := list.New(collectors, list.NewDefaultDelegate(), 0, 0)
	cl.Title = "Select Collector (Space to toggle Active Mode)"
	cl.SetShowHelp(false)

	// Case list (initialized empty, will be populated by Update)
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select Case for Collection"
	l.SetShowHelp(false)

	return runnerModel{
		state:     selectCase,
		caseList:  l,
		collList:  cl,
		textInput: ti,
	}
}

func (m runnerModel) Update(msg tea.Msg) (runnerModel, tea.Cmd) {
	var cmd tea.Cmd

	switch m.state {
	case selectCase:
		if km, ok := msg.(tea.KeyMsg); ok && km.String() == "enter" {
			if i := m.caseList.SelectedItem(); i != nil {
				m.selectedCaseID = i.(item).id
				m.state = selectCollector
				return m, nil
			}
		}
		m.caseList, cmd = m.caseList.Update(msg)
		return m, cmd

	case selectCollector:
		if km, ok := msg.(tea.KeyMsg); ok {
			switch km.String() {
			case "enter":
				if i := m.collList.SelectedItem(); i != nil {
					m.selectedColl = i.(item).title
					m.state = inputTarget
					return m, nil
				}
			case " ":
				m.activeAllowed = !m.activeAllowed
				status := "[SAFE]"
				if m.activeAllowed {
					status = "[ACTIVE/DANGEROUS]"
				}
				m.collList.Title = fmt.Sprintf("Select Collector (Space to toggle) %s", status)
				return m, nil
			}
		}
		m.collList, cmd = m.collList.Update(msg)
		return m, cmd

	case inputTarget:
		if km, ok := msg.(tea.KeyMsg); ok && km.String() == "enter" {
			target := m.textInput.Value()
			if target == "" {
				return m, nil
			}
			m.state = executing
			return m, func() tea.Msg {
				_, err := collector.Run(m.selectedColl, m.selectedCaseID, target, m.activeAllowed)
				if err != nil {
					return err
				}
				return "Collection complete!"
			}
		}
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd

	case result, executing:
		if km, ok := msg.(tea.KeyMsg); ok && km.String() == "enter" {
			m.state = selectCase
			m.message = ""
			m.err = nil
			return m, nil
		}
	}

	switch msg := msg.(type) {
	case error:
		m.err = msg
		m.state = result
	case string:
		m.message = msg
		m.state = result
	}

	return m, nil
}

func (m runnerModel) View() string {
	switch m.state {
	case selectCase:
		return m.caseList.View()
	case selectCollector:
		return m.collList.View()
	case inputTarget:
		status := "SAFE"
		if m.activeAllowed {
			status = "ACTIVE/DANGEROUS"
		}
		return fmt.Sprintf(
			"Running %s for case %s\nMode: %s\n\nEnter target:\n\n%s\n\n(enter: run • esc: cancel)",
			m.selectedColl, m.selectedCaseID, status, m.textInput.View(),
		)
	case executing:
		return fmt.Sprintf("Running %s against %s... Please wait.", m.selectedColl, m.textInput.Value())
	case result:
		if m.err != nil {
			return fmt.Errorf("Error: %v\n\n(enter: continue)", m.err).Error()
		}
		return fmt.Sprintf("%s\n\n(enter: continue)", m.message)
	}
	return ""
}
