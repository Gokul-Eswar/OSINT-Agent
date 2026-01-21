package storage

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

const schema = `
CREATE TABLE IF NOT EXISTS cases (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    status TEXT DEFAULT 'active'
);

CREATE TABLE IF NOT EXISTS entities (
    id TEXT PRIMARY KEY,
    case_id TEXT NOT NULL,
    type TEXT NOT NULL,
    value TEXT NOT NULL,
    source TEXT,
    confidence REAL DEFAULT 0.5,
    discovered_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    metadata JSON,
    FOREIGN KEY (case_id) REFERENCES cases(id),
    UNIQUE(case_id, type, value)
);

CREATE TABLE IF NOT EXISTS relationships (
    id TEXT PRIMARY KEY,
    case_id TEXT NOT NULL,
    from_entity TEXT NOT NULL,
    to_entity TEXT NOT NULL,
    rel_type TEXT NOT NULL,
    confidence REAL DEFAULT 0.5,
    evidence_id TEXT,
    discovered_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (case_id) REFERENCES cases(id),
    FOREIGN KEY (from_entity) REFERENCES entities(id),
    FOREIGN KEY (to_entity) REFERENCES entities(id),
    UNIQUE(from_entity, to_entity, rel_type)
);

CREATE TABLE IF NOT EXISTS evidence (
    id TEXT PRIMARY KEY,
    case_id TEXT NOT NULL,
    entity_id TEXT,
    collector TEXT NOT NULL,
    file_path TEXT NOT NULL,
    file_hash TEXT NOT NULL,
    collected_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    metadata JSON,
    FOREIGN KEY (case_id) REFERENCES cases(id),
    FOREIGN KEY (entity_id) REFERENCES entities(id)
);

CREATE TABLE IF NOT EXISTS analyses (
    id TEXT PRIMARY KEY,
    case_id TEXT NOT NULL,
    findings JSON,
    risks JSON,
    connections JSON,
    next_steps JSON,
    confidence REAL,
    analyzed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (case_id) REFERENCES cases(id)
);
`

// InitSchema applies the initial database schema.
func InitSchema() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	_, err := DB.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to apply schema: %w", err)
	}

	log.Info().Msg("Database schema initialized")
	return nil
}
