package analyzer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

// Request defines the structure sent to the Python analyzer.
type Request struct {
	Task     string      `json:"task"`
	CaseID   string      `json:"case_id"`
	CaseName string      `json:"case_name"`
	Context  string      `json:"context"`
	Model    string      `json:"model"`
	Data     interface{} `json:"data"` // For track 5 graph data
}

// RunPythonTask executes the Python analyzer module.
func RunPythonTask(req Request) (string, error) {
	inputJSON, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	// Execute: python -m analyzer --task <task> --input <json>
	cmd := exec.Command("python", "-m", "analyzer", "--task", req.Task, "--input", string(inputJSON))
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("python execution failed: %w\nStderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}
