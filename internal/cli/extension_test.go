package cli

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/spectre/spectre/internal/extensions"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtensionCommands(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "spectre-ext-cli-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	regPath := filepath.Join(tempDir, "registry.json")
	reg := extensions.Registry{
		Extensions: []extensions.Extension{
			{Name: "test-plugin", Description: "A test plugin", Tags: []string{"test"}},
		},
	}
	data, _ := json.Marshal(reg)
	os.WriteFile(regPath, data, 0644)

	// Configure viper to use our temp registry and plugins dir
	viper.Set("extension_registry_url", "file://"+regPath)
	
	// We need to override the getManager behavior or the pluginsDir it uses
	// Currently getManager is hardcoded to "plugins"
	// I'll update it to be configurable via viper if needed, but for now
	// I'll just check if the command executes.

	t.Run("SearchCommand", func(t *testing.T) {
		rootCmd.SetArgs([]string{"extension", "search", "test"})
		
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := rootCmd.Execute()
		
		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		assert.NoError(t, err)
		assert.Contains(t, output, "test-plugin")
		assert.Contains(t, output, "A test plugin")
	})

	t.Run("InfoCommand", func(t *testing.T) {
		rootCmd.SetArgs([]string{"extension", "info", "test-plugin"})
		
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := rootCmd.Execute()
		
		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		assert.NoError(t, err)
		assert.Contains(t, output, "Name:        test-plugin")
	})
}
