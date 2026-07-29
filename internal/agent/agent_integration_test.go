package agent

import (
	"database/sql"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/Gokul-Eswar/Spectre/internal/core"
	"github.com/Gokul-Eswar/Spectre/internal/storage"
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

	t.Run("update_hypotheses execution", func(t *testing.T) {
		tool := Registry["update_hypotheses"]
		args := map[string]interface{}{
			"hypothesis":         "Whois indicates the server host is located in the US.",
			"confidence":         0.75,
			"evidence_filenames": []interface{}{"whois_result.txt"},
			"status":             "active",
		}
		res, err := tool.Execute(caseID, args)
		if err != nil {
			t.Fatalf("update_hypotheses failed: %v", err)
		}
		if !strings.Contains(res, "recorded successfully") {
			t.Errorf("expected output to contain 'recorded successfully', got: %s", res)
		}

		// Verify database
		leads, err := storage.ListLeadsByCase(caseID)
		if err != nil {
			t.Fatalf("ListLeadsByCase failed: %v", err)
		}
		if len(leads) != 1 {
			t.Fatalf("expected 1 lead, got %d", len(leads))
		}
		if leads[0].Hypothesis != "Whois indicates the server host is located in the US." {
			t.Errorf("expected hypothesis to match, got: %s", leads[0].Hypothesis)
		}
		if leads[0].Confidence != 0.75 {
			t.Errorf("expected confidence 0.75, got %f", leads[0].Confidence)
		}
	})
}
