package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spectre/spectre/internal/core"
)

// OnEntityCreated is a hook for real-time updates
var OnEntityCreated func(*core.Entity)

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
		log.Error().Err(err).Str("entity_id", e.ID).Msg("failed to marshal entity metadata")
		return fmt.Errorf("failed to marshal metadata for entity %s: %w", e.ID, err)
	}

	query := TranslatePlaceholder(`INSERT INTO entities (id, case_id, type, value, source, confidence, discovered_at, metadata) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	_, err = DB.Exec(query, e.ID, e.CaseID, e.Type, e.Value, e.Source, e.Confidence, e.DiscoveredAt, string(metadataJSON))
	if err != nil {
		log.Error().
			Err(err).
			Str("entity_id", e.ID).
			Str("case_id", e.CaseID).
			Str("type", e.Type).
			Msg("failed to execute insert entity query")
		return fmt.Errorf("failed to create entity %s in case %s: %w", e.ID, e.CaseID, err)
	}

	if OnEntityCreated != nil {
		OnEntityCreated(e)
	}

	return nil
}

// GetEntity retrieves an entity by its ID.
func GetEntity(id string) (*core.Entity, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := TranslatePlaceholder(`SELECT id, case_id, type, value, source, confidence, discovered_at, metadata FROM entities WHERE id = ?`)
	row := DB.QueryRow(query, id)

	var e core.Entity
	var metadataStr string
	err := row.Scan(&e.ID, &e.CaseID, &e.Type, &e.Value, &e.Source, &e.Confidence, &e.DiscoveredAt, &metadataStr)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Error().Err(err).Str("entity_id", id).Msg("failed to scan entity row")
		return nil, fmt.Errorf("failed to get entity %s: %w", id, err)
	}

	if err := json.Unmarshal([]byte(metadataStr), &e.Metadata); err != nil {
		log.Error().Err(err).Str("entity_id", id).Msg("failed to unmarshal entity metadata")
		return nil, fmt.Errorf("failed to unmarshal metadata for entity %s: %w", id, err)
	}

	return &e, nil
}

// GetEntityByValue retrieves an entity by its case ID and value.
func GetEntityByValue(caseID, value string) (*core.Entity, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := TranslatePlaceholder(`SELECT id, case_id, type, value, source, confidence, discovered_at, metadata 
	          FROM entities WHERE case_id = ? AND value = ?`)
	row := DB.QueryRow(query, caseID, value)

	var e core.Entity
	var metadataStr string
	err := row.Scan(&e.ID, &e.CaseID, &e.Type, &e.Value, &e.Source, &e.Confidence, &e.DiscoveredAt, &metadataStr)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Error().Err(err).Str("case_id", caseID).Str("value", value).Msg("failed to scan entity row by value")
		return nil, fmt.Errorf("failed to get entity by value %s in case %s: %w", value, caseID, err)
	}

	if err := json.Unmarshal([]byte(metadataStr), &e.Metadata); err != nil {
		log.Error().Err(err).Str("entity_id", e.ID).Msg("failed to unmarshal entity metadata by value")
		return nil, fmt.Errorf("failed to unmarshal metadata for entity %s: %w", e.ID, err)
	}

	return &e, nil
}

// UpdateEntity updates an existing entity's fields (metadata, confidence).
func UpdateEntity(e *core.Entity) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	metadataJSON, err := json.Marshal(e.Metadata)
	if err != nil {
		log.Error().Err(err).Str("entity_id", e.ID).Msg("failed to marshal entity metadata for update")
		return fmt.Errorf("failed to marshal metadata for entity update %s: %w", e.ID, err)
	}

	query := TranslatePlaceholder(`UPDATE entities SET metadata = ?, confidence = ? WHERE id = ?`)
	_, err = DB.Exec(query, string(metadataJSON), e.Confidence, e.ID)
	if err != nil {
		log.Error().Err(err).Str("entity_id", e.ID).Msg("failed to execute update entity query")
		return fmt.Errorf("failed to update entity %s: %w", e.ID, err)
	}
	return nil
}

// ListEntitiesByCase retrieves all entities associated with a specific case.
func ListEntitiesByCase(caseID string) ([]*core.Entity, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := TranslatePlaceholder(`SELECT id, case_id, type, value, source, confidence, discovered_at, metadata FROM entities WHERE case_id = ?`)
	rows, err := DB.Query(query, caseID)
	if err != nil {
		log.Error().Err(err).Str("case_id", caseID).Msg("failed to query entities by case")
		return nil, fmt.Errorf("failed to list entities for case %s: %w", caseID, err)
	}
	defer rows.Close()

	var entities []*core.Entity
	for rows.Next() {
		var e core.Entity
		var metadataStr string
		if err := rows.Scan(&e.ID, &e.CaseID, &e.Type, &e.Value, &e.Source, &e.Confidence, &e.DiscoveredAt, &metadataStr); err != nil {
			log.Error().Err(err).Str("case_id", caseID).Msg("failed to scan entity row from list")
			return nil, fmt.Errorf("failed to scan entity for case %s: %w", caseID, err)
		}
		if err := json.Unmarshal([]byte(metadataStr), &e.Metadata); err != nil {
			log.Error().Err(err).Str("entity_id", e.ID).Msg("failed to unmarshal entity metadata in list")
			return nil, fmt.Errorf("failed to unmarshal metadata for entity %s in case %s: %w", e.ID, caseID, err)
		}
		entities = append(entities, &e)
	}

	return entities, nil
}

// ListEntitiesByType retrieves all entities of a specific type in a case.
func ListEntitiesByType(caseID, entityType string) ([]*core.Entity, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := TranslatePlaceholder(`SELECT id, case_id, type, value, source, confidence, discovered_at, metadata 
	          FROM entities WHERE case_id = ? AND type = ?`)
	rows, err := DB.Query(query, caseID, entityType)
	if err != nil {
		return nil, fmt.Errorf("failed to query entities by type: %w", err)
	}
	defer rows.Close()

	var entities []*core.Entity
	for rows.Next() {
		var e core.Entity
		var metadataStr string
		if err := rows.Scan(&e.ID, &e.CaseID, &e.Type, &e.Value, &e.Source, &e.Confidence, &e.DiscoveredAt, &metadataStr); err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(metadataStr), &e.Metadata)
		entities = append(entities, &e)
	}

	return entities, nil
}
