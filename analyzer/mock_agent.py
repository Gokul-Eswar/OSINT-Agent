import json
import sys
import argparse

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--task")
    parser.add_argument("--input")
    args = parser.parse_args()
    
    input_data = json.loads(args.input)
    messages = input_data.get("messages", [])
    
    # Check if we just received a tool result
    if messages and messages[-1]["role"] == "system" and "Observation from 'collect'" in messages[-1]["content"]:
        print(json.dumps({
            "role": "assistant",
            "content": "I have finished the DNS collection for google.com. It was successful."
        }))
        return

    # Simple logic to simulate an agent loop
    last_user_msg = ""
    for m in reversed(messages):
        if m["role"] == "user":
            last_user_msg = m["content"].lower()
            break
            
    if "run dns" in last_user_msg:
        # Simulate tool use
        print(json.dumps({
            "role": "assistant",
            "tool_use": {
                "name": "collect",
                "arguments": {"collector": "dns", "target": "google.com"}
            }
        }))
    else:
        # Default response
        print(json.dumps({
            "role": "assistant",
            "content": "Hello! I am the SPECTRE mock agent. How can I help?"
        }))

if __name__ == "__main__":
    main()
