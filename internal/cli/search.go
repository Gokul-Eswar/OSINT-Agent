package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Gokul-Eswar/Spectre/internal/analyzer"
	"github.com/Gokul-Eswar/Spectre/internal/storage"
	"github.com/spf13/cobra"
)

var semanticSearch bool

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for a string or concept across evidence in the current case",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")

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

		if semanticSearch {
			fmt.Printf("Performing semantic vector search for '%s' in case %s...\n\n", query, caseID)
			req := analyzer.Request{
				Task:   "search_evidence",
				CaseID: caseID,
				Data: map[string]interface{}{
					"case_id": caseID,
					"query":   query,
				},
			}

			output, err := analyzer.GlobalTaskRunner.Run(req)
			if err != nil {
				return fmt.Errorf("semantic search failed: %w", err)
			}

			var resp struct {
				Status  string `json:"status"`
				Error   string `json:"error"`
				Results []struct {
					ID       string                 `json:"id"`
					Content  string                 `json:"content"`
					Distance float64                `json:"distance"`
					Metadata map[string]interface{} `json:"metadata"`
				} `json:"results"`
			}

			if err := json.Unmarshal([]byte(output), &resp); err != nil {
				return fmt.Errorf("failed to parse search results: %w", err)
			}

			if resp.Status == "error" {
				return fmt.Errorf("semantic search error: %s", resp.Error)
			}

			if len(resp.Results) == 0 {
				fmt.Println("No semantically relevant matches found in vector store.")
				return nil
			}

			fmt.Printf("Found %d semantic match(es):\n\n", len(resp.Results))
			for i, r := range resp.Results {
				fmt.Printf("[%d] Chunk: %s\n", i+1, r.ID)
				fmt.Printf("    Snippet: %s\n\n", strings.TrimSpace(r.Content))
			}
			return nil
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
	searchCmd.Flags().BoolVarP(&semanticSearch, "semantic", "s", false, "Use vector-based semantic search across evidence")
	rootCmd.AddCommand(searchCmd)
}
