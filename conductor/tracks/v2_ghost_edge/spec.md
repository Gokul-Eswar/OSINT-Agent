# Specification: V2 Ghost Edge

## Goal
Enhance SPECTRE's operational security (OPSEC) by enforcing system-wide proxy support (Ghost Mode) and adding local multimodal AI to analyze evidence (images) without sending data to external cloud services.

## Components

### 1. Ghost Mode Networking
*   **Tor / SOCKS5 Integration:** The internal HTTP client must route all outgoing traffic through a defined SOCKS5 or HTTP proxy when `ghost_mode` is enabled.
*   **DNS Resolution:** Ensure DNS requests also traverse the proxy to prevent leaks.
*   **Centralized Enforcement:** All core collectors (whois, geo, github, etc.) must utilize this central HTTP client.

### 2. Local Multimodal AI (Vision)
*   **Visual Evidence Ingestion:** Ability to parse and index screenshots or image evidence.
*   **Local Processing:** Use a local vision model (e.g., `llava` via Ollama) to generate descriptions of images (e.g., detecting text, logos, or faces) and convert them to structured intelligence that the text-based vector store can index.

## Success Criteria
- [x] Enabling `ghost_mode` in settings routes all subsequent HTTP requests through the configured proxy.
- [x] Non-proxied external requests fail securely if `ghost_mode` is active but the proxy is unreachable.
- [x] Submitting an image to the local analyzer returns a synthesized textual description of its contents.
