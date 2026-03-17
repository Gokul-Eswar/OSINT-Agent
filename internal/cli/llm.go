package cli

import (
	"fmt"

	"github.com/spectre/spectre/internal/analysis"
	"github.com/spf13/cobra"
)

var llmCmd = &cobra.Command{
	Use:   "llm",
	Short: "Direct interaction with the intelligence layer",
}

var queryCmd = &cobra.Command{
	Use:   "query [question]",
	Short: "Ask a specific question about a case",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if caseID == "" {
			return fmt.Errorf("case ID is required (use --case)")
		}

		question := args[0]
		fmt.Printf("Querying case %s about: %s...\n", caseID, question)

		answer, err := analysis.QueryCase(caseID, modelName, question)
		if err != nil {
			return err
		}

		fmt.Printf("\n--- AI RESPONSE ---\n%s\n-------------------\n", answer)
		return nil
	},
}

func init() {
	queryCmd.Flags().StringVarP(&caseID, "case", "c", "", "Case ID (required)")
	queryCmd.Flags().StringVarP(&modelName, "model", "m", "llama3", "Model to use")
	
	llmCmd.AddCommand(queryCmd)
	rootCmd.AddCommand(llmCmd)
}
