package active

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/spectre/spectre/internal/collector"
	"github.com/spectre/spectre/internal/core"
	"github.com/spectre/spectre/internal/ethics"
	"github.com/spf13/viper"
)

type ScreenshotCollector struct{}

func init() {
	collector.Register(&ScreenshotCollector{})
}

func (c *ScreenshotCollector) Name() string {
	return "screenshot"
}

func (c *ScreenshotCollector) Description() string {
	return "Captures full-page screenshots of the target domain"
}

func (c *ScreenshotCollector) IsActive() bool {
	return true
}

func (c *ScreenshotCollector) Collect(caseID string, target string) ([]core.Evidence, error) {
	if err := ethics.Wait("screenshot"); err != nil {
		return nil, err
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
	)

	// Proxy Logic (Respect Ghost Mode)
	var proxy string
	if viper.GetBool("ghost_mode") {
		proxy = viper.GetString("http.tor_proxy")
		if proxy == "" {
			proxy = "socks5://127.0.0.1:9050"
		}
	} else {
		proxy = viper.GetString("http.proxy")
	}

	if proxy != "" {
		opts = append(opts, chromedp.ProxyServer(proxy))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Setup context
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	url := target
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = fmt.Sprintf("http://%s", target)
	}
	var buf []byte

	// Run tasks
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.FullScreenshot(&buf, 90),
	)
	if err != nil {
		return nil, fmt.Errorf("screenshot failed: %w", err)
	}

	// Store file
	storageDir := filepath.Join("evidence_storage", caseID)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, err
	}

	safeTarget := strings.ReplaceAll(target, "://", "_")
	safeTarget = strings.ReplaceAll(safeTarget, ":", "_")
	safeTarget = strings.ReplaceAll(safeTarget, "/", "_")
	safeTarget = strings.ReplaceAll(safeTarget, "\\", "_")

	fileName := fmt.Sprintf("screenshot_%s_%d.png", safeTarget, time.Now().Unix())
	filePath := filepath.Join(storageDir, fileName)
	if err := os.WriteFile(filePath, buf, 0644); err != nil {
		return nil, err
	}

	// Hash
	hash := sha256.Sum256(buf)
	hashStr := hex.EncodeToString(hash[:])

	evidence := core.Evidence{
		CaseID:      caseID,
		Collector:   "screenshot",
		FilePath:    filePath,
		FileHash:    hashStr,
		CollectedAt: time.Now(),
		Metadata: map[string]interface{}{
			"target": target,
			"size":   len(buf),
			"type":   "image/png",
		},
	}

	return []core.Evidence{evidence}, nil
}
