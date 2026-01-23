package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spectre/spectre/internal/core"
)

// IngestEvidence parses evidence data and populates the graph (entities/relationships).
func IngestEvidence(ev *core.Evidence) error {
	switch ev.Collector {
	case "dns":
		return ingestDNS(ev)
	default:
		return nil // No ingestion logic for this collector yet
	}
}

func ingestDNS(ev *core.Evidence) error {
	data, err := os.ReadFile(ev.FilePath)
	if err != nil {
		return fmt.Errorf("failed to read evidence file: %w", err)
	}

	var results map[string][]string
	if err := json.Unmarshal(data, &results); err != nil {
		return fmt.Errorf("failed to unmarshal DNS results: %w", err)
	}

	targetDomain := ev.Metadata["target"].(string)

	// Ensure target domain entity exists
	domainEnt := &core.Entity{
		CaseID: ev.CaseID,
		Type:   "domain",
		Value:  targetDomain,
		Source: "dns",
	}
	
	// Check if already exists to avoid errors (or use GetEntityByValue)
	existing, _ := GetEntityByValue(ev.CaseID, targetDomain)
	if existing == nil {
		if err := CreateEntity(domainEnt); err != nil {
			return err
		}
	} else {
		domainEnt = existing
	}

	// Process A records
	for _, ip := range results["A"] {
		ipEnt := &core.Entity{
			CaseID: ev.CaseID,
			Type:   "ip",
			Value:  ip,
			Source: "dns",
		}
		
		existingIP, _ := GetEntityByValue(ev.CaseID, ip)
		if existingIP == nil {
			if err := CreateEntity(ipEnt); err != nil {
				return err
			}
		} else {
			ipEnt = existingIP
		}

		// Create relationship
		rel := &core.Relationship{
			CaseID:       ev.CaseID,
			FromEntityID: domainEnt.ID,
			ToEntityID:   ipEnt.ID,
			Type:         "resolves_to",
			EvidenceID:   ev.ID,
			Confidence:   1.0,
		}
		if err := CreateRelationship(rel); err != nil {
			// Might already exist due to unique constraint, ignore error
		}
	}

	return nil
}
