# SPECTRE Phase 2: Feature Expansion Roadmap

This document outlines the engineering path to transform SPECTRE from a functional prototype into a capable intelligence tool.

## Strategic Goals
1.  **Data Depth:** Expand collection beyond DNS to include ownership (WHOIS) and code footprint (GitHub).
2.  **Narrative:** Enable analysts to view investigations chronologically (Timeline) and export results (Reports).
3.  **Safety:** Professionalize the tool with strict rate limiting and ethics checks.

---

## Track 6: Expanded Collection Capabilities
**Goal:** Integrate rich data sources that require authentication or parsing complex text.

### 6.1 Configuration Management
*   **Why:** GitHub and other APIs require tokens. We need a secure way to manage them.
*   **Tasks:**
    *   Implement `spectre config set <key> <value>` (e.g., `spectre config set keys.github gh_xxxx`).
    *   Ensure config file permissions are secure (0600).
    *   Refactor `internal/config` to expose typed accessors for secrets.

### 6.2 WHOIS Collector (Passive)
*   **Why:** Identifying domain ownership, registration dates, and abuse contacts is critical for attribution.
*   **Tasks:**
    *   Integrate `github.com/likexian/whois` library.
    *   **Parsing:** Extract "Registrant Email", "Creation Date", "Registrar".
    *   **Graphing:**
        *   Create `Person` or `Email` entities from registrant info.
        *   Link `Domain` -> `owns` -> `Person`.

### 6.3 GitHub Collector (API-Based)
*   **Why:** Threat actors often leak credentials or host malware on GitHub.
*   **Tasks:**
    *   Implement GitHub REST API client using `net/http`.
    *   **Search:** `spectre collect github <keyword/user>`.
    *   **Ingestion:**
        *   `User` entity (GitHub username).
        *   `Repository` entity (URL).
        *   Relationship: `User` -> `owns` -> `Repository`.

---

## Track 7: Timeline & Reporting
**Goal:** Transform the static "Graph" into a chronological story.

### 7.1 Timeline Engine
*   **Why:** Investigations are about *when* something happened, not just *what* exists.
*   **Tasks:**
    *   **Backend:** Create `internal/core/timeline.go` to aggregate `Evidence.CollectedAt` and `Entity.DiscoveredAt`.
    *   **CLI:** Implement `spectre timeline --case <id>`.
    *   **Visualization:** Render a CLI tree view (like `git log --graph`).
    *   **Python Bridge:** Optional extension to generate a Vis.js Timeline HTML.

### 7.2 Markdown Reporting
*   **Why:** Analysts need to deliver a final document, not a database file.
*   **Tasks:**
    *   **Template:** Create a robust Markdown template (`templates/report.md.j2`).
    *   **Generator:** Build `internal/report/generator.go`.
    *   **Content:**
        *   Executive Summary (from AI Analysis).
        *   Key Findings (High confidence entities).
        *   Investigation Timeline.
        *   Visual Graph Screenshot (placeholder or link).

---

## Track 8: Ethics Guardian & Safety
**Goal:** Ensure the tool is safe for professional use and prevents accidental "attacking".

### 8.1 Rate Limiting (The "Governor")
*   **Why:** Prevent IP bans and API quota exhaustion.
*   **Tasks:**
    *   Implement `internal/ethics/limiter.go` using Token Buckets.
    *   Enforce per-collector limits (e.g., DNS: 10/sec, WHOIS: 1/sec).

### 8.2 Scope Control (The "Fence")
*   **Why:** Prevent accidental collection against out-of-scope targets (e.g., `.gov` or explicit blacklists).
*   **Tasks:**
    *   Implement `spectre run --passive-only` flag logic.
    *   Create a blacklist configuration (`scope.blacklist: ["*.gov", "internal.corp"]`).
    *   Add `CanCollect(target)` check in the Collector Runner.

---

## Execution Order
1.  **Track 6 (Collectors)** is the priority to make the tool useful.
2.  **Track 7 (Timeline)** adds value to the data collected.
3.  **Track 8 (Ethics)** prepares it for release/distribution.
