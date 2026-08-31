package cli

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/Gokul-Eswar/Spectre/internal/analyzer"
	"github.com/Gokul-Eswar/Spectre/internal/core"
	"github.com/Gokul-Eswar/Spectre/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchCommand(t *testing.T) {
	caseID = ""
	// Setup temporary DB

	tempDir, err := os.MkdirTemp("", "spectre-search-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "spectre.db")
	os.Setenv("SPECTRE_DB_PATH", dbPath)
	defer os.Unsetenv("SPECTRE_DB_PATH")

	err = storage.InitDB()
	require.NoError(t, err)
	defer storage.CloseDB()

	// 1. Create a case
	caseID := "search-test-case"
	err = storage.CreateCase(&core.Case{
		ID:   caseID,
		Name: "Search Test Case",
	})
	require.NoError(t, err)

	// 2. Save context
	err = SaveContext(caseID)
	require.NoError(t, err)

	// 3. Create evidence file with content
	evidenceDir := filepath.Join(tempDir, "evidence", caseID)
	err = os.MkdirAll(evidenceDir, 0755)
	require.NoError(t, err)

	evidencePath := filepath.Join(evidenceDir, "evidence.json")
	err = os.WriteFile(evidencePath, []byte(`{"ip": "1.2.3.4", "target": "api.example.com"}`), 0644)
	require.NoError(t, err)

	// 4. Register evidence in DB
	err = storage.CreateEvidence(&core.Evidence{
		CaseID:    caseID,
		Collector: "test",
		FilePath:  evidencePath,
	})
	require.NoError(t, err)

	// 5. Run search command
	rootCmd.SetArgs([]string{"search", "api.example.com"})

	// Capture output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = rootCmd.Execute()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	assert.NoError(t, err)
	assert.Contains(t, output, "Found in:")
	assert.Contains(t, output, "api.example.com")
}

type mockSearchTaskRunner struct {
	response string
	err      error
}

func (m *mockSearchTaskRunner) Run(req analyzer.Request) (string, error) {
	return m.response, m.err
}

func TestSearchCommand_Semantic(t *testing.T) {
	caseID = ""
	semanticSearch = false

	tempDir, err := os.MkdirTemp("", "spectre-search-sem-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "spectre.db")
	os.Setenv("SPECTRE_DB_PATH", dbPath)
	defer os.Unsetenv("SPECTRE_DB_PATH")

	err = storage.InitDB()
	require.NoError(t, err)
	defer storage.CloseDB()

	testCaseID := "sem-search-case"
	err = storage.CreateCase(&core.Case{
		ID:   testCaseID,
		Name: "Semantic Search Case",
	})
	require.NoError(t, err)

	err = SaveContext(testCaseID)
	require.NoError(t, err)

	// Mock GlobalTaskRunner
	oldRunner := analyzer.GlobalTaskRunner
	defer func() { analyzer.GlobalTaskRunner = oldRunner }()

	analyzer.GlobalTaskRunner = &mockSearchTaskRunner{
		response: `{"status": "success", "results": [{"id": "ports.json#chunk_0", "content": "Port 80 is open running Apache", "distance": 0.12}]}`,
	}

	rootCmd.SetArgs([]string{"search", "--semantic", "open web port"})

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = rootCmd.Execute()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	assert.NoError(t, err)
	assert.Contains(t, output, "Performing semantic vector search")
	assert.Contains(t, output, "Port 80 is open running Apache")
}
