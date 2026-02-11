package whois

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWHOISCollector_Collect(t *testing.T) {
	// Skip if no internet connection or if whois servers block CI
	if os.Getenv("CI") != "" {
		t.Skip("Skipping WHOIS test in CI environment")
	}

	c := &WHOISCollector{}
	caseID := "test_case_whois"
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
	assert.Equal(t, "whois", e.Collector)
	assert.NotEmpty(t, e.FilePath)
	
	// Verify file created
	content, err := os.ReadFile(e.FilePath)
	require.NoError(t, err)
	assert.NotEmpty(t, content)
	assert.Contains(t, string(content), "Domain Name") // Common in whois output
}
