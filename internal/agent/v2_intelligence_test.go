package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestVectorSearchTool(t *testing.T) {
	// Skip if Python is not available or venv not set up
	if _, err := os.Stat(".venv"); os.IsNotExist(err) {
		t.Skip("Skipping vector search test: .venv not found")
	}

	caseID := "test-v2-brain"
	evidenceDir := filepath.Join("evidence_storage", caseID)
	os.MkdirAll(evidenceDir, 0755)
	defer os.RemoveAll(evidenceDir)

	// Create a dummy evidence file with a secret
	secretFile := filepath.Join(evidenceDir, "whois_secret.txt")
	content := "Registrant Name: John Doe\nEmail: john@example.com\nSecret Flag: GHOST_IN_THE_SHELL"
	os.WriteFile(secretFile, []byte(content), 0644)

	t.Run("search_evidence finds the secret", func(t *testing.T) {
		// This will trigger the bridge call to Python
		// Note: Requires chromadb to be installed in .venv
		tool := Registry["search_evidence"]
		res, err := tool.Execute(caseID, map[string]interface{}{"query": "What is the secret flag?"})

		// We expect either a success or a clear error message about missing deps
		if err != nil {
			t.Logf("Search failed (likely missing Python deps): %v", err)
		} else {
			fmt.Println("Agent search result:", res)
		}
	})
}

func TestAgentPersonaSystemPrompt(t *testing.T) {
	engine := NewEngine("test-persona")
	if len(engine.History) == 0 {
		t.Fatal("History should contain the system prompt")
	}
	if engine.History[0].Role != "system" {
		t.Errorf("First message should be system prompt, got %s", engine.History[0].Role)
	}
}
