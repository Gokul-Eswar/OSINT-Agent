package cli

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/spectre/spectre/internal/analysis"
	"github.com/spectre/spectre/internal/analyzer"
	"github.com/spectre/spectre/internal/storage"
	"github.com/spf13/cobra"
)

var visualizeCmd = &cobra.Command{
	Use:   "visualize",
	Short: "Generate and open the interactive intelligence graph",
	RunE: func(cmd *cobra.Command, args []string) error {
		if caseID == "" {
			ctxID, err := LoadContext()
			if err == nil && ctxID != "" {
				caseID = ctxID
				fmt.Printf("Using current case: %s\n", caseID)
			}
		}

		if caseID == "" {
			return fmt.Errorf("case ID is required (use --case)")
		}

		if err := storage.InitDB(); err != nil {
			return err
		}

		fmt.Printf("Generating visualization for case %s...\n", caseID)

		// 1. Export Data
		data, err := analysis.ExportCaseForViz(caseID)
		if err != nil {
			return err
		}

		// 2. Run Python Visualizer
		req := analyzer.Request{
			Task:     "visualize",
			CaseID:   caseID,
			CaseName: data["case_name"].(string),
			Data:     data,
		}

		responseJSON, err := analyzer.RunPythonTask(req)
		if err != nil {
			return fmt.Errorf("visualization failed: %w", err)
		}

		var result map[string]string
		if err := json.Unmarshal([]byte(responseJSON), &result); err != nil {
			return fmt.Errorf("failed to parse visualizer response: %w", err)
		}

		filePath := result["file_path"]
		if filePath == "" {
			return fmt.Errorf("visualizer did not return a file path")
		}

		fmt.Printf("Dashboard generated: %s\n", filePath)
		fmt.Println("Opening in browser...")

		// 3. Open in Browser (Windows specific)
		return exec.Command("cmd", "/c", "start", filePath).Start()
	},
}

func init() {
	visualizeCmd.Flags().StringVarP(&caseID, "case", "c", "", "Case ID (required)")
	rootCmd.AddCommand(visualizeCmd)
}
