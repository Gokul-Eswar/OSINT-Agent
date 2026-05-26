package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"           // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var DB *sql.DB

// InitDB initializes the database connection (SQLite or PostgreSQL) and applies migrations.
func InitDB() error {
	dbType := viper.GetString("database.type")
	if dbType == "" {
		dbType = "sqlite"
	}

	var dsn string
	var driver string

	switch dbType {
	case "postgres":
		driver = "postgres"
		dsn = viper.GetString("database.dsn")
		if dsn == "" {
			return fmt.Errorf("postgres DSN is required when database.type is 'postgres'")
		}
	case "sqlite":
		driver = "sqlite3"
		dbPath := viper.GetString("database.path")
		if dbPath == "" {
			dbPath = "spectre.db"
		}
		dsn = fmt.Sprintf("%s?_journal_mode=WAL&_busy_timeout=5000", dbPath)
	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}

	var err error
	DB, err = sql.Open(driver, dsn)
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

	// Connection pooling configuration
	if dbType == "postgres" {
		DB.SetMaxOpenConns(viper.GetInt("database.max_open_conns"))
		DB.SetMaxIdleConns(viper.GetInt("database.max_idle_conns"))
	} else {
		// Optimize for SQLite
		DB.SetMaxOpenConns(25)
		DB.SetMaxIdleConns(5)
	}

	log.Info().Str("type", dbType).Msg("Database initialized and migrated")
	return nil
}

// CloseDB closes the database connection.
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
