package github

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGitHubCollector_Collect(t *testing.T) {
	// Skip if no internet
	if os.Getenv("CI") != "" {
		t.Skip("Skipping GitHub test in CI")
	}

	c := &GitHubCollector{}
	caseID := "test_case_github"
	// Use a very specific target that is likely to exist but small, or just a keyword
	target := "spectre-cli-test-target" 
	
	// Reset config
	viper.Reset()

	// Cleanup
	defer os.RemoveAll("evidence_storage")

	// Execute
	evidence, err := c.Collect(caseID, target)

	// Verify
	require.NoError(t, err)
	require.Len(t, evidence, 1)

	e := evidence[0]
	assert.Equal(t, caseID, e.CaseID)
	assert.Equal(t, "github", e.Collector)
	
	// Check content
	content, err := os.ReadFile(e.FilePath)
	require.NoError(t, err)
	
	var result map[string]interface{}
	err = json.Unmarshal(content, &result)
	require.NoError(t, err)
	
	// GitHub Search API response has "items" or "total_count"
	_, hasItems := result["items"]
	_, hasTotal := result["total_count"]
	assert.True(t, hasItems || hasTotal)
}
