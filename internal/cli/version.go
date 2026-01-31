package cli

import (
	"fmt"

	"github.com/spectre/spectre/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of SPECTRE",
	Long:  `All software has versions. This is SPECTRE's.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("SPECTRE %s\n", version.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

