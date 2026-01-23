package ai

// Provider defines the interface for AI backends (e.g., Ollama, OpenAI).
type Provider interface {
	// Generate sends a prompt to the AI and returns the text response.
	Generate(prompt string) (string, error)
	// Name returns the provider name.
	Name() string
}
