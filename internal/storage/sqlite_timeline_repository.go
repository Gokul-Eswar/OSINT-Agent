package storage

import (
	"fmt"
	"sort"

	"github.com/spectre/spectre/internal/core"
)

// GetCaseTimeline aggregates entities and evidence into a sorted timeline.
func GetCaseTimeline(caseID string) ([]core.TimelineEvent, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var events []core.TimelineEvent

	// 1. Get Entities
	entities, err := ListEntitiesByCase(caseID)
	if err != nil {
		return nil, err
	}
	for _, e := range entities {
		events = append(events, core.TimelineEvent{
			Timestamp:   e.DiscoveredAt,
			Type:        "entity_discovered",
			Description: fmt.Sprintf("Discovered %s: %s", e.Type, e.Value),
			Source:      e.Source,
		})
	}

	// 2. Get Evidence
	evidence, err := ListEvidenceByCase(caseID)
	if err != nil {
		return nil, err
	}
	for _, ev := range evidence {
		events = append(events, core.TimelineEvent{
			Timestamp:   ev.CollectedAt,
			Type:        "evidence_collected",
			Description: fmt.Sprintf("Collected evidence from %s: %s", ev.Collector, ev.FilePath),
			Source:      ev.Collector,
		})
	}

	// 3. Sort by timestamp
	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp.Before(events[j].Timestamp)
	})

	return events, nil
}
