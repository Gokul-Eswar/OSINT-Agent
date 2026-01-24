package ethics

import (
	"fmt"
	"strings"
)

var (
	blacklist = []string{".gov", ".mil", "localhost", "127.0.0.1"}
	whitelist []string // If empty, everything not blacklisted is allowed
)

// IsAllowed checks if a target is within the allowed scope.
func IsAllowed(target string) (bool, error) {
	target = strings.ToLower(strings.TrimSpace(target))

	// Check Blacklist
	for _, b := range blacklist {
		if strings.Contains(target, b) {
			return false, fmt.Errorf("target '%s' is in the blacklist (matches '%s')", target, b)
		}
	}

	// Check Whitelist (if active)
	if len(whitelist) > 0 {
		allowed := false
		for _, w := range whitelist {
			if strings.Contains(target, w) {
				allowed = true
				break
			}
		}
		if !allowed {
			return false, fmt.Errorf("target '%s' is not in the whitelist", target)
		}
	}

	return true, nil
}

// SetBlacklist overrides the default blacklist.
func SetBlacklist(list []string) {
	blacklist = list
}

// SetWhitelist sets a strict whitelist.
func SetWhitelist(list []string) {
	whitelist = list
}
