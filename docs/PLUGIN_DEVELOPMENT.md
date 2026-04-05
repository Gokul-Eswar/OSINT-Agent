# Plugin Development Guide

This is the canonical guide for building, testing, and distributing SPECTRE plugins.

SPECTRE uses a decoupled execution model:
- the core orchestrates collection and policy in Go,
- plugins run as external executables,
- plugin stdout is captured as evidence.

For high-level ecosystem usage, see docs/plugins.md.

## 1. Architecture and Lifecycle

Plugin lifecycle:
1. Discovery scans plugins/* for plugin.yaml.
2. A manifest is loaded into an external collector wrapper.
3. Execution runs command + args + target.
4. Stdout is saved as raw evidence.
5. Evidence metadata is enriched from parsed JSON output.
6. Downstream ingestion and analysis consume evidence.

Relevant runtime paths:
- internal/collector/external.go
- internal/collector/registry.go
- internal/core/collector.go
- internal/core/evidence.go

## 2. Plugin Directory Layout

Minimum layout:

```text
plugins/
  my_plugin/
    plugin.yaml
    main.py
```

Recommended layout:

```text
plugins/
  my_plugin/
    plugin.yaml
    main.py
    requirements.txt
    README.md
    LICENSE
    tests/
```

## 3. Manifest Contract (plugin.yaml)

Required fields:
- name: unique collector name used by CLI.
- description: short collector summary.
- command: executable or interpreter.
- args: static argument list; target is appended by SPECTRE at runtime.
- is_active: true for active probing, false for passive collection.

Typical manifest:

```yaml
name: "dns_custom"
description: "Custom DNS intelligence collector"
version: "1.0.0"
command: "python"
args: ["main.py"]
is_active: false
author: "Example Team"
repo: "https://github.com/example/dns_custom"
tags: ["dns", "passive"]
```

Notes:
- version, author, repo, tags are optional but recommended for extension-store metadata.
- The plugin folder name does not need to match name, but name must be stable.

## 4. Runtime Execution Contract

Execution shape:
- SPECTRE runs: command args... target
- Working directory is the plugin directory.
- Stdout is treated as evidence payload.
- Stderr is used for diagnostics and failure details.

Output requirements:
- Emit exactly one valid JSON object to stdout.
- Exit code 0 on success.
- Non-zero exit code indicates failure.

If stdout is not valid JSON, raw output is still stored as evidence, but metadata enrichment will be limited.

## 5. Environment Variables Exposed to Plugins

Core execution passes environment variables to plugins.

SPECTRE-managed variables:
- SPECTRE_GHOST_MODE: 1 when ghost mode is enabled, otherwise 0.
- SPECTRE_PROXY: proxy endpoint if configured.

Behavior:
- In ghost mode, SPECTRE_PROXY defaults to socks5://127.0.0.1:9050 unless overridden.
- Outside ghost mode, SPECTRE_PROXY is set only when http.proxy is configured.

## 6. Ethics and Safety Model

Policy checks happen before plugin execution in registry flow:
1. Active consent check: active collectors require explicit active permission.
2. Scope control: blocked targets are denied.
3. Rate limiting: per-collector wait gate.

Implications for plugin authors:
- Set is_active correctly in plugin.yaml.
- Design passive collectors as default where possible.
- Do not bypass policy using alternate network channels.

## 7. Implementation Examples

### Python Example

plugin.yaml:

```yaml
name: "whois_custom"
description: "WHOIS lookup collector"
command: "python"
args: ["main.py"]
is_active: false
```

main.py:

```python
import json
import sys


def collect(target: str) -> dict:
    return {
        "target": target,
        "status": "ok",
        "source": "whois_custom",
        "registrar": "example-registrar"
    }


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("missing target", file=sys.stderr)
        sys.exit(1)

    target = sys.argv[1]
    result = collect(target)
    print(json.dumps(result))
```

### Bash Example

plugin.yaml:

```yaml
name: "http_status"
description: "Simple HTTP status probe"
command: "bash"
args: ["check.sh"]
is_active: true
```

check.sh:

```bash
#!/usr/bin/env bash
set -euo pipefail

target="$1"
code="$(curl -o /dev/null -s -w "%{http_code}" "$target")"

printf '{"target":"%s","http_status":%s}\n' "$target" "$code"
```

### Compiled Binary Example

plugin.yaml:

```yaml
name: "fast_scan"
description: "Compiled collector binary"
command: "./fast_scan"
args: []
is_active: true
```

## 8. Evidence and Metadata Enrichment

SPECTRE stores raw stdout in evidence_storage and computes a file hash.

Metadata enrichment behavior:
- Base metadata always includes target and source=external_plugin.
- If stdout parses as a JSON object, top-level keys are merged into metadata.

Guidance:
- Keep top-level keys flat and meaningful.
- Include fields that will be useful for downstream ingestion and AI synthesis.
- Avoid huge payloads in a single field; prefer concise structured values.

## 9. Testing Strategy

### Unit-level plugin testing
- Validate argument parsing.
- Validate JSON serialization.
- Validate non-zero exit on failures.

### Integration testing with SPECTRE
1. Place plugin in plugins/<name>/.
2. Run collection for a test case.
3. Confirm evidence file is generated.
4. Confirm metadata fields appear in stored evidence.

Example command:

```bash
spectre collect --case <case-id> --target example.com --scanners whois_custom
```

### Failure-mode tests
- Missing target input.
- Network timeouts.
- Invalid API credentials.
- Non-JSON stdout.

## 10. Packaging and Distribution

Recommended for publication:
- README with usage, target types, ethics profile, and data destinations.
- License file.
- Versioned releases (semantic versioning).
- Optional metadata fields in plugin.yaml for extension discovery.

Registry references:
- registry_sample.json
- docs/plugins.md

## 11. Troubleshooting Matrix

Plugin not found:
- Verify plugin.yaml exists in plugin directory.
- Verify manifest name matches scanner name.
- Restart command/session so discovery reruns.

Execution failed:
- Check stderr output from plugin.
- Validate command and args paths in plugin.yaml.
- Confirm executable permissions where applicable.

JSON parse issues:
- Ensure stdout contains only JSON.
- Send debug logs to stderr instead of stdout.

Policy blocked:
- Active collectors require active permission.
- Scope restrictions may block the target.
- Rate limiter can delay execution.

Dependency errors (Python):
- Add requirements.txt.
- Install required packages in the project environment if auto-install is unavailable.

## 12. Native Go Collectors (Core Extension Path)

For built-in or performance-critical collectors, implement the core interface directly in Go.

Reference:
- internal/core/collector.go
- internal/collector/registry.go

This guide focuses on external plugins. Use native collectors when tight integration, high throughput, or low-overhead execution is required.

## 13. Author Checklist

- Manifest is complete and correct.
- is_active classification is accurate.
- Stdout emits valid JSON object only.
- Stderr carries diagnostics.
- Exit codes are meaningful.
- README documents data collection behavior and ethics implications.
- Plugin tested through SPECTRE collect command.
