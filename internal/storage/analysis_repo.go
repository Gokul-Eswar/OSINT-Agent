package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spectre/spectre/internal/core"
)

// SaveAnalysis inserts a new analysis record into the database.
func SaveAnalysis(a *core.Analysis) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	if a.AnalyzedAt.IsZero() {
		a.AnalyzedAt = time.Now()
	}

	findingsJSON, _ := json.Marshal(a.Findings)
	risksJSON, _ := json.Marshal(a.Risks)
	connectionsJSON, _ := json.Marshal(a.Connections)
	nextStepsJSON, _ := json.Marshal(a.NextSteps)
	missingDataJSON, _ := json.Marshal(a.MissingData)
	suggestedColsJSON, _ := json.Marshal(a.SuggestedCollectors)

	query := TranslatePlaceholder(`INSERT INTO analyses (id, case_id, context_hash, findings, risks, connections, next_steps, missing_data, suggested_collectors, confidence, analyzed_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	_, err := DB.Exec(query, a.ID, a.CaseID, a.ContextHash, string(findingsJSON), string(risksJSON), string(connectionsJSON), string(nextStepsJSON), string(missingDataJSON), string(suggestedColsJSON), a.Confidence, a.AnalyzedAt)
	if err != nil {
		return fmt.Errorf("failed to save analysis: %w", err)
	}

	return nil
}

// GetLatestAnalysis retrieves the most recent analysis for a case.
func GetLatestAnalysis(caseID string) (*core.Analysis, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := TranslatePlaceholder(`SELECT id, case_id, context_hash, findings, risks, connections, next_steps, missing_data, suggested_collectors, confidence, analyzed_at 
	          FROM analyses WHERE case_id = ? ORDER BY analyzed_at DESC LIMIT 1`)
	row := DB.QueryRow(query, caseID)

	var a core.Analysis
	var contextHash sql.NullString
	var findingsStr, risksStr, connStr, nextStepsStr, missingDataStr, suggestedColsStr string

	err := row.Scan(&a.ID, &a.CaseID, &contextHash, &findingsStr, &risksStr, &connStr, &nextStepsStr, &missingDataStr, &suggestedColsStr, &a.Confidence, &a.AnalyzedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get analysis: %w", err)
	}

	if contextHash.Valid {
		a.ContextHash = contextHash.String
	}

	json.Unmarshal([]byte(findingsStr), &a.Findings)
	json.Unmarshal([]byte(risksStr), &a.Risks)
	json.Unmarshal([]byte(connStr), &a.Connections)
	json.Unmarshal([]byte(nextStepsStr), &a.NextSteps)
	
	if missingDataStr != "" {
		json.Unmarshal([]byte(missingDataStr), &a.MissingData)
	}
	if suggestedColsStr != "" {
		json.Unmarshal([]byte(suggestedColsStr), &a.SuggestedCollectors)
	}

	return &a, nil
}

// GetAnalysisByHash checks if a report already exists for this context.
func GetAnalysisByHash(caseID string, hash string) (*core.Analysis, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := TranslatePlaceholder(`SELECT id, case_id, context_hash, findings, risks, connections, next_steps, missing_data, suggested_collectors, confidence, analyzed_at 
	          FROM analyses WHERE case_id = ? AND context_hash = ? LIMIT 1`)
	row := DB.QueryRow(query, caseID, hash)

	var a core.Analysis
	var contextHash sql.NullString
	var findingsStr, risksStr, connStr, nextStepsStr, missingDataStr, suggestedColsStr string

	err := row.Scan(&a.ID, &a.CaseID, &contextHash, &findingsStr, &risksStr, &connStr, &nextStepsStr, &missingDataStr, &suggestedColsStr, &a.Confidence, &a.AnalyzedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get analysis: %w", err)
	}

	if contextHash.Valid {
		a.ContextHash = contextHash.String
	}

	json.Unmarshal([]byte(findingsStr), &a.Findings)
	json.Unmarshal([]byte(risksStr), &a.Risks)
	json.Unmarshal([]byte(connStr), &a.Connections)
	json.Unmarshal([]byte(nextStepsStr), &a.NextSteps)
	
	if missingDataStr != "" {
		json.Unmarshal([]byte(missingDataStr), &a.MissingData)
	}
	if suggestedColsStr != "" {
		json.Unmarshal([]byte(suggestedColsStr), &a.SuggestedCollectors)
	}

	return &a, nil
}