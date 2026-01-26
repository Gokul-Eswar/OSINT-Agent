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
	case "whois":
		return ingestWHOIS(ev)
	case "github":
		return ingestGitHub(ev)
	case "geo":
		return ingestGeo(ev)
	case "ports":
		return ingestPorts(ev)
	case "http":
		return ingestHTTP(ev)
	case "screenshot":
		return ingestScreenshot(ev)
	default:
		return nil // No ingestion logic for this collector yet
	}
}

func ingestScreenshot(ev *core.Evidence) error {
	target := ev.Metadata["target"].(string)

	// Ensure target entity exists (usually a domain or IP)
	targetEnt, _ := GetEntityByValue(ev.CaseID, target)
	if targetEnt == nil {
		// Try to guess type
		entityType := "domain"
		if len(target) > 0 && (target[0] >= '0' && target[0] <= '9') {
			entityType = "ip"
		}
		targetEnt = &core.Entity{CaseID: ev.CaseID, Type: entityType, Value: target, Source: "screenshot"}
		CreateEntity(targetEnt)
	}

	// Link target to the screenshot evidence
	// We don't create a new entity for the screenshot itself, 
	// but the relationship record stores the EvidenceID.
	rel := &core.Relationship{
		CaseID:       ev.CaseID,
		FromEntityID: targetEnt.ID,
		ToEntityID:   targetEnt.ID, // Self-link to represent property/evidence
		Type:         "has_screenshot",
		EvidenceID:   ev.ID,
		Confidence:   1.0,
	}
	return CreateRelationship(rel)
}

func ingestPorts(ev *core.Evidence) error {
	targetIP := ev.Metadata["target"].(string)
	
	data, err := os.ReadFile(ev.FilePath)
	if err != nil {
		return err
	}

	var results map[string]string
	if err := json.Unmarshal(data, &results); err != nil {
		return err
	}

	// Ensure IP entity exists
	ipEnt, _ := GetEntityByValue(ev.CaseID, targetIP)
	if ipEnt == nil {
		ipEnt = &core.Entity{CaseID: ev.CaseID, Type: "ip", Value: targetIP, Source: "ports"}
		CreateEntity(ipEnt)
	}

	for port, status := range results {
		if status == "open" {
			svcName := fmt.Sprintf("TCP/%s", port)
			svcEnt := &core.Entity{
				CaseID: ev.CaseID,
				Type:   "service",
				Value:  svcName,
				Source: "ports",
			}
			
			existing, _ := GetEntityByValue(ev.CaseID, svcName)
			if existing == nil {
				CreateEntity(svcEnt)
			} else {
				svcEnt = existing
			}

			// Link IP -> has -> Service
			rel := &core.Relationship{
				CaseID:       ev.CaseID,
				FromEntityID: ipEnt.ID,
				ToEntityID:   svcEnt.ID,
				Type:         "has_port",
				EvidenceID:   ev.ID,
			}
			CreateRelationship(rel)
		}
	}
	return nil
}

func ingestHTTP(ev *core.Evidence) error {
	target := ev.Metadata["target"].(string)
	server := ""
	if s, ok := ev.Metadata["server"].(string); ok {
		server = s
	}

	// Ensure target entity exists
	targetEnt, _ := GetEntityByValue(ev.CaseID, target)
	if targetEnt == nil {
		targetEnt = &core.Entity{CaseID: ev.CaseID, Type: "domain", Value: target, Source: "http"}
		CreateEntity(targetEnt)
	}

	if server != "" {
		svcEnt := &core.Entity{
			CaseID: ev.CaseID,
			Type:   "service",
			Value:  server,
			Source: "http",
		}
		
		existing, _ := GetEntityByValue(ev.CaseID, server)
		if existing == nil {
			CreateEntity(svcEnt)
		} else {
			svcEnt = existing
		}

		// Link Target -> runs -> Service
		rel := &core.Relationship{
			CaseID:       ev.CaseID,
			FromEntityID: targetEnt.ID,
			ToEntityID:   svcEnt.ID,
			Type:         "runs_service",
			EvidenceID:   ev.ID,
		}
		CreateRelationship(rel)
	}
	return nil
}

func ingestGeo(ev *core.Evidence) error {
	targetIP := ev.Metadata["target"].(string)
	
	// Ensure IP entity exists
	ipEnt, err := GetEntityByValue(ev.CaseID, targetIP)
	if err != nil {
		return err
	}
	if ipEnt == nil {
		// Create it if it doesn't exist (though rare if we collected on it)
		ipEnt = &core.Entity{
			CaseID: ev.CaseID,
			Type:   "ip",
			Value:  targetIP,
			Source: "geo",
			Metadata: make(map[string]interface{}),
		}
		if err := CreateEntity(ipEnt); err != nil {
			return err
		}
	}

	// Update metadata
	if ipEnt.Metadata == nil {
		ipEnt.Metadata = make(map[string]interface{})
	}
	
	// Copy relevant fields from evidence metadata
	fields := []string{"country", "city", "isp", "lat", "lon"}
	for _, f := range fields {
		if v, ok := ev.Metadata[f]; ok {
			ipEnt.Metadata[f] = v
		}
	}
	
	return UpdateEntity(ipEnt)
}

func ingestGitHub(ev *core.Evidence) error {
	var data []byte
	var err error

	if ev.RawData != nil {
		if b, ok := ev.RawData.([]byte); ok {
			data = b
		}
	}

	if data == nil {
		data, err = os.ReadFile(ev.FilePath)
		if err != nil {
			return err
		}
	}

	var results struct {
		Items []struct {
			FullName string `json:"full_name"`
			HTMLURL  string `json:"html_url"`
			Owner    struct {
				Login string `json:"login"`
			} `json:"owner"`
		} `json:"items"`
	}

	if err := json.Unmarshal(data, &results); err != nil {
		return err
	}

	for _, item := range results.Items {
		// Create Repo entity
		repoEnt := &core.Entity{
			CaseID: ev.CaseID,
			Type:   "repo",
			Value:  item.HTMLURL,
			Source: "github",
		}
		CreateEntity(repoEnt)

		// Create User entity
		userEnt := &core.Entity{
			CaseID: ev.CaseID,
			Type:   "username",
			Value:  item.Owner.Login,
			Source: "github",
		}
		
		existingUser, _ := GetEntityByValue(ev.CaseID, item.Owner.Login)
		if existingUser == nil {
			CreateEntity(userEnt)
		} else {
			userEnt = existingUser
		}

		// Link User -> owns -> Repo
		rel := &core.Relationship{
			CaseID:       ev.CaseID,
			FromEntityID: userEnt.ID,
			ToEntityID:   repoEnt.ID,
			Type:         "owns",
			EvidenceID:   ev.ID,
			Confidence:   1.0,
		}
		CreateRelationship(rel)
	}

	return nil
}

func ingestWHOIS(ev *core.Evidence) error {
	targetDomain := ev.Metadata["target"].(string)
	
	// Ensure domain entity exists
	domainEnt, _ := GetEntityByValue(ev.CaseID, targetDomain)
	if domainEnt == nil {
		domainEnt = &core.Entity{
			CaseID: ev.CaseID,
			Type:   "domain",
			Value:  targetDomain,
			Source: "whois",
		}
		if err := CreateEntity(domainEnt); err != nil {
			return err
		}
	}

	// If we have a registrant email, create it and link it
	if email, ok := ev.Metadata["registrant_email"].(string); ok && email != "" {
		emailEnt := &core.Entity{
			CaseID: ev.CaseID,
			Type:   "email",
			Value:  email,
			Source: "whois",
		}
		
		existingEmail, _ := GetEntityByValue(ev.CaseID, email)
		if existingEmail == nil {
			if err := CreateEntity(emailEnt); err != nil {
				return err
			}
		} else {
			emailEnt = existingEmail
		}

		// Link Domain -> owns -> Email (or registered_by)
		rel := &core.Relationship{
			CaseID:       ev.CaseID,
			FromEntityID: domainEnt.ID,
			ToEntityID:   emailEnt.ID,
			Type:         "registered_by",
			EvidenceID:   ev.ID,
			Confidence:   1.0,
		}
		CreateRelationship(rel)
	}

	return nil
}

func ingestDNS(ev *core.Evidence) error {
	var results map[string][]string

	// Try in-memory first
	if ev.RawData != nil {
		if r, ok := ev.RawData.(map[string][]string); ok {
			results = r
		}
	}

	// Fallback to disk
	if results == nil {
		data, err := os.ReadFile(ev.FilePath)
		if err != nil {
			return fmt.Errorf("failed to read evidence file: %w", err)
		}

		if err := json.Unmarshal(data, &results); err != nil {
			return fmt.Errorf("failed to unmarshal DNS results: %w", err)
		}
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
