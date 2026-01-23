# Implementation Plan - Core Data Engine (Graph Models & CLI)

**Goal:** Implement the fundamental data structures for Entities and Relationships, enabling manual graph construction via CLI.

## Phase 1: Domain Models & Storage
- [x] **Task: Define Core Structs**
    - [x] Create `internal/core/entity.go` (ID, Type, Value, Metadata).
    - [x] Create `internal/core/relationship.go` (SourceID, TargetID, RelType, Directed).
    - [x] Add JSON tags for serialization.
- [x] **Task: Implement Storage Repositories**
    - [x] Update `internal/storage/sqlite.go` (verified global DB usage).
    - [x] Create `internal/storage/entity_repo.go` (Create, Get, ListByCase).
    - [x] Create `internal/storage/relationship_repo.go` (Create, Get, ListByEntity).
    - [x] Add unit tests for repositories.

## Phase 2: CLI Commands for Manual Entry
- [x] **Task: Entity Management Commands**
    - [x] Create `internal/cli/entity.go`.
    - [x] Implement `spectre entity add <type> <value> --case <id>`.
    - [x] Implement `spectre entity list --case <id>`.
- [x] **Task: Relationship Management Commands**
    - [x] Create `internal/cli/link.go`.
    - [x] Implement `spectre link add <source_val> <target_val> --type <rel_type> --case <id>`.
    - [x] Implement `spectre link list --case <id>`.

## Phase 3: Integration Verification
- [x] **Task: Manual Graph Construction Test**
    - [x] Create a test script (or manual verification steps) to:
        1. Create a case "Operation-X".
        2. Add entity (IP) "192.168.1.1".
        3. Add entity (Domain) "malicious.com".
        4. Link them with "resolves_to".
        5. Verify the data exists in SQLite.
- [x] **Task: Conductor Checkpoint**
    - [x] Verify all tests pass.
    - [x] Ensure code style compliance.
