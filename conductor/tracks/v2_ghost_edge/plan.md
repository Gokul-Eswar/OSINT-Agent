# Implementation Plan: V2 Ghost Edge

## Phase 1: Proxy Infrastructure (Go)
- [x] Inspect and update `internal/http/client.go` to support a configured proxy URL (e.g., SOCKS5 or HTTP).
- [x] Read `ghost_mode` flag from the configuration and conditionally enforce the proxy.
- [x] Ensure that collectors (like `geo`, `github`, and `whois`) use this centralized HTTP client.

## Phase 2: Local Multimodal Vision (Python)
- [x] Ensure `llava` or a similar local vision model is accessible via the local Ollama instance configured in `llm.url`.
- [x] Implement an `analyze_image` function in `analyzer/llm.py` that accepts an image path and returns a textual description.
- [x] Add `analyze_image` logic to the agent tools in `internal/agent/tools.go`.

## Phase 3: Validation
- [ ] Enable `ghost_mode` and confirm that a network request fails when no proxy is running, proving the bypass protection.
- [ ] Pass an image to the Python analyzer and confirm it returns a text description suitable for vector indexing.
