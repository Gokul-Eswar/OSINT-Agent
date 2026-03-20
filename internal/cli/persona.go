package cli

import (
	"fmt"
	"strings"

	"github.com/spectre/spectre/internal/core"
	"github.com/spectre/spectre/internal/storage"
	"github.com/spf13/cobra"
)

var personaCmd = &cobra.Command{
	Use:   "persona",
	Short: "Manage digital personas (correlated identities)",
}

var personaMapCmd = &cobra.Command{
	Use:   "map [username]",
	Short: "Correlate a username and its linked accounts into a single persona",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]
		
		if err := storage.InitDB(); err != nil {
			return err
		}

		if caseID == "" {
			ctxID, err := LoadContext()
			if err == nil && ctxID != "" {
				caseID = ctxID
				fmt.Printf("Using current case: %s\n", caseID)
			}
		}

		if caseID == "" {
			return fmt.Errorf("case ID is required")
		}

		// 1. Get the username entity
		userEnt, err := storage.GetEntityByValue(caseID, username)
		if err != nil {
			return err
		}
		if userEnt == nil {
			return fmt.Errorf("username entity '%s' not found in case %s", username, caseID)
		}

		// 2. Check for existing persona
		personaName := "Persona: " + username
		personaEnt, err := storage.GetEntityByValue(caseID, personaName)
		if err != nil {
			return err
		}

		if personaEnt == nil {
			fmt.Printf("[*] Creating new persona: %s\n", personaName)
			personaEnt = &core.Entity{
				CaseID: caseID,
				Type:   "persona",
				Value:  personaName,
				Source: "persona_mapping",
			}
			if err := storage.CreateEntity(personaEnt); err != nil {
				return err
			}
		} else {
			fmt.Printf("[*] Using existing persona: %s\n", personaName)
		}

		// 3. Link Username to Persona
		rel := &core.Relationship{
			CaseID:       caseID,
			FromEntityID: userEnt.ID,
			ToEntityID:   personaEnt.ID,
			Type:         "part_of_persona",
		}
		storage.CreateRelationship(rel)

		// 4. Find all linked accounts and link them to persona
		relationships, err := storage.ListRelationshipsByFromEntity(userEnt.ID)
		if err != nil {
			return err
		}

		mappedCount := 0
		for _, r := range relationships {
			if r.Type == "has_account" {
				// Link Account -> part_of_persona -> Persona
				accRel := &core.Relationship{
					CaseID:       caseID,
					FromEntityID: r.ToEntityID,
					ToEntityID:   personaEnt.ID,
					Type:         "part_of_persona",
				}
				if err := storage.CreateRelationship(accRel); err == nil {
					mappedCount++
				}
			}
		}

		fmt.Printf("✅ Persona mapping complete. Correlated %d accounts to '%s'.\n", mappedCount, personaName)
		return nil
	},
}

var personaListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all mapped personas in the current case",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := storage.InitDB(); err != nil {
			return err
		}

		if caseID == "" {
			ctxID, err := LoadContext()
			if err == nil && ctxID != "" {
				caseID = ctxID
			}
		}

		if caseID == "" {
			return fmt.Errorf("case ID is required")
		}

		personas, err := storage.ListEntitiesByType(caseID, "persona")
		if err != nil {
			return err
		}

		if len(personas) == 0 {
			fmt.Println("No personas found in this case.")
			return nil
		}

		fmt.Printf("PERSONAS IN CASE %s:\n", caseID)
		fmt.Println(strings.Repeat("─", 50))
		for _, p := range personas {
			fmt.Printf("- %s (ID: %s)\n", p.Value, p.ID[:8])
		}

		return nil
	},
}

func init() {
	personaCmd.AddCommand(personaMapCmd)
	personaCmd.AddCommand(personaListCmd)
	rootCmd.AddCommand(personaCmd)
}
