-- Migration: 002_add_analyses_table
-- Description: Create the analyses table for storing AI synthesis results.

CREATE TABLE IF NOT EXISTS analyses (
    id TEXT PRIMARY KEY,
    case_id TEXT NOT NULL,
    context_hash TEXT,
    findings JSON,
    risks JSON,
    connections JSON,
    next_steps JSON,
    confidence REAL,
    analyzed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (case_id) REFERENCES cases(id)
);
