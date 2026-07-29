package cli

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/Gokul-Eswar/Spectre/internal/analysis"
	"github.com/Gokul-Eswar/Spectre/internal/analyzer"
	"github.com/Gokul-Eswar/Spectre/internal/storage"
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

		responseJSON, err := analyzer.GlobalTaskRunner.Run(req)
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

		// 3. Open in Browser (Cross-platform)
		return openURL(filePath)
	},
}

func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}
	return exec.Command(cmd, args...).Start()
}

func init() {
	visualizeCmd.Flags().StringVarP(&caseID, "case", "c", "", "Case ID (required)")
	rootCmd.AddCommand(visualizeCmd)
}
