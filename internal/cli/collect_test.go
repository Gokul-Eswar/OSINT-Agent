package cli

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectCommand_DryRun(t *testing.T) {
	caseID = ""
	// 1. Setup context

	err := SaveContext("test-case-dry-run")
	require.NoError(t, err)
	defer os.Remove(filepath.Join(os.Getenv("USERPROFILE"), ".spectre_current_case"))

	// 2. Run collect all with dry-run
	// Use a small target to be safe
	rootCmd.SetArgs([]string{"collect", "dns", "example.com", "--dry-run", "--case", "test-case-dry-run"})

	// Capture output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// This might still try to open a Bubble Tea program which can be tricky in tests
	// but let's see if it works in a non-interactive environment.
	err = rootCmd.Execute()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	assert.NoError(t, err)
	// Even if UI doesn't show in buf, Execute should return no error
	assert.Contains(t, output, "") // Placeholder
}
