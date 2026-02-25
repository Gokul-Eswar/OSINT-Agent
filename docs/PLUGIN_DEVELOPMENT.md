# 🛠️ Plugin Development Guide

SPECTRE is built on a "Decoupled Execution" model. This means you can write collectors in any language, as long as they can be executed from the command line and return JSON.

---

## 1. Plugin Architecture

A plugin consists of two main parts:
1.  **The Manifest (`plugin.yaml`)**: Metadata that tells SPECTRE how to run your tool.
2.  **The Executable**: A script or binary (Python, Bash, Go, etc.) that performs the collection.

### Directory Structure
Your plugin must live in a subfolder within the `plugins/` directory:
```text
plugins/
└── my-awesome-plugin/
    ├── plugin.yaml       # REQUIRED
    ├── main.py           # Your script
    └── requirements.txt  # Optional (Python deps)
```

---

## 2. Implementation Examples

### 🐍 Python Implementation
Python is the recommended language for most plugins due to its rich library ecosystem.

**`plugin.yaml`**
```yaml
name: "my_python_scanner"
version: "1.0.0"
description: "Fetches custom data via API"
command: "python"
args: ["main.py"]
is_active: false
```

**`main.py`**
```python
import sys
import json

def collect(target):
    # Your logic here
    return {
        "status": "success",
        "data": f"Results for {target}"
    }

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(json.dumps({"error": "No target provided"}))
        sys.exit(1)
    
    target = sys.argv[1]
    print(json.dumps(collect(target)))
```

### 🐚 Bash Implementation
Ideal for simple wrappers around existing CLI tools like `curl`, `nmap`, or `dig`.

**`plugin.yaml`**
```yaml
name: "http_check"
command: "bash"
args: ["check.sh"]
is_active: true
```

**`check.sh`**
```bash
#!/bin/bash
TARGET=$1
STATUS_CODE=$(curl -o /dev/null -s -w "%{http_code}" "$TARGET")

echo "{"target": "$TARGET", "http_status": $STATUS_CODE}"
```

### 🐹 Go Implementation (External)
You can distribute your plugin as a compiled binary.

**`plugin.yaml`**
```yaml
name: "fast_scanner"
command: "./scanner_bin"
args: []
is_active: true
```

---

## 3. Best Practices

1.  **JSON Output**: Always return a single, valid JSON object to `stdout`.
2.  **Stderr for Logs**: Send all debugging info or non-critical errors to `stderr`. SPECTRE captures `stdout` as primary evidence.
3.  **Timeout Handling**: If your tool takes a long time, ensure it handles interrupts gracefully.
4.  **Target Agnostic**: Ensure your plugin handles different target types (IP, Domain, Username) or validates that it received the correct type.

---

## 4. Plugin Marketplace Strategy

SPECTRE uses a decentralized registry system. To make your plugin available in the marketplace:

1.  **Host on GitHub**: Create a public repository for your plugin.
2.  **Registry Entry**: A `registry.json` file (hosted centrally or provided locally) maps plugin names to their Git URLs.
3.  **Submission**: Currently, new plugins are added by submitting a PR to the official SPECTRE registry repository.

### Marketplace Metadata
Your `plugin.yaml` can include extra fields for the UI:
```yaml
author: "YourName"
repo: "https://github.com/user/my-plugin"
tags: ["network", "active", "api"]
```

---

## 5. Security Guidelines

- **No Remote Execution**: Plugins should never download and execute unknown code at runtime.
- **Privacy**: If your plugin uses an external API, document what data is sent to that API.
- **Ethics**: Respect the `ethics` configuration provided by the SPECTRE core (Rate limits, Scope).
