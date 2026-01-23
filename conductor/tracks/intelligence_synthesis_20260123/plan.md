# Implementation Plan - Intelligence Synthesis (AI Analysis)

**Goal:** Integrate Local LLMs (via Ollama/OpenAI API) to analyze investigation data and generate structured intelligence reports.

## Phase 1: AI Provider & Core Models
- [ ] **Task: Define Analysis Structs**
    - Create `internal/core/analysis.go` (Findings, Risks, NextSteps).
    - Ensure alignment with the `analyses` table in the schema.
- [ ] **Task: AI Provider Interface**
    - Create `internal/ai/provider.go` (Interface: `Generate(prompt string) (string, error)`).
    - Create `internal/ai/ollama.go` (Implementation using standard HTTP client).

## Phase 2: Context Construction
- [ ] **Task: Context Builder Service**
    - Create `internal/analysis/context.go`.
    - Implement `BuildCaseContext(caseID string)` -> aggregates Name, Description, Entities, Relationships, and Evidence into a structured prompt.

## Phase 3: Analysis Engine & Storage
- [ ] **Task: Analysis Logic**
    - Create `internal/analysis/engine.go`.
    - Implement `AnalyzeCase(caseID string)`:
        1. Build Context.
        2. Send to AI with a system prompt enforcing JSON output.
        3. Parse JSON response into `core.Analysis`.
- [ ] **Task: Analysis Storage**
    - Create `internal/storage/analysis_repo.go` (Save, GetLatest).

## Phase 4: CLI Integration
- [ ] **Task: Analyze Command**
    - Create `internal/cli/analyze.go`.
    - Implement `spectre analyze --case <id> --model <name>`.
    - Display the summary/risks to the user.

## Phase 5: Verification
- [ ] **Task: Mock/Live Test**
    - Test against a local Ollama instance (or mock if unavailable) using the existing "Graph-Test" case.
