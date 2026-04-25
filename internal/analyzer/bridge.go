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

// PythonCommand allows overriding the command used to run python tasks (useful for testing)
var PythonCommand = []string{"-m", "analyzer"}

// Request defines the structure sent to the Python analyzer.
type Request struct {
	Task      string        `json:"task"`
	CaseID    string        `json:"case_id"`
	CaseName  string        `json:"case_name"`
	Context   string        `json:"context"`
	Model     string        `json:"model"`
	Data      interface{}   `json:"data"` // For track 5 graph data
	LLMConfig LLMConfig     `json:"llm_config"`
	Messages  []Message     `json:"messages,omitempty"`
	Tools     []interface{} `json:"tools,omitempty"`
}

// Message represents a single turn in a chat.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Response represents a structured response from the Python analyzer.
type Response struct {
	Role    string   `json:"role"`
	Content string   `json:"content"`
	ToolUse *ToolUse `json:"tool_use,omitempty"`
	Error   string   `json:"error,omitempty"`
}

// ToolUse represents an LLM's request to call a tool.
type ToolUse struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// Validate enforces minimum request requirements before invoking Python.
//
// CaseID is required for all current tasks except visualize-style tasks where
// graph data may be supplied directly in Data.
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

// IndexEvidence triggers the Python analyzer to index all evidence files for a case.
func IndexEvidence(caseID string) error {
	// List files in evidence_storage/<caseID>
	evidenceDir := fmt.Sprintf("evidence_storage/%s", caseID)
	files, err := os.ReadDir(evidenceDir)
	if err != nil {
		return fmt.Errorf("failed to read evidence directory: %w", err)
	}

	var filePaths []string
	for _, f := range files {
		if !f.IsDir() {
			filePaths = append(filePaths, fmt.Sprintf("%s/%s", evidenceDir, f.Name()))
		}
	}

	req := Request{
		Task:   "index_evidence",
		CaseID: caseID,
		Data: map[string]interface{}{
			"case_id": caseID,
			"files":   filePaths,
		},
	}

	_, err = RunPythonTask(req)
	return err
}

// RunPythonTask executes the Python analyzer as a subprocess with a strict
// timeout and JSON-over-CLI contract.
//
// The bridge chooses a local virtualenv Python when present, then falls back to
// PATH resolution. On failures it attempts structured error extraction from
// stdout before returning stderr-heavy execution errors.
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

	// Execute: <python> <PythonCommand...> --task <task> --input <json>
	args := append(PythonCommand, "--task", req.Task, "--input", string(inputJSON))
	cmd := exec.CommandContext(ctx, pythonPath, args...)

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
