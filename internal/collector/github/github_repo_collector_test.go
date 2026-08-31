package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGitHubCollector_Collect(t *testing.T) {
	mockResponse := map[string]interface{}{
		"total_count": 1,
		"items": []map[string]interface{}{
			{
				"full_name": "spectre/test-target",
				"html_url":  "https://github.com/spectre/test-target",
				"owner": map[string]interface{}{
					"login": "spectre",
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	c := &GitHubCollector{BaseURL: server.URL}
	caseID := "test_case_github"
	target := "spectre-cli-test-target"

	// Reset config
	viper.Reset()

	// Cleanup
	defer os.RemoveAll("evidence_storage")

	// Execute
	evidence, err := c.Collect(caseID, target, nil)

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
