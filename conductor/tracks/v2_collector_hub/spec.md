# Specification: V2 Collector Hub

## Goal
Standardize the extension of SPECTRE through a formal Plugin SDK and introduce a "Smart Search" (Automated Dorking) capability to find sensitive exposed data.

## Components

### 1. SPECTRE Python SDK (`spectre_sdk`)
*   **Base Class:** Provide a `BaseCollector` class that handles standard output, environment variables, and error reporting.
*   **Helper Functions:** Include utilities for common OSINT tasks (e.g., regex extraction, JSON formatting).
*   **Boilerplate Generator:** A simple command or template to create a new plugin folder with `plugin.yaml` and a Python script.

### 2. Automated Dorking Engine
*   **Prompt Engineering:** The agent generates complex search queries (Google Dorks) based on the target (e.g., `site:github.com "target.com" password`).
*   **Tool Integration:** A new `generate_dorks` tool for the agent to find leaked databases, `.env` files, and exposed logs.
*   **Search Integration:** (Optional) If an API key for a search engine (like Serper or Google) is provided, the engine can execute these dorks automatically.

## Success Criteria
- [x] A new plugin can be implemented using the `spectre_sdk` in under 20 lines of code.
- [x] The agent can generate at least 5 relevant dorks for a given domain target.
- [x] Plugins correctly inherit proxy settings from the environment.
