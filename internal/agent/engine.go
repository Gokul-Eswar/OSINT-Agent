package agent

import (
	"encoding/json"
	"fmt"

	"github.com/spectre/spectre/internal/analyzer"
	"github.com/spf13/viper"
)

// Engine manages the state and execution loop of the SPECTRE Agent.
type Engine struct {
	CaseID  string
	History []analyzer.Message
	Tools   []interface{}
}

// NewEngine initializes a new agent engine for a specific case.
func NewEngine(caseID string) *Engine {
	return &Engine{
		CaseID:  caseID,
		History: []analyzer.Message{},
		Tools:   GetToolDefinitions(),
	}
}

// Execute processes a user input and returns the final agent response.
// It handles the tool-use loop internally.
func (e *Engine) Execute(userInput string) (string, error) {
	e.History = append(e.History, analyzer.Message{
		Role:    "user",
		Content: userInput,
	})

	for i := 0; i < 5; i++ { // Limit iterations to prevent infinite loops
		req := analyzer.Request{
			Task:     "chat",
			CaseID:   e.CaseID,
			Messages: e.History,
			Tools:    e.Tools,
			Model:    viper.GetString("llm.model"),
			LLMConfig: analyzer.LLMConfig{
				Provider: viper.GetString("llm.provider"),
				URL:      viper.GetString("llm.url"),
				APIKey:   viper.GetString("llm.api_key"),
				Timeout:  viper.GetInt("llm.timeout"),
			},
		}

		if req.Model == "" {
			req.Model = "llama3"
		}
		if req.LLMConfig.URL == "" {
			req.LLMConfig.URL = "http://localhost:11434/api/chat"
		}

		output, err := analyzer.RunPythonTask(req)
		if err != nil {
			return "", fmt.Errorf("llm error: %w", err)
		}

		var resp analyzer.Response
		if err := json.Unmarshal([]byte(output), &resp); err != nil {
			return "", fmt.Errorf("failed to parse llm response: %w\nOutput: %s", err, output)
		}

		if resp.Error != "" {
			return "", fmt.Errorf("llm returned error: %s", resp.Error)
		}

		// Update history with assistant message
		e.History = append(e.History, analyzer.Message{
			Role:    "assistant",
			Content: resp.Content,
		})

		// If no tool use, we are done
		if resp.ToolUse == nil {
			return resp.Content, nil
		}

		// Execute Tool
		tool, ok := Registry[resp.ToolUse.Name]
		if !ok {
			toolResult := fmt.Sprintf("Error: Tool '%s' not found.", resp.ToolUse.Name)
			e.History = append(e.History, analyzer.Message{
				Role:    "system",
				Content: toolResult,
			})
			continue
		}

		result, err := tool.Execute(e.CaseID, resp.ToolUse.Arguments)
		if err != nil {
			result = fmt.Sprintf("Error executing tool '%s': %v", resp.ToolUse.Name, err)
		}

		// Add tool result to history
		e.History = append(e.History, analyzer.Message{
			Role:    "system",
			Content: fmt.Sprintf("Tool '%s' result: %s", resp.ToolUse.Name, result),
		})

		// Loop continues to let LLM process the tool result
	}

	return "", fmt.Errorf("agent reached maximum iteration limit (5) without finishing")
}
