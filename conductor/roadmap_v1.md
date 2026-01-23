# SPECTRE v1 Roadmap - End-to-End Plan

This document outlines the strategic tracks required to deliver the full SPECTRE v1 product capability, transitioning from the current foundational prototype to a complete intelligence synthesis tool.

## Current Status
- **[x] Track 1: Core Foundation** (Completed)
  - *Deliverables:* Project structure, CLI scaffold, SQLite setup, Basic Case Management.

---

## Planned Execution Tracks

### Track 2: Core Data Engine (The "Graph")
**Goal:** Enable the storage and manual management of complex investigation data (Entities & Relationships).
**Key Deliverables:**
1.  **Domain Models:** Implement `Entity` (Type, Value, Confidence) and `Relationship` (Source, Target, Type) structs.
2.  **Storage Access:** SQLite repositories for creating and querying the graph nodes/edges.
3.  **Manual CLI:** Commands to manually build the graph.
    - `spectre entity add --type ip --value 1.1.1.1`
    - `spectre link --from 1.1.1.1 --to cloudflare.com --type "resolves_to"`
**Why:** This is the prerequisite for storing anything useful from collectors.

### Track 3: Passive Collection Framework (The "Ears")
**Goal:** Automate the gathering of OSINT data to populate the graph.
**Key Deliverables:**
1.  **Collector Interface:** A Go interface for plugins (`Collect(target) -> []Evidence`).
2.  **Plugin Registry:** System to register and select collectors.
3.  **Initial Collectors:**
    - `dns`: Resolve A, MX, NS, TXT records.
    - `whois`: Fetch registrar and age data.
4.  **Ingestion Pipeline:** Convert raw Collector output -> Evidence Records -> Entities/Relationships.
**Why:** Replaces manual entry with automated discovery.

### Track 4: Intelligence Synthesis (The "Brain")
**Goal:** Integrate Local LLMs to analyze data and find patterns.
**Key Deliverables:**
1.  **AI Provider Adapter:** Interface for LLM backends (Ollama, OpenAI-compatible).
2.  **Analysis Engine:** Logic to serialize case data into a prompt context.
3.  **CLI Command:** `spectre analyze` to generate summaries and risk assessments.
**Why:** Delivers the "Synthesis" promise of the product definition.

### Track 5: Visual Intelligence Dashboard (The "Eyes")
**Goal:** Provide an interactive way to explore the investigation graph.
**Key Deliverables:**
1.  **Static Generator:** Go code to bundle case JSON into a single HTML file.
2.  **Frontend Assets:** Embedded Vis.js/D3.js for graph visualization.
3.  **Timeline View:** Visual sequence of discovered events.
4.  **CLI Command:** `spectre visualize` (opens default browser).
**Why:** Text-based graph exploration is inefficient for complex cases.

---

## Execution Strategy
- Tracks will be executed sequentially to maintain stability.
- "Conductor" checks will be performed at the start and end of each track.
- `master` branch will remain buildable at all times.
