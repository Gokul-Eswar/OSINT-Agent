package storage

import (
	"embed"
	"fmt"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// TranslatePlaceholder handles dialect differences between SQLite (?) and Postgres ($1).
func TranslatePlaceholder(query string) string {
	if viper.GetString("database.type") != "postgres" {
		return query
	}

	// Simple placeholder translation
	count := 1
	for strings.Contains(query, "?") {
		query = strings.Replace(query, "?", fmt.Sprintf("$%d", count), 1)
		count++
	}
	return query
}

// TranslateDialect handles minor SQL differences for SQLite-style migrations.
func TranslateDialect(sql string) string {
	if viper.GetString("database.type") != "postgres" {
		return sql
	}

	// Basic translations for Postgres
	replacer := strings.NewReplacer(
		"DATETIME", "TIMESTAMP",
		"REAL", "DOUBLE PRECISION",
		"JSON", "JSONB",
	)
	return replacer.Replace(sql)
}

// Migrate runs all pending migrations in the migrations directory.
func Migrate() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	// Dialect translation
	dbType := viper.GetString("database.type")

	createTableSQL := "CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY)"
	if dbType == "postgres" {
		createTableSQL = "CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY)"
	}

	// 1. Create migrations table if not exists
	_, err := DB.Exec(createTableSQL)
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
		query := TranslatePlaceholder("SELECT version FROM schema_migrations WHERE version = ?")
		err := DB.QueryRow(query, file).Scan(&existing)
		if err == nil {
			// Already applied
			continue
		}

		log.Info().Str("migration", file).Msg("applying database migration")

		content, err := migrationsFS.ReadFile("migrations/" + file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		migrationSQL := TranslateDialect(string(content))

		// Execute migration in a transaction
		tx, err := DB.Begin()
		if err != nil {
			return fmt.Errorf("failed to start transaction for %s: %w", file, err)
		}

		if _, err := tx.Exec(migrationSQL); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		insertQuery := TranslatePlaceholder("INSERT INTO schema_migrations (version) VALUES (?)")
		if _, err := tx.Exec(insertQuery, file); err != nil {
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
