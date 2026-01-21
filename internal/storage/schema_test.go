package storage

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestInitSchema(t *testing.T) {
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
	if err := InitSchema(); err != nil {
		t.Fatalf("InitSchema failed: %v", err)
	}

	// Verify tables exist
	tables := []string{"cases", "entities", "relationships", "evidence", "analyses"}
	for _, table := range tables {
		var name string
		err := DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		if err != nil {
			t.Errorf("table %s not found: %v", table, err)
		}
	}
}