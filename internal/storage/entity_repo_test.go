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

	if err := Migrate(); err != nil {
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

func TestGetEntityByValue(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	oldDB := DB
	DB = db
	defer func() { DB = oldDB }()

	if err := Migrate(); err != nil {
		t.Fatal(err)
	}

	c1 := &core.Case{ID: "case-1", Name: "Case 1"}
	CreateCase(c1)

	e1 := &core.Entity{ID: "e1", CaseID: "case-1", Type: "ip", Value: "1.1.1.1"}
	CreateEntity(e1)

	retrieved, err := GetEntityByValue("case-1", "1.1.1.1")
	if err != nil {
		t.Fatalf("GetEntityByValue failed: %v", err)
	}

	if retrieved == nil {
		t.Fatal("expected retrieved entity to not be nil")
	}

	if retrieved.ID != "e1" {
		t.Errorf("expected ID e1, got %s", retrieved.ID)
	}

	// Test non-existent value
	retrieved, err = GetEntityByValue("case-1", "2.2.2.2")
	if err != nil {
		t.Fatalf("GetEntityByValue for non-existent failed: %v", err)
	}
	if retrieved != nil {
		t.Errorf("expected nil for non-existent value, got %v", retrieved)
	}
}
