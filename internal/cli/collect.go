package cli

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/progress"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spectre/spectre/internal/collector"
	_ "github.com/spectre/spectre/internal/collector/active" // Register Active Probes
	_ "github.com/spectre/spectre/internal/collector/dns"    // Register DNS
	_ "github.com/spectre/spectre/internal/collector/geo"    // Register GeoIP
	_ "github.com/spectre/spectre/internal/collector/github" // Register GitHub
	_ "github.com/spectre/spectre/internal/collector/whois"  // Register WHOIS
	"github.com/spectre/spectre/internal/storage"
	"github.com/spf13/cobra"
)

type collectMsg struct {
	name string
	err  error
}

type collectModel struct {
	target        string
	collectors    []string
	completed     map[string]bool
	failed        map[string]error
	progress      progress.Model
	spinner       spinner.Model
	quitting      bool
	err           error
	activeAllowed bool
	caseID        string
	dryRun        bool
	threads       int
}

func (m collectModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.spinner.Tick)

	// Concurrency control
	sem := make(chan struct{}, m.threads)

	for _, name := range m.collectors {
		name := name
		cmds = append(cmds, func() tea.Msg {
			sem <- struct{}{}
			defer func() { <-sem }()

			if m.dryRun {
				return collectMsg{name: name}
			}

			_, err := collector.RunAndSave(name, m.caseID, m.target, m.activeAllowed, nil)
			return collectMsg{name: name, err: err}
		})
	}

	return tea.Batch(cmds...)
}

func (m collectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			m.quitting = true
			return m, tea.Quit
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 4
		if m.progress.Width > 80 {
			m.progress.Width = 80
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case collectMsg:
		if msg.err != nil {
			m.failed[msg.name] = msg.err
		} else {
			m.completed[msg.name] = true
		}

		if len(m.completed)+len(m.failed) == len(m.collectors) {
			m.quitting = true
			return m, tea.Quit
		}

		pct := float64(len(m.completed)+len(m.failed)) / float64(len(m.collectors))
		return m, m.progress.SetPercent(pct)

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m collectModel) View() string {
	if m.err != nil {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(fmt.Sprintf("Error: %v", m.err))
	}

	header := lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("Target: %s", m.target))
	if m.dryRun {
		header += lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render(" [DRY-RUN]")
	}
	header += "\n"

	s := header + "\n"

	for _, name := range m.collectors {
		status := " "
		if m.completed[name] {
			status = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("✔")
		} else if err, failed := m.failed[name]; failed {
			status = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("✘")
			name = fmt.Sprintf("%s (%v)", name, err)
		} else {
			status = m.spinner.View()
		}
		s += fmt.Sprintf(" %s %s\n", status, name)
	}

	s += "\n" + m.progress.View() + "\n\n"

	if m.quitting {
		s += "Done!\n"
	} else {
		s += lipgloss.NewStyle().Faint(true).Render("Press q or ctrl+c to quit") + "\n"
	}

	return s
}

var collectCmd = &cobra.Command{
	Use:   "collect [collector|all] [target]",
	Short: "Run a passive collector (or all) against a target",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Try to load context if caseID is missing
		if caseID == "" {
			ctxID, err := LoadContext()
			if err == nil && ctxID != "" {
				caseID = ctxID
				fmt.Printf("Using current case: %s\n", caseID)
			}
		}

		if caseID == "" {
			return fmt.Errorf("case ID is required (use --case or create a new case)")
		}

		collectorName := args[0]
		var target string

		if len(args) == 2 {
			target = args[1]
			// Save for next time
			SaveTarget(target)
		} else {
			// Try to load from context
			ctxTarget, err := LoadTarget()
			if err != nil || ctxTarget == "" {
				return fmt.Errorf("target is required (no active target in context)")
			}
			target = ctxTarget
			fmt.Printf("Using current target: %s\n", target)
		}

		if !dryRun {

			if err := storage.InitDB(); err != nil {
				return err
			}
		}

		var collectorsToRun []string
		if collectorName == "all" {
			for _, c := range collector.List() {
				// Skip active collectors if not allowed
				if c.IsActive() && !activeAllowed {
					continue
				}
				collectorsToRun = append(collectorsToRun, c.Name())
			}
		} else {
			collectorsToRun = []string{collectorName}
		}

		if os.Getenv("GO_TESTING") == "true" {
			fmt.Printf("Running collection in test mode for target: %s\n", target)
			for _, name := range collectorsToRun {
				if dryRun {
					fmt.Printf("[DRY-RUN] Would run %s\n", name)
					continue
				}
				_, err := collector.RunAndSave(name, caseID, target, activeAllowed, nil)
				if err != nil {
					return err
				}
			}
			return nil
		}

		m := collectModel{

			target:        target,
			collectors:    collectorsToRun,
			completed:     make(map[string]bool),
			failed:        make(map[string]error),
			progress:      progress.New(progress.WithDefaultGradient()),
			spinner:       spinner.New(spinner.WithSpinner(spinner.Dot)),
			activeAllowed: activeAllowed,
			caseID:        caseID,
			dryRun:        dryRun,
			threads:       threads,
		}

		if _, err := tea.NewProgram(m).Run(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	collectCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulate collection without making network requests")
	collectCmd.Flags().IntVar(&threads, "threads", 5, "Number of concurrent collectors to run")
	rootCmd.AddCommand(collectCmd)
}
