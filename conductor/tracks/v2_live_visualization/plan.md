# Implementation Plan: V2 Live Visualization

## Phase 1: Event Infrastructure (Go)
- [x] Add `OnRelationshipCreated` hook to `internal/storage/relationship_repo.go`.
- [x] Update `internal/server/server.go` to broadcast relationship events.
- [x] Ensure `OnEntityCreated` and `OnRelationshipCreated` are properly initialized in `server.Start`.

## Phase 2: Dashboard Core (Frontend)
- [x] Implement SSE listener in `web/index.html`.
- [x] Add `relationship_created` handler to the frontend event loop.
- [x] Implement incremental graph updates (adding to `vis.DataSet` instead of full re-render).

## Phase 3: Visual Polish
- [ ] Add pulsing animations for "Live" nodes.
- [ ] Implement a "Lead Feed" sidebar that shows the latest 5 discoveries in real-time.

## Phase 4: Validation
- [ ] Run a WHOIS collector and verify that nodes and edges appear in the browser without refreshing.
