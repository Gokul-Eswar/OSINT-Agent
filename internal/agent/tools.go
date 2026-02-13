package agent

import (
	"fmt"
	"strings"

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
			evidence, err := collector.RunAndSave(name, caseID, target, false)
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
