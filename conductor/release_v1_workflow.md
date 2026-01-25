# SPECTRE v1.0.0 Release Workflow

This document outlines the step-by-step process to validate, build, and release SPECTRE v1.0.0.

---

## 🏗️ Track 1: Release Engineering (Versioning & Build)
**Goal:** Prepare the codebase for distribution.

- [ ] **Task 1.1: Version Stamp**
    - Create `internal/version/version.go` package.
    - Define `const Version = "v1.0.0"`.
    - Implement `spectre version` command in `internal/cli/version.go`.
- [ ] **Task 1.2: Documentation Polish**
    - Update `README.md`:
        - Verify "Installation" instructions.
        - Ensure "Quick Start" commands are copy-paste ready.
    - Update `docs/features.md` to reflect final v1 feature set.
- [ ] **Task 1.3: Cross-Compilation Script**
    - Create `scripts/build_release.sh` (or `.bat`) to generate:
        - `dist/spectre_v1.0.0_windows_amd64.exe`
        - `dist/spectre_v1.0.0_linux_amd64`
        - `dist/spectre_v1.0.0_darwin_arm64`

---

## 🧪 Track 2: End-to-End Verification (The "Golden Path")
**Goal:** Prove the system works in a production-like scenario before release.

**Target:** `scanme.nmap.org` (Authorized for scanning)

- [ ] **Task 2.1: Environment Reset**
    - Archive/Delete existing `spectre.db` and `evidence_storage/`.
    - Run `spectre init`.
- [ ] **Task 2.2: The Investigation Cycle (Manual Test)**
    1.  **Create Case:** `spectre case new --name "v1_Verification_Run"`
    2.  **Passive Collect:** `spectre collect --case <ID> --target scanme.nmap.org --scanners dns,whois,geo`
    3.  **Active Collect:** `spectre collect --case <ID> --target scanme.nmap.org --scanners ports,http,screenshot --active`
    4.  **Verification Point A:** Check `evidence_storage/` for JSON files and PNG screenshot.
    5.  **Analyze:** `spectre analyze --case <ID>` (Ensure it doesn't crash even if LLM is offline/mocked).
    6.  **Report:** `spectre report --case <ID> --format markdown`.
    7.  **Verification Point B:** Check generated Markdown report for "Scanme" details.
    8.  **Web UI:** `spectre server` -> `http://localhost:8080`.
    9.  **Verification Point C:** Open Graph, click node, check "Screenshot" link loads.

---

## 🚀 Track 3: Future Roadmap (Phase 4 - "Deep Operations")
**Goal:** Define the post-v1 trajectory.

- [ ] **Feature: Ghost Mode Hardening (Strict)**
    - *Concept:* A `--strict` flag that ensures the process *cannot* make network connections if the configured proxy is unreachable.
- [ ] **Feature: PDF Reporting Engine**
    - *Concept:* High-fidelity PDF generation using `gofpdf` with embedded graphs and branding.
- [ ] **Feature: Collaborative Backend**
    - *Concept:* Optional support for PostgreSQL to allow multiple analysts to work on the same case simultaneously.
- [ ] **Feature: Plugin Marketplace**
    - *Concept:* `spectre plugin install <url>` to pull community collectors (e.g., Shodan, Censys) easily.
