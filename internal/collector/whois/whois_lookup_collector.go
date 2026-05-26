package whois

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
	"github.com/rs/zerolog/log"
	"github.com/spectre/spectre/internal/collector"
	"github.com/spectre/spectre/internal/core"
)

type WhoisClient interface {
	Whois(domain string) (string, error)
}

type DefaultWhoisClient struct{}

func (c *DefaultWhoisClient) Whois(domain string) (string, error) {
	return whois.Whois(domain)
}

type WHOISCollector struct {
	client WhoisClient
}

func Register() {
	collector.Register(&WHOISCollector{client: &DefaultWhoisClient{}})
}

func (w *WHOISCollector) Name() string {
	return "whois"
}

func (w *WHOISCollector) Description() string {
	return "Retrieve domain registration information"
}

func (w *WHOISCollector) IsActive() bool {
	return false
}

func (w *WHOISCollector) Collect(caseID string, target string, options map[string]interface{}) ([]core.Evidence, error) {
	log.Info().
		Str("collector", "whois").
		Str("case_id", caseID).
		Str("target", target).
		Msg("collection_started")

	raw, err := w.client.Whois(target)
	if err != nil {
		log.Error().Err(err).Str("target", target).Msg("whois lookup failed")
		return nil, fmt.Errorf("whois lookup failed for %s: %w", target, err)
	}

	// Parse to verify it's valid and get metadata
	result, err := whoisparser.Parse(raw)
	if err != nil {
		log.Debug().Err(err).Str("target", target).Msg("whois parsing failed, continuing with raw data")
	}

	// Store raw file
	storageDir := filepath.Join("evidence_storage", caseID)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		log.Error().Err(err).Str("dir", storageDir).Msg("failed to create storage directory")
		return nil, fmt.Errorf("failed to create storage directory %s: %w", storageDir, err)
	}

	fileName := fmt.Sprintf("whois_%s_%d.txt", target, time.Now().Unix())
	filePath := filepath.Join(storageDir, fileName)
	if err := os.WriteFile(filePath, []byte(raw), 0644); err != nil {
		log.Error().Err(err).Str("path", filePath).Msg("failed to write WHOIS evidence file")
		return nil, fmt.Errorf("failed to write WHOIS evidence file %s: %w", filePath, err)
	}

	// Hash
	hash := sha256.Sum256([]byte(raw))
	hashStr := hex.EncodeToString(hash[:])

	metadata := map[string]interface{}{
		"target": target,
	}
	if result.Registrar != nil {
		metadata["registrar"] = result.Registrar.Name
	}
	if result.Registrant != nil {
		metadata["registrant_email"] = result.Registrant.Email
		metadata["registrant_name"] = result.Registrant.Name
	}

	evidence := core.Evidence{
		CaseID:      caseID,
		Collector:   "whois",
		FilePath:    filePath,
		FileHash:    hashStr,
		CollectedAt: time.Now(),
		Metadata:    metadata,
	}

	log.Info().
		Str("collector", "whois").
		Str("case_id", caseID).
		Str("target", target).
		Msg("collection_completed")

	return []core.Evidence{evidence}, nil
}
