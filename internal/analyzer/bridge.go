package analyzer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

// Request defines the structure sent to the Python analyzer.
type Request struct {
	Task      string      `json:"task"`
	CaseID    string      `json:"case_id"`
	CaseName  string      `json:"case_name"`
	Context   string      `json:"context"`
	Model     string      `json:"model"`
	Data      interface{} `json:"data"` // For track 5 graph data
	LLMConfig LLMConfig   `json:"llm_config"`
}

// LLMConfig holds configuration for the LLM provider.
type LLMConfig struct {
	Provider string `json:"provider"`
	URL      string `json:"url"`
	APIKey   string `json:"api_key"`
	Timeout  int    `json:"timeout"`
}

// RunPythonTask executes the Python analyzer module.
func RunPythonTask(req Request) (string, error) {
	inputJSON, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// Execute: python -m analyzer --task <task> --input <json>
	cmd := exec.CommandContext(ctx, "python", "-m", "analyzer", "--task", req.Task, "--input", string(inputJSON))
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("python analysis timed out after 3 minutes")
		}
		return "", fmt.Errorf("python execution failed: %w\nStderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}
