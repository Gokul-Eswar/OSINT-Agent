# 💾 OS Installation Guide

This document provides complete, OS-specific installation instructions for deploying **SPECTRE**.

---

## 📋 System Prerequisites

Before installation, verify your environment meets the minimum version requirements:

| Dependency | Minimum Version | Required For | Verification Command |
| :--- | :--- | :--- | :--- |
| **Go** | `1.22+` | Compiling the Go core binary | `go version` |
| **Python** | `3.10+` | Intelligence analyzer and plugins | `python --version` |
| **Git** | `2.x+` | Plugin updates & codebase setup | `git --version` |

---

## 🪟 Windows Setup (Automated)

SPECTRE includes PowerShell automation wrappers for fast setup.

### 1. Configure Execution Policy
PowerShell defaults to restrictive execution limits. Open a PowerShell console as a regular user and enable script execution for the current user scope:
```powershell
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### 2. Run the Installer
Run the installer script from the root of the project directory. This compiles the Go binary, initializes a Python virtual environment (`.venv`), and installs all core ML and database libraries:
```powershell
.\install.ps1
```

### 3. (Optional) Set up Global Terminal access
To run `spectre` from any console prompt (adding its directory to your User PATH):
```powershell
.\setup_global.ps1
```
*Note: Restart your terminal window after running this script for changes to apply.*

---

## 🐧 Linux / 🍎 macOS Setup (Manual)

Follow these steps to compile and set up the platform on Unix-based systems.

### 1. Build the Go Binary
Compile the main entry point to generate the executable binary:
```bash
go build -o spectre cmd/spectre/main.go
chmod +x spectre
```

### 2. Setup Python Virtual Environment
Initialize a local virtual environment in the `.venv` directory to isolate dependencies:
```bash
python3 -m venv .venv
source .venv/bin/activate
```

### 3. Install Python Dependencies
Upgrade pip and install the NLP/vector databases:
```bash
pip install --upgrade pip
pip install -r analyzer/requirements.txt
```

---

## 🧪 Post-Installation Verification

To confirm compilation and environment setup are fully functional, run the built-in checks:

### 1. Verify Core CLI Output
```bash
# Windows
.\spectre.exe version

# Linux/macOS
./spectre version
```
Expected output: `SPECTRE v1.0.0`

### 2. Run Test Suites
Validate Go and Python test suites:
```bash
# Go unit tests
go test ./...

# Python unit tests
# (Make sure to run this using the virtual environment interpreter)
.venv/bin/python -m pytest analyzer/     # Linux/macOS
.venv/Scripts/python -m pytest analyzer/  # Windows
```

---

## 🛡️ Operational Security & Proxy Setup

SPECTRE incorporates a hardened **Ghost Mode** to ensure your OSINT queries are routed securely.

### 1. Install Tor
*   **Windows**: Download and run [Tor Browser](https://www.torproject.org/) (runs proxy on `127.0.0.1:9150`) or run a standalone Tor service.
*   **Debian/Ubuntu**: `sudo apt install tor` (runs daemon on `127.0.0.1:9050`).
*   **macOS**: `brew install tor && brew services start tor` (runs daemon on `127.0.0.1:9050`).

### 2. Configure HTTP Proxy Settings
Edit `configs/default.yaml` to specify your proxy configurations:
```yaml
http:
  proxy: "socks5://127.0.0.1:9050"   # Standard Tor Port
  insecure_skip_verify: false       # Set true only for debugging local SSL proxies
```

### 3. Harden Execution (Strict Proxy Mode)
Run commands with the `--strict` flag. This enforces **fail-closed** behavior: if the proxy goes offline or is unreachable, all network requests are blocked:
```bash
spectre collect all target.com --strict
```
