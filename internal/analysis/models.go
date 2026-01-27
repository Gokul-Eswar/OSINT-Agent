package analysis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type OllamaModel struct {
	Name string `json:"name"`
}

type OllamaTagsResponse struct {
	Models []OllamaModel `json:"models"`
}

// FetchAvailableModels connects to the configured LLM provider (Ollama)
// and retrieves the list of installed models.
func FetchAvailableModels() ([]string, error) {
	provider := viper.GetString("llm.provider")
	
	// If not Ollama, we can't auto-detect easily (OpenAI models are fixed/remote)
	if provider != "ollama" {
		return []string{viper.GetString("llm.model")}, nil
	}

	// Construct Tags URL from Config
	// Config typically: http://localhost:11434/api/generate
	// Goal: http://localhost:11434/api/tags
	genURL := viper.GetString("llm.url")
	
	var tagsURL string
	if strings.Contains(genURL, "/api/generate") {
		tagsURL = strings.Replace(genURL, "/api/generate", "/api/tags", 1)
	} else {
		// Fallback: assume the URL is just the base
		tagsURL = strings.TrimRight(genURL, "/") + "/api/tags"
	}

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := client.Get(tagsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ollama at %s: %w", tagsURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("ollama API returned status: %s", resp.Status)
	}

	var tagsResp OllamaTagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&tagsResp); err != nil {
		return nil, fmt.Errorf("failed to parse ollama response: %w", err)
	}

	var models []string
	for _, m := range tagsResp.Models {
		models = append(models, m.Name)
	}

	if len(models) == 0 {
		return nil, fmt.Errorf("no models found in Ollama registry")
	}

	return models, nil
}
