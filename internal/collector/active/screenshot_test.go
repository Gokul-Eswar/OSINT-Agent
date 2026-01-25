package active

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScreenshotCollector_Collect(t *testing.T) {
	// Skip if CI or explicitly requested (optional, but good practice)
	if os.Getenv("SKIP_HEADLESS") == "true" {
		t.Skip("Skipping headless browser tests")
	}

	// Setup test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, "<html><body style='background-color:red;'><h1>Screenshot Me</h1></body></html>")
	}))
	defer server.Close()

	collector := &ScreenshotCollector{}
	caseID := "test_case_screenshot"
	// Use the full URL from server, now that we fixed the prefix bug
	target := server.URL

	// Cleanup
	defer os.RemoveAll(filepath.Join("evidence_storage", caseID))
	os.RemoveAll(filepath.Join("evidence_storage", caseID))

	evidence, err := collector.Collect(caseID, target)
	if err != nil {
		// If chrome is missing, this might fail.
		// Check error message to decide if we should fail or skip?
		// For now, let's fail to be explicit, but log it clearly.
		if strings.Contains(err.Error(), "exec: \"google-chrome\": executable file not found") || 
		   strings.Contains(err.Error(), "executable file not found") {
			t.Skipf("Skipping screenshot test: Chrome not found: %v", err)
		}
		t.Fatalf("Collect failed: %v", err)
	}

	if len(evidence) != 1 {
		t.Fatalf("Expected 1 evidence item, got %d", len(evidence))
	}

	ev := evidence[0]
	if ev.Collector != "screenshot" {
		t.Errorf("Expected collector 'screenshot', got '%s'", ev.Collector)
	}

	// Verify File Exists
	if _, err := os.Stat(ev.FilePath); os.IsNotExist(err) {
		t.Errorf("Screenshot file does not exist: %s", ev.FilePath)
	}

	// Verify Metadata
	if _, ok := ev.Metadata["size"]; !ok {
		t.Error("Metadata 'size' missing")
	}
}
