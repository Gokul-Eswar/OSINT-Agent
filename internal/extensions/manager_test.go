package extensions

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager_ListInstalled(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "spectre-plugins-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock plugin
	pluginName := "test-plugin"
	pluginDir := filepath.Join(tempDir, pluginName)
	err = os.MkdirAll(pluginDir, 0755)
	require.NoError(t, err)

	pluginYaml := `
name: test-plugin
description: A test plugin
tags: ["test", "mock"]
`
	err = os.WriteFile(filepath.Join(pluginDir, "plugin.yaml"), []byte(pluginYaml), 0644)
	require.NoError(t, err)

	mgr := NewManager(tempDir, "")
	installed, err := mgr.ListInstalled()
	require.NoError(t, err)

	assert.Len(t, installed, 1)
	assert.Equal(t, pluginName, installed[0].Name)
	assert.Equal(t, "A test plugin", installed[0].Description)
	assert.Contains(t, installed[0].Tags, "test")
}

func TestManager_Search(t *testing.T) {
	tempReg := filepath.Join(os.TempDir(), "registry.json")
	registry := Registry{
		Extensions: []Extension{
			{Name: "shodan", Description: "Shodan search", Tags: []string{"passive"}},
			{Name: "nmap", Description: "Port scanner", Tags: []string{"active"}},
		},
	}
	data, _ := json.Marshal(registry)
	_ = os.WriteFile(tempReg, data, 0644)
	defer os.Remove(tempReg)

	mgr := NewManager("", tempReg)

	// Test search by name
	results, err := mgr.Search("shodan")
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "shodan", results[0].Name)

	// Test search by description
	results, err = mgr.Search("Port scanner")
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "nmap", results[0].Name)
}

func TestManager_GetInfo(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "spectre-plugins-info-test")
	defer os.RemoveAll(tempDir)

	tempReg := filepath.Join(tempDir, "registry.json")
	registry := Registry{
		Extensions: []Extension{
			{Name: "test-ext", Description: "Registry version", Tags: []string{"reg"}},
		},
	}
	data, _ := json.Marshal(registry)
	_ = os.WriteFile(tempReg, data, 0644)

	mgr := NewManager(tempDir, tempReg)

	// 1. Not installed, in registry
	ext, installed, err := mgr.GetInfo("test-ext")
	require.NoError(t, err)
	assert.False(t, installed)
	assert.Equal(t, "test-ext", ext.Name)

	// 2. Installed, not in registry (custom)
	pluginDir := filepath.Join(tempDir, "custom-plugin")
	_ = os.MkdirAll(pluginDir, 0755)
	_ = os.WriteFile(filepath.Join(pluginDir, "plugin.yaml"), []byte("name: custom-plugin"), 0644)

	ext, installed, err = mgr.GetInfo("custom-plugin")
	require.NoError(t, err)
	assert.True(t, installed)
	assert.Equal(t, "custom-plugin", ext.Name)
}
