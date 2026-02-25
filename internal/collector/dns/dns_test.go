package dns

import (
	"encoding/json"
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockResolver struct {
	Hosts []string
	MXs   []*net.MX
	NSs   []*net.NS
}

func (m *MockResolver) LookupHost(host string) ([]string, error) { return m.Hosts, nil }
func (m *MockResolver) LookupMX(name string) ([]*net.MX, error)  { return m.MXs, nil }
func (m *MockResolver) LookupNS(name string) ([]*net.NS, error)  { return m.NSs, nil }

func TestDNSCollector_Collect_Mocked(t *testing.T) {
	// Setup mock
	mock := &MockResolver{
		Hosts: []string{"1.2.3.4", "5.6.7.8"},
		MXs:   []*net.MX{{Host: "mail.example.com", Pref: 10}},
		NSs:   []*net.NS{{Host: "ns1.example.com"}},
	}
	
	c := &DNSCollector{resolver: mock}
	caseID := "test_case_mock_dns"
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
	
	// Check content
	content, err := os.ReadFile(e.FilePath)
	require.NoError(t, err)
	
	var results map[string][]string
	err = json.Unmarshal(content, &results)
	require.NoError(t, err)

	assert.ElementsMatch(t, mock.Hosts, results["A"])
	assert.Contains(t, results["MX"], "mail.example.com")
	assert.Contains(t, results["NS"], "ns1.example.com")
}
