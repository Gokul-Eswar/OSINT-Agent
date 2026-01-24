# Implementation Plan - Geo-Intelligence Enrichment

## Phase 1: The Collector
- [x] **Task: GeoIP Collector**
    - Create `internal/collector/geo/geoip.go`.
    - Implement `Collect(target)` where target is an IP.
    - Use `http://ip-api.com/json/{ip}`.
- [x] **Task: Registration**
    - Register the collector in `internal/collector/registry.go`.
    - Add rate limit (45/min ≈ 0.75/sec) to `configs/default.yaml`.

## Phase 2: Auto-Enrichment Logic
- [x] **Task: Event Hook**
    - Modify `internal/cli/collect.go` (or a new pipeline manager).
    - When a DNS collector produces an IP entity, automatically queue a GeoIP collection for it.
    - *Simplification for v1:* Just allow running `spectre collect geo <ip>` manually, or chain it in the `dns` collector.

## Phase 3: Visualization Update
- [x] **Task: Python Visualizer**
    - Update `analyzer/graph_viz.py`.
    - Read `metadata.country_code`.
    - Append flag emoji to the node label.

## Phase 4: Verification
- [x] **Task: Test**
    - Resolve `google.com` -> Get IP.
    - Run `spectre collect geo <IP>`.
    - Verify `spectre visualize` shows the flag.
