# ADR-0008: Local-First LLM Strategy with Optional Cloud Providers

## Status
ACCEPTED

## Context
Operators need analysis capabilities that can run offline, while still allowing optional cloud model integrations when explicitly configured.

## Decision
Default to local model routing and allow cloud providers through explicit configuration.
- Local-first execution is the default operational mode.
- Cloud providers are optional and opt-in.
- Provider details are read from runtime configuration and passed through bridge requests.

## Consequences
- Better privacy and offline capability by default.
- Flexible model/provider selection when needed.
- Requires clear provider configuration and secret management guidance.

## Alternatives Considered
- Cloud-only inference: rejected due to privacy, connectivity, and cost constraints.
- Local-only hard restriction: rejected to retain flexibility for teams with cloud LLM requirements.

## References
- Tech stack AI synthesis section: ../tech-stack.md
- LLM docs: ../../docs/llm.md
- Code paths:
  - internal/analysis/engine.go
  - internal/analyzer/bridge.go
  - analyzer/llm.py
