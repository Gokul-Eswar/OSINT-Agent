package tui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/spectre/spectre/internal/storage"
)

// FetchEntities fetches all entities for a case and converts them to table rows
func FetchEntities(caseID string) ([]table.Row, error) {
	entities, err := storage.ListEntitiesByCase(caseID)
	if err != nil {
		return nil, err
	}

	var rows []table.Row
	for _, e := range entities {
		rows = append(rows, table.Row{
			e.Type,
			e.Value,
			e.Source,
		})
	}
	return rows, nil
}

func NewEntityTable() table.Model {
	columns := []table.Column{
		{Title: "Type", Width: 10},
		{Title: "Value", Width: 30},
		{Title: "Source", Width: 15},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}
