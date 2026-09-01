# ADR-0010: YAML and Viper Configuration Hierarchy

## Status
ACCEPTED

## Context
SPECTRE must provide configurable behavior across collectors, ethics policies, storage, and LLM settings while remaining easy to operate locally.

## Decision
Use YAML defaults with Viper-based override layering.
- Base defaults live in repository configuration files.
- Runtime values can be overridden by environment variables and CLI flags.
- Core modules read from a shared configuration substrate.

## Consequences
- Predictable configuration model across commands and modules.
- Supports environment-specific deployment without source edits.
- Requires clear documentation of precedence and sensitive-value handling.

## Alternatives Considered
- Hardcoded values: rejected due to poor operational flexibility.
- Database-stored runtime config only: rejected to keep bootstrap simple and file-driven.

## References
- Tech stack configuration section: ../tech-stack.md
- Code paths:
  - configs/default.yaml
  - internal/config/config.go
  - internal/cli/config.go
