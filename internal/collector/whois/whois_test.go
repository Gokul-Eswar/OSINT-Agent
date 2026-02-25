package whois

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockWhoisClient struct {
	Response string
	Error    error
}

func (m *MockWhoisClient) Whois(domain string) (string, error) {
	return m.Response, m.Error
}

func TestWHOISCollector_Collect_Mocked(t *testing.T) {
	// Setup mock with sample WHOIS data
	sampleData := `
Domain Name: EXAMPLE.COM
Registrar: SafeNames Ltd
Registrant Email: admin@example.com
`
	mock := &MockWhoisClient{Response: sampleData}
	c := &WHOISCollector{client: mock}
	caseID := "test_case_mock_whois"
	target := "example.com"

	// Cleanup
	defer os.RemoveAll("evidence_storage")

	// Execute
	evidence, err := c.Collect(caseID, target)

	// Verify
	require.NoError(t, err)
	require.Len(t, evidence, 1)

	e := evidence[0]
	assert.Equal(t, "whois", e.Collector)
	assert.Contains(t, e.Metadata["registrar"], "SafeNames")
	assert.Equal(t, "admin@example.com", e.Metadata["registrant_email"])

	// Verify file content
	content, err := os.ReadFile(e.FilePath)
	require.NoError(t, err)
	assert.Equal(t, sampleData, string(content))
}
