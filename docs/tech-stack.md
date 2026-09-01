# Tech Stack — SPECTRE

## Core Orchestration (Authoritative)
- **Language:** Go (Golang)
- **CLI Framework:** `cobra`
- **Configuration:** `viper`
- **Concurrency:** Go routines + worker pools
- **Plugin Host:** external exec / WASM / Go plugins
- **Policy Engine:** core-enforced (permissions, ethics, rate limits)

## Intelligence, Analysis & Visualization
- **Language:** Python (sidecar)
- **Graph Analysis:** `networkx`
- **Visualization:** `pyvis` (HTML)
- **Reporting:** `jinja2` templates
- **Timeline synthesis:** `pandas`
- **Constraint:** No state mutation permitted from this layer

## AI Synthesis Layer (Local-First)
- **Default:** Ollama (local models)
- **Embeddings:** `bge` / `e5`
- **Inference:** `llama.cpp` backend
- **Mode:** offline by default
- **Optional (explicit opt-in, audited):** OpenAI, Anthropic, Gemini

## Storage & Evidence
- **Metadata:** SQLite (WAL mode)
- **Evidence:** JSONL (append-only)
- **Blobs:** structured file store
- **Hashing:** SHA-256 for all evidence
- **Constraint:** Case isolation enforced

## Networking & OSINT
- **HTTP:** `net/http` (stdlib)
- **DNS:** `miekg/dns`
- **WHOIS:** `likexian/whois`
- **TLS:** `crypto/x509`
- **Rate limiting:** Core-enforced
- **Constraint:** Passive-only by default

## Build & Distribution
- **Format:** Single Go binary
- **Dependency:** Python sidecar folder
- **Constraint:** No runtime internet required
- **Extensibility:** Optional plugins (explicit enable)
