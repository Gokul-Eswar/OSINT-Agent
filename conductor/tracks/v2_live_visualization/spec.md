# Specification: V2 Live Visualization

## Goal
Replace static investigation reports with a real-time, interactive "Command Center" that visualizes the investigation as it happens.

## Components

### 1. Real-Time Event Stream (SSE)
*   **Backend:** Broadcaster in Go that emits `entity_created` and `relationship_created` events.
*   **Frontend:** React-based listener that updates the local state without page refreshes.

### 2. Dynamic Graph Canvas (Vis.js/D3)
*   **Physics-Engine:** Smooth node placement using ForceAtlas2.
*   **Live Updates:** New nodes and edges should "pop" into existence with animations.
*   **Interactive Inspect:** Clicking a node opens a side-panel with the full evidence chain and LLM analysis.

### 3. Agent Integration
*   **Chat Overlay:** Integrated agent chat for triggering collectors directly from the visualization.
*   **Status Indicators:** Visual cues when a collector is active on a specific node (e.g., a pulsing glow).

## Success Criteria
- [x] New entities appear on the graph within 500ms of being saved to DB.
- [x] Relationships are drawn automatically between existing and new nodes.
- [x] Dashboard remains responsive with 500+ nodes.
