package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spectre/spectre/internal/report"
	"github.com/spectre/spectre/internal/storage"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a comprehensive Markdown report",
	RunE: func(cmd *cobra.Command, args []string) error {
		if caseID == "" {
			return fmt.Errorf("case ID is required (use --case)")
		}

		if err := storage.InitDB(); err != nil {
			return err
		}

		fmt.Printf("Generating report for case %s...\n", caseID)
		md, err := report.GenerateMarkdownReport(caseID)
		if err != nil {
			return err
		}

		outputDir := filepath.Join("evidence_storage", caseID)
		os.MkdirAll(outputDir, 0755)
		outputPath := filepath.Join(outputDir, "investigation_report.md")

		err = os.WriteFile(outputPath, []byte(md), 0644)
		if err != nil {
			return fmt.Errorf("failed to save report: %w", err)
		}

		fmt.Printf("Report successfully generated: %s\n", outputPath)
		fmt.Println("\n--- PREVIEW ---")
		
		// Let's just print a short header instead of complex logic
		if len(md) > 200 {
			fmt.Println(md[:200] + "...")
		} else {
			fmt.Println(md)
		}
		
		return nil
	},
}

func init() {
	reportCmd.Flags().StringVarP(&caseID, "case", "c", "", "Case ID (required)")
	rootCmd.AddCommand(reportCmd)
}
