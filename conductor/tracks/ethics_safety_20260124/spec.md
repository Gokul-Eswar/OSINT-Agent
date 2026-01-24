# Specification - Ethics Guardian & Safety

## Overview
SPECTRE must be a responsible tool. To prevent abuse, accidental "attacking" behavior, and API bans, we need a safety layer.

## Requirements

### 1. Rate Limiting (The "Governor")
- **Mechanism:** Token Bucket algorithm.
- **Granularity:** Per-collector limits.
- **Defaults:**
    - DNS: 10 requests/sec.
    - WHOIS: 1 request/sec (WHOIS servers are very sensitive).
    - GitHub: 2 requests/sec (authenticated).

### 2. Scope Control (The "Fence")
- **Blacklist:** Prevent collection on specific domains (e.g., `.gov`, `.mil`).
- **Target Verification:** A central check before any collector executes.
- **Config:** Managed via `configs/default.yaml`.

## Success Criteria
- Collectors respect defined rate limits.
- Attempting to collect on a blacklisted domain returns a safety error.
- The system prevents concurrent execution beyond safe thresholds.
