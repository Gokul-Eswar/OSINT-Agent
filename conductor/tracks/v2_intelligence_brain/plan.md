# Implementation Plan: V2 Intelligence Brain

## Phase 1: Vector Foundation (Python/Analyzer)
- [x] Install `chromadb` and `sentence-transformers` in `analyzer/requirements.txt`.
- [x] Implement `analyzer/vector_store.py` for local indexing and querying.
- [x] Add `index_evidence` task to `analyzer/__main__.py`.

## Phase 2: Agent Tooling (Go/Internal)
- [x] Implement `search_evidence` tool in `internal/agent/tools.go`.
- [x] Implement `read_evidence` tool to allow the agent to read raw files.
- [x] Update `Engine.Execute` loop to support multi-turn reasoning and planning.

## Phase 3: Automated Ingestion
- [x] Update `internal/collector/registry.go` to automatically trigger re-indexing after any collector run.
- [ ] Implement an `IntelligenceLead` model in `internal/core` to track agent-derived hypotheses.

## Phase 4: Validation
- [ ] Create an integration test where the agent must find a specific "flag" hidden in a WHOIS record file.
