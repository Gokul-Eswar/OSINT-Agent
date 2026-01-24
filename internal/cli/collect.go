package cli

import (
	"fmt"

	"github.com/spectre/spectre/internal/collector"
	_ "github.com/spectre/spectre/internal/collector/dns"    // Register DNS
	_ "github.com/spectre/spectre/internal/collector/whois"  // Register WHOIS
	_ "github.com/spectre/spectre/internal/collector/github" // Register GitHub
	_ "github.com/spectre/spectre/internal/collector/geo"    // Register GeoIP
	"github.com/spectre/spectre/internal/storage"
	"github.com/spf13/cobra"
)

var collectCmd = &cobra.Command{
	Use:   "collect [collector] [target]",
	Short: "Run a passive collector against a target",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if caseID == "" {
			return fmt.Errorf("case ID is required (use --case)")
		}

		collectorName := args[0]
		target := args[1]

		if err := storage.InitDB(); err != nil {
			return err
		}

		fmt.Printf("Running collector '%s' against target '%s'...\n", collectorName, target)
		evidenceList, err := collector.Run(collectorName, caseID, target)
		if err != nil {
			return fmt.Errorf("collection failed: %w", err)
		}

		for _, ev := range evidenceList {
			if err := storage.CreateEvidence(&ev); err != nil {
				return fmt.Errorf("failed to save evidence: %w", err)
			}
			fmt.Printf("Saved evidence: %s\n", ev.FilePath)
			
			// Auto-Ingestion (Next step)
			if err := storage.IngestEvidence(&ev); err != nil {
				fmt.Printf("Warning: auto-ingestion failed: %v\n", err)
			}
		}

		return nil
	},
}

func init() {
	collectCmd.Flags().StringVarP(&caseID, "case", "c", "", "Case ID (required)")
	rootCmd.AddCommand(collectCmd)
}
