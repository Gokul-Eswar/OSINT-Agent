package agent

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spectre/spectre/internal/analyzer"
	"github.com/spectre/spectre/internal/core"
	"github.com/spectre/spectre/internal/storage"
)

func TestFullAgentLoop(t *testing.T) {
	// Skip if python is not available
	if _, err := os.Stat("../../analyzer/mock_agent.py"); os.IsNotExist(err) {
		t.Skip("mock_agent.py not found")
	}

	setupTestDB(t)
	caseID := "test-loop-case"
	storage.CreateCase(&core.Case{
		ID:   caseID,
		Name: "Loop Test Case",
	})

	// Override Python command to use our mock script
	oldCmd := analyzer.PythonCommand
	analyzer.PythonCommand = []string{"../../analyzer/mock_agent.py"}
	defer func() { analyzer.PythonCommand = oldCmd }()

	engine := NewEngine(caseID)

	t.Run("Agent Tool Use Loop", func(t *testing.T) {
		// This should trigger:
		// 1. User: "run dns"
		// 2. LLM (mock) -> tool_use: collect(dns)
		// 3. Go -> executes dns tool
		// 4. LLM (mock) -> content: "I have finished..."

		response, err := engine.Execute("Please run dns on google.com")
		if err != nil {
			t.Fatalf("Agent loop failed: %v", err)
		}

		if !strings.Contains(response, "I have finished the DNS collection") {
			t.Errorf("Expected final response about DNS completion, got: %s", response)
		}

		// Verify history length (User, Assistant/ToolUse, System/ToolResult, Assistant/Final)
		if len(engine.History) < 4 {
			t.Errorf("Expected at least 4 messages in history, got %d", len(engine.History))
		}
	})

	t.Run("Agent Whois Flag Integration Test", func(t *testing.T) {
		// 1. Create a dummy WHOIS evidence file with a secret flag
		evidenceDir := filepath.Join("evidence_storage", caseID)
		os.MkdirAll(evidenceDir, 0755)
		defer os.RemoveAll(evidenceDir)

		secretFile := filepath.Join(evidenceDir, "whois_flag.txt")
		content := "Registrant: John Doe\nFlag: SPECTRE_FLAG_998877\n"
		if err := os.WriteFile(secretFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write mock whois flag file: %v", err)
		}

		// 2. Initialize a fresh engine
		flagEngine := NewEngine(caseID)

		// 3. Run execution
		response, err := flagEngine.Execute("Please search WHOIS flag inside evidence files")
		if err != nil {
			t.Fatalf("Agent WHOIS flag loop failed: %v", err)
		}

		if !strings.Contains(response, "SPECTRE_FLAG_998877") {
			t.Errorf("Expected response to contain 'SPECTRE_FLAG_998877', got: %s", response)
		}
	})
}
