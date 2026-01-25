import requests
import json
import re
import time
import sys

def extract_json(text):
    """Robustly extract JSON from text using regex."""
    if not text:
        return None
    text = text.strip()
    # Try to find the outermost {}
    match = re.search(r'\{.*\}', text, re.DOTALL)
    if match:
        try:
            return json.loads(match.group())
        except json.JSONDecodeError:
            pass
    return None

def analyze_case(data):
    """
    Synthesize case data using an LLM. 
    Accepts configuration for LLM provider from Go.
    """
    case_name = data.get("case_name", "Unknown")
    context = data.get("context", "")
    model = data.get("model", "llama3")
    llm_config = data.get("llm_config", {})

    api_url = llm_config.get("url", "http://localhost:11434/api/generate")
    api_key = llm_config.get("api_key", "")
    timeout = llm_config.get("timeout", 120)
    
    headers = {"Content-Type": "application/json"}
    if api_key:
        headers["Authorization"] = f"Bearer {api_key}"
    
    system_prompt = (
        "You are SPECTRE, an expert intelligence analyst. "
        "Analyze the provided case data and generate a structured report. "
        "Your output MUST be strict JSON.\n"
        "Format:\n"
        "{\n"
        "  \"findings\": [\"string\"],\n"
        "  \"risks\": [\"string\"],\n"
        "  \"connections\": [\"string\"],\n"
        "  \"next_steps\": [\"string\"],\n"
        "  \"confidence\": 0.85\n"
        "}"
    )
    
    full_prompt = f"{system_prompt}\n\nCASE DATA:\n{context}"
    
    # Payload adaptation could be improved for OpenAI vs Ollama
    # For now, sticking to Ollama default but flexible URL
    payload = {
        "model": model,
        "prompt": full_prompt,
        "stream": False
    }

    retries = 3
    last_error = ""

    for attempt in range(retries):
        try:
            resp = requests.post(
                api_url,
                json=payload,
                headers=headers,
                timeout=timeout
            )
            resp.raise_for_status()
            
            response_json = resp.json()
            # Handle Ollama (response) vs OpenAI (choices[0].message.content)
            raw_response = response_json.get("response", "")
            if not raw_response and "choices" in response_json:
                 raw_response = response_json["choices"][0]["message"]["content"]
            
            extracted = extract_json(raw_response)
            if extracted:
                return extracted
            
            last_error = "No valid JSON found in LLM response"
            # print(f"Debug: raw response: {raw_response}", file=sys.stderr)

        except Exception as e:
            last_error = str(e)
            
        if attempt < retries - 1:
            time.sleep(1)
        
    return {"error": f"LLM Analysis failed after {retries} attempts. Last error: {last_error}"}