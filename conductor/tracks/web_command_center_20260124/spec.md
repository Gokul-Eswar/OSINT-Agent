# Specification - Web Command Center

## Overview
While the CLI and TUI are powerful for speed, a Web UI is superior for complex graph exploration and high-level reporting. The Web Command Center will provide a persistent server to manage investigations across the team (or just for better local visualization).

## Requirements

### 1. API Server (Go)
- **Command:** `spectre server --port 8080`.
- **Core Endpoints:**
    - `GET /api/cases`: List all cases.
    - `GET /api/cases/{id}`: Get case details.
    - `GET /api/cases/{id}/graph`: Get nodes and edges for the graph.
    - `POST /api/cases/{id}/collect`: Trigger a collector via API.
- **Technology:** `net/http` with a standard router or lightweight framework.

### 2. Dashboard (Frontend)
- **Technology:** React + Bootstrap CSS.
- **Features:**
    - **Case Navigator:** Sidebar with active investigations.
    - **Intelligence Graph:** Interactive Cytoscape.js or Vis.js graph (real-time).
    - **Activity Feed:** Live stream of incoming evidence.
    - **Responsive Design:** Clean, modern "Dark Mode" aesthetic.

### 3. Distribution
- **Embedding:** The frontend assets MUST be compiled and embedded into the Go binary using `go:embed`.

## Success Criteria
- Running `spectre server` starts a local web server.
- Visiting `http://localhost:8080` shows a polished, functional dashboard.
- Users can view and explore the same data seen in the CLI/TUI.
