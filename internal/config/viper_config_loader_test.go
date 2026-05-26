package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestInitConfig(t *testing.T) {
	// Create a temporary config file
	tmpDir, err := os.MkdirTemp("", "spectre-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cfgPath := filepath.Join(tmpDir, "config.yaml")
	content := []byte("database:\n  path: \"test.db\"\nkeys:\n  test: \"secret\"")
	if err := os.WriteFile(cfgPath, content, 0644); err != nil {
		t.Fatal(err)
	}

	// Initialize config
	InitConfig(cfgPath)

	// Verify value
	dbPath := viper.GetString("database.path")
	if dbPath != "test.db" {
		t.Errorf("expected test.db, got %s", dbPath)
	}

	// Test GetAPIKey
	key := GetAPIKey("test")
	if key != "secret" {
		t.Errorf("expected secret, got %s", key)
	}
}

func TestInitConfigLight(t *testing.T) {
	err := InitConfigLight("")
	if err != nil {
		t.Fatalf("InitConfigLight failed: %v", err)
	}
}
