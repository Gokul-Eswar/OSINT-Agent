package agent

import (
	"testing"
)

func TestToolRegistry(t *testing.T) {
	// Verify core tools are present
	tools := []string{"collect", "list_collectors", "search_entities", "get_case_summary"}
	for _, name := range tools {
		if _, ok := Registry[name]; !ok {
			t.Errorf("Tool '%s' missing from registry", name)
		}
	}
}

func TestGetToolDefinitions(t *testing.T) {
	defs := GetToolDefinitions()
	if len(defs) == 0 {
		t.Error("Tool definitions should not be empty")
	}

	// Check one definition structure
	found := false
	for _, d := range defs {
		m := d.(map[string]interface{})
		if m["name"] == "collect" {
			found = true
			if m["description"] == "" {
				t.Error("Collect tool missing description")
			}
			break
		}
	}

	if !found {
		t.Error("Collect tool definition not found in GetToolDefinitions output")
	}
}

// Note: Testing actual tool execution requires a database mock or a temporary test database.
// Given the project structure, we can do a basic check on logic that doesn't hit the DB if we refactor,
// but for now, we'll focus on registry and definition integrity.
