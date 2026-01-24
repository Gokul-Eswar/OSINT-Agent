# Implementation Plan - Web Command Center

## Phase 1: API Foundation
- [x] **Task: Server Scaffolding**
    - Create `internal/server/server.go`.
    - Implement the basic HTTP server and router.
- [x] **Task: API Endpoints**
    - Implement `/api/cases` and `/api/cases/{id}/graph`.
    - Use `internal/storage` and `internal/analysis` to fetch data.

## Phase 2: Frontend Prototype
- [x] **Task: Basic Layout**
    - Create `web/` directory.
    - Implement a single-page React app (SPA).
- [x] **Task: Graph Integration**
    - Integrate `vis-network` (same as the HTML reports but dynamic).
- [x] **Task: Real-time Updates**
    - (Optional) Use simple polling or WebSockets to show new evidence.

## Phase 3: Embedding & CLI
- [x] **Task: Asset Embedding**
    - Use `go:embed` to serve the `web/` folder.
- [x] **Task: Server Command**
    - Create `internal/cli/server.go`.
    - Implement `spectre server`.

## Phase 4: Verification
- [x] **Task: Demo**
    - Launch server, browse to localhost, and explore a case.
