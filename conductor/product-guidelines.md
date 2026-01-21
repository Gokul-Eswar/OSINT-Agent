# Product Guidelines — SPECTRE

## Tone & Prose
- **Professional & Forensic:** All documentation, CLI output, and reports must be clinical, precise, and objective. Avoid slang or overly informal language.
- **Evidentiary Focus:** Prioritize clarity on data sources, timestamps, and confidence scores.

## Visual & Interface Design
- **Structured & Tabular:** Use tables and bordered blocks in the CLI to make data easy to scan and cross-reference.
- **Rich Terminal Experience:** employ bold typography, subtle color-coding for entity types (e.g., Blue for Domains, Green for Emails), and clear progress indicators.
- **Machine-Readable (JSON-First):** Ensure all core data outputs can be emitted as structured JSON for automation and piping.

## Engineering Principles
- **Modular Hybrid Architecture:** Maintain a strict boundary between the Go-based Core (orchestration, storage) and the Python-based Intelligence layer (AI, visualization).
- **Zero-Trust Dependencies:** Minimize external libraries. Every third-party dependency must be audited to mitigate supply-chain risks.
- **Documentation as Code:** Maintain comprehensive documentation for all modules and functions to ensure the platform remains maintainable and open-source friendly.

## Error Handling & Integrity
- **Verbose & Traceable:** Provide detailed, timestamped error logs. If a collection fails, the user must be informed exactly why (e.g., rate-limiting vs. network failure) to maintain the investigative audit trail.
- **Isolation:** Failures in one collector or the AI layer must never impact the integrity of the local SQLite database or the evidence locker.
