package analysis

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/spectre/spectre/internal/analyzer"
	"github.com/spectre/spectre/internal/core"
	"github.com/spectre/spectre/internal/storage"
	"github.com/spf13/viper"
)

// AnalyzeCase runs the AI analysis via the Python analyzer.
func AnalyzeCase(caseID string, model string) (*core.Analysis, error) {
	// 1. Fetch Case
	c, err := storage.GetCase(caseID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, fmt.Errorf("case not found")
	}

	// 2. Build Context
	contextData, err := BuildCaseContext(caseID)
	if err != nil {
		return nil, err
	}

	// Optimization: Check Cache
	hash := sha256.Sum256([]byte(contextData))
	hashStr := hex.EncodeToString(hash[:])
	
	cached, err := storage.GetAnalysisByHash(caseID, hashStr)
	if err == nil && cached != nil {
		// Cache hit
		return cached, nil
	}

	// 3. Prepare Bridge Request
	req := analyzer.Request{
		Task:     "synthesize",
		CaseName: c.Name,
		Context:  contextData,
		Model:    model,
		LLMConfig: analyzer.LLMConfig{
			Provider: viper.GetString("llm.provider"),
			URL:      viper.GetString("llm.url"),
			APIKey:   viper.GetString("llm.api_key"),
			Timeout:  viper.GetInt("llm.timeout"),
		},
	}

	// 4. Run Python Analyzer
	responseJSON, err := analyzer.RunPythonTask(req)
	if err != nil {
		return nil, fmt.Errorf("python analysis failed: %w", err)
	}

	// 5. Parse
	var result core.Analysis
	if err := json.Unmarshal([]byte(responseJSON), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w\nResponse was: %s", err, responseJSON)
	}

		result.CaseID = caseID

		result.ContextHash = hashStr

		

		// 6. Save

		if err := storage.SaveAnalysis(&result); err != nil {

	
	
				return nil, err
	
			}
	
		
	
			return &result, nil
	
		}
	
		
	
	