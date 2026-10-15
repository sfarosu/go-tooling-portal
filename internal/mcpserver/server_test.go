package mcpserver

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMCPGetExplainsEndpointRequirements(t *testing.T) {
	handler, err := NewHandler("http://127.0.0.1")
	if err != nil {
		t.Fatalf("NewHandler() error = %v", err)
	}

	request := httptest.NewRequest(http.MethodGet, "/mcp", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET /mcp returned status %d, want %d", response.Code, http.StatusMethodNotAllowed)
	}
	if got := response.Header().Get("Allow"); got != http.MethodPost {
		t.Fatalf("GET /mcp Allow = %q, want %q", got, http.MethodPost)
	}
	for _, message := range []string{
		"MCP Streamable HTTP endpoint",
		"Content-Type: application/json",
		"Accept: application/json, text/event-stream",
		"MCP-Protocol-Version",
	} {
		if !strings.Contains(response.Body.String(), message) {
			t.Errorf("GET /mcp response does not include %q: %s", message, response.Body.String())
		}
	}
}
