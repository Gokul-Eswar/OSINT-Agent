# Specification - Advanced Active Reconnaissance

## 1. Social Collector
**Goal:** Identify social media accounts associated with a target (usually a username derived from domain or manual input).

**Requirements:**
- Support at least 10 major platforms (GitHub, Twitter, etc.).
- Respect rate limits.
- **Testability:** Must be testable without external network access.

## 2. Screenshot Collector
**Goal:** Capture visual state of the target website.

**Requirements:**
- Use Headless Chrome (`chromedp`).
- Support Proxy/Tor (Ghost Mode).
- Save as PNG.
- Store metadata (resolution, size).
