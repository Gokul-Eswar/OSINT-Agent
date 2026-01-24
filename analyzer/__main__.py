import json
import sys
import argparse
from .llm import analyze_case
from .graph_viz import generate_visual_report

def main():
    parser = argparse.ArgumentParser(description="SPECTRE Analyzer (Python)")
    parser.add_argument("--task", choices=["synthesize", "visualize"], required=True)
    parser.add_argument("--input", help="JSON input data", required=True)
    
    args = parser.parse_args()
    
    try:
        input_data = json.loads(args.input)
        
        if args.task == "synthesize":
            result = analyze_case(input_data)
            print(json.dumps(result))
        elif args.task == "visualize":
            # Extract the actual graph data payload
            graph_data = input_data.get("data", {})
            result = generate_visual_report(graph_data)
            print(json.dumps(result))
            
    except Exception as e:
        print(json.dumps({"error": str(e)}), file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()
