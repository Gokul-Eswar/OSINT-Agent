package analysis

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spectre/spectre/internal/ai"
	"github.com/spectre/spectre/internal/core"
	"github.com/spectre/spectre/internal/storage"
)

const systemPrompt = `You are SPECTRE, an expert intelligence analyst. 
Analyze the provided case data (Entities, Relationships, Evidence) and generate a structured report.
Your output MUST be strict JSON in the following format:
{
  "findings": ["finding 1", "finding 2"],
  "risks": ["risk 1", "risk 2"],
  "connections": ["potential connection 1", "potential connection 2"],
  "next_steps": ["step 1", "step 2"],
  "confidence": 0.85
}
Do not include markdown formatting (like ` + "`" + `json). Just return the raw JSON string.`

// AnalyzeCase runs the AI analysis on a specific case.
func AnalyzeCase(caseID string, provider ai.Provider) (*core.Analysis, error) {
	// 1. Build Context
	contextData, err := BuildCaseContext(caseID)
	if err != nil {
		return nil, err
	}

	// 2. Construct Prompt
	fullPrompt := fmt.Sprintf("%s\n\nCASE DATA:\n%s", systemPrompt, contextData)

	// 3. Generate
	responseJSON, err := provider.Generate(fullPrompt)
	if err != nil {
		return nil, err
	}

	// Clean up potential markdown code blocks if the LLM ignores instructions
	responseJSON = strings.TrimPrefix(responseJSON, "```json")
	responseJSON = strings.TrimPrefix(responseJSON, "```")
	responseJSON = strings.TrimSuffix(responseJSON, "```")
	responseJSON = strings.TrimSpace(responseJSON)

	// 4. Parse
	var result core.Analysis
	if err := json.Unmarshal([]byte(responseJSON), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w\nResponse was: %s", err, responseJSON)
	}

	result.CaseID = caseID
	
	// 5. Save
	if err := storage.SaveAnalysis(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
