package analysis

import (
	"database/sql"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spectre/spectre/internal/core"
	"github.com/spectre/spectre/internal/storage"
)

func TestBuildCaseContext(t *testing.T) {
	// Initialize in-memory database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Temporarily set global DB
	storage.DB = db
	if err := storage.Migrate(); err != nil {
		t.Fatal(err)
	}

	caseID := "test-case-1"
	c := &core.Case{ID: caseID, Name: "Operation Midnight", Description: "Tracking mysterious activity"}
	if err := storage.CreateCase(c); err != nil {
		t.Fatal(err)
	}

	storage.CreateEntity(&core.Entity{ID: "e1", CaseID: caseID, Type: "ip", Value: "192.168.1.1", Source: "scanner"})
	storage.CreateEntity(&core.Entity{ID: "e2", CaseID: caseID, Type: "domain", Value: "evil.com", Source: "dns"})

	storage.CreateRelationship(&core.Relationship{
		ID:           "r1",
		CaseID:       caseID,
		FromEntityID: "e1",
		ToEntityID:   "e2",
		Type:         "resolves_to",
	})

	context, err := BuildCaseContext(caseID)
	if err != nil {
		t.Fatalf("BuildCaseContext failed: %v", err)
	}

	if !strings.Contains(context, "Operation Midnight") {
		t.Errorf("Expected context to contain case name")
	}
	if !strings.Contains(context, "192.168.1.1") {
		t.Errorf("Expected context to contain entity value")
	}
	if !strings.Contains(context, "resolves_to") {
		t.Errorf("Expected context to contain relationship type")
	}
}
