package core

import "time"

// TimelineEvent represents a chronological event in an investigation.
type TimelineEvent struct {
	Timestamp   time.Time `json:"timestamp"`
	Type        string    `json:"type"`        // e.g., "entity_discovered", "evidence_collected"
	Description string    `json:"description"`
	Source      string    `json:"source"`      // e.g., "dns", "whois", "manual"
}
