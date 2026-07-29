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

	"github.com/rs/zerolog/log"
	"github.com/Gokul-Eswar/Spectre/internal/core"
	"github.com/Gokul-Eswar/Spectre/internal/ethics"
	netclient "github.com/Gokul-Eswar/Spectre/internal/http"
)

type SocialCollector struct {
	Sites map[string]string
}

func NewSocialCollector() *SocialCollector {
	return &SocialCollector{
		Sites: map[string]string{
			"GitHub":         "https://github.com/%s",
			"Twitter":        "https://twitter.com/%s",
			"Instagram":      "https://www.instagram.com/%s",
			"Reddit":         "https://www.reddit.com/user/%s",
			"Facebook":       "https://www.facebook.com/%s",
			"GitLab":         "https://gitlab.com/%s",
			"Medium":         "https://medium.com/@%s",
			"YouTube":        "https://www.youtube.com/@%s",
			"Twitch":         "https://www.twitch.tv/%s",
			"TikTok":         "https://www.tiktok.com/@%s",
			"Pinterest":      "https://www.pinterest.com/%s/",
			"Snapchat":       "https://www.snapchat.com/add/%s",
			"Steam":          "https://steamcommunity.com/id/%s",
			"SoundCloud":     "https://soundcloud.com/%s",
			"Spotify":        "https://open.spotify.com/user/%s",
			"Mastodon":       "https://mastodon.social/@%s",
			"Behance":        "https://www.behance.net/%s",
			"Dribbble":       "https://dribbble.com/%s",
			"Patreon":        "https://www.patreon.com/%s",
			"Telegram":       "https://t.me/%s",
			"Archive.org":    "https://archive.org/details/@%s",
			"Keybase":        "https://keybase.io/%s",
			"Letterboxd":     "https://letterboxd.com/%s/",
			"MyAnimeList":    "https://myanimelist.net/profile/%s",
			"Duolingo":       "https://www.duolingo.com/profile/%s",
			"Chess.com":      "https://www.chess.com/member/%s",
			"Codewars":       "https://www.codewars.com/users/%s",
			"Docker Hub":     "https://hub.docker.com/u/%s",
			"Flickr":         "https://www.flickr.com/people/%s",
			"Goodreads":      "https://www.goodreads.com/user/show/%s",
			"Gumroad":        "https://gumroad.com/%s",
			"HackerOne":      "https://hackerone.com/%s",
			"IfThisThenThat": "https://ifttt.com/p/%s",
			"Issuu":          "https://issuu.com/%s",
			"Kaggle":         "https://www.kaggle.com/%s",
			"Last.fm":        "https://www.last.fm/user/%s",
			"Linktree":       "https://linktr.ee/%s",
			"MySpace":        "https://myspace.com/%s",
			"Pastebin":       "https://pastebin.com/u/%s",
			"ProductHunt":    "https://www.producthunt.com/@%s",
			"Quora":          "https://www.quora.com/profile/%s",
			"ReverbNation":   "https://www.reverbnation.com/%s",
			"Scribd":         "https://www.scribd.com/%s",
			"Slack":          "https://%s.slack.com",
			"Slideshare":     "https://www.slideshare.net/%s",
			"Substack":       "https://%s.substack.com",
			"Trakt":          "https://trakt.tv/users/%s",
			"TripAdvisor":    "https://www.tripadvisor.com/Profile/%s",
			"Vimeo":          "https://vimeo.com/%s",
			"Wattpad":        "https://www.wattpad.com/user/%s",
			"Wikipedia":      "https://en.wikipedia.org/wiki/User:%s",
			"Wordpress":      "https://%s.wordpress.com",
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

func (c *SocialCollector) Collect(caseID string, target string, options map[string]interface{}) ([]core.Evidence, error) {
	log.Info().
		Str("collector", "social").
		Str("case_id", caseID).
		Str("username", target).
		Msg("collection_started")

	// Target is assumed to be the username
	username := target

	var results []SiteResult
	var mu sync.Mutex
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5) // Limit concurrency

	client, err := netclient.NewClient()
	if err != nil {
		log.Error().Err(err).Msg("failed to create http client")
		return nil, fmt.Errorf("failed to create http client: %w", err)
	}

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
		log.Info().
			Str("collector", "social").
			Str("case_id", caseID).
			Str("username", username).
			Msg("collection_completed_no_results")
		return nil, nil // No evidence found
	}

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal social results")
		return nil, fmt.Errorf("failed to marshal social results: %w", err)
	}

	// Store file
	storageDir := filepath.Join("evidence_storage", caseID)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		log.Error().Err(err).Str("dir", storageDir).Msg("failed to create storage directory")
		return nil, fmt.Errorf("failed to create storage directory %s: %w", storageDir, err)
	}

	fileName := fmt.Sprintf("social_%s_%d.json", username, time.Now().Unix())
	filePath := filepath.Join(storageDir, fileName)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		log.Error().Err(err).Str("path", filePath).Msg("failed to write social evidence file")
		return nil, fmt.Errorf("failed to write social evidence file %s: %w", filePath, err)
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

	log.Info().
		Str("collector", "social").
		Str("case_id", caseID).
		Str("username", username).
		Int("results_count", len(results)).
		Msg("collection_completed")

	return []core.Evidence{evidence}, nil
}
