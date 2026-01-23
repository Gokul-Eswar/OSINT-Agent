package cli

import (
	"fmt"

	"github.com/spectre/spectre/internal/core"
	"github.com/spectre/spectre/internal/storage"
	"github.com/spf13/cobra"
)

var caseID string

var entityCmd = &cobra.Command{
	Use:   "entity",
	Short: "Manage entities within a case",
}

var entityAddCmd = &cobra.Command{
	Use:   "add [type] [value]",
	Short: "Add a new entity to a case",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if caseID == "" {
			return fmt.Errorf("case ID is required (use --case)")
		}

		if err := storage.InitDB(); err != nil {
			return err
		}

		entityType := args[0]
		entityValue := args[1]

		e := &core.Entity{
			CaseID: caseID,
			Type:   entityType,
			Value:  entityValue,
			Source: "manual",
		}

		if err := storage.CreateEntity(e); err != nil {
			return err
		}

		fmt.Printf("Successfully added entity: %s (%s) to case %s\n", entityValue, entityType, caseID)
		return nil
	},
}

var entityListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all entities in a case",
	RunE: func(cmd *cobra.Command, args []string) error {
		if caseID == "" {
			return fmt.Errorf("case ID is required (use --case)")
		}

		if err := storage.InitDB(); err != nil {
			return err
		}

		entities, err := storage.ListEntitiesByCase(caseID)
		if err != nil {
			return err
		}

		if len(entities) == 0 {
			fmt.Printf("No entities found for case %s\n", caseID)
			return nil
		}

		fmt.Printf("Entities for case %s:\n", caseID)
		fmt.Printf("%-36s | %-15s | %-20s\n", "ID", "TYPE", "VALUE")
		fmt.Println("--------------------------------------------------------------------------------")
		for _, e := range entities {
			fmt.Printf("%-36s | %-15s | %-20s\n", e.ID, e.Type, e.Value)
		}

		return nil
	},
}

func init() {
	entityCmd.PersistentFlags().StringVarP(&caseID, "case", "c", "", "Case ID (required)")
	
	entityCmd.AddCommand(entityAddCmd)
	entityCmd.AddCommand(entityListCmd)
	rootCmd.AddCommand(entityCmd)
}
