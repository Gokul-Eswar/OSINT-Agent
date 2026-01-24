package collector

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spectre/spectre/internal/core"
	"gopkg.in/yaml.v3"
)

// PluginMetadata defines the structure of the plugin.yaml file.
type PluginMetadata struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Command     string   `yaml:"command"`
	Args        []string `yaml:"args"`
	IsActive    bool     `yaml:"is_active"`
}

// ExternalCollector implements the core.Collector interface for external scripts.
type ExternalCollector struct {
	metadata PluginMetadata
	path     string
}

func (e *ExternalCollector) Name() string {
	return e.metadata.Name
}

func (e *ExternalCollector) Description() string {
	return e.metadata.Description
}

func (e *ExternalCollector) IsActive() bool {
	return e.metadata.IsActive
}

// Collect executes the external plugin and captures its output.
func (e *ExternalCollector) Collect(caseID string, target string) ([]core.Evidence, error) {
	// Prepare command
	args := append(e.metadata.Args, target)
	cmd := exec.Command(e.metadata.Command, args...)
	cmd.Dir = e.path

	// Run and capture output
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("plugin execution failed: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("failed to execute plugin: %w", err)
	}

	// Store file
	storageDir := filepath.Join("evidence_storage", caseID)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, err
	}

	fileName := fmt.Sprintf("%s_%s_%d.json", e.metadata.Name, target, time.Now().Unix())
	filePath := filepath.Join(storageDir, fileName)
	if err := os.WriteFile(filePath, output, 0644); err != nil {
		return nil, err
	}

	// Hash
	hash := sha256.Sum256(output)
	hashStr := hex.EncodeToString(hash[:])

	evidence := core.Evidence{
		CaseID:      caseID,
		Collector:   e.metadata.Name,
		FilePath:    filePath,
		FileHash:    hashStr,
		CollectedAt: time.Now(),
		Metadata: map[string]interface{}{
			"target": target,
			"source": "external_plugin",
		},
	}

	return []core.Evidence{evidence}, nil
}

// DiscoverPlugins scans the plugins directory for valid plugins.
func DiscoverPlugins() ([]core.Collector, error) {
	pluginsDir := "plugins"
	entries, err := os.ReadDir(pluginsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var collectors []core.Collector
	for _, entry := range entries {
		if entry.IsDir() {
			pluginDir := filepath.Join(pluginsDir, entry.Name())
			metadataPath := filepath.Join(pluginDir, "plugin.yaml")

			if _, err := os.Stat(metadataPath); err == nil {
				data, err := os.ReadFile(metadataPath)
				if err != nil {
					log.Error().Err(err).Str("path", metadataPath).Msg("Failed to read plugin metadata")
					continue
				}

				var meta PluginMetadata
				if err := yaml.Unmarshal(data, &meta); err != nil {
					log.Error().Err(err).Str("path", metadataPath).Msg("Failed to parse plugin metadata")
					continue
				}

				log.Info().Str("name", meta.Name).Msg("Loaded external plugin")
				collectors = append(collectors, &ExternalCollector{
					metadata: meta,
					path:     pluginDir,
				})
			}
		}
	}

	return collectors, nil
}
