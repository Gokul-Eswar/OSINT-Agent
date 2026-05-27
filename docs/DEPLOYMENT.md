# 🚀 Deployment and Operational Guide

This document provides guidance on deploying, daemonizing, and optimizing **SPECTRE** in shared host or enterprise settings.

---

## 1. Running the API Server as a System Service

To run SPECTRE's backend server as a background service on Linux systems, use **systemd**.

Create a systemd unit file at `/etc/systemd/system/spectre.service`:
```ini
[Unit]
Description=SPECTRE OSINT API Server & Web UI
After=network.target

[Service]
Type=simple
User=spectre-service
WorkingDirectory=/opt/spectre
ExecStart=/opt/spectre/spectre server
Restart=on-failure
Environment=SPECTRE_HOME=/var/lib/spectre
Environment=PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/opt/spectre/.venv/bin

[Install]
WantedBy=multi-user.target
```

Enable and start the service:
```bash
sudo systemctl daemon-reload
sudo systemctl enable --now spectre
sudo systemctl status spectre
```

---

## 2. Headless Browser Screenshot Dependencies (Linux)

SPECTRE can capture webpage screenshots of domain targets. Under the hood, this requires a headless browser to render JavaScript and export PNGs.

### Debian / Ubuntu Server Setup
Install Chromium and its required graphics libraries for sandbox execution:
```bash
sudo apt update
sudo apt install -y chromium-browser chromium-chromedriver
```

### Red Hat / CentOS Server Setup
```bash
sudo dnf install -y chromium chromedriver
```

### Troubleshooting Headless Render Failures
If the screenshot collector logs errors regarding browser startup, verify that Chrome can run without a display server:
```bash
chromium-browser --headless --disable-gpu --dump-html https://example.com
```
*Note: In containerized environments (Docker), add `--no-sandbox` options to the browser arguments if running as root.*

---

## 3. SQLite Concurrent Access and Data Locks

SPECTRE relies on SQLite (`spectre.db`) for lightweight, file-based graph storage. 

### Multi-User / API Caveats
*   **Concurrent Reads**: SQLite supports unlimited concurrent reads.
*   **Write Locking**: SQLite serializes write transactions. If multiple collectors or background tasks try to write to the database at the exact same millisecond, a `database is locked` error may occur.
*   **Best Practices**:
    *   Avoid running parallel collection scans for the *same* case ID across different terminal sessions.
    *   Ensure the SQLite database file resides on a local SSD drive. Avoid mounting SQLite files over network-shared filesystems (like NFS, SMB/CIFS, or AWS EFS), as network file locks do not scale and will corrupt the database.

---

## 4. Hardware Optimization for Vector Store Ingestion

Ingesting evidence documents via ChromaDB and `sentence-transformers` requires high CPU or GPU performance during vector embedding calculation.

### CPU Ingestion
By default, PyTorch and tokenizers will attempt to consume all available CPU cores. For shared server environments, restrict thread usage to prevent CPU starvation:
```bash
export OMP_NUM_THREADS=4
export MKL_NUM_THREADS=4
```

### GPU Ingestion (NVIDIA CUDA)
If your deployment server has an NVIDIA GPU:
1.  Verify the CUDA toolkit is installed.
2.  Install PyTorch with GPU support inside `.venv`:
    ```bash
    pip install torch --extra-index-url https://download.pytorch.org/whl/cu121
    ```
3.  SPECTRE checks for the presence of GPU monitoring libraries (such as `NVML`) and CUDA drivers to automatically offload embedding models to GPU memory, cutting down ingestion time by up to 90%.
