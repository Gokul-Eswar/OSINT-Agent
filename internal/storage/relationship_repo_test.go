package storage

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spectre/spectre/internal/core"
)

func TestCreateAndGetRelationship(t *testing.T) {
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

	c := &core.Case{ID: "case-1", Name: "Case 1"}
	CreateCase(c)

	e1 := &core.Entity{ID: "e1", CaseID: "case-1", Type: "ip", Value: "1.1.1.1"}
	CreateEntity(e1)
	e2 := &core.Entity{ID: "e2", CaseID: "case-1", Type: "domain", Value: "example.com"}
	CreateEntity(e2)

	rel := &core.Relationship{
		ID:           "rel-1",
		CaseID:       "case-1",
		FromEntityID: "e1",
		ToEntityID:   "e2",
		Type:         "resolves_to",
		Confidence:   1.0,
	}

	if err := CreateRelationship(rel); err != nil {
		t.Fatalf("CreateRelationship failed: %v", err)
	}

	retrieved, err := GetRelationship("rel-1")
	if err != nil {
		t.Fatalf("GetRelationship failed: %v", err)
	}

	if retrieved == nil {
		t.Fatal("expected retrieved relationship to not be nil")
	}

	if retrieved.Type != "resolves_to" {
		t.Errorf("expected type resolves_to, got %s", retrieved.Type)
	}
}

func TestListRelationshipsByCase(t *testing.T) {
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

	e1 := &core.Entity{ID: "e1", CaseID: "case-1", Type: "ip", Value: "1.1.1.1"}
	CreateEntity(e1)
	e2 := &core.Entity{ID: "e2", CaseID: "case-1", Type: "domain", Value: "example.com"}
	CreateEntity(e2)

	CreateRelationship(&core.Relationship{ID: "r1", CaseID: "case-1", FromEntityID: "e1", ToEntityID: "e2", Type: "rel1"})
	CreateRelationship(&core.Relationship{ID: "r2", CaseID: "case-1", FromEntityID: "e2", ToEntityID: "e1", Type: "rel2"})

	rels, err := ListRelationshipsByCase("case-1")
	if err != nil {
		t.Fatal(err)
	}

	if len(rels) != 2 {
		t.Errorf("expected 2 relationships for case-1, got %d", len(rels))
	}
}
