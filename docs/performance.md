# Performance Baseline

This note captures a repeatable baseline for CLI startup, collector benchmarks, and one representative investigation flow.

## Measurement Conditions

- Timestamp: 2026-04-05T12:27:42.9702258+05:30
- Workspace: D:\GOKUL_ESWAR\Codebase\CLI_AGENNT
- Binary: .\bin\spectre.exe
- Host: GOKULESWAR
- OS: Microsoft Windows 11 Home Single Language 10.0.26200
- CPU: AMD Ryzen 5 5600H with Radeon Graphics
- RAM: 15.4 GB
- Go: go1.25.4 windows/amd64
- Plugin set present during measurement: built-in collectors plus plugins/echo_test

## Repeatable Commands

Build once, then run the same binary for each measurement:

```powershell
go build -o .\bin\spectre.exe .\cmd\spectre
```

CLI startup latency commands:

```powershell
.\bin\spectre.exe version
.\bin\spectre.exe --help
.\bin\spectre.exe init
$env:GO_TESTING='true'; .\bin\spectre.exe collect dns example.com --dry-run --case perf-baseline; Remove-Item Env:GO_TESTING -ErrorAction SilentlyContinue
```

Collector benchmarks:

```powershell
go test -run '^$' -bench 'BenchmarkCollectorConcurrency|BenchmarkHTTPCollector|BenchmarkPortCollector_Collect' -benchmem ./internal/collector ./internal/collector/active
```

## CLI Baseline

| Timestamp | Machine | Command | Options / Target | Time |
| --- | --- | --- | --- | ---: |
| 2026-04-05T12:27:42.9702258+05:30 | GOKULESWAR / Windows 11 Home Single Language / Ryzen 5 5600H / 15.4 GB | version | `.in\spectre.exe version` | 144.89 ms |
| 2026-04-05T12:27:42.9702258+05:30 | GOKULESWAR / Windows 11 Home Single Language / Ryzen 5 5600H / 15.4 GB | help | `.in\spectre.exe --help` | 65.02 ms |
| 2026-04-05T12:27:42.9702258+05:30 | GOKULESWAR / Windows 11 Home Single Language / Ryzen 5 5600H / 15.4 GB | init | `.in\spectre.exe init` | 89.97 ms |
| 2026-04-05T12:27:42.9702258+05:30 | GOKULESWAR / Windows 11 Home Single Language / Ryzen 5 5600H / 15.4 GB | collect dry-run | `GO_TESTING=true .\bin\spectre.exe collect dns example.com --dry-run --case perf-baseline` | 68.33 ms |

## Collector Benchmarks

Benchmark output was captured with `-benchmem` so allocations remain comparable across future runs.

### internal/collector

| Benchmark | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| BenchmarkCollectorConcurrency/Concurrency-1-12 | 1029016500 | 130328 | 1024 |
| BenchmarkCollectorConcurrency/Concurrency-5-12 | 207452320 | 72110 | 825 |
| BenchmarkCollectorConcurrency/Concurrency-10-12 | 105579040 | 70463 | 816 |
| BenchmarkCollectorConcurrency/Concurrency-50-12 | 21486509 | 68020 | 807 |
| BenchmarkCollectorConcurrency/Concurrency-100-12 | 10873397 | 67139 | 806 |
| BenchmarkHTTPCollector/HTTP-Concurrency-5-12 | 9554681800 | 366424 | 4027 |
| BenchmarkHTTPCollector/HTTP-Concurrency-20-12 | 2279308900 | 352072 | 3841 |

### internal/collector/active

| Benchmark | ns/op | B/op | allocs/op |
| --- | ---: | ---: | ---: |
| BenchmarkPortCollector_Collect-12 | 2012340300 | 110840 | 506 |

## Representative Investigation Flow

The dry-run collect path is the safest repeatable investigation smoke test because it exercises CLI parsing, context loading, collector selection, and dry-run handling without making network requests.

| Timestamp | Machine | Flow | Command | Time |
| --- | --- | --- | --- | ---: |
| 2026-04-05T12:27:42.9702258+05:30 | GOKULESWAR / Windows 11 Home Single Language / Ryzen 5 5600H / 15.4 GB | investigation smoke flow | `GO_TESTING=true .\bin\spectre.exe collect dns example.com --dry-run --case perf-baseline` | 68.33 ms |

## Re-run Rules

Keep future runs comparable by holding these inputs steady:
- Use the same built binary path unless explicitly comparing build changes.
- Keep the same command set and arguments.
- Keep the same target set (`example.com` for the smoke flow).
- Keep the same plugin set present in the workspace when measuring startup.
- Run `init` from a clean workspace if you want cold-start database timings.
- Use `-benchmem` for all benchmark comparisons.

## Notes

- The HTTP collector benchmark currently measures the localhost fallback/failure path when no local service is listening. Treat it as a collector-overhead baseline rather than a successful remote fetch benchmark.
- Startup measurements include normal CLI initialization and plugin discovery for the current workspace.