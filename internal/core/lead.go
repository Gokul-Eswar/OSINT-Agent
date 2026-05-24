package core

import "time"

// IntelligenceLead represents an agent-generated hypothesis or lead.
type IntelligenceLead struct {
	ID          string    `json:"id"`
	CaseID      string    `json:"case_id"`
	Hypothesis  string    `json:"hypothesis"`
	Confidence  float64   `json:"confidence"`
	EvidenceIDs []string  `json:"evidence_ids"` // References to supporting evidence records
	Status      string    `json:"status"`       // "active", "verified", "refuted"
	CreatedAt   time.Time `json:"created_at"`
}
