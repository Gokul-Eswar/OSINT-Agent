package tui

import (
	"fmt"
	"strings"

	"github.com/spectre/spectre/internal/storage"
)

func RenderASCIIGraph(caseID string) string {
	if caseID == "" {
		return "No case selected."
	}

	entities, err := storage.ListEntitiesByCase(caseID)
	if err != nil {
		return "Error fetching entities: " + err.Error()
	}

	relationships, err := storage.ListRelationshipsByCase(caseID)
	if err != nil {
		return "Error fetching relationships: " + err.Error()
	}

	if len(entities) == 0 {
		return "No entities found for this case."
	}

	var s strings.Builder
	s.WriteString(fmt.Sprintf("INTELLIGENCE GRAPH — %s\n", caseID))
	s.WriteString("────────────────────────────────────\n\n")

	// Map to look up entity values by ID
	entityMap := make(map[string]string)
	for _, e := range entities {
		entityMap[e.ID] = e.Value
	}

	// Simple representation: list relationships as connected nodes
	if len(relationships) == 0 {
		s.WriteString("No relationships identified yet. Run more collectors.\n\n")
		s.WriteString("Entities:\n")
		for _, e := range entities {
			s.WriteString(fmt.Sprintf(" • [%s] %s\n", e.Type, e.Value))
		}
	} else {
		for _, r := range relationships {
			from := entityMap[r.FromEntityID]
			to := entityMap[r.ToEntityID]
			if from == "" { from = r.FromEntityID }
			if to == "" { to = r.ToEntityID }
			
			s.WriteString(fmt.Sprintf("[%s] ──(%s)──> [%s]\n", from, r.Type, to))
		}
	}

	return s.String()
}
