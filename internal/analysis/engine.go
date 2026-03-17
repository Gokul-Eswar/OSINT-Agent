package analysis

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spectre/spectre/internal/analyzer"
	"github.com/spectre/spectre/internal/core"
	"github.com/spectre/spectre/internal/storage"
	"github.com/spf13/viper"
)

// QueryCase asks a specific question about a case to the LLM.
func QueryCase(caseID string, model string, question string) (string, error) {
	log.Info().
		Str("case_id", caseID).
		Str("model", model).
		Str("question", question).
		Msg("case_query_started")

	// 1. Fetch Case
	c, err := storage.GetCase(caseID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch case %s: %w", caseID, err)
	}

	// 2. Build Context
	contextData, err := BuildCaseContext(caseID)
	if err != nil {
		return "", fmt.Errorf("failed to build context for case %s: %w", caseID, err)
	}

	// 3. Prepare Request
	req := analyzer.Request{
		Task:     "query",
		CaseID:   caseID,
		CaseName: c.Name,
		Context:  contextData,
		Model:    model,
		Data:     question, // Pass question as data
		LLMConfig: analyzer.LLMConfig{
			Provider: viper.GetString("llm.provider"),
			URL:      viper.GetString("llm.url"),
			APIKey:   viper.GetString("llm.api_key"),
			Timeout:  viper.GetInt("llm.timeout"),
		},
	}

	// 4. Run Task
	responseJSON, err := analyzer.RunPythonTask(req)
	if err != nil {
		return "", fmt.Errorf("query failed: %w", err)
	}

	// 5. Parse
	var resp struct {
		Answer string `json:"answer"`
	}
	if err := json.Unmarshal([]byte(responseJSON), &resp); err != nil {
		// If it's not JSON, return as is (fallback)
		return responseJSON, nil
	}

	return resp.Answer, nil
}

// AnalyzeCase runs the AI analysis via the Python analyzer.
func AnalyzeCase(caseID string, model string) (*core.Analysis, error) {
	log.Info().
		Str("case_id", caseID).
		Str("model", model).
		Msg("analysis_started")

	// 1. Fetch Case
	c, err := storage.GetCase(caseID)
	if err != nil {
		log.Error().Err(err).Str("case_id", caseID).Msg("failed to fetch case for analysis")
		return nil, fmt.Errorf("failed to fetch case %s: %w", caseID, err)
	}
	if c == nil {
		log.Error().Str("case_id", caseID).Msg("case not found for analysis")
		return nil, fmt.Errorf("case %s not found", caseID)
	}

	// 2. Build Context
	contextData, err := BuildCaseContext(caseID)
	if err != nil {
		log.Error().Err(err).Str("case_id", caseID).Msg("failed to build case context")
		return nil, fmt.Errorf("failed to build context for case %s: %w", caseID, err)
	}

	// Optimization: Check Cache
	hash := sha256.Sum256([]byte(contextData))
	hashStr := hex.EncodeToString(hash[:])
	
	cached, err := storage.GetAnalysisByHash(caseID, hashStr)
	if err == nil && cached != nil {
		log.Info().Str("case_id", caseID).Msg("analysis_cache_hit")
		return cached, nil
	}

	// 3. Prepare Bridge Request
	req := analyzer.Request{
		Task:     "synthesize",
		CaseID:   caseID,
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

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("invalid analyzer request")
		return nil, fmt.Errorf("invalid analyzer request: %w", err)
	}

	// 4. Run Python Analyzer
	responseJSON, err := analyzer.RunPythonTask(req)
	if err != nil {
		log.Error().Err(err).Str("case_id", caseID).Msg("python analyzer task failed")
		return nil, fmt.Errorf("python analysis failed: %w", err)
	}

	// 5. Parse
	var result core.Analysis
	if err := json.Unmarshal([]byte(responseJSON), &result); err != nil {
		log.Error().
			Err(err).
			Str("case_id", caseID).
			Str("response", responseJSON).
			Msg("failed to parse AI response")
		return nil, fmt.Errorf("failed to parse AI response: %w\nResponse was: %s", err, responseJSON)
	}

	result.CaseID = caseID
	result.ContextHash = hashStr

	// 6. Save
	if err := storage.SaveAnalysis(&result); err != nil {
		log.Error().Err(err).Str("case_id", caseID).Msg("failed to save analysis result")
		return nil, fmt.Errorf("failed to save analysis for case %s: %w", caseID, err)
	}

	log.Info().
		Str("case_id", caseID).
		Float64("confidence", result.Confidence).
		Msg("analysis_completed")

	return &result, nil
}