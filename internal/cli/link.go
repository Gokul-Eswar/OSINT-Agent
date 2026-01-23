package cli

import (
	"fmt"

	"github.com/spectre/spectre/internal/core"
	"github.com/spectre/spectre/internal/storage"
	"github.com/spf13/cobra"
)

var relType string

var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Manage relationships between entities",
}

var linkAddCmd = &cobra.Command{
	Use:   "add [source_val] [target_val]",
	Short: "Link two entities within a case",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if caseID == "" {
			return fmt.Errorf("case ID is required (use --case)")
		}

		if relType == "" {
			return fmt.Errorf("relationship type is required (use --type)")
		}

		if err := storage.InitDB(); err != nil {
			return err
		}

		sourceVal := args[0]
		targetVal := args[1]

		sourceEnt, err := storage.GetEntityByValue(caseID, sourceVal)
		if err != nil {
			return err
		}
		if sourceEnt == nil {
			return fmt.Errorf("source entity '%s' not found in case %s", sourceVal, caseID)
		}

		targetEnt, err := storage.GetEntityByValue(caseID, targetVal)
		if err != nil {
			return err
		}
		if targetEnt == nil {
			return fmt.Errorf("target entity '%s' not found in case %s", targetVal, caseID)
		}

		rel := &core.Relationship{
			CaseID:       caseID,
			FromEntityID: sourceEnt.ID,
			ToEntityID:   targetEnt.ID,
			Type:         relType,
			Confidence:   1.0,
		}

		if err := storage.CreateRelationship(rel); err != nil {
			return err
		}

		fmt.Printf("Successfully linked %s -> %s (Type: %s) in case %s\n", sourceVal, targetVal, relType, caseID)
		return nil
	},
}

var linkListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all links in a case",
	RunE: func(cmd *cobra.Command, args []string) error {
		if caseID == "" {
			return fmt.Errorf("case ID is required (use --case)")
		}

		if err := storage.InitDB(); err != nil {
			return err
		}

		rels, err := storage.ListRelationshipsByCase(caseID)
		if err != nil {
			return err
		}

		if len(rels) == 0 {
			fmt.Printf("No links found for case %s\n", caseID)
			return nil
		}

		fmt.Printf("Links for case %s:\n", caseID)
		fmt.Printf("% -36s | % -36s | % -20s\n", "FROM ID", "TO ID", "TYPE")
		fmt.Println("---------------------------------------------------------------------------------------------------")
		for _, r := range rels {
			fmt.Printf("% -36s | % -36s | % -20s\n", r.FromEntityID, r.ToEntityID, r.Type)
		}

		return nil
	},
}

func init() {
	linkCmd.PersistentFlags().StringVarP(&caseID, "case", "c", "", "Case ID (required)")
	linkAddCmd.Flags().StringVarP(&relType, "type", "t", "", "Relationship type (required)")
	
	linkCmd.AddCommand(linkAddCmd)
	linkCmd.AddCommand(linkListCmd)
	rootCmd.AddCommand(linkCmd)
}