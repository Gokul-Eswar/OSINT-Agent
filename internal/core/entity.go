package core

import "time"

// Entity represents a single intelligence node (e.g., IP, Domain, Person).
type Entity struct {
	ID           string                 `json:"id"`
	CaseID       string                 `json:"case_id"`
	Type         string                 `json:"type"`
	Value        string                 `json:"value"`
	Source       string                 `json:"source"`
	Confidence   float64                `json:"confidence"`
	DiscoveredAt time.Time              `json:"discovered_at"`
	Metadata     map[string]interface{} `json:"metadata"`
}
