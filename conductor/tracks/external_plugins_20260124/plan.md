# Implementation Plan - External Plugin System

## Phase 1: Plugin Loader
- [x] **Task: Discovery Logic**
    - Create `internal/collector/external.go`.
    - Implement logic to walk the `plugins/` directory and parse metadata.
- [x] **Task: Registry Integration**
    - Automatically register discovered plugins into the global collector registry.

## Phase 2: Execution Wrapper
- [x] **Task: External Runner**
    - Implement the `Collect(caseID, target)` method for external plugins.
    - Handle subprocess execution and capture JSON output.
- [x] **Task: Normalization**
    - Parse the plugin's output and convert it into `core.Evidence`.

## Phase 3: Developer Experience
- [x] **Task: Example Plugin**
    - Create `plugins/echo_test/` with a simple Python script.
    - Provide a `plugin.yaml` template.

## Phase 4: Verification
- [x] **Task: Test**
    - Run the echo_test plugin and verify ingestion.
