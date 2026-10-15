package metrics

import (
	"net/http"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	hostname, _ = os.Hostname()
	service     = "go-tooling-portal"
	namespace   = "go-tooling-portal"

	registry = prometheus.NewRegistry()

	APIRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "api",
			Name:      "requests_total",
			Help:      "Total number of completed API requests.",
		},
		[]string{"hostname", "service", "endpoint"},
	)
	MCPToolCallsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "mcp",
			Name:      "tool_calls_total",
			Help:      "Total number of completed MCP tool calls.",
		},
		[]string{"hostname", "service", "tool"},
	)
)

func init() {
	registry.MustRegister(APIRequestsTotal, MCPToolCallsTotal)
}

// RecordAPIRequest increments the request counter for an API endpoint.
func RecordAPIRequest(endpoint string) {
	APIRequestsTotal.WithLabelValues(hostname, service, endpoint).Inc()
}

// InitializeAPIEndpoints creates a zero-valued counter series for each API endpoint.
func InitializeAPIEndpoints(endpoints []string) {
	for _, endpoint := range endpoints {
		if strings.HasPrefix(endpoint, "/api/") {
			APIRequestsTotal.WithLabelValues(hostname, service, endpoint).Add(0)
		}
	}
}

// RecordMCPToolCall increments the call counter for an MCP tool.
func RecordMCPToolCall(tool string) {
	MCPToolCallsTotal.WithLabelValues(hostname, service, tool).Inc()
}

// InitializeMCPTool creates a zero-valued counter series for a registered MCP tool.
func InitializeMCPTool(tool string) {
	MCPToolCallsTotal.WithLabelValues(hostname, service, tool).Add(0)
}

// Handler exposes the registered metrics in Prometheus text format.
func Handler() http.Handler {
	return promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
}
