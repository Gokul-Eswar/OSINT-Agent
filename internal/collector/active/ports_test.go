package active

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestPortCollector_Collect(t *testing.T) {
	// Start a listener on a random port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start listener: %v", err)
	}
	defer listener.Close()

	// Get the port
	addr := listener.Addr().(*net.TCPAddr)
	port := addr.Port

	// Configure Viper
	viper.Set("collectors.ports.mode", "custom")
	viper.Set("collectors.ports.custom_ports", []int{port, port + 1}) // Scan open port and one closed port

	collector := &PortCollector{}
	caseID := "test_case_ports"
	target := "127.0.0.1"

	// Cleanup
	defer os.RemoveAll(filepath.Join("evidence_storage", caseID))
	os.RemoveAll(filepath.Join("evidence_storage", caseID))

	// Run Collect
	evidence, err := collector.Collect(caseID, target)
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if len(evidence) != 1 {
		t.Fatalf("Expected 1 evidence item, got %d", len(evidence))
	}

	ev := evidence[0]
	if ev.Collector != "ports" {
		t.Errorf("Expected collector 'ports', got '%s'", ev.Collector)
	}

	// Verify File Content
	content, err := os.ReadFile(ev.FilePath)
	if err != nil {
		t.Fatalf("Failed to read evidence file: %v", err)
	}

	var results map[string]string // Key is stringified int in JSON
	if err := json.Unmarshal(content, &results); err != nil {
		t.Fatalf("Failed to parse evidence JSON: %v", err)
	}

	// Verify open port found
	if status, ok := results[fmt.Sprintf("%d", port)]; !ok || status != "open" {
		t.Errorf("Expected port %d to be open, got %v", port, results[fmt.Sprintf("%d", port)])
	}

	// Verify closed port not found (or closed if logic changes, but current logic only adds "open" ports)
	if _, ok := results[fmt.Sprintf("%d", port+1)]; ok {
		t.Errorf("Expected port %d to be closed/missing, but it was in results", port+1)
	}
}
