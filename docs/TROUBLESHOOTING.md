# ❓ Troubleshooting & FAQ

Stuck? Find your issue below. Most problems have simple fixes.

---

## Installation Problems

### "Python not found" or "Go not found"
**What it means:** SPECTRE can't find Python or Go on your computer.

**How to fix:**
1. Check if Python is installed: `python --version` (or `python3 --version`)
2. Check if Go is installed: `go version`
3. If not installed, download them:
   - Go: https://go.dev/dl/
   - Python: https://www.python.org/downloads/
4. Restart your terminal after installing

---

### "Permission denied" on Linux/macOS
**What it means:** The SPECTRE binary can't be run.

**How to fix:**
```bash
chmod +x spectre
```

---

### "Execution Policy" error on Windows
**What it means:** PowerShell is blocking the installer script.

**How to fix:**
```powershell
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
```

Then try running the installer again.

---

## Running SPECTRE Problems

### "Database is locked" error
**What it means:** Multiple operations are trying to write to the database at the same time.

**How to fix:**
- Run commands one at a time, not in parallel
- Don't run multiple `spectre collect` on the same case simultaneously
- Wait for one command to finish before starting another

---

### "Collector not found"
**What it means:** SPECTRE can't find a custom collector you're trying to use.

**How to fix:**
1. Make sure the plugin folder is in the `plugins/` directory
2. Check that `plugin.yaml` file exists
3. Restart SPECTRE (it rescans plugins on startup)

---

### "ImportError: No module named..."
**What it means:** A Python library is missing.

**How to fix:**

**Windows:**
```powershell
.\.venv\Scripts\activate
pip install -r analyzer/requirements.txt
```

**Linux/macOS:**
```bash
source .venv/bin/activate
pip install -r analyzer/requirements.txt
```

---

### "Connection refused" when starting web dashboard
**What it means:** Port 8080 is already in use by another program.

**How to fix:**
- Option 1: Close the other program using port 8080
- Option 2: Use a different port:
  ```bash
  spectre server --port 8888
  ```
- Option 3: Check what's using the port:
  - Windows: `netstat -ano | findstr 8080`
  - Linux/macOS: `lsof -i :8080`

---

## Common Tasks

### "How do I reset a case?"
Delete it and create a new one:
```bash
spectre case delete <CASE_ID>
spectre case new "New Investigation"
```

### "Where is my data stored?"
Everything is in your project folder:
- `spectre.db` — the main database
- `evidence_storage/` — all raw files and screenshots

### "How do I back up my investigations?"
Copy these two folders:
```bash
# Windows
copy spectre.db spectre.db.backup
xcopy evidence_storage evidence_storage.backup /E /I

# Linux/macOS
cp spectre.db spectre.db.backup
cp -r evidence_storage evidence_storage.backup
```

### "Can I run SPECTRE through a proxy?"
Yes! SPECTRE respects standard proxy settings. Either:
1. Set environment variables:
   ```bash
   export HTTP_PROXY=http://proxy.example.com:8080
   export HTTPS_PROXY=http://proxy.example.com:8080
   ```

2. Or edit `configs/default.yaml` and set the proxy there.

3. To enable "strict mode" (fail if proxy is down):
   ```bash
   spectre collect all example.com --strict
   ```

---

## Performance Issues

### "SPECTRE is slow"
Try these fixes:

1. **For data collection:** Run collectors in parallel
   ```bash
   spectre collect all example.com --case <CASE_ID>
   ```

2. **For AI analysis:** If running slow LLM queries:
   - Make sure you have enough RAM (8GB+ recommended)
   - If you have a GPU, configure Ollama to use it
   - Try with a smaller model

3. **For large cases:** Create smaller, focused cases instead of one massive investigation

---

## Windows-Specific Tips

- **Terminal looks weird:** Try using Windows Terminal instead of cmd.exe
- **Paths with spaces:** Put them in quotes: `spectre collect "My Target.com"`
- **Backslashes:** Always use backslashes (`\`) in Windows paths

---

## Linux-Specific Tips

- Make sure `python3` and `pip` are installed:
  ```bash
  sudo apt install python3 python3-pip  # Debian/Ubuntu
  sudo dnf install python3 python3-pip  # Fedora
  ```

---

## macOS-Specific Tips

- If you installed Go/Python via Homebrew, they should work automatically
- If something isn't found, reinstall:
  ```bash
  brew install go python@3.10
  ```

---

## Still Stuck?

1. Check the main docs: [docs/README.md](README.md)
2. Look at [docs/INSTALLATION.md](INSTALLATION.md) for setup details
3. Open an issue on GitHub with:
   - What you were trying to do
   - The exact error message
   - Your operating system and Go/Python versions
