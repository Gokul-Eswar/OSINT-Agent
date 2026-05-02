package agent

import (
	"encoding/json"
	"fmt"

	"github.com/spectre/spectre/internal/analyzer"
	"github.com/spf13/viper"
)

const AgentPersona = `You are the SPECTRE Autonomous Intelligence Agent. 
Your goal is to assist in OSINT (Open Source Intelligence) investigations.
You have access to a variety of tools to collect data, search entities, and analyze evidence.

Guidelines:
1. Be methodical. Plan your steps before executing them.
2. If you find a new lead (e.g., an email in WHOIS), follow it using relevant collectors.
3. Use 'search_evidence' to look for specific details within the files you've collected.
4. When you have enough information or have reached a dead end, summarize your findings clearly.
5. Always maintain operational security and follow ethical guidelines.

Respond in a structured way:
Thought: <Your reasoning about the current state>
Plan: <What you intend to do next>
Tool: <If calling a tool, specify the tool name and arguments>
`

// Engine manages the state and execution loop of the SPECTRE Agent.
type Engine struct {
	CaseID  string
	History []analyzer.Message
	Tools   []interface{}
}

// NewEngine initializes a new agent engine for a specific case.
func NewEngine(caseID string) *Engine {
	return &Engine{
		CaseID: caseID,
		History: []analyzer.Message{
			{Role: "system", Content: AgentPersona},
		},
		Tools: GetToolDefinitions(),
	}
}

// Execute runs a bounded tool-use loop for agent chat.
func (e *Engine) Execute(userInput string) (string, error) {
	e.History = append(e.History, analyzer.Message{
		Role:    "user",
		Content: userInput,
	})

	for i := 0; i < 8; i++ { // Increased iteration limit for autonomous loops
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

		output, err := analyzer.GlobalTaskRunner.Run(req)
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

		fmt.Printf("\n[Agent] Decision: %s\n", resp.Content)
		fmt.Printf("[Agent] Action: Running tool '%s' with arguments: %v\n", resp.ToolUse.Name, resp.ToolUse.Arguments)
		
		result, err := tool.Execute(e.CaseID, resp.ToolUse.Arguments)
		if err != nil {
			result = fmt.Sprintf("Error executing tool '%s': %v", resp.ToolUse.Name, err)
		}

		// Add tool result to history
		e.History = append(e.History, analyzer.Message{
			Role:    "system",
			Content: fmt.Sprintf("Observation from '%s': %s", resp.ToolUse.Name, result),
		})

		// Loop continues to let LLM process the tool result
	}

	return "Agent reached maximum iteration limit without finishing the task.", nil
}
