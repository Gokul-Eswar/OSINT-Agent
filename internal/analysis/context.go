package analysis

import (
	"fmt"
	"strings"

	"github.com/spectre/spectre/internal/core"
	"github.com/spectre/spectre/internal/storage"
)

// PromptBuilder materializes a stable text context from graph state.
//
// The generated format is designed for model readability and deterministic
// hashing so AnalyzeCase can safely cache synthesis results.
type PromptBuilder struct {
	Case          *core.Case
	Entities      []*core.Entity
	Relationships []*core.Relationship
	Evidence      []*core.Evidence
}

// NewPromptBuilder loads all case-side primitives required to build a prompt.
//
// Keeping data collection centralized here prevents drift between chat/query
// and synthesis flows that depend on the same contextual snapshot.
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

// Build serializes case metadata, entities, relationships, and evidence into a
// single prompt string consumed by synthesis and question-answer tasks.
//
// Relationship rendering first resolves entity IDs to readable labels, then
// falls back to raw IDs when entities are missing so partial graph corruption
// does not block analysis.
func (pb *PromptBuilder) Build() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("CASE NAME: %s\n", pb.Case.Name))
	sb.WriteString(fmt.Sprintf("CASE DESCRIPTION: %s\n", pb.Case.Description))
	sb.WriteString(fmt.Sprintf("CREATED AT: %s\n\n", pb.Case.CreatedAt.Format("2006-01-02")))

	sb.WriteString("=== ENTITIES INVOLVED ===\n")
	entityMap := make(map[string]string)
	for _, e := range pb.Entities {
		sb.WriteString(fmt.Sprintf("- [%s] %s (Source: %s)\n", e.Type, e.Value, e.Source))
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

// BuildCaseContext is a convenience entrypoint used by analysis workflows to
// produce the canonical prompt payload for the analyzer bridge.
func BuildCaseContext(caseID string) (string, error) {
	pb, err := NewPromptBuilder(caseID)
	if err != nil {
		return "", err
	}
	return pb.Build(), nil
}

// ExportCaseForViz exports case primitives in a bridge-friendly map for
// visualization tasks.
//
// The payload keys are consumed by the Python visualizer module and should be
// treated as a compatibility contract across Go/Python boundaries.
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
