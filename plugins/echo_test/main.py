import sys
import json
import time

def main():
    if len(sys.argv) < 2:
        print(json.dumps({"error": "No target provided"}))
        sys.exit(1)

    target = sys.argv[1]
    
    result = {
        "plugin": "echo_test",
        "target": target,
        "timestamp": int(time.time()),
        "message": f"Successfully echoed target: {target}",
        "status": "ok"
    }
    
    print(json.dumps(result, indent=2))

if __name__ == "__main__":
    main()
