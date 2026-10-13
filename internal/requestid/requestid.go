package requestid

import "context"

type contextKey struct{}

// WithContext returns a copy of ctx containing the request ID.
func WithContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextKey{}, id)
}

// FromContext returns the request ID stored in ctx, or an empty string if absent.
func FromContext(ctx context.Context) string {
	id, _ := ctx.Value(contextKey{}).(string)
	return id
}
