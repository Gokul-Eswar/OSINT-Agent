# 🕵️ SPECTRE — Local-First OSINT Intelligence Platform

> **Turn raw internet noise into structured, verifiable intelligence — fast, repeatable, and 100% local.**

"SPECTRE is like having a smart detective assistant that gathers all online clues about a target in one place, figures out connections automatically, and keeps everything private and locked down—all without needing the internet or cloud services."


## 🎯 The Problem SPECTRE Solves

Imagine you're investigating a suspicious website or tracking down who's behind an online account. Right now, you'd have to use 10+ different tools—check who owns the domain, look up the server location, scan for open ports, take screenshots, search social media manually. Then you'd copy-paste everything into a spreadsheet and try to connect the dots yourself. It's slow, messy, and your data gets scattered across cloud services you don't fully trust.

SPECTRE brings everything together in one place. Tell it a website or username, and it automatically gathers all the information you need—who owns it, where the servers are, what services are running, and what else it's connected to. Then it uses AI to connect the dots for you and show you what matters. Everything stays on your computer (no cloud), and every finding is locked down so it can be used as evidence in court or investigations. What used to take hours now takes minutes.





**SPECTRE solves this by providing a unified, local-first intelligence workbench:**

```
                                  THE SPECTRE PIPELINE
                                  
  [ Target ] ──► 1. COLLECT ──► 2. LINK & INGEST ──► 3. SYNTHESIZE ──► 4. DELIVER
  (Domain,       (DNS, WHOIS,     (Extract entities   (Local LLM AI    (Interactive Web Graph,
   IP, Email,     GeoIP, Ports,    into knowledge      finds risks &    TUI Console, PDF /
   Username)      Screenshots)     graph in SQLite)    connections)     Markdown Reports)
```

## 🚀 Setup

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



## 📄 License

This project is licensed under the **MIT License** — see the [LICENSE](LICENSE) file for details.
