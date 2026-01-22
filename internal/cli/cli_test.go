package cli

import (
	"os"
	"testing"

	"github.com/spectre/spectre/internal/config"
	"github.com/spectre/spectre/internal/storage"
)

func TestCLICommands(t *testing.T) {
	// Setup temporary database for testing
	dbPath := "test_cli.db"
	defer os.Remove(dbPath)

	// Mock config
	config.InitConfig("")
	
	// Temporarily set DB path in config if possible, 
	// or just rely on InitDB using default and then Cleanup.
	// For simplicity in this E2E-style unit test:
	
	t.Run("InitCommand", func(t *testing.T) {
		rootCmd.SetArgs([]string{"init"})
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("init command failed: %v", err)
		}
		
		if _, err := os.Stat("spectre.db"); os.IsNotExist(err) {
			t.Error("spectre.db was not created by init command")
		}
		defer os.Remove("spectre.db")
	})

	t.Run("CaseNewCommand", func(t *testing.T) {
		// Ensure DB is initialized
		storage.InitDB()
		storage.InitSchema()
		defer storage.CloseDB()
		defer os.Remove("spectre.db")

		rootCmd.SetArgs([]string{"case", "new", "test-cli-case"})
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("case new command failed: %v", err)
		}
		
		// Verify case exists
		_, err := storage.GetCase("some-id") // We don't know the ID easily here without extra logic
		// But we can check if any case exists
		rows, err := storage.DB.Query("SELECT name FROM cases WHERE name='test-cli-case'")
		if err != nil {
			t.Fatal(err)
		}
		defer rows.Close()
		if !rows.Next() {
			t.Error("case was not found in database")
		}
	})
}
