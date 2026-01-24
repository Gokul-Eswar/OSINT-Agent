# Specification - Geo-Intelligence Enrichment

## Overview
Intelligence is spatial. Knowing where a server resides (e.g., a bulletproof host in a specific jurisdiction vs. a corporate HQ) is critical context.

## Requirements

### 1. GeoIP Enrichment
- **Source:** Use a free, privacy-respecting API (e.g., `ip-api.com` or similar) to avoid large database downloads for now, or support a local MMDB if available.
- **Trigger:** When an `IP` entity is added or discovered (e.g., from DNS resolution).
- **Data Points:**
    - Country Code (US, DE, CN)
    - City
    - ISP / Organization
    - Latitude / Longitude

### 2. Storage
- **Metadata:** Store this data in the `metadata` JSON field of the `Entity`.

### 3. Visualization
- **Graph:** Display flag emojis (🇺🇸, 🇩🇪) in the graph node labels.
- **Map:** (Future) Plot on a world map.

## Constraints
- **Rate Limits:** `ip-api.com` has rate limits (45/min). We must respect them via our `ethics` package.
- **Privacy:** Do not send case-sensitive data, only the IP address.

## Success Criteria
- Running `spectre collect dns ...` automatically enriches resulting IPs with location data.
- `spectre visualize` shows country flags on IP nodes.
