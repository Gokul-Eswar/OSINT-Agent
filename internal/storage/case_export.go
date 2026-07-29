package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Gokul-Eswar/Spectre/internal/core"
)

// CaseBundle represents a fully exported case containing its metadata, entities, relationships, evidence, and timeline.
type CaseBundle struct {
	Case          *core.Case           `json:"case"`
	Entities      []*core.Entity       `json:"entities"`
	Relationships []*core.Relationship `json:"relationships"`
	Evidence      []*core.Evidence     `json:"evidence"`
	Timeline      []core.TimelineEvent `json:"timeline"`
}

// ExportCaseBundle exports all data for a case into a JSON file at targetPath.
func ExportCaseBundle(caseID string, targetPath string) error {
	c, err := GetCase(caseID)
	if err != nil {
		return fmt.Errorf("failed to get case: %w", err)
	}
	if c == nil {
		return fmt.Errorf("case not found")
	}

	entities, err := ListEntitiesByCase(caseID)
	if err != nil {
		entities = []*core.Entity{}
	}

	rels, err := ListRelationshipsByCase(caseID)
	if err != nil {
		rels = []*core.Relationship{}
	}

	ev, err := ListEvidenceByCase(caseID)
	if err != nil {
		ev = []*core.Evidence{}
	}

	tl, err := GetCaseTimeline(caseID)
	if err != nil {
		tl = []core.TimelineEvent{}
	}

	bundle := CaseBundle{
		Case:          c,
		Entities:      entities,
		Relationships: rels,
		Evidence:      ev,
		Timeline:      tl,
	}

	data, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal case bundle: %w", err)
	}

	return os.WriteFile(targetPath, data, 0644)
}

// ImportCaseBundle imports a case bundle from a JSON file into the database.
func ImportCaseBundle(sourcePath string) (*core.Case, error) {
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read import file: %w", err)
	}

	var bundle CaseBundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		return nil, fmt.Errorf("failed to parse case bundle JSON: %w", err)
	}

	if bundle.Case == nil {
		return nil, fmt.Errorf("invalid bundle: missing case object")
	}

	// Create or update case
	existingCase, _ := GetCase(bundle.Case.ID)
	if existingCase == nil {
		if err := CreateCase(bundle.Case); err != nil {
			return nil, fmt.Errorf("failed to import case: %w", err)
		}
	}

	// Import entities
	for _, ent := range bundle.Entities {
		_ = CreateEntity(ent)
	}

	// Import relationships
	for _, rel := range bundle.Relationships {
		_ = CreateRelationship(rel)
	}

	// Import evidence
	for _, ev := range bundle.Evidence {
		_ = CreateEvidence(ev)
	}

	return bundle.Case, nil
}
