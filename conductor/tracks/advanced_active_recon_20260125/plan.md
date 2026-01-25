# Implementation Plan - Advanced Active Reconnaissance

## Phase 1: Social Intelligence (Sherlock)
- [x] **Task: Refactor Social Collector**
    - Modify `internal/collector/active/social.go` to use a configurable `SiteMap` or `template` system to allow testing without hitting real social networks.
    - Ensure `netclient` is injected or mockable.
- [x] **Task: Social Tests**
    - Create `internal/collector/active/social_test.go`.
    - Verify it correctly identifies "found" vs "not found" profiles using a local test server.

## Phase 2: Visual Evidence (Screenshot)
- [x] **Task: Screenshot Tests**
    - Create `internal/collector/active/screenshot_test.go`.
    - Use a simple HTTP test server as the target.
    - **Note:** This test might require a specialized environment (Docker/CI) for `chromedp`. We will implement a skip-if-headless-missing check.
- [x] **Task: Integration Check**
    - Verify `screenshot` collector works with the `web` dashboard (path correctness).
    - *Action:* Added `/evidence/` handler to `internal/server/server.go`.

## Phase 3: Documentation
- [x] **Task: Docs Update**
    - Update `docs/features.md` to include Social and Screenshot capabilities.
