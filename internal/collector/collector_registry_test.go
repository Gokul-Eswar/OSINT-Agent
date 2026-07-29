package collector

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/Gokul-Eswar/Spectre/internal/core"
	"github.com/Gokul-Eswar/Spectre/internal/storage"
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

func TestRegistry_RegisterDuplicate(t *testing.T) {
	r := NewRegistry()
	c := &mockCollector{}
	if err := r.Register(c); err != nil {
		t.Fatal(err)
	}
	if err := r.Register(c); err == nil {
		t.Error("expected error when registering duplicate collector")
	}
}

func TestRegistry_Get(t *testing.T) {
	r := NewRegistry()
	c := &mockCollector{}
	r.Register(c)

	got, err := r.Get("mock")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name() != "mock" {
		t.Errorf("got %s, want mock", got.Name())
	}

	_, err = r.Get("nonexistent")
	if err == nil {
		t.Error("expected error when getting nonexistent collector")
	}
}

func TestRegistry_List(t *testing.T) {
	r := NewRegistry()
	c := &mockCollector{}
	r.Register(c)

	list := r.List()
	if len(list) != 1 {
		t.Errorf("expected 1 collector, got %d", len(list))
	}
}

func TestRun_ActiveConsent(t *testing.T) {
	r := NewRegistry()
	active := &activeMockCollector{}
	r.Register(active)

	// Save the old global registry and restore it after
	oldRegistry := globalRegistry
	globalRegistry = r
	defer func() { globalRegistry = oldRegistry }()

	// Try to run active collector without consent
	_, err := Run("active_mock", "case-1", "example.com", false, nil)
	if err == nil {
		t.Error("expected error when running active collector without consent")
	}

	// Run with consent
	_, err = Run("active_mock", "case-1", "example.com", true, nil)
	if err != nil {
		t.Fatalf("Run failed with consent: %v", err)
	}
}

type activeMockCollector struct{}

func (m *activeMockCollector) Name() string        { return "active_mock" }
func (m *activeMockCollector) Description() string { return "Active Mock" }
func (m *activeMockCollector) IsActive() bool      { return true }
func (m *activeMockCollector) Collect(caseID string, target string, options map[string]interface{}) ([]core.Evidence, error) {
	return []core.Evidence{}, nil
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
