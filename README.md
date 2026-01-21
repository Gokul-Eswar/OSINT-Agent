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
│ • Evidence   │ │ • Certs      │ │ • GraphML    │ │ • Synthesis  │
│ • Logs       │ │ • GitHub     │ │ • pyvis Viz  │ │ • Reports    │
└──────────────┘ └──────────────┘ └──────────────┘ └──────────────┘
```

---

## 🚀 Installation

### Prerequisites
*   Go 1.21+
*   Python 3.11+

### Build from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/spectre.git
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
spectre new-case "acme-breach-2024"

# 2. Add initial entities (passive collection targets)
spectre add domain acme.com
spectre add email security@acme.com

# 3. Run collectors to gather intelligence
spectre run --case acme-breach-2024

# 4. Visualize the entity graph
spectre graph --case acme-breach-2024
# (Opens an interactive HTML graph in your browser)

# 5. Generate an AI synthesis report
spectre analyze --case acme-breach-2024
```

---

## 🎨 Command Reference

### Case Management
```bash
spectre init                      # Initialize system
spectre new-case "name"           # Create a new investigation
spectre list                      # List all cases
spectre show --case "name"        # Show case details
```

### Entity Management
```bash
spectre add domain example.com    # Add a domain
spectre add email user@site.com   # Add an email
spectre entities --case "name"    # List entities in a case
```

### Collection & Analysis
```bash
spectre run --case "name"         # Run all enabled collectors
spectre run --passive-only        # Enforce passive collection
spectre analyze --case "name"     # Run AI synthesis
spectre timeline --case "name"    # Generate investigation timeline
```

---

## 📁 Project Structure

```
spectre/
├── cmd/spectre/       # Main entry point
├── internal/
│   ├── cli/           # Cobra CLI commands
│   ├── core/          # Core domain logic (Case, Entity, Evidence)
│   ├── collectors/    # OSINT collectors (DNS, Whois, etc.)
│   ├── storage/       # SQLite and file storage
│   └── analyzer/      # Go bridge to Python analyzer
├── analyzer/          # Python intelligence module (LLM, Graph Viz)
├── configs/           # Configuration files
└── cases/             # Local data storage (created at runtime)
```

---

## 🛡️ Ethics & Safety

Spectre is built for **ethical investigation**.
*   **Passive by Design:** Default collectors do not probe target infrastructure aggressively.
*   **Rate Limiting:** Built-in safeguards preventing API abuse and scanning detection.
*   **Audit Trail:** All actions are logged to a local, immutable chain of evidence.

---

## 🤝 Contributing

Contributions are welcome! Please check the `conductor` folder for detailed product guidelines and architectural specs.

## 📄 License

MIT License. See `LICENSE` for details.
