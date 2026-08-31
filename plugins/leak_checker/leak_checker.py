import sys
import os
import requests
from datetime import datetime, timezone

# Add project root to path to find spectre_sdk
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../../")))
from lib.python.spectre_sdk import BaseCollector

class LeakChecker(BaseCollector):
    def __init__(self):
        super().__init__(
            name="leak_checker",
            description="Checks target for exposed sensitive files (.env, .git, backups)."
        )

    def collect(self, target):
        self.log(f"Starting leak check for {target}...")
        
        proxies = None
        if self.proxy:
            proxies = {"http": self.proxy, "https": self.proxy}

        sensitive_paths = [
            (".env", "Environment Config"),
            (".git/config", "Git Repository Config"),
            ("config.json", "Application Configuration"),
            ("backup.sql", "Database Backup"),
            (".htaccess", "Apache Access Rules")
        ]

        base_url = target if target.startswith(("http://", "https://")) else f"http://{target}"
        potential_leaks = []

        for path, leak_type in sensitive_paths:
            url = f"{base_url.rstrip('/')}/{path}"
            try:
                resp = requests.head(url, proxies=proxies, timeout=5, allow_redirects=True)
                status = "exposed" if resp.status_code == 200 else "restricted"
                potential_leaks.append({
                    "type": leak_type,
                    "url": url,
                    "status": status,
                    "status_code": resp.status_code
                })
            except Exception as e:
                potential_leaks.append({
                    "type": leak_type,
                    "url": url,
                    "status": "unreachable",
                    "error": str(e)
                })

        return {
            "target": target,
            "scan_time": datetime.now(timezone.utc).isoformat(),
            "potential_leaks": potential_leaks,
            "ghost_mode_active": self.ghost_mode,
            "proxy_used": self.proxy
        }

if __name__ == "__main__":
    collector = LeakChecker()
    collector.run()

