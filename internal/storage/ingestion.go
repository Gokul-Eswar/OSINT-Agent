package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spectre/spectre/internal/core"
)

// IngestEvidence dispatches collector-specific parsers that convert raw evidence
// into graph entities and relationships.
//
// The collector name is the routing key for ingestion behavior. Unknown
// collectors are ignored so newly added plugins can still persist evidence even
// before a dedicated ingestor exists.
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
	case "social":
		return ingestSocial(ev)
	default:
		return nil // No ingestion logic for this collector yet
	}
}

// ingestSocial maps discovered social profiles into account entities linked to
// the seed username.
func ingestSocial(ev *core.Evidence) error {
	username := ev.Metadata["target"].(string)

	data, err := os.ReadFile(ev.FilePath)
	if err != nil {
		return err
	}

	var results []struct {
		Site string `json:"site"`
		URL  string `json:"url"`
	}
	if err := json.Unmarshal(data, &results); err != nil {
		return err
	}

	// Ensure username entity exists
	userEnt, err := EnsureEntity(ev.CaseID, "username", username, "social")
	if err != nil {
		return err
	}

	for _, res := range results {
		// Create site entity
		siteEnt, err := EnsureEntity(ev.CaseID, "account", res.URL, "social")
		if err != nil {
			return err
		}

		// Update metadata if it's a new entity or we want to ensure it has platform info
		if siteEnt.Metadata == nil {
			siteEnt.Metadata = make(map[string]interface{})
		}
		siteEnt.Metadata["platform"] = res.Site
		UpdateEntity(siteEnt)

		// Link Username -> has_account -> Site
		rel := &core.Relationship{
			CaseID:       ev.CaseID,
			FromEntityID: userEnt.ID,
			ToEntityID:   siteEnt.ID,
			Type:         "has_account",
			EvidenceID:   ev.ID,
		}
		CreateRelationship(rel)
	}
	return nil
}

// ingestScreenshot records screenshot evidence as a self-referential
// relationship on the target entity.
//
// This models screenshots as supporting evidence for an entity without creating
// a second synthetic node type just for image artifacts.
func ingestScreenshot(ev *core.Evidence) error {
	target := ev.Metadata["target"].(string)

	// Ensure target entity exists (usually a domain or IP)
	entityType := "domain"
	if len(target) > 0 && (target[0] >= '0' && target[0] <= '9') {
		entityType = "ip"
	}
	targetEnt, err := EnsureEntity(ev.CaseID, entityType, target, "screenshot")
	if err != nil {
		return err
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

// ingestPorts creates service entities for open ports and links them to the
// target IP.
//
// Closed ports are intentionally ignored to keep the graph focused on actionable
// attack surface relationships.
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
	ipEnt, err := EnsureEntity(ev.CaseID, "ip", targetIP, "ports")
	if err != nil {
		return err
	}

	for port, status := range results {
		if status == "open" {
			svcName := fmt.Sprintf("TCP/%s", port)
			svcEnt, err := EnsureEntity(ev.CaseID, "service", svcName, "ports")
			if err != nil {
				return err
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

// ingestHTTP enriches a target domain with server/software relationship data
// when available from HTTP collector metadata.
func ingestHTTP(ev *core.Evidence) error {
	target := ev.Metadata["target"].(string)
	server := ""
	if s, ok := ev.Metadata["server"].(string); ok {
		server = s
	}

	// Ensure target entity exists
	targetEnt, err := EnsureEntity(ev.CaseID, "domain", target, "http")
	if err != nil {
		return err
	}

	if server != "" {
		svcEnt, err := EnsureEntity(ev.CaseID, "service", server, "http")
		if err != nil {
			return err
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

// ingestGeo applies geolocation metadata directly to the IP entity instead of
// creating additional geo nodes.
//
// This keeps the graph compact while still making location attributes queryable
// through entity metadata.
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
			CaseID:   ev.CaseID,
			Type:     "ip",
			Value:    targetIP,
			Source:   "geo",
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

// ingestGitHub builds repository and owner entities from search results and
// links them via ownership relationships.
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
		repoEnt, err := EnsureEntity(ev.CaseID, "repo", item.HTMLURL, "github")
		if err != nil {
			return err
		}

		// Create User entity
		userEnt, err := EnsureEntity(ev.CaseID, "username", item.Owner.Login, "github")
		if err != nil {
			return err
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

// ingestWHOIS links a domain to registrant contact artifacts extracted from
// WHOIS metadata.
func ingestWHOIS(ev *core.Evidence) error {
	targetDomain := ev.Metadata["target"].(string)

	// Ensure domain entity exists
	domainEnt, err := EnsureEntity(ev.CaseID, "domain", targetDomain, "whois")
	if err != nil {
		return err
	}

	// If we have a registrant email, create it and link it
	if email, ok := ev.Metadata["registrant_email"].(string); ok && email != "" {
		emailEnt, err := EnsureEntity(ev.CaseID, "email", email, "whois")
		if err != nil {
			return err
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

// ingestDNS maps DNS A-record resolution edges from domain entities to IP
// entities, preserving evidence provenance on each relationship.
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
	domainEnt, err := EnsureEntity(ev.CaseID, "domain", targetDomain, "dns")
	if err != nil {
		return err
	}

	// Process A records
	for _, ip := range results["A"] {
		ipEnt, err := EnsureEntity(ev.CaseID, "ip", ip, "dns")
		if err != nil {
			return err
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
