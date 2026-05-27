import os
import sys

def validate_plugins():
    plugins_dir = "plugins"
    if not os.path.exists(plugins_dir):
        print(f"Plugins directory '{plugins_dir}' does not exist.")
        return 0

    try:
        import yaml
    except ImportError:
        print("Error: 'yaml' module (PyYAML) is required. Please install it using 'pip install pyyaml'.")
        return 1

    errors = []
    
    for entry in os.scandir(plugins_dir):
        if entry.is_dir():
            plugin_yaml = os.path.join(entry.path, "plugin.yaml")
            if not os.path.exists(plugin_yaml):
                # Check if this subdirectory is empty or contains code files
                # Some folders like __pycache__ or other metadata might exist and we want to ignore them if they aren't real plugin dirs.
                # However, if there are files (like python files) we warn or check.
                files = os.listdir(entry.path)
                if not entry.name.startswith((".", "__")) and files:
                    errors.append(f"Directory '{entry.path}' contains files but is missing 'plugin.yaml'")
                continue

            print(f"Validating plugin manifest: {plugin_yaml}")
            try:
                with open(plugin_yaml, "r", encoding="utf-8") as f:
                    data = yaml.safe_load(f)
            except Exception as e:
                errors.append(f"Invalid YAML in {plugin_yaml}: {e}")
                continue

            if not isinstance(data, dict):
                errors.append(f"Plugin manifest in {plugin_yaml} must be a YAML dictionary")
                continue

            required_fields = ["name", "description", "command", "args"]
            for field in required_fields:
                if field not in data:
                    errors.append(f"Missing required field '{field}' in {plugin_yaml}")
            
            if "command" in data and "args" in data:
                command = data["command"]
                args = data["args"]
                if not isinstance(args, list):
                    errors.append(f"'args' must be a list in {plugin_yaml}")
                
                # Check if referenced scripts exist
                if command in ["python", "python3", "bash", "sh"] and args:
                    script_file = os.path.join(entry.path, args[0])
                    if not os.path.exists(script_file):
                        errors.append(f"Referenced script file '{script_file}' does not exist for plugin '{data.get('name', entry.name)}'")

    if errors:
        print("\nPlugin Validation Failures:")
        for err in errors:
            print(f"- {err}")
        return 1

    print("\nAll plugins validated successfully!")
    return 0

if __name__ == "__main__":
    sys.exit(validate_plugins())
