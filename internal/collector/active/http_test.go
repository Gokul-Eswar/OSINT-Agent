package active

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHTTPCollector_Collect(t *testing.T) {
	// Setup test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "SpectreTestServer")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "<html><head><title>Test Title</title></head><body>Hello</body></html>")
	}))
	defer server.Close()

	// Extract host (ip:port) from URL
	// server.URL is like "http://127.0.0.1:12345"
	target := strings.TrimPrefix(server.URL, "http://")

	collector := &HTTPCollector{}
	caseID := "test_case_http"

	// Cleanup before and after
	defer os.RemoveAll(filepath.Join("evidence_storage", caseID))
	os.RemoveAll(filepath.Join("evidence_storage", caseID))

	evidence, err := collector.Collect(caseID, target)
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if len(evidence) != 1 {
		t.Fatalf("Expected 1 evidence item, got %d", len(evidence))
	}

	ev := evidence[0]
	if ev.Collector != "http" {
		t.Errorf("Expected collector 'http', got '%s'", ev.Collector)
	}

	// Verify Metadata
	if title, ok := ev.Metadata["title"]; !ok || title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got '%v'", title)
	}
	if srv, ok := ev.Metadata["server"]; !ok || srv != "SpectreTestServer" {
		t.Errorf("Expected server 'SpectreTestServer', got '%v'", srv)
	}

	// Verify File Content
	content, err := os.ReadFile(ev.FilePath)
	if err != nil {
		t.Fatalf("Failed to read evidence file: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("Failed to parse evidence JSON: %v", err)
	}

	if result["status_code"].(float64) != 200 {
		t.Errorf("Expected status 200, got %v", result["status_code"])
	}
}
