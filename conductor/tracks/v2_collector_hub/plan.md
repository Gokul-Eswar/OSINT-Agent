# Implementation Plan: V2 Collector Hub

## Phase 1: Python SDK Foundation
- [x] Create `lib/python/spectre_sdk.py` containing the `BaseCollector` and utility functions.
- [x] Create a sample plugin in `plugins/leak_checker` that uses the new SDK.
- [x] Update `internal/collector/external.go` to ensure it passes all necessary context to Python plugins.

## Phase 2: Automated Dorking Engine
- [x] Add `generate_dorks` capability to `analyzer/llm.py`.
- [x] Implement `generate_dorks` tool in `internal/agent/tools.go`.
- [x] Create a library of dork templates (Cloud buckets, Git leaks, DB dumps).

## Phase 3: Validation
- [x] Run the `leak_checker` plugin and verify output.
- [x] Ask the agent to "Find potential data leaks for example.com" and verify it generates valid dorks.
