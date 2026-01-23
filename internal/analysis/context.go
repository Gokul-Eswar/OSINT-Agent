package analysis

import (
	"fmt"
	"strings"

	"github.com/spectre/spectre/internal/storage"
)

// BuildCaseContext aggregates all case data into a prompt-ready string.
func BuildCaseContext(caseID string) (string, error) {
	c, err := storage.GetCase(caseID)
	if err != nil {
		return "", fmt.Errorf("failed to get case: %w", err)
	}
	if c == nil {
		return "", fmt.Errorf("case not found")
	}

	entities, err := storage.ListEntitiesByCase(caseID)
	if err != nil {
		return "", fmt.Errorf("failed to list entities: %w", err)
	}

	rels, err := storage.ListRelationshipsByCase(caseID)
	if err != nil {
		return "", fmt.Errorf("failed to list relationships: %w", err)
	}

	evidence, err := storage.ListEvidenceByCase(caseID)
	if err != nil {
		return "", fmt.Errorf("failed to list evidence: %w", err)
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("CASE: %s\n", c.Name))
	sb.WriteString(fmt.Sprintf("DESCRIPTION: %s\n\n", c.Description))

	sb.WriteString("ENTITIES:\n")
	entityMap := make(map[string]string)
	for _, e := range entities {
		sb.WriteString(fmt.Sprintf("- [%s] %s (Source: %s)\n", e.Type, e.Value, e.Source))
		entityMap[e.ID] = fmt.Sprintf("%s (%s)", e.Value, e.Type)
	}
	sb.WriteString("\n")

	sb.WriteString("RELATIONSHIPS:\n")
	for _, r := range rels {
		from := entityMap[r.FromEntityID]
		to := entityMap[r.ToEntityID]
		if from == "" { from = r.FromEntityID }
		if to == "" { to = r.ToEntityID }
		
		sb.WriteString(fmt.Sprintf("- %s --[%s]--> %s\n", from, r.Type, to))
	}
	sb.WriteString("\n")

	sb.WriteString("EVIDENCE:\n")
	for _, ev := range evidence {
		sb.WriteString(fmt.Sprintf("- %s (Collector: %s)\n", ev.FilePath, ev.Collector))
	}

	return sb.String(), nil
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

	return map[string]interface{}{
		"case_id":       caseID,
		"case_name":     c.Name,
		"entities":      entities,
		"relationships": rels,
	}, nil
}
