package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spectre/spectre/internal/core"
)

// CreateCase inserts a new case into the database.
func CreateCase(c *core.Case) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
	}
	if c.UpdatedAt.IsZero() {
		c.UpdatedAt = time.Now()
	}
	if c.Status == "" {
		c.Status = "active"
	}

	query := `INSERT INTO cases (id, name, description, created_at, updated_at, status) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := DB.Exec(query, c.ID, c.Name, c.Description, c.CreatedAt, c.UpdatedAt, c.Status)
	if err != nil {
		return fmt.Errorf("failed to create case: %w", err)
	}

	return nil
}

// GetCase retrieves a case by its ID.
func GetCase(id string) (*core.Case, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	query := `SELECT id, name, description, created_at, updated_at, status FROM cases WHERE id = ?`
	row := DB.QueryRow(query, id)

	var c core.Case
	err := row.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt, &c.Status)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get case: %w", err)
	}

	return &c, nil
}
