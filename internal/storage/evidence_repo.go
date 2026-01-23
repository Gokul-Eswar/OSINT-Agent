package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spectre/spectre/internal/core"
)

// CreateEvidence inserts a new evidence record into the database.
func CreateEvidence(ev *core.Evidence) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	if ev.ID == "" {
		ev.ID = uuid.New().String()
	}
	if ev.CollectedAt.IsZero() {
		ev.CollectedAt = time.Now()
	}

	metadataJSON, err := json.Marshal(ev.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `INSERT INTO evidence (id, case_id, entity_id, collector, file_path, file_hash, collected_at, metadata) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = DB.Exec(query, ev.ID, ev.CaseID, ev.EntityID, ev.Collector, ev.FilePath, ev.FileHash, ev.CollectedAt, string(metadataJSON))
	if err != nil {
		return fmt.Errorf("failed to create evidence: %w", err)
	}

	return nil
}

// ListEvidenceByCase retrieves all evidence for a specific case.
func ListEvidenceByCase(caseID string) ([]*core.Evidence, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := `SELECT id, case_id, entity_id, collector, file_path, file_hash, collected_at, metadata FROM evidence WHERE case_id = ?`
	rows, err := DB.Query(query, caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list evidence: %w", err)
	}
	defer rows.Close()

	var evidenceList []*core.Evidence
	for rows.Next() {
		var ev core.Evidence
		var entityID sql.NullString
		var metadataStr string
		if err := rows.Scan(&ev.ID, &ev.CaseID, &entityID, &ev.Collector, &ev.FilePath, &ev.FileHash, &ev.CollectedAt, &metadataStr); err != nil {
			return nil, fmt.Errorf("failed to scan evidence: %w", err)
		}
		
		if entityID.Valid {
			ev.EntityID = entityID.String
		}

		if err := json.Unmarshal([]byte(metadataStr), &ev.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
		evidenceList = append(evidenceList, &ev)
	}

	return evidenceList, nil
}
