# SPECTRE - Project Context

## Project Overview
**SPECTRE** is a local-first OSINT (Open Source Intelligence) platform designed to collect, analyze, and synthesize intelligence from various sources. It prioritizes operational security ("Ghost Mode") and offline capability.

*   **Type:** Hybrid Go (System/CLI) and Python (Analysis/AI) application.
*   **Goal:** Turn raw internet data into structured intelligence graphs and reports.
*   **Key Features:**
    *   **TUI & CLI:** Professional terminal interface and command-line tools.
    *   **Collection:** Active (scanning) and Passive (DNS, Whois) data gathering.
    *   **Analysis:** Local LLM integration (via Python) for synthesis.
    *   **Visualization:** Interactive node-link diagrams.

## Architecture
The project follows a hybrid architecture:

1.  **Core (Go):** Located in `internal/`. Handles the CLI framework (`cobra`), data collection, orchestration, TUI (`bubbletea`), and SQLite storage.
2.  **Intelligence Layer (Python):** Located in `analyzer/`. Handles AI analysis (`llm.py`) and graph visualization logic (`graph_viz.py`).
3.  **Bridge:** Go invokes the Python analyzer as a subprocess, passing data via JSON arguments and capturing stdout.
4.  **Storage:**
    *   **Structured:** SQLite (`spectre.db`) for cases, entities, and relationships.
    *   **Unstructured:** Raw evidence files stored in `evidence_storage/`.

## Directory Structure
*   `cmd/spectre/`: Entry point for the Go binary.
*   `internal/`: Core Go logic (CLI, Collectors, Storage, TUI).
*   `analyzer/`: Python module for AI analysis and visualization.
*   `conductor/`: Product guidelines, roadmaps, and documentation.
*   `configs/`: Configuration files (`default.yaml`).
*   `scripts/`: Utility scripts (e.g., build releases).
*   `web/`: Static assets for the web dashboard.

## Building and Running

### Prerequisites
*   **Go:** v1.25.4+
*   **Python:** 3.x with `pip`
*   **Make** (optional, for using the Makefile)

### Build Commands
Use the provided `Makefile` for standard operations:

```bash
# Build the binary
make build

# Install Python dependencies
make install-python

# Build and run
make run

# Clean build artifacts
make clean
```

**Windows Users:**
*   One-click setup: `.\install.ps1`
*   Run wrapper: `.\spectre.bat`

### Manual Build
```bash
go build -o spectre cmd/spectre/main.go
pip install -r analyzer/requirements.txt
```

## Development Conventions

### Go (Core)
*   **Frameworks:** Uses `cobra` for CLI, `bubbletea` for TUI, `zerolog` for logging.
*   **Pattern:** Follows standard Go project layout (`cmd`, `internal`).
*   **Testing:** Run tests via `go test ./...`.

### Python (Analyzer)
*   **Entry Point:** `analyzer/__main__.py` handles CLI arguments (`--task`, `--input`).
*   **Input/Output:** Communicates with Go via JSON over stdin/stdout.

### Configuration
*   Controlled by `configs/default.yaml`.
*   Supports overriding via environment variables or CLI flags (`--config`).

### External Plugins
*   Spectre supports a plugin system (Python/Bash/Go) for custom collectors.
*   Plugins are discovered and registered at startup (see `internal/collector/registry.go`).
