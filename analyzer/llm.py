import requests
import json

def analyze_case(data):
    """
    Synthesize case data using an LLM. 
    Currently defaults to Ollama, mirroring the Go implementation but in Python.
    """
    case_name = data.get("case_name", "Unknown")
    context = data.get("context", "")
    model = data.get("model", "llama3")
    
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
    
    # Ollama integration
    try:
        resp = requests.post(
            "http://localhost:11434/api/generate",
            json={
                "model": model,
                "prompt": full_prompt,
                "stream": False
            },
            timeout=120
        )
        resp.raise_for_status()
        raw_response = resp.json().get("response", "")
        
        # Simple JSON extraction in case of LLM noise
        start = raw_response.find("{")
        end = raw_response.rfind("}") + 1
        if start != -1 and end != -1:
            return json.loads(raw_response[start:end])
        
        raise ValueError("No valid JSON found in LLM response")
        
    except Exception as e:
        return {"error": f"LLM Analysis failed: {str(e)}"}