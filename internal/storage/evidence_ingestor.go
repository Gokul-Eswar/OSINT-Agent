package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Gokul-Eswar/Spectre/internal/core"
)

// IngestEvidence is the central routing function for data ingestion.
// It inspects the Evidence's Collector field and delegates the parsing to the appropriate
// collector-specific ingestion handler.
// Unknown collectors are ignored (returning nil), which allows third-party extensions
// to write files to evidence_storage without causing failures before ingestors are written.
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
		return nil // Graceful skip for collectors without specialized ingestion rules
	}
}

// ingestSocial parses social media account presence checking results.
// It extracts the seed username and the sites where matches were found,
// creates account entities, and maps them back to the seed username.
func ingestSocial(ev *core.Evidence) error {
	// Extract target username from the collection metadata.
	username := ev.Metadata["target"].(string)

	// Read the JSON evidence file containing search results.
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

	// 1. Ensure the root username entity exists in the database.
	userEnt, err := EnsureEntity(ev.CaseID, "username", username, "social")
	if err != nil {
		return err
	}

	// 2. Iterate over positive site matches.
	for _, res := range results {
		// Create a specific entity representing the profile URL.
		siteEnt, err := EnsureEntity(ev.CaseID, "account", res.URL, "social")
		if err != nil {
			return err
		}

		// Update or ensure platform metadata (e.g. "platform": "GitHub") is attached to the entity.
		if siteEnt.Metadata == nil {
			siteEnt.Metadata = make(map[string]interface{})
		}
		siteEnt.Metadata["platform"] = res.Site
		if err := UpdateEntity(siteEnt); err != nil {
			return err
		}

		// 3. Link Username -> has_account -> Profile URL entity.
		rel := &core.Relationship{
			CaseID:       ev.CaseID,
			FromEntityID: userEnt.ID,
			ToEntityID:   siteEnt.ID,
			Type:         "has_account",
			EvidenceID:   ev.ID,
		}
		if err := CreateRelationship(rel); err != nil {
			// Log or handle error
		}
	}
	return nil
}

// ingestScreenshot records headful screenshot capture evidence.
// Rather than creating a dedicated "image node" in the graph, it represents
// screenshots as a self-referential relationship on the target entity (e.g., domain/IP).
// This keeps the intelligence graph clean while linking the visual evidence to the target.
func ingestScreenshot(ev *core.Evidence) error {
	target := ev.Metadata["target"].(string)

	// Determine if target target is an IP or domain based on the first character.
	entityType := "domain"
	if len(target) > 0 && (target[0] >= '0' && target[0] <= '9') {
		entityType = "ip"
	}

	// Ensure the target node (IP/Domain) exists in the database.
	targetEnt, err := EnsureEntity(ev.CaseID, entityType, target, "screenshot")
	if err != nil {
		return err
	}

	// Create a self-loop relationship: Target -> has_screenshot -> Target.
	// We bind the relationship to the screenshot's EvidenceID to preserve forensic linkage.
	rel := &core.Relationship{
		CaseID:       ev.CaseID,
		FromEntityID: targetEnt.ID,
		ToEntityID:   targetEnt.ID,
		Type:         "has_screenshot",
		EvidenceID:   ev.ID,
		Confidence:   1.0,
	}
	return CreateRelationship(rel)
}

// ingestPorts parses TCP port scan results.
// For every port marked as "open", it creates a "service" entity (e.g. TCP/80)
// and links the service to the host IP. Closed ports are skipped to keep the database tidy.
func ingestPorts(ev *core.Evidence) error {
	targetIP := ev.Metadata["target"].(string)

	data, err := os.ReadFile(ev.FilePath)
	if err != nil {
		return err
	}

	// Port scan results are marshaled as map[port_string]status_string (e.g. {"22": "open"})
	var results map[string]string
	if err := json.Unmarshal(data, &results); err != nil {
		return err
	}

	// Ensure the base host IP entity exists.
	ipEnt, err := EnsureEntity(ev.CaseID, "ip", targetIP, "ports")
	if err != nil {
		return err
	}

	// Process each scanned port.
	for port, status := range results {
		if status == "open" {
			// Format the service name (e.g., "TCP/22")
			svcName := fmt.Sprintf("TCP/%s", port)

			// Ensure the Service service entity exists.
			svcEnt, err := EnsureEntity(ev.CaseID, "service", svcName, "ports")
			if err != nil {
				return err
			}

			// Link Host IP -> has_port -> Service entity.
			rel := &core.Relationship{
				CaseID:       ev.CaseID,
				FromEntityID: ipEnt.ID,
				ToEntityID:   svcEnt.ID,
				Type:         "has_port",
				EvidenceID:   ev.ID,
			}
			_ = CreateRelationship(rel)
		}
	}
	return nil
}

// ingestHTTP parses basic HTTP inspection metadata (such as server banner headers).
// If a server banner is found, it registers a service and links the target domain to it.
func ingestHTTP(ev *core.Evidence) error {
	target := ev.Metadata["target"].(string)
	server := ""
	if s, ok := ev.Metadata["server"].(string); ok {
		server = s
	}

	// Ensure the parent domain entity exists.
	targetEnt, err := EnsureEntity(ev.CaseID, "domain", target, "http")
	if err != nil {
		return err
	}

	// If the HTTP response included a 'Server' header, register it.
	if server != "" {
		svcEnt, err := EnsureEntity(ev.CaseID, "service", server, "http")
		if err != nil {
			return err
		}

		// Link Domain -> runs_service -> Service.
		rel := &core.Relationship{
			CaseID:       ev.CaseID,
			FromEntityID: targetEnt.ID,
			ToEntityID:   svcEnt.ID,
			Type:         "runs_service",
			EvidenceID:   ev.ID,
		}
		_ = CreateRelationship(rel)
	}
	return nil
}

// ingestGeo processes IP geolocation data.
// Rather than cluttering the graph with distinct "City" or "Country" nodes,
// we append the geolocation properties directly to the parent IP entity's Metadata JSON store.
func ingestGeo(ev *core.Evidence) error {
	targetIP := ev.Metadata["target"].(string)

	// Fetch the IP entity from the database.
	ipEnt, err := GetEntityByValue(ev.CaseID, targetIP)
	if err != nil {
		return err
	}
	if ipEnt == nil {
		// If the IP node is missing, create it.
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

	// Ensure metadata map is allocated.
	if ipEnt.Metadata == nil {
		ipEnt.Metadata = make(map[string]interface{})
	}

	// Map geolocation keys from the evidence metadata directly to the entity.
	fields := []string{"country", "city", "isp", "lat", "lon"}
	for _, f := range fields {
		if v, ok := ev.Metadata[f]; ok {
			ipEnt.Metadata[f] = v
		}
	}

	// Persist the updated metadata properties to the SQLite database.
	return UpdateEntity(ipEnt)
}

// ingestGitHub extracts repository details, repository owners,
// creates the nodes, and links them via ownership relationships.
func ingestGitHub(ev *core.Evidence) error {
	var data []byte
	var err error

	// Prefer reading raw bytes in memory if pre-loaded.
	if ev.RawData != nil {
		if b, ok := ev.RawData.([]byte); ok {
			data = b
		}
	}

	// Fallback to reading the saved JSON evidence file from disk.
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

	// Iterate through matching GitHub repositories.
	for _, item := range results.Items {
		// Create Repository entity node.
		repoEnt, err := EnsureEntity(ev.CaseID, "repo", item.HTMLURL, "github")
		if err != nil {
			return err
		}

		// Create Owner Username entity node.
		userEnt, err := EnsureEntity(ev.CaseID, "username", item.Owner.Login, "github")
		if err != nil {
			return err
		}

		// Link Owner -> owns -> Repository.
		rel := &core.Relationship{
			CaseID:       ev.CaseID,
			FromEntityID: userEnt.ID,
			ToEntityID:   repoEnt.ID,
			Type:         "owns",
			EvidenceID:   ev.ID,
			Confidence:   1.0,
		}
		_ = CreateRelationship(rel)
	}

	return nil
}

// ingestWHOIS extracts registry contact details (specifically email addresses) from WHOIS text responses.
// It links the domain node to the email node.
func ingestWHOIS(ev *core.Evidence) error {
	targetDomain := ev.Metadata["target"].(string)

	// Ensure target domain entity exists.
	domainEnt, err := EnsureEntity(ev.CaseID, "domain", targetDomain, "whois")
	if err != nil {
		return err
	}

	// If registrant_email metadata is present, map it.
	if email, ok := ev.Metadata["registrant_email"].(string); ok && email != "" {
		emailEnt, err := EnsureEntity(ev.CaseID, "email", email, "whois")
		if err != nil {
			return err
		}

		// Link Domain -> registered_by -> Email.
		rel := &core.Relationship{
			CaseID:       ev.CaseID,
			FromEntityID: domainEnt.ID,
			ToEntityID:   emailEnt.ID,
			Type:         "registered_by",
			EvidenceID:   ev.ID,
			Confidence:   1.0,
		}
		_ = CreateRelationship(rel)
	}

	return nil
}

// ingestDNS parses passive DNS query records (A, MX, NS).
// For A-records, it sets up link relationships from the domain node to resolved IP nodes.
func ingestDNS(ev *core.Evidence) error {
	var results map[string][]string

	// Read in-memory buffer if present.
	if ev.RawData != nil {
		if r, ok := ev.RawData.(map[string][]string); ok {
			results = r
		}
	}

	// Or load from the saved JSON file on disk.
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

	// Ensure target domain entity node exists.
	domainEnt, err := EnsureEntity(ev.CaseID, "domain", targetDomain, "dns")
	if err != nil {
		return err
	}

	// Process resolved A records.
	for _, ip := range results["A"] {
		// Ensure IP node exists.
		ipEnt, err := EnsureEntity(ev.CaseID, "ip", ip, "dns")
		if err != nil {
			return err
		}

		// Link Domain -> resolves_to -> IP.
		rel := &core.Relationship{
			CaseID:       ev.CaseID,
			FromEntityID: domainEnt.ID,
			ToEntityID:   ipEnt.ID,
			Type:         "resolves_to",
			EvidenceID:   ev.ID,
			Confidence:   1.0,
		}
		if err := CreateRelationship(rel); err != nil {
			// Failures usually mean relationship constraints already exist. We skip.
		}
	}

	return nil
}
