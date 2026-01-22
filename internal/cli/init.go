package cli

import (
	"github.com/spectre/spectre/internal/storage"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the SPECTRE database",
	Long:  `Initialize the SQLite database and apply necessary schemas for SPECTRE.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := storage.InitDB(); err != nil {
			return err
		}
		return storage.InitSchema()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
