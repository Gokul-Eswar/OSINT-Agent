# Implementation Plan - Ethics Guardian & Safety

## Phase 1: The Governor (Rate Limiting)
- [x] **Task: Core Limiter Engine**
    - Create `internal/ethics/limiter.go`.
    - Implement a Token Bucket manager that can hold different rates for different collector IDs.
- [x] **Task: Registry Integration**
    - Update `internal/collector/registry.go` to include rate limit definitions for each collector.
- [x] **Task: Runner Enforcement**
    - Wrap collector execution in `internal/collector/registry.go` with ethics checks.

## Phase 2: The Fence (Scope Control)
- [x] **Task: Scope Checker**
    - Create `internal/ethics/scope.go`.
    - Implement `IsAllowed(target string)` logic using glob matching.
- [x] **Task: Configuration**
    - Update `internal/config/config.go` to support `Ethics` block in YAML.
    - Add default blacklist to `configs/default.yaml`.

## Phase 3: CLI & UX
- [x] **Task: Safety Warnings**
    - Ensure the CLI prints a clear, professional message when a target is blocked or rate limiting is active.

## Phase 4: Verification
- [x] **Task: Stress Test**
    - Run 50 DNS queries rapidly and verify they are throttled (Verified with WHOIS test).
- [x] **Task: Blacklist Test**
    - Attempt `spectre collect dns nasa.gov` and verify it is blocked.
