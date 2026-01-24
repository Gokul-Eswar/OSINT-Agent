package tui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/spectre/spectre/internal/storage"
)

func FetchRelationships(caseID string) ([]table.Row, error) {
	rels, err := storage.ListRelationshipsByCase(caseID)
	if err != nil {
		return nil, err
	}

	// We need a map for entity values to make the table readable
	entities, _ := storage.ListEntitiesByCase(caseID)
	entMap := make(map[string]string)
	for _, e := range entities {
		entMap[e.ID] = e.Value
	}

	var rows []table.Row
	for _, r := range rels {
		from := entMap[r.FromEntityID]
		if from == "" { from = r.FromEntityID }
		to := entMap[r.ToEntityID]
		if to == "" { to = r.ToEntityID }

		rows = append(rows, table.Row{
			from,
			r.Type,
			to,
		})
	}
	return rows, nil
}

func NewRelationshipTable() table.Model {
	columns := []table.Column{
		{Title: "Source", Width: 20},
		{Title: "Type", Width: 15},
		{Title: "Target", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.Bold(true)
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57"))
	t.SetStyles(s)

	return t
}
