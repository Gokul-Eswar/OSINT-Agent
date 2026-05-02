# Implementation Plan: Codebase Refinements

## Objective
Enhance the maintainability, performance, and robustness of the SPECTRE platform by addressing key areas identified during the codebase analysis. This plan outlines a phased approach to implementing these refinements.

## Key Files & Context
*   **Go Core:** `internal/storage/entity_repo.go`, `internal/storage/ingestion.go`, `internal/collector/registry.go`, `internal/agent/tools.go`, `internal/analyzer/bridge.go`, `internal/agent/engine.go`, `internal/cli/chat.go`
*   **Python Analyzer:** `analyzer/llm.py`, `analyzer/vector_store.py`

## Implementation Steps

### Phase 1: Architectural DRYness (Storage & Ingestion)
1.  **Introduce `EnsureEntity` helper:**
    *   Implement `EnsureEntity(caseID, entityType, entityValue, source string) (*core.Entity, error)` in `internal/storage/entity_repo.go` (or `database.go`).
    *   This function will check if an entity exists by value. If not, it creates and returns it; otherwise, it returns the existing entity.
2.  **Refactor Ingestion Logic:**
    *   Update all `ingest*` functions in `internal/storage/ingestion.go` (e.g., `ingestDNS`, `ingestWHOIS`) to use `EnsureEntity` instead of the manual `GetEntityByValue` check followed by `CreateEntity`.

### Phase 2: Performance (Vector Store)
1.  **Batch Vector Re-indexing:**
    *   In `internal/collector/registry.go`, the `RunAndSave` function currently calls `analyzer.IndexEvidence(caseID)` in a goroutine for *every* evidence record processed.
    *   Modify `RunAndSave` to accumulate the saved evidence and call `analyzer.IndexEvidence` only *once* after the loop completes.

### Phase 3: Core Feature Completion (Agent Pipeline)
1.  **Implement `read_evidence` Tool:**
    *   In `internal/agent/tools.go`, implement the logic for the `read_evidence` tool.
    *   Construct the file path (`evidence_storage/<case_id>/<filename>`), verify existence, read contents, and return the raw text (truncate if excessively large).
2.  **Improve Bridge Robustness:**
    *   Enhance `RunPythonTask` in `internal/analyzer/bridge.go` to better handle malformed JSON responses or unexpected stderr outputs, ensuring graceful degradation.

### Phase 4: Testing & UX Improvements
1.  **Mockable Analyzer Interface (Optional Testing Polish):**
    *   Extract `RunPythonTask` logic into an interface in `internal/analyzer/bridge.go` to allow injecting a mock analyzer during `agent.Engine` unit testing.
2.  **Real-time Agent Feedback:**
    *   Improve the "Agent is thinking..." state in `internal/cli/chat.go` and `internal/agent/engine.go` to provide real-time updates when tools are executed (e.g., "Agent decided to run WHOIS on example.com...").

## Verification & Testing
*   **Unit Tests:** Ensure all existing tests pass (`go test ./...`). Add a unit test for the new `EnsureEntity` helper.
*   **Integration Tests:** Verify that running a collector (e.g., `spectre collect dns example.com`) correctly creates entities and relationships using the refactored ingestion logic.
*   **Performance Check:** Run a high-volume collector and monitor the number of Python subprocesses spawned to confirm the vector store indexing optimization.

## Alternatives Considered
*   *Complex Event Bus for Indexing:* Considered using a publish-subscribe system for triggering re-indexing, but simply batching the call at the end of `RunAndSave` is far simpler and meets the immediate performance needs without adding overhead.

## Migration & Rollback
These changes are primarily internal refactorings and non-breaking feature additions.
*   **Rollback:** Standard `git checkout` to revert. No database schema migrations are introduced.