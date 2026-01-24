# SPECTRE Phase 3: Active Operations & Web Intelligence

This phase shifts SPECTRE from a passive CLI tool into an active, interactive intelligence platform.

## Strategic Goals
1.  **Active Reconnaissance:** safely introduce active scanning capabilities (Ports, HTTP) to validate passive findings.
2.  **Web Command Center:** Replace static reports with a persistent web server (`spectre server`) for real-time case management and graph exploration.
3.  **Geo-Intelligence:** Enrich infrastructure data with physical location to map threat landscapes.
4.  **Extensibility:** Formalize the plugin system to support external tools.

---

## Track 9: Active Reconnaissance (The "Probe")
**Goal:** Validate infrastructure existence and identify running services.
**Constraint:** MUST require explicit user consent via `--active` flag.

### 9.1 Port Scanner
*   **Why:** Knowing open ports (80, 443, 22, 3389) identifies server function.
*   **Tasks:**
    *   Implement `internal/collector/active/ports.go`.
    *   Scan top 20 common ports by default.
    *   Strict timeout and rate limiting.

### 9.2 HTTP Probe
*   **Why:** Headers and status codes reveal technology stacks (Nginx, IIS, PHP).
*   **Tasks:**
    *   Implement `internal/collector/active/http.go`.
    *   Fetch Title, Server header, and StatusCode.
    *   Screenshot capability (optional, via headless chrome if available).

---

## Track 10: Geo-Intelligence (The "Map")
**Goal:** Visualize the physical footprint of the target.

### 10.1 GeoIP Enrichment
*   **Why:** Knowing an IP is in "Pyongyang" vs "Ashburn" changes the threat profile.
*   **Tasks:**
    *   Integrate a GeoIP provider (MaxMind GeoLite2 or a privacy-respecting API).
    *   Enrich `IP` entities with `lat`, `long`, `country`, `city`.

### 10.2 Map Visualization
*   **Why:** Analysts think spatially.
*   **Tasks:**
    *   Update `spectre visualize` (or the new Web UI) to plot IPs on a world map.

---

## Track 11: Web Command Center (The "Console")
**Goal:** A persistent, interactive UI to manage investigations.

### 11.1 The API Server
*   **Why:** Decouple the UI from the CLI.
*   **Tasks:**
    *   Implement `spectre server` using `net/http` or `Gin/Echo`.
    *   Endpoints: `/api/cases`, `/api/graph`, `/api/query`.

### 11.2 The Frontend
*   **Why:** HTML reports are static. We need dynamic filtering and exploration.
*   **Tasks:**
    *   Build a simple SPA (Single Page App) or server-rendered HTML.
    *   Embed the assets into the Go binary (`embed`).
    *   Features: Search, Graph View, Timeline View, Map View.

---

## Track 12: External Plugin System (The "Ecosystem")
**Goal:** Allow users to write collectors in any language.

### 12.1 External Runner
*   **Why:** We can't write every collector in Go. Community tools are often Python/Bash.
*   **Tasks:**
    *   Define a JSON-over-Stdin protocol for plugins.
    *   Implement a loader that finds executables in `~/.spectre/plugins/`.
    *   Map external output to internal `Entity` and `Evidence` types.

---

## Track 13: Interactive TUI Console (The "Terminal")
**Goal:** A rich, interactive terminal interface using Bubble Tea.

### 13.1 The Dashboard
*   **Why:** Running individual commands is slow. A persistent TUI allows rapid exploration.
*   **Tasks:**
    *   Integrate `github.com/charmbracelet/bubbletea`.
    *   Implement `spectre console` command.
    *   Views: Case List, Entity Table, Log Stream.

### 13.2 Interactive Runner
*   **Why:** Trigger collectors without remembering flags.
*   **Tasks:**
    *   Form-based input for running collectors.
    *   Real-time progress bars for collection tasks.

---

## Execution Order
1.  **Track 10 (GeoIP)** is low hanging fruit and high value.
2.  **Track 13 (TUI)** provides the "cool factor" and usability.
3.  **Track 9 (Active)** adds a new dimension of data.
4.  **Track 11 (Web UI)** is a major architectural addition.
5.  **Track 12 (Plugins)** ensures long-term scalability.
