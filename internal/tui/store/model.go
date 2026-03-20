package store

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spectre/spectre/internal/extensions"
)

var (
	appStyle   = lipgloss.NewStyle().Padding(1, 2)
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type item struct {
	ext extensions.Extension
}

func (i item) Title() string       { return i.ext.Name }
func (i item) Description() string { return fmt.Sprintf("[%s] %s", i.ext.Version, i.ext.Description) }
func (i item) FilterValue() string { return i.ext.Name }

type Model struct {
	list          list.Model
	manager       *extensions.Manager
	err           error
	installing    string
	spinner       spinner.Model
	width, height int
}

func NewModel(mgr *extensions.Manager) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Spectre Extension Store"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle

	return Model{
		list:    l,
		manager: mgr,
		spinner: s,
	}
}

type extensionsLoadedMsg []extensions.Extension
type installFinishedMsg struct{ err error }

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.loadExtensions,
	)
}

func (m Model) loadExtensions() tea.Msg {
	exts, err := m.manager.ListRemote()
	if err != nil {
		return err
	}
	return extensionsLoadedMsg(exts)
}

func (m Model) installExtension(name string) tea.Cmd {
	return func() tea.Msg {
		err := m.manager.Install(name)
		return installFinishedMsg{err: err}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case extensionsLoadedMsg:
		var items []list.Item
		for _, e := range msg {
			items = append(items, item{ext: e})
		}
		m.list.SetItems(items)

	case tea.KeyMsg:
		if m.installing != "" {
			return m, nil // Block input while installing
		}
		switch msg.String() {
		case "enter":
			if i, ok := m.list.SelectedItem().(item); ok {
				m.installing = i.ext.Name
				return m, m.installExtension(i.ext.Name)
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case installFinishedMsg:
		if msg.err != nil {
			m.list.NewStatusMessage(statusMessageStyle("Error: " + msg.err.Error()))
		} else {
			m.list.NewStatusMessage(statusMessageStyle("Successfully installed " + m.installing))
		}
		m.installing = ""
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	// Update list only if not installing (or handle key blocks better)
	if m.installing == "" {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.installing != "" {
		return fmt.Sprintf("\n\n   %s Installing %s...\n\n", m.spinner.View(), m.installing)
	}
	if m.err != nil {
		return fmt.Sprintf("\nError: %v\n", m.err)
	}
	return appStyle.Render(m.list.View())
}
