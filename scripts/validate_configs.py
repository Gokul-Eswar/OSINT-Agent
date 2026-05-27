import os
import sys

def validate_configs():
    configs_dir = "configs"
    if not os.path.exists(configs_dir):
        print(f"Configs directory '{configs_dir}' does not exist.")
        return 0

    # We try to import yaml. If not available, we warn and skip or install it.
    try:
        import yaml
    except ImportError:
        print("Error: 'yaml' module (PyYAML) is required. Please install it using 'pip install pyyaml'.")
        return 1

    errors = []
    
    for root, _, files in os.walk(configs_dir):
        for file in files:
            if file.endswith((".yaml", ".yml")):
                config_path = os.path.join(root, file)
                print(f"Validating YAML config: {config_path}")
                try:
                    with open(config_path, "r", encoding="utf-8") as f:
                        yaml.safe_load(f)
                except Exception as e:
                    errors.append(f"Invalid YAML in {config_path}: {e}")

    if errors:
        print("\nConfig Validation Failures:")
        for err in errors:
            print(f"- {err}")
        return 1

    print("\nAll config files validated successfully!")
    return 0

if __name__ == "__main__":
    sys.exit(validate_configs())
