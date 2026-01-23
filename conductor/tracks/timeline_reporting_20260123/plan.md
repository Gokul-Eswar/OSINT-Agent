# Implementation Plan - Timeline & Reporting

**Goal:** Provide a chronological view of investigation events and generate exportable Markdown reports.

## Phase 1: Chronological Timeline
- [x] **Task: Timeline Engine**
    - [x] Create `internal/core/timeline.go` to define `TimelineEvent` (Timestamp, Type, Description).
    - [x] Implement `GetCaseTimeline(caseID string)` in `internal/storage` to aggregate discovery times from Entities and Evidence.
- [x] **Task: Timeline CLI**
    - [x] Create `internal/cli/timeline.go`.
    - [x] Implement `spectre timeline --case <id>`.
    - [x] Render events in a clean, sorted list.

## Phase 2: Report Generation
- [x] **Task: Markdown Generator**
    - [x] Create `internal/report/markdown.go`.
    - [x] Implement a generator that combines Case details, AI Analysis, Entity Graph summary, and the Timeline.
- [x] **Task: Report CLI**
    - [x] Create `internal/cli/report.go`.
    - [x] Implement `spectre report --case <id>`.
    - [x] Output to stdout or save to `evidence_storage/<case_id>/summary.md`.

## Phase 3: Verification
- [x] **Task: End-to-End Investigation Story**
    - [x] Run collectors.
    - [x] Run AI Analysis.
    - [x] Generate Timeline.
    - [x] Generate Report and verify all findings are present.
