# Specification - Active Reconnaissance (Probes)

## Overview
Passive intelligence provides the *identity* of infrastructure; active probes provide the *state*. SPECTRE will now support basic service discovery.

## Requirements

### 1. Port Scanner
- **Target:** IP entities.
- **Scope:** Top 20 common ports (80, 443, 22, 21, 25, 53, 445, 3389, 8080, etc.).
- **Mechanism:** SYN or Connect scan (TCP).
- **Flag:** MUST require `--active` to run.

### 2. HTTP Probe
- **Target:** Domain or IP.
- **Data Points:**
    - Server header (Nginx, Apache, IIS).
    - Page Title.
    - Status Code (200, 403, 404).
    - Powered-By headers (X-Powered-By).

## Constraints
- **Consent:** Any active collector MUST check a global `AllowActive` flag.
- **Safety:** Stricter rate limiting than passive collectors.

## Success Criteria
- Running `spectre collect ports <IP> --active` identifies open ports.
- Running `spectre collect http <Domain> --active` identifies the web server stack.
