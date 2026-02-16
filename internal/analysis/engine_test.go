package analysis

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spectre/spectre/internal/analyzer"
	"github.com/spectre/spectre/internal/core"
	"github.com/spectre/spectre/internal/storage"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyzeCase_Flow(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "spectre-analysis-flow-test")
	defer os.RemoveAll(tempDir)

	// Mock DB
	viper.Set("database.path", filepath.Join(tempDir, "test.db"))
	err := storage.InitDB()
	require.NoError(t, err)
	defer storage.CloseDB()

	caseID := "test-case"
	storage.CreateCase(&core.Case{ID: caseID, Name: "Test Case"})

	// Create a mock python script
	mockScript := filepath.Join(tempDir, "mock_analyzer.py")
	scriptContent := `
import json
import sys
result = {
    "confidence": 0.95,
    "findings": ["Everything is fine"],
    "risks": ["None"],
    "next_steps": ["Relax"]
}
print(json.dumps(result))
`
	os.WriteFile(mockScript, []byte(scriptContent), 0644)

	// Override PythonCommand to run our mock script
	oldCmd := analyzer.PythonCommand
	analyzer.PythonCommand = []string{mockScript}
	defer func() { analyzer.PythonCommand = oldCmd }()

	// 1. First run (Cache Miss)
	res, err := AnalyzeCase(caseID, "mock-model")
	require.NoError(t, err)
	assert.Equal(t, 0.95, res.Confidence)
	assert.Equal(t, "Everything is fine", res.Findings[0])

	// 2. Second run (Cache Hit)
	// Even if we change the mock script, it should return the cached result
	os.WriteFile(mockScript, []byte(`print(json.dumps({"confidence": 0.1}))`), 0644)
	
	res2, err := AnalyzeCase(caseID, "mock-model")
	require.NoError(t, err)
	assert.Equal(t, 0.95, res2.Confidence, "Should have returned cached result")
}
