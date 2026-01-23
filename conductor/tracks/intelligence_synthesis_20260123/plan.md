# Implementation Plan - Intelligence Synthesis (AI Analysis)

**Goal:** Integrate Local LLMs (via Ollama/OpenAI API) to analyze investigation data and generate structured intelligence reports.

## Phase 1: AI Provider & Core Models
- [x] **Task: Define Analysis Structs**
    - [x] Create `internal/core/analysis.go` (Findings, Risks, NextSteps).
    - [x] Ensure alignment with the `analyses` table in the schema.
- [x] **Task: AI Provider Interface**
    - [x] Create `internal/ai/provider.go` (Interface: `Generate(prompt string) (string, error)`).
    - [x] Create `internal/ai/ollama.go` (Implementation using standard HTTP client).

## Phase 2: Context Construction
- [x] **Task: Context Builder Service**
    - [x] Create `internal/analysis/context.go`.
    - [x] Implement `BuildCaseContext(caseID string)` -> aggregates Name, Description, Entities, Relationships, and Evidence into a structured prompt.

## Phase 3: Analysis Engine & Storage
- [x] **Task: Analysis Logic**
    - [x] Create `internal/analysis/engine.go`.
    - [x] Implement `AnalyzeCase(caseID string)`:
        1. Build Context.
        2. Send to AI with a system prompt enforcing JSON output.
        3. Parse JSON response into `core.Analysis`.
- [x] **Task: Analysis Storage**
    - [x] Create `internal/storage/analysis_repo.go` (Save, GetLatest).

## Phase 4: CLI Integration
- [x] **Task: Analyze Command**
    - [x] Create `internal/cli/analyze.go`.
    - [x] Implement `spectre analyze --case <id> --model <name>`.
    - [x] Display the summary/risks to the user.

## Phase 5: Verification
- [x] **Task: Mock/Live Test**
    - [x] Test against a local Ollama instance (or mock if unavailable) using the existing "Graph-Test" case.
