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
	"github.com/spectre/spectre/internal/collector"
	"github.com/spectre/spectre/internal/core"
)

type WHOISCollector struct{}

func init() {
	collector.Register(&WHOISCollector{})
}

func (w *WHOISCollector) Name() string {
	return "whois"
}

func (w *WHOISCollector) Description() string {
	return "Retrieve domain registration information"
}

func (w *WHOISCollector) Collect(caseID string, target string) ([]core.Evidence, error) {
	raw, err := whois.Whois(target)
	if err != nil {
		return nil, fmt.Errorf("whois lookup failed: %w", err)
	}

	// Parse to verify it's valid and get metadata
	result, err := whoisparser.Parse(raw)
	if err != nil {
		// Even if parsing fails, we keep the raw data as evidence
	}

	// Store raw file
	storageDir := filepath.Join("evidence_storage", caseID)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, err
	}

	fileName := fmt.Sprintf("whois_%s_%d.txt", target, time.Now().Unix())
	filePath := filepath.Join(storageDir, fileName)
	if err := os.WriteFile(filePath, []byte(raw), 0644); err != nil {
		return nil, err
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

	return []core.Evidence{evidence}, nil
}
