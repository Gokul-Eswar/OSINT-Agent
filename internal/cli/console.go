package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spectre/spectre/internal/server"
	"github.com/spectre/spectre/internal/storage"
	"github.com/spectre/spectre/internal/tui"
	"github.com/spf13/cobra"
)

var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Launch the interactive SPECTRE console",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := storage.InitDB(); err != nil {
			return err
		}

		// Start Server in Background
		go func() {
			if err := server.Start(8080); err != nil {
				fmt.Printf("Server error: %v\n", err)
			}
		}()

		p := tea.NewProgram(tui.InitialModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("failed to start TUI: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(consoleCmd)
}
