package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/spectre/spectre/internal/storage"
)

type item struct {
	id, title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// FetchCases gets all cases from the DB and converts them to list items
func FetchCases() ([]list.Item, error) {
	cases, err := storage.ListCases()
	if err != nil {
		return nil, err
	}

	var items []list.Item
	for _, c := range cases {
		items = append(items, item{
			id:    c.ID,
			title: c.Name,
			desc:  fmt.Sprintf("ID: %s | Status: %s", c.ID, c.Status),
		})
	}
	return items, nil
}