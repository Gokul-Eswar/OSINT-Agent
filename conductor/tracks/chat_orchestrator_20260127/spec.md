# Specification - Chat Orchestrator (The "Agent")

## Overview
The **Chat Orchestrator** transforms Spectre from a tool that *follows orders* into an agent that *understands intent*. It allows users to control the entire intelligence cycle using natural language.

Instead of typing:
`spectre collect --case 123 --target example.com --scanners dns,whois,ports`

The user types:
`"Check example.com for open ports and tell me who owns the domain."`

## Core Requirements

### 1. Intent Recognition (The "Brain")
*   **LLM-Driven:** Uses the existing LLM infrastructure (Ollama/Python Bridge).
*   **Tool Definitions:** The Agent must know about available capabilities:
    *   `run_collector(name, target)`
    *   `create_case(name)`
    *   `search_knowledge(query)`
*   **Prompt Engineering:** A specialized system prompt that forces the LLM to output structured "Actions" (JSON) instead of just text.

### 2. The Agent Loop (ReAct Pattern)
The system operates in a loop:
1.  **Observe:** Receive user input.
2.  **Think:** LLM decides which tool to use.
3.  **Act:** Go executes the tool (e.g., runs a port scan).
4.  **Observe:** Capture the tool output (e.g., "Port 80 Open").
5.  **Think:** LLM synthesizes the result or decides on the next step.
6.  **Response:** Final answer to the user.

### 3. Interfaces

#### A. CLI Mode (`spectre chat`)
A simple REPL (Read-Eval-Print Loop) in the terminal.
```text
spectre> investigate google.com
[Agent] Running DNS collector...
[Agent] Running Whois collector...
[Agent] Found 4 IPs and the domain is owned by Google LLC.
```

#### B. TUI Integration
A new tab in the TUI (`Chat`) that displays the conversation history and allows input.

## Success Criteria
*   Users can trigger standard collectors using natural language variations.
*   The agent correctly handles missing parameters (e.g., "Scan ports" -> "Which target?").
*   The agent provides a summary of actions taken.
