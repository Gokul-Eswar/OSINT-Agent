package github

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spectre/spectre/internal/collector"
	"github.com/spectre/spectre/internal/config"
	"github.com/spectre/spectre/internal/core"
	netclient "github.com/spectre/spectre/internal/http"
)

type GitHubCollector struct {
	Client *http.Client
}

func Register() {
	collector.Register(&GitHubCollector{})
}

func (g *GitHubCollector) Name() string {
	return "github"
}

func (g *GitHubCollector) Description() string {
	return "Search GitHub for repositories and users"
}

func (g *GitHubCollector) IsActive() bool {
	return false
}

func (g *GitHubCollector) Collect(caseID string, target string, options map[string]interface{}) ([]core.Evidence, error) {
	log.Info().
		Str("collector", "github").
		Str("case_id", caseID).
		Str("target", target).
		Msg("collection_started")

	apiKey := config.GetAPIKey("github")
	client, err := netclient.NewClient()
	if err != nil {
		log.Error().Err(err).Msg("failed to create http client")
		return nil, fmt.Errorf("failed to create http client: %w", err)
	}

	// Search repositories
	url := fmt.Sprintf("https://api.github.com/search/repositories?q=%s", target)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("failed to create github request")
		return nil, fmt.Errorf("failed to create github request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if apiKey != "" {
		req.Header.Set("Authorization", "token "+apiKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Str("target", target).Msg("github search failed")
		return nil, fmt.Errorf("github search failed for %s: %w", target, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read github response body")
		return nil, fmt.Errorf("failed to read github response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Error().
			Int("status_code", resp.StatusCode).
			Str("body", string(body)).
			Msg("github api error")
		return nil, fmt.Errorf("github API returned %d: %s", resp.StatusCode, string(body))
	}

	// Store raw evidence
	storageDir := filepath.Join("evidence_storage", caseID)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		log.Error().Err(err).Str("dir", storageDir).Msg("failed to create storage directory")
		return nil, fmt.Errorf("failed to create storage directory %s: %w", storageDir, err)
	}

	fileName := fmt.Sprintf("github_%s_%d.json", target, time.Now().Unix())
	filePath := filepath.Join(storageDir, fileName)
	if err := os.WriteFile(filePath, body, 0644); err != nil {
		log.Error().Err(err).Str("path", filePath).Msg("failed to write github evidence file")
		return nil, fmt.Errorf("failed to write github evidence file %s: %w", filePath, err)
	}

	hash := sha256.Sum256(body)
	hashStr := hex.EncodeToString(hash[:])

	evidence := core.Evidence{
		CaseID:      caseID,
		Collector:   "github",
		FilePath:    filePath,
		FileHash:    hashStr,
		CollectedAt: time.Now(),
		Metadata: map[string]interface{}{
			"target": target,
		},
		RawData: body,
	}

	log.Info().
		Str("collector", "github").
		Str("case_id", caseID).
		Str("target", target).
		Msg("collection_completed")

	return []core.Evidence{evidence}, nil
}
