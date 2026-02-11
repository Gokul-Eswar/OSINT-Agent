package storage

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestMigrate(t *testing.T) {
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

	// Run migrations
	if err := Migrate(); err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	// Verify tables exist
	tables := []string{"cases", "entities", "relationships", "evidence", "analyses", "schema_migrations"}
	for _, table := range tables {
		var name string
		err := DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		if err != nil {
			t.Errorf("table %s not found: %v", table, err)
		}
	}

	// Run migrations again (should be idempotent)
	if err := Migrate(); err != nil {
		t.Fatalf("Second Migrate failed: %v", err)
	}
}
