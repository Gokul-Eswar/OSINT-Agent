package analysis

import (
	"fmt"
	"strings"

	"github.com/Gokul-Eswar/Spectre/internal/core"
	"github.com/Gokul-Eswar/Spectre/internal/storage"
)

// PersonaCluster represents a correlated identity profile spanning multiple online platforms.
type PersonaCluster struct {
	Handle     string         `json:"handle"`
	PersonaID  string         `json:"persona_id"`
	Platforms  []string       `json:"platforms"`
	Accounts   []*core.Entity `json:"accounts"`
	Emails     []string       `json:"emails"`
	Confidence float64        `json:"confidence"`
}

// CorrelatePersonas analyzes all discovered account and username entities in a case to map personas.
func CorrelatePersonas(caseID string) ([]*PersonaCluster, error) {
	entities, err := storage.ListEntitiesByCase(caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list case entities: %w", err)
	}

	clusters := make(map[string]*PersonaCluster)

	for _, ent := range entities {
		if ent.Type == "username" || ent.Type == "account" {
			handle := strings.ToLower(strings.TrimSpace(ent.Value))
			if idx := strings.LastIndex(handle, "/"); idx != -1 {
				handle = handle[idx+1:]
			}
			handle = strings.TrimPrefix(handle, "@")

			if handle == "" {
				continue
			}

			cluster, exists := clusters[handle]
			if !exists {
				cluster = &PersonaCluster{
					Handle:     handle,
					PersonaID:  fmt.Sprintf("persona_%s", handle),
					Platforms:  []string{},
					Accounts:   []*core.Entity{},
					Emails:     []string{},
					Confidence: 0.8,
				}
				clusters[handle] = cluster
			}

			cluster.Accounts = append(cluster.Accounts, ent)
			if platform, ok := ent.Metadata["platform"].(string); ok && platform != "" {
				cluster.Platforms = append(cluster.Platforms, platform)
			}
		}
	}

	// Create Persona Entities & Relationships in DB
	var result []*PersonaCluster
	for _, cluster := range clusters {
		if len(cluster.Accounts) >= 1 {
			personaEnt, _ := storage.EnsureEntity(caseID, "person", cluster.Handle, "persona_engine")
			if personaEnt != nil {
				for _, acc := range cluster.Accounts {
					rel := &core.Relationship{
						CaseID:       caseID,
						FromEntityID: personaEnt.ID,
						ToEntityID:   acc.ID,
						Type:         "owns_account",
						Confidence:   cluster.Confidence,
					}
					_ = storage.CreateRelationship(rel)
				}
			}
			result = append(result, cluster)
		}
	}

	return result, nil
}
