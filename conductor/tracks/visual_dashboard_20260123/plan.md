# Implementation Plan - Visual Intelligence Dashboard (Hybrid)

**Goal:** Generate interactive, offline HTML reports using Python's `pyvis` and `networkx`, orchestrated via the Go CLI.

## Phase 1: Python Visualization Module
- [x] **Task: Setup Analyzer Module** (Done in refactor)
- [ ] **Task: Implement Graph Viz Logic**
    - Create `analyzer/graph_viz.py`.
    - Use `networkx` to build the graph from JSON input.
    - Use `pyvis` to generate an interactive HTML file.
    - Color-code nodes by entity type.

## Phase 2: Go Bridge Extension
- [ ] **Task: Update Bridge**
    - Ensure `internal/analyzer/bridge.go` can handle the `visualize` task.
- [ ] **Task: Case Serializer**
    - Implement `analysis.ExportCaseForViz(caseID)` to prepare the JSON for Python.

## Phase 3: CLI Integration
- [ ] **Task: Visualize Command**
    - Create `internal/cli/visualize.go`.
    - Implement `spectre visualize --case <id>`.
    - It should: Call Go Serializer -> Call Python Bridge -> Open generated HTML.

## Phase 4: Verification
- [ ] **Task: Final v1 Demo**
    - Run `spectre visualize` on the "Graph-Test" case.
    - Verify the graph shows google.com connected to its IPs in a professional-grade dashboard.
