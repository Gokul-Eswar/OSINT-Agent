# 🕵️ SPECTRE — Local-First OSINT Intelligence Platform

[![CI Status](https://github.com/Gokul-Eswar/Spectre/actions/workflows/ci.yml/badge.svg)](https://github.com/Gokul-Eswar/Spectre/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?logo=go)](go.mod)
[![Python Version](https://img.shields.io/badge/Python-3.10%2B-3776AB?logo=python)](analyzer/requirements.txt)
[![Local-First](https://img.shields.io/badge/Privacy-100%25%20Local--First-success)](#-core-principles)

> **Turn raw internet noise into structured, verifiable intelligence — fast, repeatable, and 100% local.**

---

## 📖 What is SPECTRE? (Explain It In Plain English)

Imagine you are investigating a suspicious website, a malicious domain, a leaked email, or an online persona. 

Normally, an investigator has to manually juggle a dozen disconnected tools: running `whois`, querying DNS records, checking IP geolocations, running port scanners, capturing website screenshots, searching GitHub for leaked API keys, and searching username availability across 50+ social platforms. Then, they have to manually piece together messy notes to figure out how everything is connected.

**SPECTRE solves this by providing a unified, local-first intelligence workbench:**

```
                                  THE SPECTRE PIPELINE
                                  
  [ Target ] ──► 1. COLLECT ──► 2. LINK & INGEST ──► 3. SYNTHESIZE ──► 4. DELIVER
  (Domain,       (DNS, WHOIS,     (Extract entities   (Local LLM AI    (Interactive Web Graph,
   IP, Email,     GeoIP, Ports,    into knowledge      finds risks &    TUI Console, PDF /
   Username)      Screenshots)     graph in SQLite)    connections)     Markdown Reports)
```

1. **Target In:** You give SPECTRE a seed target (a domain like `example.com`, an IP, or a username).
2. **Automated Collection:** SPECTRE runs passive and active collectors simultaneously, capturing raw evidence files.
3. **Graph Ingestion:** It parses the evidence, extracts entities (People, Emails, Domains, IPs, Ports, Repos, Accounts), and links them into an interconnected knowledge graph.
4. **AI-Powered Synthesis:** A local AI model (via Ollama or LLaVA) reads the gathered evidence, uncovers hidden relationships, identifies risks, and suggests next investigation leads.
5. **Multi-Interface Delivery:** You can explore findings through an interactive terminal interface (TUI), a live web dashboard with real-time physics-based graphs, or export forensic-grade PDF/HTML/Markdown reports.

---

## 🎯 Why SPECTRE? (Core Principles)

- 🔒 **100% Local-First & Private:** All investigation data, SQLite databases, and evidence files live exclusively on your machine (`./evidence_storage/` & `spectre.db`). No telemetry, no mandatory cloud subscriptions, and no leaks to third-party servers.
- ⚖️ **Forensic Auditability:** Every piece of evidence is hashed using SHA-256 and immutably timestamped so investigations can be audited or presented as evidence.
- 👻 **Ghost Mode (Operational Security):** Built-in routing through Tor, SOCKS5, or HTTP proxies with automatic rate limiting and scope safety rules (refuses to accidentally scan `.gov` or `.mil` targets).
- ⚡ **Hybrid Architecture:** High-speed orchestration, database storage, and CLI/TUI written in **Go**; AI synthesis, vector search, and graph visualization written in **Python**.

---

## 🌟 Key Features

| Capability | Description |
| :--- | :--- |
| **Passive Reconnaissance** | Query DNS records (A, MX, NS, TXT), WHOIS registrations, IP Geolocation (City, Country, ISP), and public GitHub repositories without touching target servers. |
| **Active Reconnaissance** | Multi-threaded TCP port scanning (common top-100 or custom ports), automated headless browser full-page screenshot capture, and username availability checking across 50+ platforms. |
| **Autonomous AI Agent** | Interactive chat REPL where an AI investigator autonomously decides which tools to run, searches evidence, analyzes screenshots, and records hypotheses. |
| **Semantic Vector Search** | Search through unstructured case evidence using vector embeddings (via ChromaDB & Sentence Transformers) or literal substring matching. |
| **Persona Correlation** | Cluster related usernames, social media accounts, and email handles into unified persona profiles. |
| **Interactive Web Command Center** | Real-time web application featuring a physics-based force-directed graph (Vis.js), Server-Sent Events (SSE) live updates, and an embedded agent chat drawer. |
| **Terminal User Interface (TUI)** | Keyboard-driven terminal console with ASCII graph visualizers, entity data tables, and live logs. |
| **Multi-Format Export** | Export investigations into executive-ready PDF documents, standalone HTML graphs, CSV data sheets, JSON bundles, or Markdown notes. |
| **Plugin Ecosystem & Store** | Create custom collectors in Python, Bash, or Go in under 20 lines of code using the built-in `spectre_sdk` and manage extensions via the marketplace. |

---

## 🚀 Quick Start

### 1. Prerequisites
- **Go:** `1.21+` ([Download Go](https://go.dev/dl/))
- **Python:** `3.10+` with `pip` ([Download Python](https://www.python.org/downloads/))
- **Ollama (Optional, for local AI):** ([Download Ollama](https://ollama.ai/)) — Run `ollama pull llama3` and `ollama pull llava`.

---

### 2. Installation

#### Windows (One-Click Setup)
Run the automated installation script in PowerShell:
```powershell
.\install.ps1
```
To enable global command line access from any terminal:
```powershell
.\setup_global.ps1
```

#### Linux & macOS / Manual Setup
```bash
# 1. Clone repository
git clone https://github.com/Gokul-Eswar/Spectre.git
cd Spectre

# 2. Build the Go binary
go build -o spectre cmd/spectre/main.go

# 3. Setup Python virtual environment & dependencies
python3 -m venv .venv
source .venv/bin/activate
pip install -r analyzer/requirements.txt
```

---

## 🕹️ Three Ways to Use SPECTRE

### 1️⃣ Quick Scan (One Command)
Want to quickly investigate a target? Get results instantly:
```bash
spectre investigate scanme.nmap.org
```
SPECTRE will:
- Create a case automatically
- Gather DNS, WHOIS, GeoIP info
- Scan ports
- Analyze everything with AI
- Show you a summary

Perfect for: Quick checks, when you don't need a detailed investigation.

---

### 2️⃣ Terminal Dashboard (TUI)
Prefer working in your terminal? Launch the interactive dashboard:
```bash
spectre
```

You'll get:
- A menu to manage cases
- View all collected data
- Search through findings
- See interactive graphs
- Run commands with keyboard shortcuts

Perfect for: Advanced investigators who live in the terminal.

---

### 3️⃣ Web Dashboard
Want a modern browser interface? Start the web server:
```bash
spectre server
```
Then open `http://localhost:8080` in your browser.

You'll get:
- Interactive graphs you can drag around
- Search and filter data
- Chat with AI analysis
- Export reports
- Real-time updates

Perfect for: Visual exploration and presenting findings.

---

## 💻 Common Commands

Here are the commands you'll use most often. Don't worry about memorizing them all — just reference this when you need it.

### Creating & Managing Cases
```bash
# Create a new investigation
spectre case new "Operation Name"

# List all your investigations
spectre case list

# See details of one case
spectre case show <CASE_ID>

# Delete a case
spectre case delete <CASE_ID>

# Set a case as your default (won't need to type --case)
spectre context set <CASE_ID>
```

### Collecting Information
```bash
# Quick passive scan (DNS, WHOIS, IP info)
spectre collect all example.com --case <CASE_ID>

# Scan for open ports (requires --active permission)
spectre collect ports example.com --case <CASE_ID> --active

# Check if username exists on social media
spectre collect accounts username123 --case <CASE_ID>
```

### Searching & Asking Questions
```bash
# Search through collected evidence
spectre search "admin@example.com"

# Ask the AI a question
spectre chat --case <CASE_ID>
# Then type your questions interactively
```

### Exporting Results
```bash
# Export as PDF (professional looking)
spectre report pdf --case <CASE_ID>

# Export as Markdown (easy to edit)
spectre report markdown --case <CASE_ID>

# Export as interactive HTML graph
spectre report html --case <CASE_ID>
```

### Other Useful Commands
```bash
# View SPECTRE version
spectre version

# Get help on any command
spectre help
spectre <command> help

# Start the web dashboard
spectre server

# Visualize the graph
spectre visualize --case <CASE_ID>
```

---

## 🏛️ System Architecture

SPECTRE utilizes a clean hybrid decoupled architecture:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              USER INTERFACES                                │
│       CLI (Cobra)       │    TUI (BubbleTea)    │    Web GUI (React/Vis.js) │
└─────────────────────────────────────────────────────────────────────────────┘
                                       │
                                       ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                            CORE ENGINE (Go)                                 │
│  ┌───────────────────────┬───────────────────────┬───────────────────────┐  │
│  │ Collector Orchestrator│  Storage & Repository │  Ethics & Rate Guard  │  │
│  │ (Active & Passive)    │  (SQLite WAL Mode)    │  (Ghost Mode / Tor)   │  │
│  └───────────────────────┴───────────────────────┴───────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
                                       │ (JSON Subprocess Bridge)
                                       ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                       INTELLIGENCE LAYER (Python)                           │
│  ┌───────────────────────┬───────────────────────┬───────────────────────┐  │
│  │ LLM Engine (Ollama)   │ ChromaDB Vector Store │ Graph Viz (PyVis)     │  │
│  │ (Llama3, LLaVA, Dorks)│ (Dense Embeddings)    │ (NetworkX HTML)       │  │
│  └───────────────────────┴───────────────────────┴───────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Directory Structure
```
├── cmd/
│   ├── spectre/          # Main application entry point & embedded web assets
│   └── installer/        # Single-binary self-extracting installer
├── internal/
│   ├── agent/            # Autonomous ReAct agent loop and tool registry
│   ├── analysis/         # Context compilation and AI synthesis coordinator
│   ├── analyzer/         # Go-to-Python subprocess bridge contract
│   ├── cli/              # Cobra CLI command definitions
│   ├── collector/        # Active & passive collectors (DNS, WHOIS, Ports, etc.)
│   ├── config/           # Viper configuration hierarchy
│   ├── core/             # Domain entities, models, relationships, and evidence
│   ├── ethics/           # Rate limiting, proxying, and blacklist scope guards
│   ├── extensions/       # Plugin installer and marketplace manager
│   ├── report/           # PDF, HTML, Markdown, CSV, and JSON report generators
│   ├── server/           # REST API & Server-Sent Events (SSE) server
│   ├── storage/          # SQLite database migrations and repositories
│   └── tui/              # BubbleTea terminal dashboard
├── analyzer/             # Python analysis engine (LLM, Vision, Vector Store)
├── configs/              # Default configuration (default.yaml)
├── docs/                 # Architectural specifications and feature guides
├── lib/                  # Frontend assets & Python SDK (spectre_sdk.py)
├── plugins/              # Built-in and third-party extension collectors
└── scripts/              # Build, packaging, and validation scripts
```

---

## ⚙️ Configuration

Configuration is managed via [`configs/default.yaml`](configs/default.yaml) or overridden using environment variables and CLI flags.

```yaml
# Database Configuration
database:
  type: "sqlite"          # "sqlite" (default) or "postgres"
  path: "spectre.db"

# LLM Backend Settings
llm:
  provider: "ollama"      # "ollama", "openai", or "local"
  url: "http://localhost:11434/api/generate"
  model: "llama3"
  timeout: 120

# Operational Security & Ghost Mode
ghost_mode: false         # Set true to route traffic through proxies
http:
  proxy: ""
  tor_proxy: "socks5://127.0.0.1:9050"
  strict: false           # If true, fails closed if proxy is unreachable

# Safety & Scope Safeguards
ethics:
  blacklist:
    - ".gov"
    - ".mil"
    - "localhost"
    - "127.0.0.1"
```

---

## 🔌 Writing Custom Plugins in 2 Minutes

SPECTRE supports custom collectors written in Python, Bash, or Go.

Create `plugins/my_collector/main.py`:
```python
import sys
from lib.python.spectre_sdk import BaseCollector

class MyCollector(BaseCollector):
    def __init__(self):
        super().__init__(name="my_collector", description="Finds custom assets")

    def collect(self, target):
        self.log(f"Scanning target: {target}")
        # Your collection logic here:
        return {
            "target": target,
            "status": "active",
            "custom_data": ["asset1.example.com", "asset2.example.com"]
        }

if __name__ == "__main__":
    MyCollector().run()
```

Create `plugins/my_collector/plugin.yaml`:
```yaml
name: my_collector
version: 1.0.0
description: Custom asset collector
author: Your Name
type: collector
entrypoint: main.py
active: false
```

SPECTRE will automatically discover and register your plugin on the next run!

---

## 🧪 Development & Quality Gates

Run all quality checks locally before submitting contributions:

```bash
# Run all Go unit tests
go test ./...

# Run Python unit tests
pytest

# Validate plugin manifests and YAML configs
python scripts/validate_plugins.py
python scripts/validate_configs.py

# Check test coverage threshold
go test -coverprofile=coverage.out ./...
go run ./scripts/check_coverage.go -profile coverage.out -min 35
```

---

## 📄 License

This project is licensed under the **MIT License** — see the [LICENSE](LICENSE) file for details.