package analysis

import (
	"fmt"
	"strings"

	"github.com/spectre/spectre/internal/storage"
	"github.com/spectre/spectre/internal/core"
)

// PromptBuilder helps in constructing complex, accurate LLM prompts
type PromptBuilder struct {
	Case          *core.Case
	Entities      []core.Entity
	Relationships []core.Relationship
	Evidence      []core.Evidence
}

// NewPromptBuilder creates a new PromptBuilder for a case
func NewPromptBuilder(caseID string) (*PromptBuilder, error) {
	c, err := storage.GetCase(caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get case: %w", err)
	}
	if c == nil {
		return nil, fmt.Errorf("case not found")
	}

	entities, err := storage.ListEntitiesByCase(caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list entities: %w", err)
	}

	rels, err := storage.ListRelationshipsByCase(caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list relationships: %w", err)
	}

	evidence, err := storage.ListEvidenceByCase(caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list evidence: %w", err)
	}

	return &PromptBuilder{
		Case:          c,
		Entities:      entities,
		Relationships: rels,
		Evidence:      evidence,
	}, nil
}

// Build creates the final context string for the LLM
func (pb *PromptBuilder) Build() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("CASE NAME: %s\n", pb.Case.Name))
	sb.WriteString(fmt.Sprintf("CASE DESCRIPTION: %s\n", pb.Case.Description))
	sb.WriteString(fmt.Sprintf("CREATED AT: %s\n\n", pb.Case.CreatedAt.Format("2006-01-02")))

	sb.WriteString("=== ENTITIES INVOLVED ===\n")
	entityMap := make(map[string]string)
	for _, e := range pb.Entities {
		sb.WriteString(fmt.Sprintf("- [%s] %s (Source: %s)\n", e.Type, e.Value, e.Source))
		if e.Notes != "" {
			sb.WriteString(fmt.Sprintf("  Notes: %s\n", e.Notes))
		}
		entityMap[e.ID] = fmt.Sprintf("%s (%s)", e.Value, e.Type)
	}
	sb.WriteString("\n")

	sb.WriteString("=== RELATIONSHIPS & CONNECTIONS ===\n")
	for _, r := range pb.Relationships {
		from := entityMap[r.FromEntityID]
		to := entityMap[r.ToEntityID]
		if from == "" {
			from = r.FromEntityID
		}
		if to == "" {
			to = r.ToEntityID
		}
		sb.WriteString(fmt.Sprintf("- %s --[%s]--> %s\n", from, r.Type, to))
	}
	sb.WriteString("\n")

	sb.WriteString("=== COLLECTED EVIDENCE ===\n")
	for _, ev := range pb.Evidence {
		sb.WriteString(fmt.Sprintf("- %s (Collector: %s)\n", ev.FilePath, ev.Collector))
	}
	
	sb.WriteString("\nANALYSIS INSTRUCTIONS: Please identify missing intelligence gaps and suggest additional collectors in the missing_data and suggested_collectors fields.")

	return sb.String()
}

// BuildCaseContext aggregates all case data into a prompt-ready string using PromptBuilder.
func BuildCaseContext(caseID string) (string, error) {
	pb, err := NewPromptBuilder(caseID)
	if err != nil {
		return "", err
	}
	return pb.Build(), nil
}

// ExportCaseForViz gathers all case data into a map for JSON export to the visualizer.
func ExportCaseForViz(caseID string) (map[string]interface{}, error) {
	c, err := storage.GetCase(caseID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, fmt.Errorf("case not found")
	}

	entities, err := storage.ListEntitiesByCase(caseID)
	if err != nil {
		return nil, err
	}

	rels, err := storage.ListRelationshipsByCase(caseID)
	if err != nil {
		return nil, err
	}

	evidence, err := storage.ListEvidenceByCase(caseID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"case_id":       caseID,
		"case_name":     c.Name,
		"entities":      entities,
		"relationships": rels,
		"evidence":      evidence,
	}, nil
}
