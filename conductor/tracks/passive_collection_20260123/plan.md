# Implementation Plan - Passive Collection Framework

**Goal:** Implement a modular system for OSINT collectors and build a DNS plugin that automatically populates the graph.

## Phase 1: Evidence & Collector Core
- [x] **Task: Define Evidence & Interface**
    - [x] Create `internal/core/evidence.go` (ID, CaseID, Collector, FilePath, Metadata).
    - [x] Define `Collector` interface in `internal/core/collector.go`.
- [x] **Task: Evidence Storage**
    - [x] Create `internal/storage/evidence_repo.go` (Create, ListByCase).
    - [x] Update schema (verified in Core Foundation).

## Phase 2: Plugin Registry & DNS Collector
- [x] **Task: Collector Registry**
    - [x] Create `internal/collector/registry.go` to register and retrieve plugins.
- [x] **Task: DNS Plugin**
    - [x] Create `internal/collector/dns/dns.go`.
    - [x] Support A, MX, and NS record lookups.
    - [x] Standardize output to `Evidence` objects.

## Phase 3: CLI & Auto-Ingestion
- [x] **Task: Collection CLI**
    - [x] Create `internal/cli/collect.go`.
    - [x] Implement `spectre collect <collector> <target> --case <id>`.
- [x] **Task: Graph Ingestion Logic**
    - [x] Build a service to parse `Evidence` and automatically create `Entities` and `Relationships`.
    - [x] Example: DNS A-record -> "domain" -(resolves_to)-> "ip".

## Phase 4: Verification
- [x] **Task: End-to-End Collection Test**
    - [x] Run `spectre collect dns google.com --case <id>`.
    - [x] Verify that both entities (domain/ip) and the relationship are created automatically.
