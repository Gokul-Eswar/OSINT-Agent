package storage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/Gokul-Eswar/Spectre/internal/core"
)

// OnLeadCreated is a hook for real-time updates
var OnLeadCreated func(*core.IntelligenceLead)

// CreateLead inserts a new intelligence lead into the database.
func CreateLead(l *core.IntelligenceLead) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	if l.CreatedAt.IsZero() {
		l.CreatedAt = time.Now()
	}
	if l.Confidence == 0 {
		l.Confidence = 0.5
	}
	if l.Status == "" {
		l.Status = "active"
	}
	if l.EvidenceIDs == nil {
		l.EvidenceIDs = []string{}
	}

	evidenceIDsJSON, err := json.Marshal(l.EvidenceIDs)
	if err != nil {
		log.Error().Err(err).Str("lead_id", l.ID).Msg("failed to marshal lead evidence IDs")
		return fmt.Errorf("failed to marshal evidence IDs for lead %s: %w", l.ID, err)
	}

	query := TranslatePlaceholder(`INSERT INTO intelligence_leads (id, case_id, hypothesis, confidence, evidence_ids, status, created_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`)
	_, err = DB.Exec(query, l.ID, l.CaseID, l.Hypothesis, l.Confidence, string(evidenceIDsJSON), l.Status, l.CreatedAt)
	if err != nil {
		log.Error().
			Err(err).
			Str("lead_id", l.ID).
			Str("case_id", l.CaseID).
			Msg("failed to execute insert lead query")
		return fmt.Errorf("failed to create lead %s in case %s: %w", l.ID, l.CaseID, err)
	}

	if OnLeadCreated != nil {
		OnLeadCreated(l)
	}

	return nil
}

// UpdateLead updates an existing intelligence lead in the database.
func UpdateLead(l *core.IntelligenceLead) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	evidenceIDsJSON, err := json.Marshal(l.EvidenceIDs)
	if err != nil {
		log.Error().Err(err).Str("lead_id", l.ID).Msg("failed to marshal lead evidence IDs")
		return fmt.Errorf("failed to marshal evidence IDs for lead %s: %w", l.ID, err)
	}

	query := TranslatePlaceholder(`UPDATE intelligence_leads SET hypothesis = ?, confidence = ?, evidence_ids = ?, status = ? WHERE id = ?`)
	_, err = DB.Exec(query, l.Hypothesis, l.Confidence, string(evidenceIDsJSON), l.Status, l.ID)
	if err != nil {
		log.Error().
			Err(err).
			Str("lead_id", l.ID).
			Msg("failed to execute update lead query")
		return fmt.Errorf("failed to update lead %s: %w", l.ID, err)
	}

	if OnLeadCreated != nil {
		OnLeadCreated(l)
	}

	return nil
}

// GetLead retrieves a single lead by its ID.
func GetLead(id string) (*core.IntelligenceLead, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := TranslatePlaceholder(`SELECT id, case_id, hypothesis, confidence, evidence_ids, status, created_at FROM intelligence_leads WHERE id = ?`)
	row := DB.QueryRow(query, id)

	var l core.IntelligenceLead
	var evidenceIDsStr string
	err := row.Scan(&l.ID, &l.CaseID, &l.Hypothesis, &l.Confidence, &evidenceIDsStr, &l.Status, &l.CreatedAt)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(evidenceIDsStr), &l.EvidenceIDs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal evidence IDs for lead %s: %w", l.ID, err)
	}

	return &l, nil
}

// ListLeadsByCase returns all intelligence leads associated with a specific case.
func ListLeadsByCase(caseID string) ([]core.IntelligenceLead, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := TranslatePlaceholder(`SELECT id, case_id, hypothesis, confidence, evidence_ids, status, created_at FROM intelligence_leads WHERE case_id = ? ORDER BY created_at DESC`)
	rows, err := DB.Query(query, caseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leads []core.IntelligenceLead
	for rows.Next() {
		var l core.IntelligenceLead
		var evidenceIDsStr string
		err := rows.Scan(&l.ID, &l.CaseID, &l.Hypothesis, &l.Confidence, &evidenceIDsStr, &l.Status, &l.CreatedAt)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(evidenceIDsStr), &l.EvidenceIDs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal evidence IDs: %w", err)
		}

		leads = append(leads, l)
	}

	return leads, nil
}
