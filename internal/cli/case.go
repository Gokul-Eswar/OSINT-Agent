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

var caseListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all investigation cases",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := storage.InitDB(); err != nil {
			return err
		}

		cases, err := storage.ListCases()
		if err != nil {
			return err
		}

		if len(cases) == 0 {
			fmt.Println("No cases found.")
			return nil
		}

		fmt.Printf("%-36s | %-20s | %-10s\n", "ID", "NAME", "STATUS")
		fmt.Println("----------------------------------------------------------------------")
		for _, c := range cases {
			fmt.Printf("%-36s | %-20s | %-10s\n", c.ID, c.Name, c.Status)
		}

		return nil
	},
}

func init() {
	caseCmd.AddCommand(newCaseCmd)
	caseCmd.AddCommand(caseListCmd)
	rootCmd.AddCommand(caseCmd)
}
