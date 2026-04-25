import json
import os
import sys
import argparse

class BaseCollector:
    """
    Standard base class for SPECTRE Python plugins.
    """
    def __init__(self, name, description):
        self.name = name
        self.description = description
        self.ghost_mode = os.environ.get("SPECTRE_GHOST_MODE") == "1"
        self.proxy = os.environ.get("SPECTRE_PROXY")
        self.target = None

    def setup_args(self):
        parser = argparse.ArgumentParser(description=self.description)
        parser.add_argument("target", help="The target for collection")
        return parser.parse_args()

    def run(self):
        """Main execution entry point."""
        args = self.setup_args()
        self.target = args.target
        
        try:
            results = self.collect(self.target)
            self.output(results)
        except Exception as e:
            self.error(str(e))

    def collect(self, target):
        """Override this method to implement collection logic."""
        raise NotImplementedError("Collectors must implement collect()")

    def output(self, data):
        """Standard JSON output to stdout."""
        print(json.dumps(data, indent=2))

    def error(self, message):
        """Standard error reporting."""
        print(json.dumps({"error": message}), file=sys.stderr)
        sys.exit(1)

    def log(self, message):
        """Optional logging to stderr."""
        print(f"[{self.name}] {message}", file=sys.stderr)
