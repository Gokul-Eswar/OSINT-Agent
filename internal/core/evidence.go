package core

import "time"

// Evidence represents a piece of data collected by a collector.
type Evidence struct {
	ID          string                 `json:"id"`
	CaseID      string                 `json:"case_id"`
	EntityID    string                 `json:"entity_id"`
	Collector   string                 `json:"collector"`
	FilePath    string                 `json:"file_path"`
	FileHash    string                 `json:"file_hash"`
	CollectedAt time.Time              `json:"collected_at"`
	Metadata    map[string]interface{} `json:"metadata"`
	RawData     interface{}            `json:"-"`
}
