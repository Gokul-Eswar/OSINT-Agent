# 🕵️ SPECTRE Documentation Index

Welcome to the documentation suite for **SPECTRE**, the local-first Open Source Intelligence (OSINT) platform. Use this index to navigate guides, API schemas, tutorials, and operational references.

---

## 📚 Guides and Walkthroughs

*   **[Installation Guide](INSTALLATION.md)**: Detailed platform-specific installation instructions (Windows, Linux, macOS) covering prerequisites, execution policies, and environment setup.
*   **[Getting Started Tutorial](GETTING_STARTED.md)**: A step-by-step walkthrough detailing how to create your first case, execute collectors, interact with the conversational AI agent, and generate reports.
*   **[Troubleshooting & FAQ](TROUBLESHOOTING.md)**: A collection of troubleshooting steps for common issues, platform-specific bugs, database locks, and general questions.

---

## 🛠️ Architecture and Extensibility

*   **[System Architecture](ARCHITECTURE.md)**: Conceptual overview of Go/Python split of labor, core internal domains, ethics guardian, and databases.
*   **[Bridge API & Architecture Schema](API_DOCUMENTATION.md)**: Details the design of the Go-to-Python subprocess bridge, including request/response JSON contracts and internal Go package layers.
*   **[Plugin Development Guide](PLUGIN_DEVELOPMENT.md)**: The developer's guide to building custom intelligence collectors using Python, Bash, or compiled binaries.
*   **[External Plugins Model](plugins.md)**: Context on how external collectors are discovered, registered, and run safely.
*   **[LLM Integration](llm.md)**: Design notes on the local intelligence/caching system, providers (Ollama), and offline fallbacks.

---

## 🚀 Operational guides

*   **[Deployment Guide](DEPLOYMENT.md)**: Operational guidance for running SPECTRE on shared servers or enterprise settings, setting up headless screenshot drivers, daemonizing backend servers, and multi-user configurations.
*   **[Performance Baselines](performance.md)**: System performance benchmarks, parallel scanning, and database optimization indices.
*   **[Feature Catalog](features.md)**: Complete list of native collectors, active scanners, and system-wide capabilities.
