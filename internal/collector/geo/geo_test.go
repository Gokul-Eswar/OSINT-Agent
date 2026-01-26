package geo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGeoIPCollector_Collect(t *testing.T) {
	// Mock ip-api response
	mockResponse := map[string]interface{}{
		"status":      "success",
		"country":     "United States",
		"countryCode": "US",
		"region":      "VA",
		"city":        "Ashburn",
		"lat":         39.0438,
		"lon":         -77.4874,
		"isp":         "Google LLC",
		"query":       "8.8.8.8",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	c := &GeoIPCollector{BaseURL: server.URL + "/"}
	evs, err := c.Collect("test-case", "8.8.8.8")
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if len(evs) != 1 {
		t.Fatalf("Expected 1 evidence, got %d", len(evs))
	}

	ev := evs[0]
	if ev.Metadata["lat"] != 39.0438 {
		t.Errorf("Expected lat 39.0438, got %v", ev.Metadata["lat"])
	}
	if ev.Metadata["lon"] != -77.4874 {
		t.Errorf("Expected lon -77.4874, got %v", ev.Metadata["lon"])
	}
	if ev.Metadata["city"] != "Ashburn" {
		t.Errorf("Expected city Ashburn, got %v", ev.Metadata["city"])
	}
}
