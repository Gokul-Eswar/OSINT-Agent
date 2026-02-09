package extensions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Manager struct {
	pluginsDir  string
	registryURL string
	cachedList  []Extension
}

func NewManager(pluginsDir string, registryURL string) *Manager {
	return &Manager{
		pluginsDir:  pluginsDir,
		registryURL: registryURL,
	}
}

// ensureRegistryFetched ensures we have the list of extensions.
func (m *Manager) ensureRegistryFetched() error {
	if m.cachedList != nil {
		return nil
	}

	var extensions []Extension

	// Check if URL is actually a local file path
	// e.g. "file://C:/path/to/registry.json" or just "registry.json"
	isLocal := strings.HasPrefix(m.registryURL, "file://") || 
              (strings.HasSuffix(m.registryURL, ".json") && !strings.HasPrefix(m.registryURL, "http"))

	if isLocal {
		cleanPath := strings.TrimPrefix(m.registryURL, "file://")
		data, err := os.ReadFile(cleanPath)
		if err != nil {
			return fmt.Errorf("failed to read local registry file '%s': %w", cleanPath, err)
		}
		
		var reg Registry
		if err := json.Unmarshal(data, &reg); err != nil {
			return fmt.Errorf("failed to parse local registry JSON: %w", err)
		}
		extensions = reg.Extensions

	} else {
		// HTTP Fetch
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(m.registryURL)
		if err != nil {
			return fmt.Errorf("failed to fetch registry from %s: %w", m.registryURL, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("registry returned status %d", resp.StatusCode)
		}

		var reg Registry
		if err := json.NewDecoder(resp.Body).Decode(&reg); err != nil {
			return fmt.Errorf("failed to decode registry from URL: %w", err)
		}
		extensions = reg.Extensions
	}

	m.cachedList = extensions
	return nil
}

func (m *Manager) Search(query string) ([]Extension, error) {
	if err := m.ensureRegistryFetched(); err != nil {
		return nil, err
	}

	var results []Extension
	query = strings.ToLower(query)

	for _, ext := range m.cachedList {
		if strings.Contains(strings.ToLower(ext.Name), query) ||
			strings.Contains(strings.ToLower(ext.Description), query) {
			results = append(results, ext)
		}
	}
	return results, nil
}

func (m *Manager) ListRemote() ([]Extension, error) {
	if err := m.ensureRegistryFetched(); err != nil {
		return nil, err
	}
	return m.cachedList, nil
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
	if err := m.ensureRegistryFetched(); err != nil {
		return err
	}

	// 1. Find the extension
	var target *Extension
	for _, ext := range m.cachedList {
		if ext.Name == extName {
			target = &ext
			break
		}
	}

	// Because we are iterating over a slice of structs, target is a pointer to the loop variable.
	// But since we break immediately, it's fine. 
	// However, correct Go idiom is:
	if target == nil {
		// Re-check just to be safe if loop finished
		found := false
		for i := range m.cachedList {
			if m.cachedList[i].Name == extName {
				target = &m.cachedList[i]
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("extension '%s' not found in registry", extName)
		}
	}

	// 2. Check if already installed
	destPath := filepath.Join(m.pluginsDir, target.Name)
	if _, err := os.Stat(destPath); err == nil {
		return fmt.Errorf("extension '%s' is already installed", target.Name)
	}

	// 3. Real Git Clone
	fmt.Printf("Installing %s from %s...\n", target.Name, target.URL)

    // Ensure git is installed
    if _, err := exec.LookPath("git"); err != nil {
        return fmt.Errorf("git is not installed or not in PATH. Please install git to fetch extensions")
    }

	cmd := exec.Command("git", "clone", "--depth", "1", target.URL, destPath)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("git clone failed: %w", err)
    }

	// 4. Verification (Optional but recommended)
    // Check if plugin.yaml exists in the cloned repo
    if _, err := os.Stat(filepath.Join(destPath, "plugin.yaml")); err != nil {
        // Rollback
        // os.RemoveAll(destPath) // DISABLED for safety, user might want to inspect
        return fmt.Errorf("warning: 'plugin.yaml' missing in repository. Extension may not load")
    }

	return nil
}

func (m *Manager) Remove(extName string) error {
    destPath := filepath.Join(m.pluginsDir, extName)
    if _, err := os.Stat(destPath); os.IsNotExist(err) {
        return fmt.Errorf("extension '%s' is not installed", extName)
    }
    return os.RemoveAll(destPath)
}