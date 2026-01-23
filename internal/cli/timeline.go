package cli

import (
	"fmt"

	"github.com/spectre/spectre/internal/storage"
	"github.com/spf13/cobra"
)

var timelineCmd = &cobra.Command{
	Use:   "timeline",
	Short: "Show the chronological timeline of an investigation",
	RunE: func(cmd *cobra.Command, args []string) error {
		if caseID == "" {
			return fmt.Errorf("case ID is required (use --case)")
		}

		if err := storage.InitDB(); err != nil {
			return err
		}

		events, err := storage.GetCaseTimeline(caseID)
		if err != nil {
			return err
		}

		if len(events) == 0 {
			fmt.Printf("No events found for case %s\n", caseID)
			return nil
		}

		fmt.Printf("Timeline for case %s:\n", caseID)
		fmt.Println("--------------------------------------------------------------------------------")
		for _, ev := range events {
			timestamp := ev.Timestamp.Format("2006-01-02 15:04:05")
			fmt.Printf("[%s] %-20s | %s (%s)\n", timestamp, ev.Type, ev.Description, ev.Source)
		}

		return nil
	},
}

func init() {
	timelineCmd.Flags().StringVarP(&caseID, "case", "c", "", "Case ID (required)")
	rootCmd.AddCommand(timelineCmd)
}
