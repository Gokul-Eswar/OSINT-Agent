package active

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spectre/spectre/internal/collector"
	"github.com/spectre/spectre/internal/core"
)

type HTTPCollector struct{}

func init() {
	collector.Register(&HTTPCollector{})
}

func (c *HTTPCollector) Name() string {
	return "http"
}

func (c *HTTPCollector) Description() string {
	return "Active HTTP service discovery (Headers, Title)"
}

func (c *HTTPCollector) IsActive() bool {
	return true
}

func (c *HTTPCollector) Collect(caseID string, target string) ([]core.Evidence, error) {
	// Try HTTPS first
	url := target
	if !strings.HasPrefix(url, "http") {
		url = "https://" + target
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		// Fallback to HTTP
		url = "http://" + target
		resp, err = client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("http request failed: %w", err)
		}
	}
	defer resp.Body.Close()

	results := make(map[string]interface{})
	results["url"] = url
	results["status_code"] = resp.StatusCode
	
	headers := make(map[string]string)
	for k, v := range resp.Header {
		headers[k] = strings.Join(v, ", ")
	}
	results["headers"] = headers

	// Extract Title
	bodyBytes := make([]byte, 4096) // Read first 4KB
	n, _ := resp.Body.Read(bodyBytes)
	body := string(bodyBytes[:n])
	
	re := regexp.MustCompile(`(?i)<title>(.*?)</title>`)
	match := re.FindStringSubmatch(body)
	if len(match) > 1 {
		results["title"] = match[1]
	}

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return nil, err
	}

	// Store file
	storageDir := filepath.Join("evidence_storage", caseID)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, err
	}

	fileName := fmt.Sprintf("http_%s_%d.json", target, time.Now().Unix())
	filePath := filepath.Join(storageDir, fileName)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return nil, err
	}

	// Hash
	hash := sha256.Sum256(data)
	hashStr := hex.EncodeToString(hash[:])

	evidence := core.Evidence{
		CaseID:      caseID,
		Collector:   "http",
		FilePath:    filePath,
		FileHash:    hashStr,
		CollectedAt: time.Now(),
		Metadata: map[string]interface{}{
			"target": target,
			"server": headers["Server"],
			"title":  results["title"],
		},
	}

	return []core.Evidence{evidence}, nil
}
