# ADR-0003: SQLite WAL and Append-Only Evidence Storage

## Status
ACCEPTED

## Context
Investigations require local-first persistence with reliable metadata queries and immutable raw evidence artifacts.

## Decision
Use a dual-store model.
- SQLite in WAL mode stores structured investigation metadata.
- Evidence files are written to the filesystem as append-only artifacts per case.

## Consequences
- Reliable local operation without external infrastructure.
- Fast metadata queries while preserving full-fidelity raw evidence.
- Requires consistency discipline between database records and file paths.

## Alternatives Considered
- Full document store only: rejected for weak relational graph querying.
- External database service: rejected to preserve offline and single-user setup simplicity.

## References
- Architecture storage section: ../../docs/ARCHITECTURE.md
- Tech stack storage section: ../tech-stack.md
- Code paths:
  - internal/storage
  - evidence_storage/
