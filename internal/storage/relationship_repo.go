package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spectre/spectre/internal/core"
)

// CreateRelationship inserts a new relationship into the database.
func CreateRelationship(r *core.Relationship) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	if r.DiscoveredAt.IsZero() {
		r.DiscoveredAt = time.Now()
	}
	if r.Confidence == 0 {
		r.Confidence = 0.5
	}

	query := TranslatePlaceholder(`INSERT INTO relationships (id, case_id, from_entity, to_entity, rel_type, confidence, evidence_id, discovered_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	_, err := DB.Exec(query, r.ID, r.CaseID, r.FromEntityID, r.ToEntityID, r.Type, r.Confidence, r.EvidenceID, r.DiscoveredAt)
	if err != nil {
		return fmt.Errorf("failed to create relationship: %w", err)
	}

	return nil
}

// GetRelationship retrieves a relationship by its ID.
func GetRelationship(id string) (*core.Relationship, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := TranslatePlaceholder(`SELECT id, case_id, from_entity, to_entity, rel_type, confidence, evidence_id, discovered_at FROM relationships WHERE id = ?`)
	row := DB.QueryRow(query, id)

	var r core.Relationship
	err := row.Scan(&r.ID, &r.CaseID, &r.FromEntityID, &r.ToEntityID, &r.Type, &r.Confidence, &r.EvidenceID, &r.DiscoveredAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get relationship: %w", err)
	}

	return &r, nil
}

// ListRelationshipsByFromEntity retrieves all relationships originating from a specific entity.
func ListRelationshipsByFromEntity(entityID string) ([]*core.Relationship, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := TranslatePlaceholder(`SELECT id, case_id, from_entity, to_entity, rel_type, confidence, evidence_id, discovered_at 
	          FROM relationships WHERE from_entity = ?`)
	rows, err := DB.Query(query, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to query relationships by from_entity: %w", err)
	}
	defer rows.Close()

	var relationships []*core.Relationship
	for rows.Next() {
		var r core.Relationship
		if err := rows.Scan(&r.ID, &r.CaseID, &r.FromEntityID, &r.ToEntityID, &r.Type, &r.Confidence, &r.EvidenceID, &r.DiscoveredAt); err != nil {
			return nil, err
		}
		relationships = append(relationships, &r)
	}

	return relationships, nil
}

// ListRelationshipsByType retrieves all relationships of a specific type in a case.
func ListRelationshipsByType(caseID, relType string) ([]*core.Relationship, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := TranslatePlaceholder(`SELECT id, case_id, from_entity, to_entity, rel_type, confidence, evidence_id, discovered_at 
	          FROM relationships WHERE case_id = ? AND rel_type = ?`)
	rows, err := DB.Query(query, caseID, relType)
	if err != nil {
		return nil, fmt.Errorf("failed to query relationships by type: %w", err)
	}
	defer rows.Close()

	var relationships []*core.Relationship
	for rows.Next() {
		var r core.Relationship
		if err := rows.Scan(&r.ID, &r.CaseID, &r.FromEntityID, &r.ToEntityID, &r.Type, &r.Confidence, &r.EvidenceID, &r.DiscoveredAt); err != nil {
			return nil, err
		}
		relationships = append(relationships, &r)
	}

	return relationships, nil
}

// ListRelationshipsByCase retrieves all relationships associated with a specific case.
func ListRelationshipsByCase(caseID string) ([]*core.Relationship, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := TranslatePlaceholder(`SELECT id, case_id, from_entity, to_entity, rel_type, confidence, evidence_id, discovered_at 
	          FROM relationships WHERE case_id = ?`)
	rows, err := DB.Query(query, caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query relationships by case: %w", err)
	}
	defer rows.Close()

	var relationships []*core.Relationship
	for rows.Next() {
		var r core.Relationship
		if err := rows.Scan(&r.ID, &r.CaseID, &r.FromEntityID, &r.ToEntityID, &r.Type, &r.Confidence, &r.EvidenceID, &r.DiscoveredAt); err != nil {
			return nil, err
		}
		relationships = append(relationships, &r)
	}

	return relationships, nil
}
