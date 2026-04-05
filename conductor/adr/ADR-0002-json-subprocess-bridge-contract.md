# ADR-0002: JSON Subprocess Bridge Contract

## Status
ACCEPTED

## Context
The Go core must invoke Python analysis tasks without introducing a long-lived RPC dependency or service daemon.

## Decision
Use subprocess execution with JSON payload exchange.
- Go serializes analyzer requests to JSON and passes them as CLI input arguments.
- Python returns JSON responses through stdout.
- The bridge applies hard timeouts and structured error extraction.

## Consequences
- No mandatory background service to run analysis.
- Easy local/offline execution model.
- CLI payload size and stdout/stderr discipline become a contract requirement.

## Alternatives Considered
- gRPC/HTTP service: rejected to avoid process lifecycle complexity and daemon dependency.
- Shared library embedding: rejected for language/runtime coupling and portability risk.

## References
- Architecture bridge section: ../../architecture.md
- Code paths:
  - internal/analyzer/bridge.go
  - internal/analysis/engine.go
  - analyzer/__main__.py
