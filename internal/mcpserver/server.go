package mcpserver

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sfarosu/go-tooling-portal/internal/requestid"
	"github.com/sfarosu/go-tooling-portal/internal/util/httpclientutil"
	"github.com/sfarosu/go-tooling-portal/internal/version"
)

// NewHandler creates a stateless Streamable HTTP MCP handler that proxies tools to the API.
func NewHandler(apiBaseURL string) (http.Handler, error) {
	baseURL, err := url.Parse(apiBaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse API base URL: %w", err)
	}
	if (baseURL.Scheme != "http" && baseURL.Scheme != "https") || baseURL.Host == "" {
		return nil, fmt.Errorf("API base URL must be an absolute HTTP or HTTPS URL")
	}
	baseURL.Path = strings.TrimRight(baseURL.Path, "/")

	client, err := httpclientutil.NewClient(httpclientutil.Config{})
	if err != nil {
		return nil, fmt.Errorf("create API HTTP client: %w", err)
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "go-tooling-portal",
		Version: version.Version,
	}, nil)
	if err := registerTools(server, client, baseURL.String()); err != nil {
		return nil, fmt.Errorf("register MCP tools: %w", err)
	}

	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{
		Stateless:                    true,
		PropagateRequestCancellation: true,
	})
	return ensureMCPRequestID(explainUnsupportedMethods(handler)), nil
}

func explainUnsupportedMethods(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(
				w,
				"This is an MCP Streamable HTTP endpoint, not a web page. Use an MCP client that sends JSON-RPC POST requests with Content-Type: application/json and Accept: application/json, text/event-stream. MCP-Protocol-Version is required after initialization.",
				http.StatusMethodNotAllowed,
			)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ensureMCPRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mcpRequestID := requestid.MCPFromContext(r.Context())
		if mcpRequestID == "" {
			mcpRequestID = r.Header.Get("MCP-Request-ID")
		}
		if _, err := uuid.Parse(mcpRequestID); err != nil {
			mcpRequestID = uuid.NewString()
		}

		w.Header().Set("MCP-Request-ID", mcpRequestID)
		ctx := requestid.WithContext(r.Context(), mcpRequestID)
		ctx = requestid.WithMCPContext(ctx, mcpRequestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
