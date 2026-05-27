# ADR-0001: Hybrid Go and Python Architecture

## Status
ACCEPTED

## Context
SPECTRE needs a fast, portable CLI core for orchestration, policy enforcement, and storage, while also supporting rapid iteration for AI synthesis and visualization.

## Decision
Use Go for the system core and Python for the intelligence sidecar.
- Go owns CLI workflows, collector orchestration, ethics enforcement, persistence, and graph mutation.
- Python owns LLM synthesis, specialized analysis routines, and visualization generation.
- Cross-language communication happens through a strict request/response bridge.

## Consequences
- Strong operational reliability for collection and storage paths.
- Faster experimentation in AI and analysis logic.
- Added bridge complexity and contract management across languages.

## Alternatives Considered
- Go only: rejected due to slower iteration for ML/AI ecosystem features.
- Python only: rejected due to weaker single-binary distribution and concurrency ergonomics for core orchestration.

## References
- Architecture: ../../docs/ARCHITECTURE.md
- Tech stack: ../tech-stack.md
- Code paths:
  - internal/cli
  - internal/collector
  - internal/storage
  - analyzer/
