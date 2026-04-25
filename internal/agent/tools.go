package agent

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spectre/spectre/internal/analyzer"
	"github.com/spectre/spectre/internal/collector"
	"github.com/spectre/spectre/internal/storage"
)

// Tool represents a capability the LLM can invoke.
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  interface{} `json:"parameters"`
	Execute     func(caseID string, args map[string]interface{}) (string, error)
}

// Registry stores all available tools for the agent.
var Registry = map[string]Tool{
	"collect": {
		Name:        "collect",
		Description: "Run a specific collector on a target to gather intelligence. Use 'list_collectors' to see available collectors.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"collector": map[string]interface{}{
					"type":        "string",
					"description": "Name of the collector to run (e.g., 'dns', 'whois').",
				},
				"target": map[string]interface{}{
					"type":        "string",
					"description": "The target to collect data for (domain, IP, email, etc.).",
				},
			},
			"required": []string{"collector", "target"},
		},
		Execute: func(caseID string, args map[string]interface{}) (string, error) {
			name, ok := args["collector"].(string)
			if !ok {
				return "", fmt.Errorf("collector name is required")
			}
			target, ok := args["target"].(string)
			if !ok {
				return "", fmt.Errorf("target is required")
			}

			// By default, we don't allow active recon via chat unless we add an 'active' parameter
			evidence, err := collector.RunAndSave(name, caseID, target, false, nil)
			if err != nil {
				return "", err
			}

			return fmt.Sprintf("Collector '%s' finished on target '%s'. Generated %d evidence records.", name, target, len(evidence)), nil
		},
	},

	"list_collectors": {
		Name:        "list_collectors",
		Description: "List all available intelligence collectors and their descriptions.",
		Parameters: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
		Execute: func(caseID string, args map[string]interface{}) (string, error) {
			collectors := collector.List()
			var sb strings.Builder
			sb.WriteString("Available Collectors:\n")
			for _, c := range collectors {
				sb.WriteString(fmt.Sprintf("- %s: %s (Active: %v)\n", c.Name(), c.Description(), c.IsActive()))
			}
			return sb.String(), nil
		},
	},

	"search_entities": {
		Name:        "search_entities",
		Description: "Search for existing entities in the current case by value or type.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "The search term (e.g., 'google.com' or 'email').",
				},
			},
			"required": []string{"query"},
		},
		Execute: func(caseID string, args map[string]interface{}) (string, error) {
			query, ok := args["query"].(string)
			if !ok {
				return "", fmt.Errorf("query is required")
			}

			entities, err := storage.ListEntitiesByCase(caseID)
			if err != nil {
				return "", err
			}

			var matched []string
			for _, e := range entities {
				if strings.Contains(strings.ToLower(e.Value), strings.ToLower(query)) ||
					strings.Contains(strings.ToLower(e.Type), strings.ToLower(query)) {
					matched = append(matched, fmt.Sprintf("[%s] %s (Source: %s)", e.Type, e.Value, e.Source))
				}
			}

			if len(matched) == 0 {
				return "No matching entities found in this case.", nil
			}

			return "Found entities:\n" + strings.Join(matched, "\n"), nil
		},
	},

	"get_case_summary": {
		Name:        "get_case_summary",
		Description: "Get a high-level summary of the current investigation, including entity counts.",
		Parameters: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
		Execute: func(caseID string, args map[string]interface{}) (string, error) {
			c, err := storage.GetCase(caseID)
			if err != nil {
				return "", err
			}
			if c == nil {
				return "Case not found.", nil
			}

			entities, _ := storage.ListEntitiesByCase(caseID)
			rels, _ := storage.ListRelationshipsByCase(caseID)

			return fmt.Sprintf("Case: %s\nDescription: %s\nStats: %d entities, %d relationships.",
				c.Name, c.Description, len(entities), len(rels)), nil
		},
	},

	"search_evidence": {
		Name:        "search_evidence",
		Description: "Perform a semantic search across all collected evidence files to find specific information.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "The search query (e.g., 'registrant email' or 'vulnerabilities').",
				},
			},
			"required": []string{"query"},
		},
		Execute: func(caseID string, args map[string]interface{}) (string, error) {
			query, _ := args["query"].(string)
			req := analyzer.Request{
				Task:   "search_evidence",
				CaseID: caseID,
				// Add extra fields for the Python task
				Data: map[string]interface{}{
					"case_id": caseID,
					"query":   query,
				},
			}

			output, err := analyzer.RunPythonTask(req)
			if err != nil {
				return "", err
			}

			var resp struct {
				Results []struct {
					ID      string `json:"id"`
					Content string `json:"content"`
				} `json:"results"`
			}
			if err := json.Unmarshal([]byte(output), &resp); err != nil {
				return "", err
			}

			if len(resp.Results) == 0 {
				return "No relevant information found in evidence files.", nil
			}

			var sb strings.Builder
			sb.WriteString("Found relevant evidence snippets:\n")
			for _, r := range resp.Results {
				sb.WriteString(fmt.Sprintf("\n--- Source: %s ---\n%s\n", r.ID, r.Content))
			}
			return sb.String(), nil
		},
	},

	"read_evidence": {
		Name:        "read_evidence",
		Description: "Read the full content of a specific evidence file.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"filename": map[string]interface{}{
					"type":        "string",
					"description": "The name of the file to read (e.g., 'whois_scanme.nmap.org.txt').",
				},
			},
			"required": []string{"filename"},
		},
		Execute: func(caseID string, args map[string]interface{}) (string, error) {
			filename, _ := args["filename"].(string)
			// Implementation would use storage or direct file access
			// For brevity, assuming we have a way to get the path
			return fmt.Sprintf("Full content of %s requested. (Implementation pending file path resolution)", filename), nil
		},
	},

	"analyze_image": {
		Name:        "analyze_image",
		Description: "Perform visual analysis on an image file (e.g., screenshots) to extract text, logos, or descriptions using a local vision model.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"filename": map[string]interface{}{
					"type":        "string",
					"description": "The name of the image file to analyze (must be in evidence_storage/<case_id>/).",
				},
				"prompt": map[string]interface{}{
					"type":        "string",
					"description": "Specific instructions or questions about the image (e.g., 'What text is visible?').",
				},
			},
			"required": []string{"filename"},
		},
		Execute: func(caseID string, args map[string]interface{}) (string, error) {
			filename, _ := args["filename"].(string)
			prompt, ok := args["prompt"].(string)
			if !ok {
				prompt = "Describe this image in detail. Focus on text, logos, or identifying features."
			}

			filePath := fmt.Sprintf("evidence_storage/%s/%s", caseID, filename)
			imgData, err := os.ReadFile(filePath)
			if err != nil {
				return "", fmt.Errorf("failed to read image file: %w", err)
			}

			base64Str := base64.StdEncoding.EncodeToString(imgData)

			req := analyzer.Request{
				Task:    "vision",
				CaseID:  caseID,
				Data:    prompt,
				Context: base64Str,
			}

			output, err := analyzer.RunPythonTask(req)
			if err != nil {
				return "", fmt.Errorf("vision analysis failed: %w", err)
			}

			var resp struct {
				Answer string `json:"answer"`
			}
			if err := json.Unmarshal([]byte(output), &resp); err != nil {
				return "", fmt.Errorf("failed to parse vision response: %w", err)
			}

			return resp.Answer, nil
		},
	},

	"generate_dorks": {
		Name:        "generate_dorks",
		Description: "Generate a list of specialized Google Dorks for a target domain to find sensitive or leaked information.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"target": map[string]interface{}{
					"type":        "string",
					"description": "The target domain (e.g., 'example.com').",
				},
			},
			"required": []string{"target"},
		},
		Execute: func(caseID string, args map[string]interface{}) (string, error) {
			target, _ := args["target"].(string)

			req := analyzer.Request{
				Task:   "generate_dorks",
				CaseID: caseID,
				Data:   target,
			}

			output, err := analyzer.RunPythonTask(req)
			if err != nil {
				return "", fmt.Errorf("dork generation failed: %w", err)
			}

			var resp struct {
				Dorks []string `json:"dorks"`
			}
			if err := json.Unmarshal([]byte(output), &resp); err != nil {
				return "", fmt.Errorf("failed to parse dorks: %w", err)
			}

			if len(resp.Dorks) == 0 {
				return "No dorks generated.", nil
			}

			return "Suggested Google Dorks:\n" + strings.Join(resp.Dorks, "\n"), nil
		},
	},
}

// GetToolDefinitions returns the JSON-serializable tool definitions for the LLM.
func GetToolDefinitions() []interface{} {
	var defs []interface{}
	for _, t := range Registry {
		defs = append(defs, map[string]interface{}{
			"name":        t.Name,
			"description": t.Description,
			"parameters":  t.Parameters,
		})
	}
	return defs
}
