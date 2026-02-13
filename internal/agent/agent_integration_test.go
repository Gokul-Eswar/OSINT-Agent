package agent

import (
	"database/sql"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spectre/spectre/internal/core"
	"github.com/spectre/spectre/internal/storage"
)

func setupTestDB(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	storage.DB = db
	if err := storage.Migrate(); err != nil {
		t.Fatal(err)
	}
}

func TestAgentToolIntegration(t *testing.T) {
	setupTestDB(t)
	caseID := "test-case-agent"
	
	// Seed database
	storage.CreateCase(&core.Case{
		ID:   caseID,
		Name: "Agent Test Case",
	})
	
	storage.CreateEntity(&core.Entity{
		CaseID: caseID,
		Type:   "domain",
		Value:  "spectre.local",
		Source: "manual",
	})

	t.Run("list_collectors execution", func(t *testing.T) {
		tool := Registry["list_collectors"]
		res, err := tool.Execute(caseID, nil)
		if err != nil {
			t.Fatalf("list_collectors failed: %v", err)
		}
		if res == "" {
			t.Error("expected output from list_collectors")
		}
	})

	t.Run("search_entities execution", func(t *testing.T) {
		tool := Registry["search_entities"]
		res, err := tool.Execute(caseID, map[string]interface{}{"query": "spectre"})
		if err != nil {
			t.Fatalf("search_entities failed: %v", err)
		}
		if !strings.Contains(res, "spectre.local") {
			t.Errorf("expected search result to contain 'spectre.local', got: %s", res)
		}
	})

	t.Run("get_case_summary execution", func(t *testing.T) {
		tool := Registry["get_case_summary"]
		res, err := tool.Execute(caseID, nil)
		if err != nil {
			t.Fatalf("get_case_summary failed: %v", err)
		}
		if !strings.Contains(res, "Agent Test Case") {
			t.Errorf("expected summary to contain case name, got: %s", res)
		}
	})
}
