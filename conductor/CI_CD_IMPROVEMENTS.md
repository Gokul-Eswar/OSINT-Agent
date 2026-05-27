## CI/CD Pipeline Improvements for SPECTRE

This document outlines the adaptive CI/CD strategy implemented to handle your hybrid Go/Python architecture and plugin system.

### Problems with the Original Pipeline

1. **Fixed action versions**: `actions/setup-go@v5` pinned SHA could become stale
2. **No conditional execution**: All steps ran regardless of what changed
3. **Mixed concerns**: Go and Python tests ran together without isolation
4. **No plugin validation**: New collectors/plugins weren't automatically validated
5. **No config validation**: YAML config changes weren't checked
6. **No coverage reporting**: Test coverage wasn't tracked over time

### New Architecture

#### 1. **Change Detection (`detect-changes` job)**
   - Analyzes which files changed in the commit
   - Outputs flags: `go-changed`, `python-changed`, `config-changed`, `collector-added`
   - Downstream jobs use these to skip unnecessary work

   ```yaml
   if: needs.detect-changes.outputs.go-changed == 'true'
   ```

#### 2. **Go Checks (`go-checks` job)**
   - **Always runs** (Go is your core system)
   - Added `go mod tidy` verification
   - Lint step uses `golangci-lint` (flexible, modern linter)
   - Coverage gate enforces 40% minimum
   - Builds and validates the CLI binary
   - Reports coverage to Codecov

#### 3. **Python Checks (`python-checks` job)**
   - **Conditional** - only runs if Python files changed
   - Formats check with `black`
   - Lints with `flake8`
   - Runs pytest with coverage
   - Reports to Codecov

#### 4. **Plugin Validation (`plugin-checks` job)**
   - **Conditional** - only runs if plugins/collectors changed
   - Discovers new plugins automatically
   - Validates plugin Go code structure
   - Ensures plugin builds without errors

#### 5. **Config Validation (`config-validation` job)**
   - **Conditional** - only runs if YAML configs changed
   - Validates all YAML syntax
   - Catches config errors before deployment

#### 6. **Integration Tests (`integration-tests` job)**
   - Optional tier: only runs on `main` branch after successful unit tests
   - Marked `continue-on-error: true` so it doesn't block merges
   - Tests real CLI flows

#### 7. **Concurrency Control**
   ```yaml
   concurrency:
     group: ${{ github.workflow }}-${{ github.ref }}
     cancel-in-progress: true
   ```
   - Cancels previous runs on the same branch to save CI minutes

---

### How This Adapts to New Changes

#### Scenario 1: You add a new Go package
- ✅ `go-changed` flag triggers
- ✅ Linting, testing, coverage checks run
- ❌ Python tests skipped (not needed)
- ❌ Plugin validation skipped

#### Scenario 2: You add a new Python analyzer function
- ✅ `python-changed` flag triggers
- ✅ Python tests, linting run
- ✅ Go tests still run (core always tested)
- ❌ Plugin validation skipped

#### Scenario 3: You create a new plugin
- ✅ `collector-added` flag triggers
- ✅ Plugin structure validated
- ✅ Go and Python tests run
- ✅ Ensures plugin follows conventions

#### Scenario 4: You update `configs/default.yaml`
- ✅ `config-changed` flag triggers
- ✅ YAML validation runs
- ✅ All other tests run normally

---

### Recommended Additions

#### 1. Pre-commit hooks (local)
Create `.git-hooks/pre-commit`:
```bash
#!/bin/bash
go fmt ./...
go mod tidy
cd analyzer && black . && flake8 .
```

#### 2. GitHub Branch Protection
In repo settings > Branches > main:
- ✅ Require status checks to pass before merging
- ✅ Require code reviews (at least 1)
- ✅ Dismiss stale PR approvals on new commits
- ✅ Require branches to be up to date

#### 3. Codecov Integration
- Automatically reports coverage trends
- Flags PRs with coverage drops
- Add badge to README:
  ```markdown
  [![codecov](https://codecov.io/gh/Gokul-Eswar/Spectre/branch/main/graph/badge.svg)](https://codecov.io/gh/Gokul-Eswar/Spectre)
  ```

#### 4. Dependabot Integration
Create `.github/dependabot.yml`:
```yaml
version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
  - package-ecosystem: "pip"
    directory: "/analyzer"
    schedule:
      interval: "weekly"
  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "weekly"
```

---

### Testing Locally Before Commit

#### Go
```bash
go fmt ./...
go mod tidy
go vet ./...
go test -race ./...
go build -o bin/spectre ./cmd/spectre
```

#### Python
```bash
cd analyzer
pip install -r requirements.txt
black .
flake8 .
pytest -v --cov=.
```

#### Full
```bash
make test  # if you add this target
```

---

### Troubleshooting

**Q: Why is my PR failing only in CI but passing locally?**
- ✅ Try running tests with `-race` flag: `go test -race ./...`
- ✅ Ensure `go mod tidy` was run
- ✅ Check Python version (tests run on 3.11)

**Q: How do I force a full CI run?**
- Push an empty commit: `git commit --allow-empty -m "Trigger CI"`
- Or use GitHub UI to re-run workflow

**Q: My new collector isn't being validated**
- Ensure it's in `internal/collector/plugins/` or `plugins/` directory
- Ensure it has `.go` extension
- Run locally: `go vet ./internal/collector/...`

---

### Metrics to Track

1. **Build time**: Should stay < 5 min for Go + Python
2. **Coverage trend**: Aim for 50%+ on core logic
3. **Test flakiness**: If tests fail randomly, investigate race conditions
4. **CI queue**: If backing up, increase runner concurrency

