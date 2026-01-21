# Implementation Plan - Initialize SPECTRE Core Foundation

## Phase 1: Project Scaffolding & CLI Setup
- [x] Task: Initialize Go Module and Directory Structure [commit: e882835]
    - [x] Initialize `go.mod`.
    - [x] Create directories: `cmd/spectre`, `internal/cli`, `internal/core`, `internal/storage`, `internal/config`.
    - [x] Create `.gitignore` for Go projects.
- [ ] Task: Implement Basic CLI with Cobra
    - [ ] Create `internal/cli/root.go` with the root command.
    - [ ] Create `cmd/spectre/main.go` to execute the root command.
    - [ ] Verify compilation with `go build ./cmd/spectre`.
- [ ] Task: Integrate Configuration with Viper
    - [ ] Create `internal/config/config.go` to handle loading.
    - [ ] Add flags to the root command (config path).
    - [ ] Create a default `configs/default.yaml`.
- [ ] Task: Setup Structured Logging
    - [ ] Implement `internal/logger` using `zerolog`.
    - [ ] Configure console output for dev and JSON for files.
    - [ ] Integrate logger into the `main` entry point.
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Project Scaffolding & CLI Setup' (Protocol in workflow.md)

## Phase 2: Storage Layer & Schema
- [ ] Task: Setup SQLite Infrastructure
    - [ ] Create `internal/storage/sqlite.go` for database connection.
    - [ ] Ensure `CGO_ENABLED=1` support (required for `go-sqlite3`).
- [ ] Task: Define Database Schema (Migrations)
    - [ ] Create SQL migration files or const strings for `cases`, `entities`, `relationships`, `evidence`, `analyses`.
    - [ ] Implement a migration runner in `internal/storage/schema.go` that runs on startup.
- [ ] Task: Implement Case Management (Proof of Concept)
    - [ ] Define `Case` struct in `internal/core/case.go`.
    - [ ] Implement `CreateCase` and `GetCase` in `internal/storage/case_repo.go`.
    - [ ] Write a unit test for `CreateCase` using an in-memory DB.
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Storage Layer & Schema' (Protocol in workflow.md)

## Phase 3: Integration & Wiring
- [ ] Task: Wire Storage to CLI
    - [ ] Add an `init` sub-command to `internal/cli/init.go` that initializes the DB.
    - [ ] Add a `new-case` sub-command to `internal/cli/case.go` that uses the storage layer.
- [ ] Task: End-to-End Verification
    - [ ] Run `spectre init` to create the DB.
    - [ ] Run `spectre new-case "test-case"` and verify it exists in the SQLite file.
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Integration & Wiring' (Protocol in workflow.md)
