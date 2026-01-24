# 🕵️ SPECTRE

**Local-First OSINT Intelligence Platform**

> Turn raw internet noise into structured intelligence — fast, repeatable, and local.

Spectre is a CLI-based OSINT agent that collects passive intelligence, builds entity graphs, generates timelines, and synthesizes findings using AI. It is designed for security researchers, journalists, and threat analysts who need professional-grade intelligence synthesis without cloud dependencies or active scanning.

---

## 🎯 Core Principles

*   **Local-First:** No cloud dependency; all data stays on your disk.
*   **Passive-Only:** No active scanning by default (ethical OSINT).
*   **Case-Based:** Every investigation is isolated, auditable, and stored in a local SQLite database.
*   **Evidence Chain:** Forensic-grade provenance and integrity with SHA-256 hashing.
*   **AI-Augmented:** Intelligence synthesis (findings, risks, connections) using local or API-based LLMs.
*   **Extensible:** Plugin architecture for custom collectors.

---

## 🏗️ Architecture

Spectre utilizes a hybrid architecture to leverage the best of both worlds:
*   **Go (System Core):** Handles orchestration, CLI framework (`cobra`), concurrent collection, and SQLite storage.
*   **Python (Intelligence Layer):** Manages AI analysis, graph visualization (`pyvis`), and report generation.

```
┌─────────────────────────────────────────────────────────────┐
│                      SPECTRE CLI (Go)                       │
│  ┌──────────┬──────────┬──────────┬──────────┬──────────┐   │
│  │   Cases  │Collectors│   Graph  │ Timeline │ Analysis │   │
│  └──────────┴──────────┴──────────┴──────────┴──────────┘   │
└─────────────────────────────────────────────────────────────┘
          │              │              │              │
          ▼              ▼              ▼              ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│    Storage   │ │  Collectors  │ │     Graph    │ │   Analyzer   │
│              │ │              │ │              │ │   (Python)   │
│ • SQLite     │ │ • DNS        │ │ • SQLite     │ │ • LLM API    │
│ • Files      │ │ • WHOIS      │ │   Edges      │ │ • Timeline   │
│ • Evidence   │ │ • GitHub     │ │ • GraphML    │ │ • Synthesis  │
│ • Logs       │ │ • Certs      │ │ • pyvis Viz  │ │ • Reports    │
└──────────────┘ └──────────────┘ └──────────────┘ └──────────────┘
```

---

## 🚀 Installation

### Prerequisites
*   Go 1.25+
*   Python 3.11+
*   Git

### Build from Source

```bash
# Clone the repository
git clone https://github.com/spectre/spectre.git
cd spectre

# Build the Go binary
make build

# Install Python dependencies for the analyzer
make install-python

# Full setup (builds binary and installs dependencies)
make install
```

---

## ⚡ Quick Start

Start a new investigation in seconds:

```bash
# 1. Initialize a new case
spectre case new "acme-breach-2026"

# 2. Run collectors to gather intelligence
spectre collect dns acme.com --case <case-id>
spectre collect whois acme.com --case <case-id>

# 3. Visualize the entity graph
spectre visualize --case <case-id>
# (Opens an interactive HTML graph in your browser)

# 4. Generate an AI synthesis report
spectre analyze --case <case-id>
```

---

## 📊 Visual Intelligence Dashboard

Spectre transforms your collected data into an interactive visual graph.

```bash
spectre visualize --case <case-id>
```

*   **Interactive HTML:** Zoom, pan, and drag nodes to explore relationships.
*   **Color-Coded Entities:**
    *   🔵 Domain
    *   🟢 Email
    *   🟠 IP
    *   🟣 Username
    *   🔴 Repository
    *   🩷 Person
*   **Offline:** The dashboard is a standalone HTML file generated in your `evidence_storage` folder.

---

## 🛡️ Ethics & Safety

Spectre is built for **ethical investigation**.

### The Governor (Rate Limiting)
Prevent API bans and reduce footprint.
*   **DNS:** 10 requests/sec
*   **WHOIS:** 1 request/sec (strict enforcement)
*   **GitHub:** 2 requests/sec

### The Fence (Scope Control)
Prevent accidental collection against sensitive targets.
*   **Blacklist:** Automatically blocks collection on `.gov`, `.mil`, `localhost`, and `127.0.0.1`.
*   **Whitelist:** Optional strict mode to only allow specific domains.
*   **Configurable:** Manage rules in `configs/default.yaml`.

---

## 🎨 Command Reference

### Case Management
```bash
spectre case new "name"           # Create a new investigation
spectre case list                 # List all cases
```

### Collection
```bash
spectre collect dns example.com --case <id>    # Collect DNS records
spectre collect whois example.com --case <id>  # Collect WHOIS data
spectre collect github user --case <id>        # Search GitHub user
```

### Visualization & Analysis
```bash
spectre visualize --case <id>     # Generate interactive graph
spectre analyze --case <id>       # Run AI synthesis
```

---

## 📁 Project Structure

```
spectre/
├── cmd/spectre/       # Main entry point
├── internal/
│   ├── cli/           # Cobra CLI commands
│   ├── core/          # Core domain logic
│   ├── collector/     # OSINT collectors (DNS, Whois, GitHub)
│   ├── storage/       # SQLite and file storage
│   ├── ethics/        # Rate limiting and scope control
│   └── analyzer/      # Go bridge to Python analyzer
├── analyzer/          # Python intelligence module (LLM, Graph Viz)
├── configs/           # Configuration files
├── evidence_storage/  # Local data storage (created at runtime)
└── spectre.db         # SQLite database
```

---

## 🤝 Contributing

Contributions are welcome! Please check the `conductor` folder for detailed product guidelines and architectural specs.

## 📄 License

MIT License. See `LICENSE` for details.