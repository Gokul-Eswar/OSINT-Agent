package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spectre/spectre/internal/agent"
	"github.com/spf13/viper"
)

func resetServerGlobals() {
	clientsMu.Lock()
	clients = make(map[chan interface{}]bool)
	clientsMu.Unlock()

	enginesMu.Lock()
	engines = make(map[string]*agent.Engine)
	enginesMu.Unlock()

	viper.Reset()
}

func TestBroadcast_DeliversToBufferedClient(t *testing.T) {
	resetServerGlobals()

	ch := make(chan interface{}, 1)
	clientsMu.Lock()
	clients[ch] = true
	clientsMu.Unlock()

	msg := map[string]string{"type": "ping"}
	Broadcast(msg)

	select {
	case got := <-ch:
		if got.(map[string]string)["type"] != "ping" {
			t.Fatalf("unexpected message payload: %#v", got)
		}
	default:
		t.Fatal("expected broadcast message to be delivered")
	}
}

func TestHandleCases_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/cases", nil)
	rr := httptest.NewRecorder()

	handleCases(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestHandleCaseDetail_InvalidPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/cases", nil)
	rr := httptest.NewRecorder()

	handleCaseDetail(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleSettings_Get(t *testing.T) {
	resetServerGlobals()
	viper.Set("ghost_mode", true)
	viper.Set("logging.level", "debug")
	viper.Set("collectors", map[string]interface{}{"dns": true})
	viper.Set("ethics", map[string]interface{}{"allow_active": false})

	req := httptest.NewRequest(http.MethodGet, "/api/settings", nil)
	rr := httptest.NewRecorder()

	handleSettings(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rr.Code)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if payload["ghost_mode"] != true {
		t.Fatalf("expected ghost_mode=true, got %#v", payload["ghost_mode"])
	}
	if payload["logging"] != "debug" {
		t.Fatalf("expected logging=debug, got %#v", payload["logging"])
	}
}

func TestHandleSettings_PostInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/settings", strings.NewReader("{"))
	rr := httptest.NewRecorder()

	handleSettings(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleSettings_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPatch, "/api/settings", nil)
	rr := httptest.NewRecorder()

	handleSettings(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestHandleChat_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/chat", nil)
	rr := httptest.NewRecorder()

	handleChat(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestHandleChat_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/chat", strings.NewReader("{"))
	rr := httptest.NewRecorder()

	handleChat(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleChat_MissingCaseID(t *testing.T) {
	body := `{"message":"hello"}`
	req := httptest.NewRequest(http.MethodPost, "/api/chat", strings.NewReader(body))
	rr := httptest.NewRecorder()

	handleChat(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rr.Code)
	}
}
