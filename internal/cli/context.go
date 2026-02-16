package cli

import (
	"os"
	"path/filepath"
	"strings"
)

const contextFile = ".spectre_current_case"
const targetFile = ".spectre_current_target"

func getSpectreHome() string {
	if h := os.Getenv("SPECTRE_HOME"); h != "" {
		return h
	}
	home, _ := os.UserHomeDir()
	return home
}

// SaveContext saves the current case ID to a file in the user's home directory.
func SaveContext(caseID string) error {
	path := filepath.Join(getSpectreHome(), contextFile)
	return os.WriteFile(path, []byte(caseID), 0644)
}

// LoadContext retrieves the most recent case ID.
func LoadContext() (string, error) {
	path := filepath.Join(getSpectreHome(), contextFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

// SaveTarget saves the current target to a file in the user's home directory.
func SaveTarget(target string) error {
	path := filepath.Join(getSpectreHome(), targetFile)
	return os.WriteFile(path, []byte(target), 0644)
}

// LoadTarget retrieves the most recent target.
func LoadTarget() (string, error) {
	path := filepath.Join(getSpectreHome(), targetFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}


