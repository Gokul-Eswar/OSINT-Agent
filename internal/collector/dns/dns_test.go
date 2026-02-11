package dns

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDNSCollector_Collect(t *testing.T) {
	// Setup
	c := &DNSCollector{}
	caseID := "test_case_dns"
	target := "example.com"

	// Cleanup
	defer os.RemoveAll("evidence_storage")

	// Execute
	evidence, err := c.Collect(caseID, target)

	// Verify
	require.NoError(t, err)
	require.Len(t, evidence, 1)

	e := evidence[0]
	assert.Equal(t, caseID, e.CaseID)
	assert.Equal(t, "dns", e.Collector)
	assert.NotEmpty(t, e.FilePath)
	assert.NotEmpty(t, e.FileHash)
	
	// Check content
	var results map[string][]string
	data, ok := e.RawData.(map[string][]string)
	if !ok {
		// If it was unmarshaled from JSON during a reload, it might be map[string]interface{}
		// But here we return it directly.
		// However, let's verify the file content too.
		content, err := os.ReadFile(e.FilePath)
		require.NoError(t, err)
		err = json.Unmarshal(content, &results)
		require.NoError(t, err)
	} else {
		results = data
	}

	assert.NotEmpty(t, results["A"])
	// example.com might not have MX or NS depending on the resolver, but usually has A.
}
