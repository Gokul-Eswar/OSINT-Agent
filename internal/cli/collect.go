package cli

import (
	"fmt"
	"sync"

	"github.com/spectre/spectre/internal/collector"
	_ "github.com/spectre/spectre/internal/collector/dns"    // Register DNS
	_ "github.com/spectre/spectre/internal/collector/whois"  // Register WHOIS
	_ "github.com/spectre/spectre/internal/collector/github" // Register GitHub
	_ "github.com/spectre/spectre/internal/collector/geo"    // Register GeoIP
	_ "github.com/spectre/spectre/internal/collector/active" // Register Active Probes
	"github.com/spectre/spectre/internal/storage"
	"github.com/spf13/cobra"
)

var collectCmd = &cobra.Command{
	Use:   "collect [collector|all] [target]",
	Short: "Run a passive collector (or all) against a target",
	Args:  cobra.ExactArgs(2),
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
		target := args[1]

		if err := storage.InitDB(); err != nil {
			return err
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

		var wg sync.WaitGroup
		var printMu sync.Mutex

		fmt.Printf("Starting collection against '%s' with %d collectors...\n", target, len(collectorsToRun))

		for _, name := range collectorsToRun {
			wg.Add(1)
			go func(cName string) {
				defer wg.Done()

				// Execute
				evidenceList, err := collector.Run(cName, caseID, target, activeAllowed)
				
				printMu.Lock()
				defer printMu.Unlock()

				if err != nil {
					fmt.Printf("[X] %s: Failed - %v\n", cName, err)
					return
				}

				fmt.Printf("[+] %s: Completed (%d evidence items)\n", cName, len(evidenceList))

				for _, ev := range evidenceList {
					if err := storage.CreateEvidence(&ev); err != nil {
						fmt.Printf("    - Failed to save evidence: %v\n", err)
						continue
					}
					// fmt.Printf("    - Saved: %s\n", ev.FilePath) // Reduced verbosity for "fluidity"
					
					if err := storage.IngestEvidence(&ev); err != nil {
						fmt.Printf("    - Warning: ingestion failed: %v\n", err)
					}
				}
			}(name)
		}

		wg.Wait()
		fmt.Println("Collection complete.")
		return nil
	},
}

func init() {
	collectCmd.Flags().StringVarP(&caseID, "case", "c", "", "Case ID (required)")
	rootCmd.AddCommand(collectCmd)
}
