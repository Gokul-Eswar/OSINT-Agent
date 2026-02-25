package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var DB *sql.DB

// InitDB initializes the SQLite database connection and applies migrations.
func InitDB() error {
	dbPath := viper.GetString("database.path")
	if dbPath == "" {
		dbPath = "spectre.db"
	}

	var err error
	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_busy_timeout=5000", dbPath)
	DB, err = sql.Open("sqlite3", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Verify connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations
	if err := Migrate(); err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}

	// Optimize connection pooling for SQLite
	// SQLite supports multiple readers but only one writer at a time.
	// WAL mode and busy_timeout handle most cases, but we can also limit connections.
	DB.SetMaxOpenConns(25) // Allow some concurrency for readers
	DB.SetMaxIdleConns(5)

	log.Info().Str("path", dbPath).Msg("SQLite database initialized and migrated")
	return nil
}

// CloseDB closes the database connection.
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
