package storage

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spectre/spectre/internal/core"
)

func TestCreateAndGetEntity(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	oldDB := DB
	DB = db
	defer func() { DB = oldDB }()

	if err := InitSchema(); err != nil {
		t.Fatal(err)
	}

	// Need a case first due to foreign key
	c := &core.Case{ID: "case-1", Name: "Case 1"}
	if err := CreateCase(c); err != nil {
		t.Fatal(err)
	}

	entity := &core.Entity{
		ID:     "ent-1",
		CaseID: "case-1",
		Type:   "ip",
		Value:  "1.1.1.1",
		Source: "manual",
		Metadata: map[string]interface{}{
			"note": "test entity",
		},
	}

	if err := CreateEntity(entity); err != nil {
		t.Fatalf("CreateEntity failed: %v", err)
	}

	retrieved, err := GetEntity("ent-1")
	if err != nil {
		t.Fatalf("GetEntity failed: %v", err)
	}

	if retrieved == nil {
		t.Fatal("expected retrieved entity to not be nil")
	}

	if retrieved.Value != "1.1.1.1" {
		t.Errorf("expected value 1.1.1.1, got %s", retrieved.Value)
	}

	if retrieved.Metadata["note"] != "test entity" {
		t.Errorf("expected metadata note 'test entity', got %v", retrieved.Metadata["note"])
	}
}

func TestListEntitiesByCase(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	oldDB := DB
	DB = db
	defer func() { DB = oldDB }()

	if err := InitSchema(); err != nil {
		t.Fatal(err)
	}

	c1 := &core.Case{ID: "case-1", Name: "Case 1"}
	CreateCase(c1)
	c2 := &core.Case{ID: "case-2", Name: "Case 2"}
	CreateCase(c2)

	CreateEntity(&core.Entity{ID: "e1", CaseID: "case-1", Type: "ip", Value: "1.1.1.1"})
	CreateEntity(&core.Entity{ID: "e2", CaseID: "case-1", Type: "domain", Value: "example.com"})
	CreateEntity(&core.Entity{ID: "e3", CaseID: "case-2", Type: "ip", Value: "8.8.8.8"})

	entities, err := ListEntitiesByCase("case-1")
	if err != nil {
		t.Fatal(err)
	}

	if len(entities) != 2 {
		t.Errorf("expected 2 entities for case-1, got %d", len(entities))
	}
}
