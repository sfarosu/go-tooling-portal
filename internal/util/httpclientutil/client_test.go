package httpclientutil

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
)

func TestNewClientUsesDefaults(t *testing.T) {
	client, err := NewClient(Config{})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if got := client.GetClient().Timeout; got != 60*time.Second {
		t.Errorf("timeout = %v, want %v", got, 60*time.Second)
	}
	if client.RetryCount != 3 {
		t.Errorf("retry count = %d, want 3", client.RetryCount)
	}
	if client.RetryWaitTime != time.Second {
		t.Errorf("retry wait time = %v, want %v", client.RetryWaitTime, time.Second)
	}
	if client.RetryMaxWaitTime != 3*time.Second {
		t.Errorf("maximum retry wait time = %v, want %v", client.RetryMaxWaitTime, 3*time.Second)
	}
	if client.ResponseBodyLimit != 100*1024*1024 {
		t.Errorf("maximum response body size = %d, want %d", client.ResponseBodyLimit, 100*1024*1024)
	}
}

func TestNewClientUsesConfiguredValues(t *testing.T) {
	config := Config{
		Timeout:          10 * time.Second,
		RetryCount:       5,
		RetryWaitTime:    2 * time.Second,
		RetryMaxWaitTime: 8 * time.Second,
		MaxBodyBytes:     1024,
	}
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if got := client.GetClient().Timeout; got != config.Timeout {
		t.Errorf("timeout = %v, want %v", got, config.Timeout)
	}
	if client.RetryCount != config.RetryCount {
		t.Errorf("retry count = %d, want %d", client.RetryCount, config.RetryCount)
	}
	if client.RetryWaitTime != config.RetryWaitTime {
		t.Errorf("retry wait time = %v, want %v", client.RetryWaitTime, config.RetryWaitTime)
	}
	if client.RetryMaxWaitTime != config.RetryMaxWaitTime {
		t.Errorf("maximum retry wait time = %v, want %v", client.RetryMaxWaitTime, config.RetryMaxWaitTime)
	}
	if client.ResponseBodyLimit != config.MaxBodyBytes {
		t.Errorf("maximum response body size = %d, want %d", client.ResponseBodyLimit, config.MaxBodyBytes)
	}
}

func TestNewClientRejectsNegativeValues(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{name: "timeout", config: Config{Timeout: -time.Second}},
		{name: "retry count", config: Config{RetryCount: -1}},
		{name: "retry wait", config: Config{RetryWaitTime: -time.Second}},
		{name: "maximum retry wait", config: Config{RetryMaxWaitTime: -time.Second}},
		{name: "maximum body size", config: Config{MaxBodyBytes: -1}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := NewClient(test.config); err == nil {
				t.Fatal("NewClient() error = nil, want negative setting error")
			}
		})
	}
}

func TestNewClientEnforcesMaximumResponseBodySize(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("response exceeds limit"))
	}))
	defer server.Close()

	client, err := NewClient(Config{
		RetryCount:       1,
		RetryWaitTime:    time.Nanosecond,
		RetryMaxWaitTime: time.Nanosecond,
		MaxBodyBytes:     4,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	_, err = client.R().Get(server.URL)
	if !errors.Is(err, resty.ErrResponseBodyTooLarge) {
		t.Fatalf("request error = %v, want %v", err, resty.ErrResponseBodyTooLarge)
	}
}
