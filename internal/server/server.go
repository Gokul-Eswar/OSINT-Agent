package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/spectre/spectre/internal/analysis"
	"github.com/spectre/spectre/internal/core"
	"github.com/spectre/spectre/internal/storage"
	"github.com/spf13/viper"
)

var (
	webAssets embed.FS
	clients   = make(map[chan interface{}]bool)
	clientsMu sync.Mutex
)

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
	mux.HandleFunc("/api/events", handleEvents)
	mux.HandleFunc("/api/settings", handleSettings)

	// Static Assets
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api") && !strings.HasPrefix(r.URL.Path, "/evidence") {
			data, err := webAssets.ReadFile("web/index.html")
			if err != nil {
				http.Error(w, "Web assets not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.Write(data)
			return
		}
		if !strings.HasPrefix(r.URL.Path, "/evidence") {
			http.NotFound(w, r)
		}
	})

	// Serve Evidence Files
	fs := http.FileServer(http.Dir("evidence_storage"))
	mux.Handle("/evidence/", http.StripPrefix("/evidence/", fs))

	// Hook into storage events
	storage.OnEntityCreated = func(e *core.Entity) {
		Broadcast(map[string]interface{}{
			"type": "entity_created",
			"data": e,
		})
	}

	fmt.Printf("SPECTRE API Server starting on :%d...\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

func handleEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	messageChan := make(chan interface{})
	clientsMu.Lock()
	clients[messageChan] = true
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients, messageChan)
		clientsMu.Unlock()
		close(messageChan)
	}()

	notify := r.Context().Done()

	for {
		select {
		case msg := <-messageChan:
			data, err := json.Marshal(msg)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "data: %s\n\n", data)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case <-notify:
			return
		}
	}
}

// Broadcast sends a message to all connected SSE clients
func Broadcast(msg interface{}) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for client := range clients {
		select {
		case client <- msg:
		default:
			// Client channel full, skip to avoid blocking
		}
	}
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

func handleSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		settings := map[string]interface{}{
			"ghost_mode": viper.GetBool("ghost_mode"),
			"logging":    viper.GetString("logging.level"),
			"collectors": viper.GetStringMap("collectors"),
			"ethics":     viper.GetStringMap("ethics"),
		}
		json.NewEncoder(w).Encode(settings)
		return
	}

	if r.Method == http.MethodPost {
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Update Settings
		if val, ok := payload["ghost_mode"]; ok {
			viper.Set("ghost_mode", val)
		}
		if val, ok := payload["logging"]; ok {
			viper.Set("logging.level", val)
		}
		// Handle nested collectors if needed, for now simplified to top-level knowns
		// Deep merging map structures with viper can be tricky, so we might need more specific handling
		// if we allow editing complex objects. For now, simple toggles.

		// Save Config
		if err := viper.WriteConfig(); err != nil {
			// If no config file exists yet, safe write
			if err = viper.SafeWriteConfig(); err != nil {
				http.Error(w, "Failed to save config: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
