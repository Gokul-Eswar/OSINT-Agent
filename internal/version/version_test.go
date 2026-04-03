package version

import (
	"strings"
	"testing"
)

func TestVersion_IsNonEmptyAndTagged(t *testing.T) {
	if Version == "" {
		t.Fatal("Version must not be empty")
	}

	if !strings.HasPrefix(Version, "v") {
		t.Fatalf("Version should be prefixed with 'v', got %q", Version)
	}
}
