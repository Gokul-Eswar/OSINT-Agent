# ADR-0007: Ethics Enforcement Layer

## Status
ACCEPTED

## Context
Collection operations can perform active probing and target sensitive infrastructure. The platform needs hard safety gates independent of collector implementation quality.

## Decision
Enforce ethics checks in core orchestration before collector execution.
- Active collectors require explicit user consent.
- Scope controls validate target eligibility.
- Rate limiting is applied per collector.

## Consequences
- Centralized, auditable safety posture.
- Reduced risk from unsafe plugin behavior.
- Some scans may be blocked unless users opt in and configure scope correctly.

## Alternatives Considered
- Collector self-enforcement only: rejected as unenforceable and inconsistent.
- Post-execution auditing only: rejected because it does not prevent unsafe execution.

## References
- Architecture ethics section: ../../docs/ARCHITECTURE.md
- Code paths:
  - internal/ethics/limiter.go
  - internal/ethics/scope.go
  - internal/collector/registry.go
