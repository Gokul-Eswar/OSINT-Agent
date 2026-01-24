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

	"github.com/spectre/spectre/internal/collector"
	"github.com/spectre/spectre/internal/config"
	"github.com/spectre/spectre/internal/core"
)

type GitHubCollector struct {
	Client *http.Client
}

func init() {
	collector.Register(&GitHubCollector{
		Client: &http.Client{Timeout: 30 * time.Second},
	})
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

func (g *GitHubCollector) Collect(caseID string, target string) ([]core.Evidence, error) {
	apiKey := config.GetAPIKey("github")
	
	// Search repositories
	url := fmt.Sprintf("https://api.github.com/search/repositories?q=%s", target)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if apiKey != "" {
		req.Header.Set("Authorization", "token "+apiKey)
	}

	resp, err := g.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github search failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github API returned %d: %s", resp.StatusCode, string(body))
	}

	// Store raw evidence
	storageDir := filepath.Join("evidence_storage", caseID)
	os.MkdirAll(storageDir, 0755)
	fileName := fmt.Sprintf("github_%s_%d.json", target, time.Now().Unix())
	filePath := filepath.Join(storageDir, fileName)
	os.WriteFile(filePath, body, 0644)

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
	}

	return []core.Evidence{evidence}, nil
}
