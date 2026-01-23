package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OllamaProvider struct {
	Endpoint string
	Model    string
	Client   *http.Client
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// NewOllamaProvider creates a new instance of the Ollama provider.
func NewOllamaProvider(endpoint, model string) *OllamaProvider {
	if endpoint == "" {
		endpoint = "http://localhost:11434/api/generate"
	}
	if model == "" {
		model = "llama3" // Default to a common model, can be overridden
	}
	return &OllamaProvider{
		Endpoint: endpoint,
		Model:    model,
		Client:   &http.Client{Timeout: 120 * time.Second}, // Generous timeout for analysis
	}
}

func (p *OllamaProvider) Name() string {
	return "ollama"
}

func (p *OllamaProvider) Generate(prompt string) (string, error) {
	reqBody := ollamaRequest{
		Model:  p.Model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", p.Endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	var ollamaResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return ollamaResp.Response, nil
}
