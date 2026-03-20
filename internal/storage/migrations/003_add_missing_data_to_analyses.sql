-- Migration: 003_add_missing_data_to_analyses
-- Description: Add missing_data and suggested_collectors columns to analyses table

ALTER TABLE analyses ADD COLUMN missing_data JSON DEFAULT '[]';
ALTER TABLE analyses ADD COLUMN suggested_collectors JSON DEFAULT '[]';
