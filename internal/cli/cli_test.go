package cli

import (
	"os"
	"testing"

	"github.com/Gokul-Eswar/Spectre/internal/config"
	"github.com/Gokul-Eswar/Spectre/internal/storage"
)

func TestMain(m *testing.M) {
	tempDir, _ := os.MkdirTemp("", "spectre-cli-test-main")
	os.Setenv("SPECTRE_HOME", tempDir)
	os.Setenv("GO_TESTING", "true")

	code := m.Run()

	os.RemoveAll(tempDir)
	os.Exit(code)
}

func TestCLICommands(t *testing.T) {

	// Mock config
	config.InitConfig("")

	t.Run("InitCommand", func(t *testing.T) {
		// Cleanup if exists
		os.Remove("spectre.db")
		os.Remove("spectre.db-wal")
		os.Remove("spectre.db-shm")

		rootCmd.SetArgs([]string{"init"})
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("init command failed: %v", err)
		}

		if _, err := os.Stat("spectre.db"); os.IsNotExist(err) {
			t.Error("spectre.db was not created by init command")
		}
		defer os.Remove("spectre.db")
		defer os.Remove("spectre.db-wal")
		defer os.Remove("spectre.db-shm")
	})

	t.Run("CaseNewCommand", func(t *testing.T) {
		// Use in-memory DB for logic test if possible, but CLI uses InitDB
		// which is hardcoded to spectre.db if not overridden in config.
		// For unit test, we just cleanup.
		os.Remove("spectre.db")
		os.Remove("spectre.db-wal")
		os.Remove("spectre.db-shm")

		if err := storage.InitDB(); err != nil {
			t.Fatal(err)
		}
		defer storage.CloseDB()
		defer os.Remove("spectre.db")
		defer os.Remove("spectre.db-wal")
		defer os.Remove("spectre.db-shm")

		rootCmd.SetArgs([]string{"case", "new", "test-cli-case"})
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("case new command failed: %v", err)
		}

		// Verify case exists
		rows, err := storage.DB.Query("SELECT name FROM cases WHERE name='test-cli-case'")
		if err != nil {
			t.Fatal(err)
		}
		defer rows.Close()
		if !rows.Next() {
			t.Error("case was not found in database")
		}
	})

	t.Run("VersionCommand", func(t *testing.T) {
		rootCmd.SetArgs([]string{"version"})
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("version command failed: %v", err)
		}
	})

	t.Run("HelpGuideFunction", func(t *testing.T) {
		printHelpGuide()
	})
}
