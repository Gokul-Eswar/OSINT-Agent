# Implementation Plan - Chat Orchestrator

## Phase 1: The Agent Engine
**Goal:** Build the backend logic to parse intent and execute tools.

- [ ] **Define Tool Registry:**
    - Create `internal/agent/tools.go` to map strings ("dns", "scan") to Go functions.
    - Define JSON schemas for each tool for the LLM.

- [ ] **Python Bridge Update:**
    - Update `analyzer/llm.py` (or create `analyzer/agent.py`) to handle "chat" tasks.
    - Implement a "Tool Use" system prompt.

- [ ] **The Agent Loop (Go):**
    - Create `internal/agent/engine.go`.
    - Implement `agent.Execute(input string) (string, error)`.
    - Handle the `User -> LLM -> Tool -> LLM -> User` flow.

## Phase 2: CLI Interface
**Goal:** Quick way to test the agent.

- [ ] **Command:** Implement `cmd/spectre/chat.go`.
- [ ] **REPL:** Build a simple loop reading from Stdin.

## Phase 3: TUI Integration
**Goal:** Visual chat within the dashboard.

- [ ] **Bubble Tea Model:** Create `internal/tui/chat.go`.
- [ ] **UI Components:**
    - `viewport` for history.
    - `textinput` for user messages.
    - `spinner` for "Agent is thinking...".

## Phase 4: Web Integration (Future)
- [ ] Add `/api/chat` endpoint.
- [ ] Add Chat Widget to React frontend.
