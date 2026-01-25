# 🕵️ SPECTRE

**Local-First OSINT Intelligence Platform**

> Turn raw internet noise into structured intelligence — fast, repeatable, and local.

Spectre is a commercial-grade OSINT platform that collects intelligence, builds entity graphs, generates timelines, and synthesizes findings using AI. It is designed for security researchers, journalists, and threat analysts who need professional-grade intelligence synthesis without cloud dependencies.

---

## 🚀 Key Features

### 🔍 Deep Collection
*   **Passive Recon:** DNS, Whois, GeoIP, GitHub.
*   **Active Recon:** Port Scanning (Top-100/Custom), Screenshot Capture (Headless Browser), Social Media Username Checks.
*   **Plugin System:** Extensible architecture for custom collectors (Python/Bash/Go).

### 🛡️ Operational Security (Ghost Mode)
*   **Proxy Support:** Route all HTTP, Social, and Screenshot traffic through Tor/SOCKS5/HTTP proxies.
*   **Rate Limiting:** Intelligent throttling to prevent IP bans.
*   **Scope Control:** Built-in safeguards to prevent scanning `.gov` or `.mil` targets.

### 🧠 AI-Augmented Analysis
*   **Local Intelligence:** Uses LLMs (Llama3, Mistral) to analyze case data.
*   **Automatic Synthesis:** Generates Findings, Risks, and Next Steps automatically.
*   **Caching:** Instant re-analysis for unchanged data contexts.

### 📊 Advanced Visualization
*   **TUI Dashboard:** Keyboard-driven terminal interface with Tables, ASCII Graphs, and Timelines.
*   **Web Dashboard:** Interactive React/D3.js node-link diagrams (`spectre web`).
*   **Forensic Timeline:** Chronological view of all discovered entities and collected evidence.

### 📄 Professional Reporting
*   **Markdown:** Instant export for developer documentation.
*   **PDF:** Branded, executive-ready reports with cover pages and summaries.

---

## ⚡ Quick Start

### Installation

**Windows (One-Click):**
```powershell
.\install.ps1
```

**Global Access:**
To run `spectre` from any terminal:
```powershell
.\setup_global.ps1
```

### The "One-Shot" Investigation
The fastest way to start. Auto-creates a case, scans, analyzes, and reports in one command.

```powershell
spectre investigate malicious-site.com
```

### Manual Workflow
For granular control over the intelligence cycle.

1.  **Launch Dashboard:**
    ```powershell
    spectre
    ```
2.  **Create Case:** Follow the TUI prompts (`n` to new case).
3.  **Run Collectors:**
    ```powershell
    spectre collect --case <ID> --target example.com --scanners dns,whois,ports,screenshot
    ```
4.  **Analyze:**
    In the TUI, navigate to **Analysis** and press `2`.
5.  **Visualize:**
    In the TUI, navigate to **Web Dashboard**, press `s` (start) then `o` (open).

---

## 🎮 TUI Controls

Spectre features a professional terminal user interface.

| Key | Action |
| :--- | :--- |
| **TAB** | Toggle focus between Sidebar and Main Content |
| **Arrows** | Navigate selection |
| **Enter** | Confirm selection / Open Case |
| **1-7** | Jump to specific views (Cases, Analysis, Evidence, Graph, Timeline, Reports, Settings) |
| **q** | Quit |

**View-Specific:**
*   **Reports:** Press `1` to generate a Markdown report.
*   **Settings:** Press `s` to switch AI models.
*   **Web Dashboard:** Press `s` to start server, `o` to open browser.

---

## ⚙️ Configuration

Edit `configs/default.yaml` to customize your experience.

### Ghost Mode (Proxy)
```yaml
http:
  proxy: "socks5://127.0.0.1:9050" # Enable for Tor
  insecure_skip_verify: false
```

### Collector Settings
```yaml
collectors:
  ports:
    mode: "top-100" # or "custom"
  screenshot:
    enabled: true
```

### AI Configuration
```yaml
llm:
  provider: "ollama"
  model: "llama3"
```

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      SPECTRE PLATFORM                       │
│  ┌──────────┬──────────┬──────────┬──────────┬──────────┐   │
│  │   TUI    │   CLI    │ Web GUI  │ Analysis │ Reports  │   │
│  └──────────┴──────────┴──────────┴──────────┴──────────┘   │
└─────────────────────────────────────────────────────────────┘
          │              │              │              │
          ▼              ▼              ▼              ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│    Storage   │ │  Collectors  │ │     Graph    │ │   Engine     │
│              │ │              │ │              │ │              │
│ • SQLite     │ │ • Active     │ │ • SQLite     │ │ • Caching    │
│ • Files      │ │ • Passive    │ │   Edges      │ │ • LLM Bridge │
│ • Evidence   │ │ • Plugins    │ │ • JSON/D3    │ │ • PDF Gen    │
└──────────────┘ └──────────────┘ └──────────────┘ └──────────────┘
```

---

## 🤝 Contributing

Contributions are welcome! Please check the `conductor` folder for detailed product guidelines.

## 📄 License

MIT License. See `LICENSE` for details.
