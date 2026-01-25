package active

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestSocialCollector_Collect(t *testing.T) {
	// Setup test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/found":
			w.WriteHeader(http.StatusOK)
		case "/not_found":
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	// Initialize collector with test sites
	collector := NewSocialCollector()
	collector.Sites = map[string]string{
		"TestSiteFound":    server.URL + "/%s",      // Will hit /found if target is "found"
		"TestSiteNotFound": server.URL + "/not_%s",  // Will hit /not_found if target is "found"
	}

	caseID := "test_case_social"
	target := "found" // This makes the URLs: /found and /not_found

	// Cleanup
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
	if ev.Collector != "social" {
		t.Errorf("Expected collector 'social', got '%s'", ev.Collector)
	}

	// Verify File Content
	content, err := os.ReadFile(ev.FilePath)
	if err != nil {
		t.Fatalf("Failed to read evidence file: %v", err)
	}

	var results []SiteResult
	if err := json.Unmarshal(content, &results); err != nil {
		t.Fatalf("Failed to parse evidence JSON: %v", err)
	}

	// Logic: Collect only appends "found" results to the list?
	// Let's check social.go implementation:
	// if status == "found" { append }
	
	// So we expect 1 result (TestSiteFound)
	if len(results) != 1 {
		t.Errorf("Expected 1 found result, got %d", len(results))
	} else {
		if results[0].Site != "TestSiteFound" {
			t.Errorf("Expected found site to be 'TestSiteFound', got '%s'", results[0].Site)
		}
		if results[0].Status != "found" {
			t.Errorf("Expected status 'found', got '%s'", results[0].Status)
		}
	}
}
