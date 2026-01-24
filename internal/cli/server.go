package cli

import (
	"fmt"

	"github.com/spectre/spectre/internal/server"
	"github.com/spectre/spectre/internal/storage"
	"github.com/spf13/cobra"
)

var port int

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the SPECTRE API server and Web UI",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := storage.InitDB(); err != nil {
			return err
		}

		fmt.Printf("Launching Web Command Center...\n")
		return server.Start(port)
	},
}

func init() {
	serverCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to listen on")
	rootCmd.AddCommand(serverCmd)
}
