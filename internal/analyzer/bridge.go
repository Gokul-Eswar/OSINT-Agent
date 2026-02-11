package analyzer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"
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

// Validate checks if the request has all required fields.
func (r *Request) Validate() error {
	if r.Task == "" {
		return fmt.Errorf("task is required")
	}
	if r.CaseID == "" && r.Task != "visualize" { // CaseID might be optional for some visualization tasks if Data is provided
		return fmt.Errorf("case_id is required")
	}
	return nil
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
	if err := req.Validate(); err != nil {
		return "", fmt.Errorf("invalid request: %w", err)
	}

	log.Info().
		Str("task", req.Task).
		Str("case_id", req.CaseID).
		Msg("python_task_started")

	inputJSON, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// Determine Python executable
	pythonPath := "python" // Default to system path
	
	// Check for local venv (Windows)
	if _, err := os.Stat(".venv/Scripts/python.exe"); err == nil {
		pythonPath = ".venv/Scripts/python.exe"
	} else if _, err := os.Stat(".venv/bin/python"); err == nil {
		// Unix/Mac
		pythonPath = ".venv/bin/python"
	}

	// Execute: <python> -m analyzer --task <task> --input <json>
	cmd := exec.CommandContext(ctx, pythonPath, "-m", "analyzer", "--task", req.Task, "--input", string(inputJSON))
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Error().Err(err).Str("task", req.Task).Msg("python_task_timeout")
			return "", fmt.Errorf("python analysis timed out after 3 minutes")
		}
		
		// Try to parse error from stdout if stderr is empty (some python errors might go there)
		var errResp struct {
			Error string `json:"error"`
		}
		if json.Unmarshal(stdout.Bytes(), &errResp) == nil && errResp.Error != "" {
			log.Error().Str("error", errResp.Error).Str("task", req.Task).Msg("python_task_failed")
			return "", fmt.Errorf("python execution failed: %s", errResp.Error)
		}

		log.Error().
			Err(err).
			Str("task", req.Task).
			Str("stderr", stderr.String()).
			Msg("python_task_execution_error")
		return "", fmt.Errorf("python execution failed: %w\nStderr: %s", err, stderr.String())
	}

	log.Debug().Str("task", req.Task).Msg("python_task_completed")
	return stdout.String(), nil
}
