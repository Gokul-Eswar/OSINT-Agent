package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/spectre/spectre/internal/agent"
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

	// Auth Middleware
	withAuth := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			apiKey := viper.GetString("server.api_key")
			if apiKey != "" {
				providedKey := r.Header.Get("X-API-Key")
				if providedKey == "" {
					providedKey = r.URL.Query().Get("api_key")
				}

				if providedKey != apiKey {
					http.Error(w, "Unauthorized: Invalid API Key", http.StatusUnauthorized)
					return
				}
			}
			h(w, r)
		}
	}

	// API Routes (Protected)
	mux.HandleFunc("/api/cases", withAuth(handleCases))
	mux.HandleFunc("/api/cases/", withAuth(handleCaseDetail))
	mux.HandleFunc("/api/events", withAuth(handleEvents))
	mux.HandleFunc("/api/settings", withAuth(handleSettings))
	mux.HandleFunc("/api/chat", withAuth(handleChat))

	// Static Assets (Public)
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

	// Serve Evidence Files (Protected)
	fs := http.FileServer(http.Dir("evidence_storage"))
	mux.Handle("/evidence/", withAuth(func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/evidence/", fs).ServeHTTP(w, r)
	}))

	// Hook into storage events
	storage.OnEntityCreated = func(e *core.Entity) {
		Broadcast(map[string]interface{}{
			"type": "entity_created",
			"data": e,
		})
	}

	fmt.Printf("SPECTRE API Server starting on 127.0.0.1:%d...\n", port)
	return http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), mux)
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

var (
	engines   = make(map[string]*agent.Engine)
	enginesMu sync.Mutex
)

func handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		CaseID  string `json:"case_id"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if payload.CaseID == "" {
		http.Error(w, "case_id is required", http.StatusBadRequest)
		return
	}

	enginesMu.Lock()
	engine, ok := engines[payload.CaseID]
	if !ok {
		engine = agent.NewEngine(payload.CaseID)
		engines[payload.CaseID] = engine
	}
	enginesMu.Unlock()

	response, err := engine.Execute(payload.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"response": response,
	})
}
