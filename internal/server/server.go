package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/spectre/spectre/internal/analysis"
	"github.com/spectre/spectre/internal/storage"
)

var webAssets embed.FS

// SetAssets sets the embedded assets for the server
func SetAssets(assets embed.FS) {
	webAssets = assets
}

// Start starts the API server
func Start(port int) error {
	mux := http.NewServeMux()

	// API Routes
	mux.HandleFunc("/api/cases", handleCases)
	mux.HandleFunc("/api/cases/", handleCaseDetail) // /api/cases/{id} and /api/cases/{id}/graph

	// Static Assets
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api") {
			data, err := webAssets.ReadFile("web/index.html")
			if err != nil {
				http.Error(w, "Web assets not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.Write(data)
			return
		}
		http.NotFound(w, r)
	})

	fmt.Printf("SPECTRE API Server starting on :%d...\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

func handleCases(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cases, err := storage.ListCases()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cases)
}

func handleCaseDetail(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	caseID := parts[3]

	// Check if it's a graph request
	if len(parts) > 4 && parts[4] == "graph" {
		data, err := analysis.ExportCaseForViz(caseID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
		return
	}

	// Normal case detail
	c, err := storage.GetCase(caseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if c == nil {
		http.Error(w, "Case not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}