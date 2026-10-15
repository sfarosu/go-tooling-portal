package requestid

import "context"

type contextKey string

const (
	requestIDKey    contextKey = "request-id"
	mcpRequestIDKey contextKey = "mcp-request-id"
)

// WithContext returns a copy of ctx containing the request ID.
func WithContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// FromContext returns the request ID stored in ctx, or an empty string if absent.
func FromContext(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey).(string)
	return id
}

// WithMCPContext returns a copy of ctx containing the MCP request ID.
func WithMCPContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, mcpRequestIDKey, id)
}

// MCPFromContext returns the MCP request ID stored in ctx, or an empty string if absent.
func MCPFromContext(ctx context.Context) string {
	id, _ := ctx.Value(mcpRequestIDKey).(string)
	return id
}
