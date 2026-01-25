# 🔌 Plugin System

Spectre is designed to be extensible. You can write your own collectors in any language (Python, Bash, Go, Node.js) and integrate them seamlessly into the platform.

## How it Works

Plugins are external executables that Spectre calls via CLI arguments.
1. Spectre runs: `./plugin_executable <target>`
2. Plugin prints JSON to **Stdout**.
3. Spectre captures output, hashes it, and stores it as evidence.
4. Spectre parses the JSON to enrich the Evidence Metadata.

## Creating a Plugin

### 1. Directory Structure
Create a new folder in `plugins/`:
```
plugins/
└── my_cool_scanner/
    ├── plugin.yaml   # Manifest
    └── scan.py       # Your script
```

### 2. The Manifest (`plugin.yaml`)
Tell Spectre how to run your tool.
```yaml
name: "cool_scanner"
description: "Checks for coolness factor"
command: "python"
args: ["scan.py"]
is_active: true  # false = passive only
```

### 3. The Script
Your script must accept the **target** as the first argument and print **JSON** to stdout.

**Example (Python):**
```python
import sys
import json
import time

target = sys.argv[1]

# Do your scanning logic here...
result = {
    "target": target,
    "status": "vulnerable",
    "score": 98,
    "details": ["Found open backdoor", "Default creds active"]
}

print(json.dumps(result))
```

### 4. Running It
Once your plugin is in place, Spectre automatically detects it.

```bash
spectre collect --case <ID> --target example.com --scanners cool_scanner
```

## Advanced Integration

### Structured Data Ingestion
If your plugin returns a valid JSON object, Spectre puts that entire object into the `Evidence.Metadata` field.
This means you can query it later or use it in AI analysis.

### Error Handling
If your plugin fails, exit with a non-zero status code. Print the error message to **Stderr**.
Spectre will capture this log and report the failure in the CLI/TUI.