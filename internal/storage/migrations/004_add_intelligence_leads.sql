-- Migration: 004_add_intelligence_leads
-- Description: Create the intelligence_leads table to store hypotheses.

CREATE TABLE IF NOT EXISTS intelligence_leads (
    id TEXT PRIMARY KEY,
    case_id TEXT NOT NULL,
    hypothesis TEXT NOT NULL,
    confidence REAL DEFAULT 0.5,
    evidence_ids JSON NOT NULL DEFAULT '[]',
    status TEXT DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (case_id) REFERENCES cases(id)
);

CREATE INDEX IF NOT EXISTS idx_leads_case_id ON intelligence_leads(case_id);
