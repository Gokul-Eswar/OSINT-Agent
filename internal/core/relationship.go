package core

import "time"

// Relationship represents a connection between two entities.
type Relationship struct {
	ID           string    `json:"id"`
	CaseID       string    `json:"case_id"`
	FromEntityID string    `json:"from_entity_id"`
	ToEntityID   string    `json:"to_entity_id"`
	Type         string    `json:"type"`
	Confidence   float64   `json:"confidence"`
	EvidenceID   string    `json:"evidence_id"`
	DiscoveredAt time.Time `json:"discovered_at"`
}
