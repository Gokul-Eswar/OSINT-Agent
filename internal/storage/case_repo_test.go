package storage

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spectre/spectre/internal/core"
)

func TestCreateAndGetCase(t *testing.T) {
	// Initialize in-memory database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Temporarily set global DB
	oldDB := DB
	DB = db
	defer func() { DB = oldDB }()

	// Apply schema
	if err := Migrate(); err != nil {
		t.Fatal(err)
	}

	// Create a test case
	testCase := &core.Case{
		ID:          "test-123",
		Name:        "Test Case",
		Description: "A test case for unit testing",
	}

	if err := CreateCase(testCase); err != nil {
		t.Fatalf("CreateCase failed: %v", err)
	}

	// Retrieve the case
	retrieved, err := GetCase("test-123")
	if err != nil {
		t.Fatalf("GetCase failed: %v", err)
	}

	if retrieved == nil {
		t.Fatal("expected retrieved case to not be nil")
	}

	if retrieved.Name != testCase.Name {
		t.Errorf("expected name %s, got %s", testCase.Name, retrieved.Name)
	}
}

func TestListCases(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	defer db.Close()
	oldDB := DB
	DB = db
	defer func() { DB = oldDB }()
	Migrate()

	CreateCase(&core.Case{ID: "c1", Name: "Case 1"})
	CreateCase(&core.Case{ID: "c2", Name: "Case 2"})

	cases, err := ListCases()
	if err != nil {
		t.Fatal(err)
	}

	if len(cases) != 2 {
		t.Errorf("expected 2 cases, got %d", len(cases))
	}
}
