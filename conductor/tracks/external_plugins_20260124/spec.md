# Specification - External Plugin System

## Overview
To become a true platform, SPECTRE must allow users to add their own collectors without recompiling the Go binary. We will implement a "Sidecar" plugin system.

## Requirements

### 1. Plugin Discovery
- SPECTRE will scan a specific directory (e.g., `plugins/`) for subdirectories containing a `plugin.yaml` or `plugin.json` metadata file.

### 2. Plugin Contract
- **Metadata:** Name, description, input type (domain, ip, etc.), and the command to execute.
- **Protocol:**
    - **Input:** Passed via environment variables or CLI arguments.
    - **Output:** The plugin MUST print a valid SPECTRE Evidence JSON to `stdout`.

### 3. Execution Engine
- A new `ExternalCollector` wrapper that implements the `core.Collector` interface.
- It will use `os/exec` to run the plugin script/binary.

## Success Criteria
- Creating a Python script in `plugins/my_tool/` and adding a `plugin.yaml` makes it appear in `spectre collectors`.
- Running `spectre collect my_tool <target>` executes the script and ingests its findings.
