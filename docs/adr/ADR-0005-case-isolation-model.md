# ADR-0005: Case Isolation Model

## Status
ACCEPTED

## Context
Multiple investigations may run on the same installation; data crossover would compromise evidence quality and operator trust.

## Decision
Enforce case-scoped isolation for entities, relationships, analysis, and evidence paths.
- Structured records include case identifiers.
- Raw artifacts are stored under case-specific directories.
- Analysis context and cache keys are resolved per case.

## Consequences
- Safer multi-case workflows and clearer evidence boundaries.
- Cleaner export and archival operations.
- Requires every ingestion and query path to preserve case scoping.

## Alternatives Considered
- Global pooled graph: rejected due to contamination risk across investigations.
- Separate database per case: rejected to avoid heavy lifecycle management overhead.

## References
- Architecture case/evidence storage flow: ../../docs/ARCHITECTURE.md
- Code paths:
  - internal/core/case.go
  - internal/storage
  - internal/analysis/engine.go
