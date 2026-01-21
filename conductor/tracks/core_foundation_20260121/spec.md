# Specification: Initialize SPECTRE Core Foundation

## 1. Goal
To establish the foundational software architecture for SPECTRE, enabling future development of collectors, analysis modules, and the user interface. This involves setting up the Go module, organizing the directory structure, implementing the core data models (SQLite), and creating the CLI entry point.

## 2. Core Requirements

### 2.1 Project Scaffolding
- Initialize a new Go module (`github.com/spectre/spectre` or similar placeholder).
- Create the standard Go project layout:
    - `cmd/spectre/`: Application entry point.
    - `internal/`: Private application code.
    - `pkg/`: Library code (if necessary, though `internal` is preferred for this phase).
    - `configs/`: Configuration files.

### 2.2 CLI Infrastructure
- Implement the root command using `cobra`.
- Integrate `viper` for configuration management (reading `config.yaml`, environment variables).
- Define global flags (e.g., `--config`, `--verbose`).

### 2.3 Storage Layer
- Implement the SQLite database connection using `mattn/go-sqlite3`.
- Define and apply the initial database schema (migrations):
    - `cases`: Track investigations.
    - `entities`: Store discovering items (domains, emails, etc.).
    - `relationships`: Link entities.
    - `evidence`: Log provenance.
    - `analyses`: Store AI results.
- Implement basic CRUD operations for the `Case` model to verify the storage layer.

### 2.4 Logging & Error Handling
- Set up structured logging using `zerolog`.
- Ensure logs are written to both `stderr` (human-readable) and a file (JSON).

## 3. Non-Functional Requirements
- **Modularity:** Ensure clear separation between CLI (presentation), Core (business logic), and Storage.
- **Testability:** The storage layer must be testable (using in-memory SQLite for unit tests).
- **Compliance:** Adhere to the `go.md` style guide (idiomatic Go, error handling, formatting).

## 4. Out of Scope
- Implementing actual collectors (DNS, WHOIS, etc.).
- The Python/AI bridge.
- The web-based graph visualization.
- Complete feature implementation beyond "Init" and basic "Case" creation.
