import requests
import json
import re
import time
import sys

# Global session for connection pooling and latency reduction
session = requests.Session()

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

def chat(data):
    """
    Handle an interactive chat session with tool use capabilities.
    """
    messages = data.get("messages", [])
    tools = data.get("tools", [])
    model = data.get("model", "llama3")
    llm_config = data.get("llm_config", {})

    api_url = llm_config.get("url", "http://localhost:11434/api/chat")
    api_key = llm_config.get("api_key", "")
    timeout = llm_config.get("timeout", 120)

    headers = {"Content-Type": "application/json"}
    if api_key:
        headers["Authorization"] = f"Bearer {api_key}"

    # System prompt for agent behavior
    system_msg = {
        "role": "system",
        "content": (
            "You are SPECTRE Agent, an OSINT automation assistant. "
            "You help users gather intelligence by using available tools. "
            "When you need to use a tool, respond with ONLY a JSON object in this format:\n"
            "{\"tool_use\": {\"name\": \"tool_name\", \"arguments\": { ... }}}\n"
            "Available tools:\n" + json.dumps(tools, indent=2) + "\n"
            "If you have enough information to answer the user, just provide a normal text response."
        )
    }

    # Ensure system message is at the start
    if not messages or messages[0].get("role") != "system":
        messages.insert(0, system_msg)

    # Payload for Chat API (Ollama style by default)
    payload = {
        "model": model,
        "messages": messages,
        "stream": False
    }

    try:
        resp = session.post(
            api_url,
            json=payload,
            headers=headers,
            timeout=timeout
        )
        resp.raise_for_status()
        
        response_json = resp.json()
        
        # Handle different API response structures
        content = ""
        if "message" in response_json:
            content = response_json["message"].get("content", "")
        elif "choices" in response_json:
            content = response_json["choices"][0]["message"].get("content", "")
        elif "response" in response_json:
            content = response_json["response"]

        # Check for tool use in the content
        tool_call = extract_json(content)
        if tool_call and "tool_use" in tool_call:
            return {"role": "assistant", "tool_use": tool_call["tool_use"]}
        
        return {"role": "assistant", "content": content}

    except Exception as e:
        # Fallback for chat
        last_user_msg = ""
        for m in reversed(messages):
            if m.get("role") == "user":
                last_user_msg = m.get("content", "").lower()
                break
        
        if "hello" in last_user_msg or "hi" in last_user_msg:
            return {"role": "assistant", "content": "Hello! (Fallback Mode) I'm currently unable to reach my LLM backend, but I can still help with basic commands."}
        
        return {
            "role": "assistant", 
            "content": f"I'm sorry, I encountered an error and the LLM is unavailable: {str(e)}. Please check if Ollama is running."
        }

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
        "  \"missing_data\": [\"string\"],\n"
        "  \"suggested_collectors\": [\"string\"],\n"
        "  \"confidence\": 0.85\n"
        "}"
    )
    
    full_prompt = f"{system_prompt}\n\nCASE DATA:\n{context}"
    
    # Payload adaptation could be improved for OpenAI vs Ollama
    payload = {
        "model": model,
        "prompt": full_prompt,
        "stream": False
    }

    retries = 3
    last_error = ""

    for attempt in range(retries):
        try:
            resp = session.post(
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

        except Exception as e:
            last_error = str(e)
            
        if attempt < retries - 1:
            time.sleep(1)
                
    # Final fallback for analyze_case
    return {
        "findings": ["LLM synthesis unavailable. Raw data review required."],
        "risks": ["Unable to perform AI-assisted risk assessment."],
        "connections": [],
        "next_steps": ["Check LLM backend connectivity.", "Manually review collected evidence."],
        "confidence": 0.0,
        "error": last_error
    }

def query_case(data):
    """
    Answer a specific question about a case using the LLM.
    """
    case_name = data.get("case_name", "Unknown")
    context = data.get("context", "")
    question = data.get("data", "")
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
        "Use the provided case context to answer the user's specific question. "
        "Be concise, professional, and focus on the evidence. "
        "If the information is not in the context, state that clearly."
    )
    
    full_prompt = f"{system_prompt}\n\nCASE CONTEXT:\n{context}\n\nQUESTION: {question}"
    
    payload = {
        "model": model,
        "prompt": full_prompt,
        "stream": False
    }

    try:
        resp = session.post(
            api_url,
            json=payload,
            headers=headers,
            timeout=timeout
        )
        resp.raise_for_status()
        
        response_json = resp.json()
        raw_response = response_json.get("response", "")
        if not raw_response and "choices" in response_json:
             raw_response = response_json["choices"][0]["message"]["content"]
        
        return {"answer": raw_response}

    except Exception as e:
        return {"answer": f"Error querying LLM: {str(e)}"}

def analyze_image(data):
    """
    Perform visual analysis on an image using a vision-capable LLM (e.g., llava).
    """
    model = data.get("model", "llava") # Default to llava for vision
    prompt = data.get("data", "Describe this image in detail. Focus on text, logos, or identifying features.")
    image_base64 = data.get("context", "") # We use context field to pass base64
    llm_config = data.get("llm_config", {})

    api_url = llm_config.get("url", "http://localhost:11434/api/generate")
    api_key = llm_config.get("api_key", "")
    timeout = llm_config.get("timeout", 120)
    
    headers = {"Content-Type": "application/json"}
    if api_key:
        headers["Authorization"] = f"Bearer {api_key}"
    
    payload = {
        "model": model,
        "prompt": prompt,
        "images": [image_base64],
        "stream": False
    }

    try:
        resp = session.post(
            api_url,
            json=payload,
            headers=headers,
            timeout=timeout
        )
        resp.raise_for_status()
        
        response_json = resp.json()
        raw_response = response_json.get("response", "")
        if not raw_response and "choices" in response_json:
             raw_response = response_json["choices"][0]["message"]["content"]
        
        return {"answer": raw_response}

    except Exception as e:
        return {"answer": f"Error performing visual analysis: {str(e)}"}
