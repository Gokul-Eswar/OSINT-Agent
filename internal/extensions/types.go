package extensions

// Extension represents a purchasable/installable unit in the store.
type Extension struct {
	Name        string   `json:"name" yaml:"name"`
	Description string   `json:"description" yaml:"description"`
	Author      string   `json:"author" yaml:"author"`
	Version     string   `json:"version" yaml:"version"`
	URL         string   `json:"url" yaml:"url"`   // Git URL or Download URL
	Type        string   `json:"type" yaml:"type"` // e.g., "collector", "reporter"
	Tags        []string `json:"tags" yaml:"tags"`
}

// Registry is a collection of extensions.
type Registry struct {
	Extensions []Extension `json:"extensions" yaml:"extensions"`
}
