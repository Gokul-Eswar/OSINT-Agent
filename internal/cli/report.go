package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spectre/spectre/internal/report"
	"github.com/spectre/spectre/internal/storage"
	"github.com/spf13/cobra"
)

var (
	format     string
	outputFile string
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a comprehensive report (Markdown or PDF)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if caseID == "" {
			return fmt.Errorf("case ID is required (use --case)")
		}

		if err := storage.InitDB(); err != nil {
			return err
		}

		outputDir := filepath.Join("evidence_storage", caseID)
		os.MkdirAll(outputDir, 0755)

		switch format {
		case "pdf":
			fmt.Printf("Generating PDF report for case %s...\n", caseID)
			targetPath := outputFile
			if targetPath == "" {
				targetPath = filepath.Join(outputDir, "investigation_report.pdf")
			}
			err := report.GeneratePDFReport(caseID, targetPath)
			if err != nil {
				return err
			}
			fmt.Printf("PDF Report successfully generated: %s\n", targetPath)
		case "html":
			fmt.Printf("Generating HTML report for case %s...\n", caseID)
			targetPath := outputFile
			if targetPath == "" {
				targetPath = filepath.Join(outputDir, "investigation_report.html")
			}
			err := report.GenerateHTMLReport(caseID, targetPath)
			if err != nil {
				return err
			}
			fmt.Printf("HTML Report successfully generated: %s\n", targetPath)
		case "json":
			fmt.Printf("Generating JSON report for case %s...\n", caseID)
			targetPath := outputFile
			if targetPath == "" {
				targetPath = filepath.Join(outputDir, "investigation_report.json")
			}
			err := report.GenerateJSONReport(caseID, targetPath)
			if err != nil {
				return err
			}
			fmt.Printf("JSON Report successfully generated: %s\n", targetPath)
		case "csv":
			fmt.Printf("Generating CSV report for case %s...\n", caseID)
			targetPath := outputFile
			if targetPath == "" {
				targetPath = filepath.Join(outputDir, "entities.csv")
			}
			err := report.GenerateCSVReport(caseID, targetPath)
			if err != nil {
				return err
			}
			fmt.Printf("CSV Report successfully generated: %s\n", targetPath)
		default:
			fmt.Printf("Generating Markdown report for case %s...\n", caseID)
			md, err := report.GenerateMarkdownReport(caseID)
			if err != nil {
				return err
			}
			targetPath := outputFile
			if targetPath == "" {
				targetPath = filepath.Join(outputDir, "investigation_report.md")
			}
			err = os.WriteFile(targetPath, []byte(md), 0644)
			if err != nil {
				return fmt.Errorf("failed to save report: %w", err)
			}
			fmt.Printf("Markdown Report successfully generated: %s\n", targetPath)
			fmt.Println("\n--- PREVIEW ---")
			if len(md) > 200 {
				fmt.Println(md[:200] + "...")
			} else {
				fmt.Println(md)
			}
		}

		return nil
	},
}

func init() {
	reportCmd.Flags().StringVarP(&caseID, "case", "c", "", "Case ID (required)")
	reportCmd.Flags().StringVarP(&format, "format", "f", "markdown", "Report format (markdown, pdf, html, json, csv)")
	reportCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Custom output file path")
	rootCmd.AddCommand(reportCmd)
}
