package storage

import (
	"embed"
	"fmt"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Migrate runs all pending migrations in the migrations directory.
func Migrate() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	// 1. Create migrations table if not exists
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY)`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// 2. Read migration files
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	// 3. Apply pending migrations
	for _, file := range files {
		var existing string
		err := DB.QueryRow("SELECT version FROM schema_migrations WHERE version = ?", file).Scan(&existing)
		if err == nil {
			// Already applied
			continue
		}

		log.Info().Str("migration", file).Msg("applying database migration")

		content, err := migrationsFS.ReadFile("migrations/" + file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		// Execute migration in a transaction
		tx, err := DB.Begin()
		if err != nil {
			return fmt.Errorf("failed to start transaction for %s: %w", file, err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", file); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", file, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", file, err)
		}

		log.Info().Str("migration", file).Msg("migration applied successfully")
	}

	return nil
}
