package httpclientutil

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	defaultTimeout       = 60 * time.Second
	defaultRetryCount    = 3
	defaultRetryWaitTime = time.Second
	defaultRetryMaxWait  = 3 * time.Second
	defaultMaxBodyBytes  = 100 * 1024 * 1024
)

// Config configures a Resty HTTP client. Zero-valued fields use their defaults.
type Config struct {
	Timeout          time.Duration // Default: 60 seconds.
	RetryCount       int           // Default: 3 retries.
	RetryWaitTime    time.Duration // Default: 1 second.
	RetryMaxWaitTime time.Duration // Default: 3 seconds.
	MaxBodyBytes     int           // Maximum response body size; default: 100 MB.
}

// NewClient creates a Resty client configured with the supplied settings.
// Zero-valued settings use defaults; negative settings are rejected.
func NewClient(config Config) (*resty.Client, error) {
	if config.Timeout < 0 {
		return nil, fmt.Errorf("timeout cannot be negative")
	}
	if config.RetryCount < 0 {
		return nil, fmt.Errorf("retry count cannot be negative")
	}
	if config.RetryWaitTime < 0 {
		return nil, fmt.Errorf("retry wait time cannot be negative")
	}
	if config.RetryMaxWaitTime < 0 {
		return nil, fmt.Errorf("maximum retry wait time cannot be negative")
	}
	if config.MaxBodyBytes < 0 {
		return nil, fmt.Errorf("maximum body size cannot be negative")
	}

	if config.Timeout == 0 {
		config.Timeout = defaultTimeout
	}
	if config.RetryCount == 0 {
		config.RetryCount = defaultRetryCount
	}
	if config.RetryWaitTime == 0 {
		config.RetryWaitTime = defaultRetryWaitTime
	}
	if config.RetryMaxWaitTime == 0 {
		config.RetryMaxWaitTime = defaultRetryMaxWait
	}
	if config.MaxBodyBytes == 0 {
		config.MaxBodyBytes = defaultMaxBodyBytes
	}

	return resty.New().
		SetTimeout(config.Timeout).
		SetRetryCount(config.RetryCount).
		SetRetryWaitTime(config.RetryWaitTime).
		SetRetryMaxWaitTime(config.RetryMaxWaitTime).
		SetResponseBodyLimit(config.MaxBodyBytes), nil
}
