package geo

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/spectre/spectre/internal/collector"
	"github.com/spectre/spectre/internal/core"
)

type GeoIPCollector struct{}

func init() {
	collector.Register(&GeoIPCollector{})
}

func (c *GeoIPCollector) Name() string {
	return "geo"
}

func (c *GeoIPCollector) Description() string {
	return "Enrich IP addresses with geolocation data via ip-api.com"
}

func (c *GeoIPCollector) IsActive() bool {
	return false
}

func (c *GeoIPCollector) Collect(caseID string, target string) ([]core.Evidence, error) {
	// API Request
	url := fmt.Sprintf("http://ip-api.com/json/%s", target)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("geoip request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geoip api returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse for metadata extraction
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	if result["status"] == "fail" {
		return nil, fmt.Errorf("geoip api error: %v", result["message"])
	}

	// Store Evidence File
	storageDir := filepath.Join("evidence_storage", caseID)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, err
	}

	fileName := fmt.Sprintf("geo_%s_%d.json", target, time.Now().Unix())
	filePath := filepath.Join(storageDir, fileName)
	if err := os.WriteFile(filePath, body, 0644); err != nil {
		return nil, err
	}

	// Hash
	hash := sha256.Sum256(body)
	hashStr := hex.EncodeToString(hash[:])

	evidence := core.Evidence{
		CaseID:      caseID,
		Collector:   "geo",
		FilePath:    filePath,
		FileHash:    hashStr,
		CollectedAt: time.Now(),
		Metadata: map[string]interface{}{
			"target":  target,
			"country": result["countryCode"], // US, DE
			"city":    result["city"],
			"isp":     result["isp"],
			"lat":     result["lat"],
			"lon":     result["lon"],
		},
	}

	return []core.Evidence{evidence}, nil
}
