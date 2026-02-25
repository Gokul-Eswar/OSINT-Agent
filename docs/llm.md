# SPECTRE Intelligence Layer (AI)

SPECTRE uses Large Language Models (LLMs) to synthesize collected evidence, identify risks, and provide an interactive chat interface for investigation orchestration.

## Backend Options

### 1. Ollama (Recommended / Default)
Ollama is the preferred backend for SPECTRE because it runs locally, ensuring that sensitive investigation data never leaves your machine.

**Setup:**
1. Download and install Ollama from [ollama.com](https://ollama.com).
2. Pull the default model:
   ```bash
   ollama pull llama3
   ```
3. Ensure the Ollama server is running (usually on port 11434).

**Configuration:**
SPECTRE looks for Ollama at `http://localhost:11434` by default. You can change this in `configs/default.yaml`:
```yaml
ai:
  model: "llama3"
  url: "http://localhost:11434/api"
```

### 2. OpenAI / Compatible APIs
You can use OpenAI's GPT-4 or any OpenAI-compatible API (like LocalAI or vLLM).

**Configuration:**
```yaml
ai:
  model: "gpt-4"
  url: "https://api.openai.com/v1"
  api_key: "your-api-key-here"
```

## How It Works

1. **Synthesize:** When you run `spectre analyze`, the Go core extracts entities and evidence from the SQLite database and passes them to the Python layer. The Python layer prompts the LLM to generate a structured JSON report.
2. **Chat Agent:** The `spectre chat` command (and Web UI) uses an agentic loop. The LLM can "decide" to use tools (collectors) by responding with a specific JSON format:
   ```json
   {"tool_use": {"name": "dns", "arguments": {"target": "example.com"}}}
   ```

## Fallback Mode
If the LLM backend is unavailable, SPECTRE will:
- Return a structured "Synthesis Unavailable" report instead of crashing.
- Provide basic greeting and error guidance in the chat interface.
- Allow manual evidence review via the CLI and Web Dashboard.

## Troubleshooting
- **Connection Refused:** Ensure Ollama is running (`ollama serve`).
- **Timeouts:** Complex cases with lots of evidence may take longer to process. You can increase the `timeout` in the config.
- **Garbage Output:** Ensure you are using a capable model like `llama3`, `mistral`, or `gpt-4`. Smaller models may fail to produce the required JSON format.
