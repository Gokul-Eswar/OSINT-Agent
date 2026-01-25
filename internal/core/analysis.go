package core

import "time"

// Analysis represents the AI-generated intelligence report for a case.
type Analysis struct {
	ID          string    `json:"id"`
	CaseID      string    `json:"case_id"`
	ContextHash string    `json:"context_hash"`
	Findings    []string  `json:"findings"`
	Risks       []string  `json:"risks"`
	Connections []string  `json:"connections"` // Suggested potential connections
	NextSteps   []string  `json:"next_steps"`
	Confidence  float64   `json:"confidence"`
	AnalyzedAt  time.Time `json:"analyzed_at"`
}
