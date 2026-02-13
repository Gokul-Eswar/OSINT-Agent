package agent

import (
	"os"
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
}
