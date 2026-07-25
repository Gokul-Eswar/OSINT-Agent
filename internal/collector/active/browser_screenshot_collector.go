package active

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/spectre/spectre/internal/core"
	"github.com/spectre/spectre/internal/ethics"
	"github.com/spf13/viper"
)

type ScreenshotCollector struct{}

func (c *ScreenshotCollector) Name() string {
	return "screenshot"
}

func (c *ScreenshotCollector) Description() string {
	return "Captures full-page screenshots of the target domain"
}

func (c *ScreenshotCollector) IsActive() bool {
	return true
}

func (c *ScreenshotCollector) Collect(caseID string, target string, options map[string]interface{}) ([]core.Evidence, error) {
	if err := ethics.Wait("screenshot"); err != nil {
		return nil, err
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
	)

	if browserBin := findBrowserExecutable(); browserBin != "" {
		opts = append(opts, chromedp.ExecPath(browserBin))
	}

	if os.Getenv("CI") != "" {
		opts = append(opts,
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-setuid-sandbox", true),
			chromedp.Flag("disable-dev-shm-usage", true),
		)
	}

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
	timeoutDuration := 30 * time.Second
	if t, ok := options["timeout"].(int); ok && t > 0 {
		timeoutDuration = time.Duration(t) * time.Second
	} else if cfgTimeout := viper.GetInt("collectors.screenshot.timeout"); cfgTimeout > 0 {
		timeoutDuration = time.Duration(cfgTimeout) * time.Second
	}
	ctx, cancel = context.WithTimeout(ctx, timeoutDuration)
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

func findBrowserExecutable() string {
	// 1. Check ENV Overrides
	for _, env := range []string{"CHROME_BIN", "BROWSER_BIN", "CHROMIUM_BIN"} {
		if val := os.Getenv(env); val != "" {
			if _, err := os.Stat(val); err == nil {
				return val
			}
		}
	}

	// 2. Check Viper Config Override
	if cfgPath := viper.GetString("browser.executable"); cfgPath != "" {
		if _, err := os.Stat(cfgPath); err == nil {
			return cfgPath
		}
	}

	// 3. Search system PATH for common browser binaries
	browsers := []string{"google-chrome", "chromium", "chromium-browser", "chrome", "msedge", "brave"}
	for _, b := range browsers {
		if p, err := exec.LookPath(b); err == nil {
			return p
		}
	}

	// 4. Standard OS installation paths fallback
	switch runtime.GOOS {
	case "windows":
		paths := []string{
			`C:\Program Files\Google\Chrome\Application\chrome.exe`,
			`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
			`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,
			filepath.Join(os.Getenv("LocalAppData"), `Google\Chrome\Application\chrome.exe`),
		}
		for _, p := range paths {
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
	case "darwin":
		paths := []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
			"/Applications/Brave Browser.app/Contents/MacOS/Brave Browser",
			"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
		}
		for _, p := range paths {
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
	default:
		paths := []string{
			"/usr/bin/google-chrome",
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
			"/usr/bin/brave-browser",
			"/snap/bin/chromium",
		}
		for _, p := range paths {
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
	}

	return ""
}
