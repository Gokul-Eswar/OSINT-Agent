package cli

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/spectre/spectre/internal/storage"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for a string across all evidence in the current case",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		// Try to load context if caseID is missing
		if caseID == "" {
			ctxID, err := LoadContext()
			if err == nil && ctxID != "" {
				caseID = ctxID
			}
		}

		if caseID == "" {
			return fmt.Errorf("case ID is required (use --case or create a new case)")
		}

		if err := storage.InitDB(); err != nil {
			return err
		}

		fmt.Printf("Searching for '%s' in case %s...\n\n", query, caseID)

		evidenceList, err := storage.ListEvidenceByCase(caseID)
		if err != nil {
			return err
		}

		foundCount := 0
		for _, ev := range evidenceList {
			if ev.FilePath == "" {
				continue
			}

			data, err := os.ReadFile(ev.FilePath)
			if err != nil {
				// File might have been moved or deleted
				continue
			}

			if bytes.Contains(data, []byte(query)) {
				foundCount++
				fmt.Printf("[+] Found in: %s (%s)\n", ev.FilePath, ev.Collector)

				// Show a snippet of the context
				lines := strings.Split(string(data), "\n")
				for _, line := range lines {
					if strings.Contains(line, query) {
						trimmed := strings.TrimSpace(line)
						if len(trimmed) > 100 {
							idx := strings.Index(trimmed, query)
							start := idx - 20
							if start < 0 {
								start = 0
							}
							end := idx + len(query) + 40
							if end > len(trimmed) {
								end = len(trimmed)
							}
							fmt.Printf("    ...%s...\n", trimmed[start:end])
						} else {
							fmt.Printf("    %s\n", trimmed)
						}
					}
				}
				fmt.Println()
			}
		}

		if foundCount == 0 {
			fmt.Println("No matches found.")
		} else {
			fmt.Printf("Total files with matches: %d\n", foundCount)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
