# ADR-0009: Static HTML Visualization Output

## Status
ACCEPTED

## Context
Investigators need interactive graph visualization without requiring a persistent web backend or service deployment.

## Decision
Generate static HTML visualizations from case exports.
- Go exports case graph payloads.
- Python visualization logic renders interactive HTML artifacts.
- Output is opened locally in the user environment.

## Consequences
- Simple, local-first visualization workflow.
- Easy sharing of generated artifacts.
- Very large graphs may require additional filtering for browser performance.

## Alternatives Considered
- Embedded web server: rejected to avoid runtime service management overhead.
- Native GUI dependency: rejected to keep CLI-first portability.

## References
- Architecture visualization flow: ../../docs/ARCHITECTURE.md
- Code paths:
  - internal/analysis/context.go
  - analyzer/graph_viz.py
  - internal/cli/visualize.go
