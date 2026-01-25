package active

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spectre/spectre/internal/collector"
	"github.com/spectre/spectre/internal/core"
	"github.com/spectre/spectre/internal/ethics"
	"github.com/spectre/spectre/internal/http" // Import path is directory, not package name
)

type SocialCollector struct {
	Sites map[string]string
}

func init() {
	collector.Register(NewSocialCollector())
}

func NewSocialCollector() *SocialCollector {
	return &SocialCollector{
		Sites: map[string]string{
			"GitHub":    "https://github.com/%s",
			"Twitter":   "https://twitter.com/%s",
			"Instagram": "https://www.instagram.com/%s",
			"Reddit":    "https://www.reddit.com/user/%s",
			"Facebook":  "https://www.facebook.com/%s",
			"GitLab":    "https://gitlab.com/%s",
			"Medium":    "https://medium.com/@%s",
			"YouTube":   "https://www.youtube.com/@%s",
			"Twitch":    "https://www.twitch.tv/%s",
			"TikTok":    "https://www.tiktok.com/@%s",
		},
	}
}

func (c *SocialCollector) Name() string {
	return "social"
}

func (c *SocialCollector) Description() string {
	return "Checks for username availability across social media sites"
}

func (c *SocialCollector) IsActive() bool {
	return true
}

type SiteResult struct {
	Site   string `json:"site"`
	URL    string `json:"url"`
	Status string `json:"status"` // "found", "not_found", "error"
}

func (c *SocialCollector) Collect(caseID string, target string) ([]core.Evidence, error) {
	// Target is assumed to be the username
	username := target
	
	var results []SiteResult
	var mu sync.Mutex
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5) // Limit concurrency

	client := netclient.NewClient()

	for site, urlTmpl := range c.Sites {
		wg.Add(1)
		go func(site, urlTmpl string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := ethics.Wait("social"); err != nil {
				return
			}

			checkURL := fmt.Sprintf(urlTmpl, username)
			status := "error"

			req, err := http.NewRequest("GET", checkURL, nil)
			if err == nil {
				req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
				resp, err := client.Do(req)
				if err == nil {
					defer resp.Body.Close()
					if resp.StatusCode == 200 {
						status = "found"
					} else if resp.StatusCode == 404 {
						status = "not_found"
					} else {
						status = fmt.Sprintf("http_%d", resp.StatusCode)
					}
				} else {
                    status = "connection_error"
                }
			}

			if status == "found" {
				mu.Lock()
				results = append(results, SiteResult{
					Site:   site,
					URL:    checkURL,
					Status: status,
				})
				mu.Unlock()
			}
		}(site, urlTmpl)
	}

	wg.Wait()

	if len(results) == 0 {
		return nil, nil // No evidence found
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

	fileName := fmt.Sprintf("social_%s_%d.json", username, time.Now().Unix())
	filePath := filepath.Join(storageDir, fileName)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return nil, err
	}

	// Hash
	hash := sha256.Sum256(data)
	hashStr := hex.EncodeToString(hash[:])

	evidence := core.Evidence{
		CaseID:      caseID,
		Collector:   "social",
		FilePath:    filePath,
		FileHash:    hashStr,
		CollectedAt: time.Now(),
		Metadata: map[string]interface{}{
			"target": username,
			"count":  len(results),
		},
	}

	return []core.Evidence{evidence}, nil
}
