# Spectre Plugin Development Guide

Spectre's plugin system allows you to extend the collector framework using any programming language (Python, Node.js, Bash, Go, Rust, etc.).

## How it Works

1.  **Discovery:** On startup, Spectre scans the `plugins/` directory.
2.  **Registration:** It reads `plugin.yaml` from each subdirectory.
3.  **Execution:** When you run `spectre collect <plugin_name> <target>`, Spectre executes the command defined in your YAML.
4.  **Ingestion:** The standard output (stdout) of your script is captured, hashed, and stored as evidence.

## Directory Structure

To create a new plugin, create a folder in `plugins/`:

```
plugins/
└── my_cool_plugin/
    ├── plugin.yaml      # Metadata configuration
    ├── main.py          # Your script (or binary)
    └── requirements.txt # Dependencies (optional)
```

## Configuration (`plugin.yaml`)

```yaml
name: my_cool_plugin
description: Fetches data from My Cool API
command: python          # The executable to run
args: ["main.py"]        # Arguments (target is appended automatically)
is_active: false         # Set to true if this performs active scanning
```

## Writing the Script

Your script must accept the **target** as the last command-line argument and print **JSON** to stdout.

### Python Example (`main.py`)

```python
import sys
import json
import time

def main():
    if len(sys.argv) < 2:
        print(json.dumps({"error": "No target provided"}))
        sys.exit(1)

    target = sys.argv[1]

    # ... Perform your OSINT logic here ...
    
    result = {
        "plugin": "my_cool_plugin",
        "target": target,
        "data": {
            "key": "value",
            "found": True
        }
    }
    
    # Print JSON to stdout
    print(json.dumps(result, indent=2))

if __name__ == "__main__":
    main()
```

## Evidence Storage

Spectre automatically:
1.  Captures the stdout.
2.  Calculates a SHA-256 hash.
3.  Saves it to `evidence_storage/<case_id>/<plugin>_<target>_<timestamp>.json`.
4.  Creates an Evidence record in the SQLite database.

## Advanced: Auto-Ingestion

Currently, external plugins produce "raw evidence". Future versions of Spectre will support a standardized JSON schema to automatically create Entities and Relationships from plugin output.

For now, your data is safely stored and hash-verified, ready for manual review or custom processing.
