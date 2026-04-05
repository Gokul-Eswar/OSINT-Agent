# 🔍 Troubleshooting & FAQ

## Common Issues

### 1. Plugin Not Found
- **Issue:** You run a collector but get "Collector not found".
- **Solution:** 
    - Ensure the plugin folder is in the `plugins/` directory.
    - Check that `plugin.yaml` exists and the `name` field matches what you are calling.
    - Restart SPECTRE to trigger a plugin rescan.

### 2. Python Dependency Errors
- **Issue:** `ImportError: No module named ...`
- **Solution:** 
    - SPECTRE attempts to install dependencies from `requirements.txt` automatically. If this fails, manually activate the SPECTRE virtual environment and install them:
      ```bash
      # Windows
      .\.venv\Scripts\activate
      pip install -r plugins/your-plugin/requirements.txt
      ```
      ```bash
      # Linux/macOS
      source .venv/bin/activate
      pip install -r plugins/your-plugin/requirements.txt
      ```

### 3. Database is Locked
- **Issue:** `database is locked` error when running multiple commands.
- **Solution:** SPECTRE uses SQLite. While it supports concurrent reads, concurrent writes can occasionally cause locks. Avoid running multiple `collect` commands simultaneously for the *same* case ID.

---

## Platform-Specific Guides

### 🪟 Windows
- **Execution Policy:** If `install.ps1` fails, run `Set-ExecutionPolicy RemoteSigned -Scope CurrentUser` in PowerShell.
- **Paths:** Always use backslashes (``) in CLI arguments or quote your paths.
- **TUI Rendering:** If the TUI looks garbled, try using **Windows Terminal** instead of the legacy `cmd.exe`.

### 🐧 Linux
- **Permissions:** Ensure the `spectre` binary has execute permissions (`chmod +x spectre`).
- **Dependencies:** Ensure `python3` and `pip` are installed via your package manager (`apt`, `dnf`, etc.).

---

## Performance Tuning

### 1. Parallel Collection
By default, SPECTRE runs collectors sequentially. To speed up investigations, you can use the comma-separated scanner list:
```bash
spectre collect --target example.com --scanners dns,whois,geo
```

Current baseline measurements and re-run rules are documented in [performance.md](./performance.md).

### 2. LLM Response Times
If `spectre analyze` is slow:
- **Local Models:** Ensure you have enough RAM (8GB+ for 7B models). If you have a GPU, ensure Ollama is configured to use it.
- **Context Size:** Very large cases can slow down the LLM. Try to limit the amount of evidence by creating focused cases.

### 3. SQLite Optimization
For very large investigations (thousands of entities), SPECTRE automatically manages indexes. If performance degrades, consider archiving old cases to keep the active database small.

---

## FAQ

**Q: Where is my data stored?**
A: Everything is in your project directory:
- `spectre.db`: The relationship graph and case metadata.
- `evidence_storage/`: Raw files, screenshots, and logs.

**Q: Can I use SPECTRE behind a proxy?**
A: Yes. SPECTRE respects standard environment variables: `HTTP_PROXY`, `HTTPS_PROXY`, and `NO_PROXY`.

**Q: How do I backup my investigations?**
A: Simply copy the `spectre.db` file and the `evidence_storage/` folder.
