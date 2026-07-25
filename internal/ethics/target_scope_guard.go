package ethics

import (
	"fmt"
	"strings"
)

var (
	// blacklist holds target substrings that are strictly forbidden from scanning/collection.
	// Defaults cover government, military, and local system addresses to prevent unauthorized local/institutional scanning.
	blacklist = []string{".gov", ".mil", "localhost", "127.0.0.1"}

	// whitelist holds target substrings that are explicitly allowed.
	// If the whitelist is populated (non-empty), then only targets matching a whitelisted substring can be scanned.
	whitelist []string
)

// IsAllowed validates if a target scan target conforms to the ethics policy.
// It checks the input against the active blacklist and whitelist.
// Returns (true, nil) if validation passes, or (false, error) if the target violates scope limits.
func IsAllowed(target string) (bool, error) {
	// Standardize casing and spaces to prevent bypasses like ".GOV" or leading/trailing whitespace tricks.
	target = strings.ToLower(strings.TrimSpace(target))

	// 1. Blacklist Evaluation:
	// Loop over all prohibited strings and block the target if it contains any blacklisted substring.
	// NOTE: Because this relies on strings.Contains, it can trigger false positives on innocent domains
	// (e.g. "mygovernmentblog.com" contains ".gov" but is not a .gov TLD; "milicious.com" contains ".mil").
	for _, b := range blacklist {
		if strings.Contains(target, b) {
			return false, fmt.Errorf("target '%s' is in the blacklist (matches '%s')", target, b)
		}
	}

	// 2. Whitelist Evaluation:
	// If a whitelist has been explicitly set in configurations, enforce it strictly.
	if len(whitelist) > 0 {
		allowed := false
		for _, w := range whitelist {
			// If the target contains any of the whitelisted strings, it passes the whitelist check.
			if strings.Contains(target, w) {
				allowed = true
				break
			}
		}
		// If the target didn't match any whitelisted entry, reject it.
		if !allowed {
			return false, fmt.Errorf("target '%s' is not in the whitelist", target)
		}
	}

	// Target passed both checks and is allowed to proceed.
	return true, nil
}

// SetBlacklist overrides the default blacklist configurations at runtime.
// Typically triggered during configuration bootstrapping from Viper values.
func SetBlacklist(list []string) {
	blacklist = list
}

// SetWhitelist enforces a strict target scope by overwriting the whitelist values.
// If active, only targets that contain these strings will pass the IsAllowed check.
func SetWhitelist(list []string) {
	whitelist = list
}
