# 🏗️ Architecture

Spectre utilizes a hybrid architecture to leverage the best of both worlds:
*   **Go (System Core):** Handles orchestration, CLI framework (`cobra`), concurrent collection, and SQLite storage.
*   **Python (Intelligence Layer):** Manages AI analysis, graph visualization (`pyvis`), and report generation.

## Component Overview

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

## 1. System Core (Go)

The `internal` directory contains the core logic of the application.

### **CLI (`internal/cli`)**
*   Built using `cobra`.
*   Handles command parsing, flags, and user interaction.
*   **Key Commands:** `case`, `collect`, `visualize`, `analyze`.

### **Core Domain (`internal/core`)**
*   Defines the primary data structures: `Case`, `Entity`, `Evidence`, `Relationship`.
*   **Collector Interface:** Defines the contract for all data collectors.

### **Collectors (`internal/collector`)**
*   **Registry:** Manages available collectors.
*   **Implementations:**
    *   `dns`: Uses `net` package for resolution.
    *   `whois`: Uses `github.com/likexian/whois` for parsing.
    *   `github`: Uses `net/http` to query GitHub API.

### **Ethics Guardian (`internal/ethics`)**
*   **Rate Limiter:** Token Bucket algorithm (`golang.org/x/time/rate`) ensures collectors respect API limits.
*   **Scope Control:** Validates targets against blacklists/whitelists before execution.

### **Storage (`internal/storage`)**
*   **SQLite:** Stores structured metadata (cases, entities, relationships).
*   **File System:** Stores raw evidence files (JSON, text) in `evidence_storage/`.

## 2. Intelligence Layer (Python)

The `analyzer` directory contains the Python module responsible for high-level synthesis.

### **Bridge (`internal/analyzer`)**
*   Go code that marshals case data into JSON and executes the Python module via `exec.Command`.
*   Handles data transfer via `stdin/stdout`.

### **Visualizer (`analyzer/graph_viz.py`)**
*   **Input:** JSON graph data (nodes, edges).
*   **Processing:** Uses `networkx` to build the graph structure.
*   **Output:** Generates an interactive HTML file using `pyvis` with physics-based layout.

### **LLM Synthesis (`analyzer/llm.py`)**
*   **Input:** Textual context of the case (entities, relationships).
*   **Processing:** Sends prompts to an LLM (currently supports local Ollama).
*   **Output:** Structured JSON report (findings, risks, connections).

## 3. Data Flow

### **Collection Flow**
1.  User runs `spectre collect <collector> <target>`.
2.  **Ethics Check:** Target is checked against blacklist.
3.  **Rate Limit:** Collector waits for available token.
4.  **Execution:** Collector fetches data.
5.  **Storage:** Raw data saved to `evidence_storage/`.
6.  **Ingestion:** Entities and relationships extracted and saved to SQLite.

### **Visualization Flow**
1.  User runs `spectre visualize`.
2.  **Export:** Go exports all case entities/relationships to JSON.
3.  **Bridge:** Go calls `python -m analyzer --task visualize`.
4.  **Generation:** Python builds the `pyvis` graph and saves `report.html`.
5.  **Display:** Go opens the HTML file in the default browser.

## 4. Configuration

Managed via `configs/default.yaml` and loaded by `internal/config`.

*   **Database Path:** Location of `spectre.db`.
*   **Collector Settings:** Rate limits and enabled status.
*   **Ethics:** Global blacklists/whitelists.