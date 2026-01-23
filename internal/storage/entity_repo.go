package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spectre/spectre/internal/core"
)

// CreateEntity inserts a new entity into the database.
func CreateEntity(e *core.Entity) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	if e.DiscoveredAt.IsZero() {
		e.DiscoveredAt = time.Now()
	}
	if e.Confidence == 0 {
		e.Confidence = 0.5
	}

	metadataJSON, err := json.Marshal(e.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `INSERT INTO entities (id, case_id, type, value, source, confidence, discovered_at, metadata) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = DB.Exec(query, e.ID, e.CaseID, e.Type, e.Value, e.Source, e.Confidence, e.DiscoveredAt, string(metadataJSON))
	if err != nil {
		return fmt.Errorf("failed to create entity: %w", err)
	}

	return nil
}

// GetEntity retrieves an entity by its ID.
func GetEntity(id string) (*core.Entity, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := `SELECT id, case_id, type, value, source, confidence, discovered_at, metadata FROM entities WHERE id = ?`
	row := DB.QueryRow(query, id)

	var e core.Entity
	var metadataStr string
	err := row.Scan(&e.ID, &e.CaseID, &e.Type, &e.Value, &e.Source, &e.Confidence, &e.DiscoveredAt, &metadataStr)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get entity: %w", err)
	}

	if err := json.Unmarshal([]byte(metadataStr), &e.Metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &e, nil
}

// ListEntitiesByCase retrieves all entities associated with a specific case.
func ListEntitiesByCase(caseID string) ([]*core.Entity, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := `SELECT id, case_id, type, value, source, confidence, discovered_at, metadata FROM entities WHERE case_id = ?`
	rows, err := DB.Query(query, caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list entities: %w", err)
	}
	defer rows.Close()

	var entities []*core.Entity
	for rows.Next() {
		var e core.Entity
		var metadataStr string
		if err := rows.Scan(&e.ID, &e.CaseID, &e.Type, &e.Value, &e.Source, &e.Confidence, &e.DiscoveredAt, &metadataStr); err != nil {
			return nil, fmt.Errorf("failed to scan entity: %w", err)
		}
		if err := json.Unmarshal([]byte(metadataStr), &e.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
		entities = append(entities, &e)
	}

	return entities, nil
}

// GetEntityByValue retrieves an entity by its value and case ID.
func GetEntityByValue(caseID, value string) (*core.Entity, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := `SELECT id, case_id, type, value, source, confidence, discovered_at, metadata 
	          FROM entities WHERE case_id = ? AND value = ?`
	row := DB.QueryRow(query, caseID, value)

	var e core.Entity
	var metadataStr string
	err := row.Scan(&e.ID, &e.CaseID, &e.Type, &e.Value, &e.Source, &e.Confidence, &e.DiscoveredAt, &metadataStr)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get entity by value: %w", err)
	}

	if err := json.Unmarshal([]byte(metadataStr), &e.Metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &e, nil
}
