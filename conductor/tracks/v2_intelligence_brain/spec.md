# Specification: V2 Intelligence Brain

## Goal
Transform the SPECTRE Agent from a manual tool-triggering bot into an autonomous investigator that can "read" its findings, reason over evidence, and suggest next steps.

## Components

### 1. Evidence Vector Store (Local-First)
*   **Engine:** Integration with a local vector database (e.g., `ChromaDB` via the Python analyzer).
*   **Indexing:** Automatic background indexing of all `evidence_storage/` files (JSON, TXT, MD).
*   **Embeddings:** Use a local embedding model (e.g., `all-MiniLM-L6-v2`) via Ollama or Sentence-Transformers.

### 2. Recursive Agent "Think-Step" Loop
*   **Memory:** Contextual memory of past tool results and their analysis.
*   **Strategy Planning:** The agent must generate a "Plan" before execution.
*   **Recursion:** The agent can decide to run follow-up collectors based on findings without user intervention (bounded by a limit).

### 3. Analyst Toolset
*   **`search_evidence`:** Semantic search across all evidence files in the case.
*   **`analyze_document`:** Targeted LLM extraction from a specific evidence file.
*   **`update_hypotheses`:** A way for the agent to track and present "Intelligence Leads" or "Hypotheses" to the user.

## Success Criteria
- [x] Agent can answer: "Based on the WHOIS records we found, who is the most likely registrant and where are they located?"
- [x] Agent can automatically follow a lead (e.g., finding an email in DNS and then checking that email on GitHub).
- [x] Zero-cloud dependency for indexing and reasoning.
