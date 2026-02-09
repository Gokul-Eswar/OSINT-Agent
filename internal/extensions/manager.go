package extensions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Mock Registry for now - in production this would be fetched from a remote JSON/YAML
var defaultRegistry = Registry{
	Extensions: []Extension{
		{
			Name:        "shodan_lookup",
			Description: "Queries Shodan API for IP information",
			Author:      "SpectreTeam",
			Version:     "1.0.0",
			URL:         "https://github.com/spectre-plugins/shodan_lookup",
			Type:        "collector",
			Tags:        []string{"osint", "ip", "passive"},
		},
		{
			Name:        "whois_advanced",
			Description: "Deep WHOIS lookup with historical data",
			Author:      "Community",
			Version:     "0.5.0",
			URL:         "https://github.com/spectre-plugins/whois_advanced",
			Type:        "collector",
			Tags:        []string{"whois", "domain"},
		},
		{
			Name:        "git_leaks",
			Description: "Scans repositories for secrets",
			Author:      "SecurityResearch",
			Version:     "2.1.0",
			URL:         "https://github.com/spectre-plugins/git_leaks",
			Type:        "collector",
			Tags:        []string{"git", "secrets", "active"},
		},
		{
			Name:        "subdomain_brute",
			Description: "Fast subdomain brute-forcing tool",
			Author:      "RedTeamOps",
			Version:     "1.2.0",
			URL:         "https://github.com/spectre-plugins/subdomain_brute",
			Type:        "collector",
			Tags:        []string{"dns", "active", "subdomains"},
		},
	},
}

type Manager struct {
	pluginsDir string
}

func NewManager(pluginsDir string) *Manager {
	return &Manager{pluginsDir: pluginsDir}
}

func (m *Manager) Search(query string) ([]Extension, error) {
	// In a real scenario, this would fetch from a remote URL.
	// For now, filtering the defaultRegistry.
	var results []Extension
	query = strings.ToLower(query)

	for _, ext := range defaultRegistry.Extensions {
		if strings.Contains(strings.ToLower(ext.Name), query) ||
			strings.Contains(strings.ToLower(ext.Description), query) {
			results = append(results, ext)
		}
	}
	return results, nil
}

func (m *Manager) ListRemote() ([]Extension, error) {
	return defaultRegistry.Extensions, nil
}

func (m *Manager) ListInstalled() ([]string, error) {
	entries, err := os.ReadDir(m.pluginsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var installed []string
	for _, e := range entries {
		if e.IsDir() {
			// Check for plugin.yaml
			if _, err := os.Stat(filepath.Join(m.pluginsDir, e.Name(), "plugin.yaml")); err == nil {
				installed = append(installed, e.Name())
			}
		}
	}
	return installed, nil
}

func (m *Manager) Install(extName string) error {
	// 1. Find the extension
	var target *Extension
	for _, ext := range defaultRegistry.Extensions {
		if ext.Name == extName {
			target = &ext
			break
		}
	}

	if target == nil {
		return fmt.Errorf("extension '%s' not found in registry", extName)
	}

	// 2. Check if already installed
	destPath := filepath.Join(m.pluginsDir, target.Name)
	if _, err := os.Stat(destPath); err == nil {
		return fmt.Errorf("extension '%s' is already installed", target.Name)
	}

	// 3. Clone/Download
	// Since we are simulating, and these repos don't actually exist, 
	// I will mock the installation by creating the folder and a dummy plugin.yaml/script.
	// IN REALITY: cmd := exec.Command("git", "clone", target.URL, destPath)

	fmt.Printf("Installing %s from %s...\n", target.Name, target.URL)

	// MOCK INSTALLATION START
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return err
	}

	isActive := false
	for _, tag := range target.Tags {
		if tag == "active" {
			isActive = true
			break
		}
	}

	dummyYaml := fmt.Sprintf(`name: "%s"
description: "%s"
command: "python"
args: ["main.py"]
is_active: %t
`, target.Name, target.Description, isActive)

	if err := os.WriteFile(filepath.Join(destPath, "plugin.yaml"), []byte(dummyYaml), 0644); err != nil {
		return err
	}

	dummyPy := `import sys
import json

# This is a placeholder script for the installed extension.
# In a real scenario, this would be the actual tool logic.

try:
    target = sys.argv[1]
except IndexError:
    target = "unknown"

print(json.dumps({
    "source": "extension_store",
    "status": "success", 
    "message": "Extension executed successfully",
    "target": target
}))
`
	if err := os.WriteFile(filepath.Join(destPath, "main.py"), []byte(dummyPy), 0644); err != nil {
		return err
	}
	// MOCK INSTALLATION END

	return nil
}

func (m *Manager) Remove(extName string) error {
    destPath := filepath.Join(m.pluginsDir, extName)
    if _, err := os.Stat(destPath); os.IsNotExist(err) {
        return fmt.Errorf("extension '%s' is not installed", extName)
    }
    return os.RemoveAll(destPath)
}