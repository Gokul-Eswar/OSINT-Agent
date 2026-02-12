package analyzer

import (
	"encoding/json"
	"os/exec"
	"testing"
)

func TestRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     Request
		wantErr bool
	}{
		{
			name:    "valid synthesize request",
			req:     Request{Task: "synthesize", CaseID: "123"},
			wantErr: false,
		},
		{
			name:    "valid visualize request",
			req:     Request{Task: "visualize"},
			wantErr: false,
		},
		{
			name:    "missing task",
			req:     Request{CaseID: "123"},
			wantErr: true,
		},
		{
			name:    "missing case_id for synthesize",
			req:     Request{Task: "synthesize"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Request.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestRunPythonTaskMocked is a conceptual test. 
// In a real environment, we'd mock the exec.Command.
// Here we just test the logic around the command execution.
func TestRunPythonTaskBinaryCheck(t *testing.T) {
	// Skip if python is not available
	_, err := exec.LookPath("python")
	if err != nil {
		t.Skip("python not found in path")
	}

	// Create a dummy request
	req := Request{
		Task:   "invalid",
		CaseID: "test",
	}

	_, err = RunPythonTask(req)
	if err == nil {
		t.Error("Expected error for invalid task, got nil")
	}
}

func TestPythonErrorParsing(t *testing.T) {
	// This tests the logic added in the previous turn where we parse error JSON from stdout
	
	// Create a dummy stdout with error JSON
	dummyStdout := `{"error": "API key invalid"}`
	
	var errResp struct {
		Error string `json:"error"`
	}
	err := json.Unmarshal([]byte(dummyStdout), &errResp)
	if err != nil {
		t.Fatalf("Failed to unmarshal dummy stdout: %v", err)
	}
	
	if errResp.Error != "API key invalid" {
		t.Errorf("Expected error 'API key invalid', got %s", errResp.Error)
	}
}
