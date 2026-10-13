package requestid

import (
	"context"
	"testing"
)

func TestContextRoundTrip(t *testing.T) {
	ctx := WithContext(context.Background(), "test-request-id")
	if got := FromContext(ctx); got != "test-request-id" {
		t.Fatalf("FromContext() = %q, want %q", got, "test-request-id")
	}
}

func TestFromContextWithoutID(t *testing.T) {
	if got := FromContext(context.Background()); got != "" {
		t.Fatalf("FromContext() = %q, want empty", got)
	}
}
