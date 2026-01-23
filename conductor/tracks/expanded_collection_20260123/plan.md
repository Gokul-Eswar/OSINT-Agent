# Implementation Plan - Expanded Collection (WHOIS, GitHub & Config)

**Goal:** Implement secure API key management and add advanced collectors for WHOIS and GitHub to enrich the entity graph.

## Phase 1: API Key & Config Management
- [x] **Task: Config CLI Commands**
    - [x] Implement `spectre config set <key> <value>` in `internal/cli/config.go`.
    - [x] Support nested keys (e.g., `keys.github`).
    - [x] Verify keys are saved correctly in `configs/default.yaml` or user-defined path.
- [x] **Task: Secure Accessors**
    - [x] Update `internal/config/config.go` to provide a helper for retrieving secrets safely.

## Phase 2: WHOIS Collector
- [x] **Task: WHOIS Implementation**
    - [x] Create `internal/collector/whois/whois.go`.
    - [x] Use `github.com/likexian/whois` for raw data retrieval.
    - [x] Implement a parser to extract: Registrant Email, Creation Date, Registrar.
- [x] **Task: WHOIS Ingestion**
    - [x] Update `internal/storage/ingestion.go` to handle WHOIS evidence.
    - [x] Logic: Create `Email` or `Person` entities and link them to the `Domain` via `owns` or `registered_by`.

## Phase 3: GitHub Collector
- [x] **Task: GitHub Client**
    - [x] Create `internal/collector/github/github.go`.
    - [x] Implement authenticated search for users and repositories.
- [x] **Task: GitHub Ingestion**
    - [x] Update `internal/storage/ingestion.go` to handle GitHub evidence.
    - [x] Logic: Create `User` (username) and `Repository` (URL) entities. Link `User` -> `owns` -> `Repository`.

## Phase 4: Verification
- [x] **Task: Multi-Source Investigation Test**
    - [x] Create a case.
    - [x] Add a domain.
    - [x] Run `collect dns`, `collect whois`, and `collect github`.
    - [x] Verify the graph shows the domain, its IPs, its owner (email), and related GitHub assets.
