# ADR-0006: Collector Registry and Plugin Contract

## Status
ACCEPTED

## Context
SPECTRE must support built-in and external collectors with uniform execution, policy checks, and evidence persistence.

## Decision
Use a registry-based collector abstraction.
- All collectors implement a shared Collector contract.
- External plugins are discovered from plugin manifests and wrapped as collectors.
- Registry run paths enforce common policy and persistence flow.

## Consequences
- Consistent execution semantics across native and external collectors.
- Easier extensibility and testability.
- Plugin manifest correctness becomes critical to runtime behavior.

## Alternatives Considered
- Ad hoc command-specific integrations: rejected due to inconsistent policy handling.
- Separate plugin runtime subsystem: rejected to avoid divergence from core collector lifecycle.

## References
- Plugin docs: ../../docs/plugins.md
- Code paths:
  - internal/core/collector.go
  - internal/collector/registry.go
  - internal/collector/external.go
