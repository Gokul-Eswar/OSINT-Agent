package collector

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spectre/spectre/internal/core"
	"github.com/spectre/spectre/internal/storage"
)

type mockCollector struct{}

func (m *mockCollector) Name() string        { return "mock" }
func (m *mockCollector) Description() string { return "Mock Collector" }
func (m *mockCollector) IsActive() bool      { return false }
func (m *mockCollector) Collect(caseID string, target string, options map[string]interface{}) ([]core.Evidence, error) {
	return []core.Evidence{
		{
			CaseID:    caseID,
			Collector: "mock",
			FilePath:  "mock.json",
			FileHash:  "abc",
			Metadata:  map[string]interface{}{"target": target},
		},
	}, nil
}

func TestRunAndSave(t *testing.T) {
	// Setup DB
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	storage.DB = db
	storage.Migrate()

	caseID := "test-case"
	storage.CreateCase(&core.Case{ID: caseID, Name: "Test"})

	// Register mock
	Register(&mockCollector{})

	// Run
	evList, err := RunAndSave("mock", caseID, "example.com", false, nil)
	if err != nil {
		t.Fatalf("RunAndSave failed: %v", err)
	}

	if len(evList) != 1 {
		t.Errorf("Expected 1 evidence item, got %d", len(evList))
	}

	// Verify it's in DB
	ev, err := storage.ListEvidenceByCase(caseID)
	if err != nil {
		t.Fatal(err)
	}
	if len(ev) != 1 {
		t.Errorf("Expected 1 evidence item in DB, got %d", len(ev))
	}
}
