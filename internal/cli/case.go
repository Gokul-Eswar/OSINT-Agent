package cli

import (
	"fmt"

	"github.com/spectre/spectre/internal/core"
	"github.com/spectre/spectre/internal/storage"
	"github.com/spf13/cobra"
)

var caseCmd = &cobra.Command{
	Use:   "case",
	Short: "Manage investigation cases",
}

var newCaseCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new investigation case",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		
		if err := storage.InitDB(); err != nil {
			return err
		}
		
		c := &core.Case{
			Name: name,
		}
		
		if err := storage.CreateCase(c); err != nil {
			return err
		}
		
		// Save context
		if err := SaveContext(c.ID); err != nil {
			fmt.Printf("Warning: failed to save context: %v\n", err)
		}

		fmt.Printf("Successfully created case: %s (ID: %s)\n", c.Name, c.ID)
		return nil
	},
}

func init() {
	caseCmd.AddCommand(newCaseCmd)
	rootCmd.AddCommand(caseCmd)
}
