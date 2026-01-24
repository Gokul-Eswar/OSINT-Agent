# Implementation Plan - Active Reconnaissance (Probes)

## Phase 1: The Active Port Collector
- [x] **Task: Port Scanner**
    - Create `internal/collector/active/ports.go`.
    - Implement TCP connect scan for common ports.
- [x] **Task: Registration**
    - Register `ports` collector.
    - Add rate limit (2 ports/sec) to `configs/default.yaml`.

## Phase 2: The HTTP Stack Collector
- [x] **Task: HTTP Prober**
    - Create `internal/collector/active/http.go`.
    - Use `net/http` to fetch headers and title.
- [x] **Task: Registration**
    - Register `http` collector.

## Phase 3: Active Consent Enforcement
- [x] **Task: Global Active Flag**
    - Add `--active` flag to `rootCmd`.
    - Update `collector.Run` to fail if an active collector is called without the flag.

## Phase 4: Verification
- [x] **Task: Test**
    - Scan a safe target (e.g., `scanme.nmap.org` or `google.com`).
    - Verify ports and HTTP headers are saved as evidence.
