package core

import (
	"fmt"
	"regexp"
	"strings"
)

var validCaseIDRe = regexp.MustCompile(`^[A-Za-z0-9\-_]{1,64}$`)
var invalidFileCharsRe = regexp.MustCompile(`[^A-Za-z0-9._-]+`)

// ValidateCaseID ensures the case ID is safe against directory traversal and contains allowed characters.
func ValidateCaseID(caseID string) error {
	if caseID == "" {
		return fmt.Errorf("case ID cannot be empty")
	}
	if !validCaseIDRe.MatchString(caseID) {
		return fmt.Errorf("invalid case ID %q: must contain only alphanumeric characters, hyphens, or underscores (max 64 chars)", caseID)
	}
	return nil
}

// SanitizeFileName replaces unsafe path and filename characters with underscores.
func SanitizeFileName(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "file"
	}
	s = invalidFileCharsRe.ReplaceAllString(s, "_")
	if len(s) > 128 {
		s = s[:128]
	}
	return s
}
