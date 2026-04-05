# ADR-0004: Evidence Hashing and Auditability

## Status
ACCEPTED

## Context
Investigation outputs must be traceable from synthesized intelligence back to raw evidence with tamper-detection support.

## Decision
Hash all raw evidence and persist hash references with evidence metadata.
- Evidence files receive SHA-256 hashes at collection time.
- Relationships and analysis artifacts maintain linkage to evidence IDs.

## Consequences
- Stronger forensic traceability and reproducibility.
- Hash checks can detect content drift.
- Slight overhead on collection write path.

## Alternatives Considered
- No hashing: rejected due to weak audit guarantees.
- Optional hashing: rejected because partial coverage weakens trust in outputs.

## References
- Architecture storage and auditability notes: ../../architecture.md
- Code paths:
  - internal/collector/external.go
  - internal/core/evidence.go
  - internal/storage
