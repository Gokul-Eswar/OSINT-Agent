package collector

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spectre/spectre/internal/core"
	"go.yaml.in/yaml/v3"
)

// ExternalPluginMetadata defines the structure of plugin.yaml
type ExternalPluginMetadata struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Command     string `yaml:"command"`     // e.g. "python scanner.py"
	IsActive    bool   `yaml:"is_active"`   // Whether this is an active probe
}

// ExternalCollector wraps an external script/binary
type ExternalCollector struct {
	Meta ExternalPluginMetadata
	Path string // Directory containing the plugin
}

func (e *ExternalCollector) Name() string {
	return e.Meta.Name
}

func (e *ExternalCollector) Description() string {
	return e.Meta.Description
}

func (e *ExternalCollector) IsActive() bool {
	return e.Meta.IsActive
}

func (e *ExternalCollector) Collect(caseID string, target string) ([]core.Evidence, error) {
	// Execute command: [command] [target]
	args := append(strings.Split(e.Meta.Command, " "), target)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = e.Path

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("plugin execution failed: %w\nOutput: %s", err, string(output))
	}

	// The plugin must print JSON evidence to stdout
	var evidence []core.Evidence
	if err := json.Unmarshal(output, &evidence); err != nil {
		return nil, fmt.Errorf("failed to parse plugin output as evidence JSON: %w\nRaw: %s", err, string(output))
	}

	// Ensure case ID is set correctly for all evidence
	for i := range evidence {
		evidence[i].CaseID = caseID
		evidence[i].Collector = e.Name()
	}

	return evidence, nil
}

// DiscoverPlugins scans the plugins/ directory
func DiscoverPlugins() ([]core.Collector, error) {
	pluginsDir := "plugins"
	if _, err := os.Stat(pluginsDir); os.IsNotExist(err) {
		return nil, nil // No plugins directory
	}

	dirs, err := os.ReadDir(pluginsDir)
	if err != nil {
		return nil, err
	}

	var collectors []core.Collector
	for _, d := range dirs {
		if d.IsDir() {
			pluginPath := filepath.Join(pluginsDir, d.Name())
			metaPath := filepath.Join(pluginPath, "plugin.yaml")
			
			if _, err := os.Stat(metaPath); err == nil {
				data, err := os.ReadFile(metaPath)
				if err != nil {
					continue
				}

				var meta ExternalPluginMetadata
				if err := yaml.Unmarshal(data, &meta); err != nil {
					continue
				}

				collectors = append(collectors, &ExternalCollector{
					Meta: meta,
					Path: pluginPath,
				})
			}
		}
	}

	return collectors, nil
}
