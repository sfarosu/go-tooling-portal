package mcpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/metrics"
)

func TestMCPToolsProxyToAPIWithRequestIDs(t *testing.T) {
	previousLogger := logger.Logger
	var logOutput bytes.Buffer
	logger.Logger = slog.New(slog.NewTextHandler(&logOutput, nil))
	t.Cleanup(func() {
		logger.Logger = previousLogger
	})

	var receivedBody map[string]any
	receivedMCPRequestIDs := make(map[string]string)
	receivedAPIRequestIDs := make(map[string]string)
	receivedAPICalls := 0
	var receivedMu sync.Mutex
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMu.Lock()
		receivedAPICalls++
		receivedMCPRequestIDs[r.URL.Path] = r.Header.Get("MCP-Request-ID")
		receivedAPIRequestIDs[r.URL.Path] = r.Header.Get("X-Request-ID")
		defer receivedMu.Unlock()
		if r.URL.Path == "/api/formatconverter" {
			http.Error(w, `{"detail":"invalid input"}`, http.StatusBadRequest)
			return
		}

		if r.URL.Path != "/api/base64converter" {
			t.Errorf("API path = %q, want /api/base64converter", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&receivedBody); err != nil {
			t.Errorf("decode API request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"YWJj"}`))
	}))
	defer api.Close()

	handler, err := NewHandler(api.URL)
	if err != nil {
		t.Fatalf("NewHandler() error = %v", err)
	}
	mcpHTTP := httptest.NewServer(handler)
	defer mcpHTTP.Close()

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "1.0.0"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{
		Endpoint:             mcpHTTP.URL,
		DisableStandaloneSSE: true,
	}, nil)
	if err != nil {
		t.Fatalf("connect MCP client: %v", err)
	}
	defer session.Close()

	tools, err := session.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatalf("list MCP tools: %v", err)
	}
	names := make([]string, 0, len(tools.Tools))
	for _, tool := range tools.Tools {
		names = append(names, tool.Name)
	}
	sort.Strings(names)
	wantNames := []string{
		"base64-converter",
		"format-converter",
		"generate-htpasswd",
		"generate-password",
		"generate-ssh-key",
		"get-version",
		"json-pretty",
		"time-converter",
		"url-decoder",
	}
	if !reflect.DeepEqual(names, wantNames) {
		t.Fatalf("MCP tools = %v, want %v", names, wantNames)
	}
	metricsResponse := httptest.NewRecorder()
	metrics.Handler().ServeHTTP(metricsResponse, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	for _, name := range wantNames {
		if !strings.Contains(metricsResponse.Body.String(), `tool="`+name+`"`) {
			t.Errorf("MCP metrics missing zero-value tool series %q: %s", name, metricsResponse.Body.String())
		}
	}
	urlDecoderSchemaFound := false
	for _, tool := range tools.Tools {
		if tool.Name != "url-decoder" {
			continue
		}
		urlDecoderSchemaFound = true
		schema, ok := tool.InputSchema.(map[string]any)
		if !ok {
			t.Fatalf("url-decoder input schema has type %T, want map[string]any", tool.InputSchema)
		}
		properties, ok := schema["properties"].(map[string]any)
		if !ok {
			t.Fatalf("url-decoder properties have type %T, want map[string]any", schema["properties"])
		}
		urlSchema, ok := properties["url"].(map[string]any)
		if !ok {
			t.Fatalf("url-decoder url property has type %T, want map[string]any", properties["url"])
		}
		if got, want := urlSchema["default"], "https://example.com/foo+bar%21?service=ab%2Bc&product=de%24f&service=second"; got != want {
			t.Errorf("url-decoder url default = %v, want %q", got, want)
		}
	}
	if !urlDecoderSchemaFound {
		t.Fatal("url-decoder tool was not returned by tools/list")
	}
	formatConverterSchemaFound := false
	for _, tool := range tools.Tools {
		if tool.Name != "format-converter" {
			continue
		}
		formatConverterSchemaFound = true
		schema, ok := tool.InputSchema.(map[string]any)
		if !ok {
			t.Fatalf("format-converter input schema has type %T, want map[string]any", tool.InputSchema)
		}
		properties, ok := schema["properties"].(map[string]any)
		if !ok {
			t.Fatalf("format-converter properties have type %T, want map[string]any", schema["properties"])
		}
		inputSchema, ok := properties["input"].(map[string]any)
		if !ok {
			t.Fatalf("format-converter input property has type %T, want map[string]any", properties["input"])
		}
		if got, want := inputSchema["default"], `{"replace":"me"}`; got != want {
			t.Errorf("format-converter input default = %v, want %q", got, want)
		}
	}
	if !formatConverterSchemaFound {
		t.Fatal("format-converter tool was not returned by tools/list")
	}
	htpasswdSchemaFound := false
	for _, tool := range tools.Tools {
		if tool.Name != "generate-htpasswd" {
			continue
		}
		htpasswdSchemaFound = true
		schema, ok := tool.InputSchema.(map[string]any)
		if !ok {
			t.Fatalf("generate-htpasswd input schema has type %T, want map[string]any", tool.InputSchema)
		}
		properties, ok := schema["properties"].(map[string]any)
		if !ok {
			t.Fatalf("generate-htpasswd properties have type %T, want map[string]any", schema["properties"])
		}
		for property, want := range map[string]any{
			"username":  "replace-me",
			"password":  "replace-me",
			"algorithm": "apr1",
		} {
			propertySchema, ok := properties[property].(map[string]any)
			if !ok {
				t.Fatalf("generate-htpasswd property %q has type %T, want map[string]any", property, properties[property])
			}
			if got := propertySchema["default"]; got != want {
				t.Errorf("generate-htpasswd property %q default = %v, want %v", property, got, want)
			}
		}
	}
	if !htpasswdSchemaFound {
		t.Fatal("generate-htpasswd tool was not returned by tools/list")
	}
	sshKeySchemaFound := false
	for _, tool := range tools.Tools {
		if tool.Name != "generate-ssh-key" {
			continue
		}
		sshKeySchemaFound = true
		schema, ok := tool.InputSchema.(map[string]any)
		if !ok {
			t.Fatalf("generate-ssh-key input schema has type %T, want map[string]any", tool.InputSchema)
		}
		properties, ok := schema["properties"].(map[string]any)
		if !ok {
			t.Fatalf("generate-ssh-key properties have type %T, want map[string]any", schema["properties"])
		}
		for property, want := range map[string]any{
			"algorithm": "ed25519",
			"ecdsaBits": "521",
			"rsaBits":   "4096",
			"email":     "replace-me@example.com",
		} {
			propertySchema, ok := properties[property].(map[string]any)
			if !ok {
				t.Fatalf("generate-ssh-key property %q has type %T, want map[string]any", property, properties[property])
			}
			if got := propertySchema["default"]; got != want {
				t.Errorf("generate-ssh-key property %q default = %v, want %v", property, got, want)
			}
		}
	}
	if !sshKeySchemaFound {
		t.Fatal("generate-ssh-key tool was not returned by tools/list")
	}
	jsonPrettySchemaFound := false
	for _, tool := range tools.Tools {
		if tool.Name != "json-pretty" {
			continue
		}
		jsonPrettySchemaFound = true
		schema, ok := tool.InputSchema.(map[string]any)
		if !ok {
			t.Fatalf("json-pretty input schema has type %T, want map[string]any", tool.InputSchema)
		}
		properties, ok := schema["properties"].(map[string]any)
		if !ok {
			t.Fatalf("json-pretty properties have type %T, want map[string]any", schema["properties"])
		}
		for property, want := range map[string]any{
			"input":     `{"key":"value"}`,
			"operation": "pretty",
			"spaces":    float64(2),
		} {
			propertySchema, ok := properties[property].(map[string]any)
			if !ok {
				t.Fatalf("json-pretty property %q has type %T, want map[string]any", property, properties[property])
			}
			if got := propertySchema["default"]; got != want {
				t.Errorf("json-pretty property %q default = %v, want %v", property, got, want)
			}
		}
	}
	if !jsonPrettySchemaFound {
		t.Fatal("json-pretty tool was not returned by tools/list")
	}
	for _, tool := range tools.Tools {
		if tool.Name != "base64-converter" {
			continue
		}
		schema, ok := tool.InputSchema.(map[string]any)
		if !ok {
			t.Fatalf("base64-converter input schema has type %T, want map[string]any", tool.InputSchema)
		}
		properties, ok := schema["properties"].(map[string]any)
		if !ok {
			t.Fatalf("base64-converter properties have type %T, want map[string]any", schema["properties"])
		}
		for property, want := range map[string]any{
			"input":     "replace-me",
			"operation": "encode",
			"format":    "standard",
		} {
			propertySchema, ok := properties[property].(map[string]any)
			if !ok {
				t.Fatalf("base64-converter property %q has type %T, want map[string]any", property, properties[property])
			}
			if got := propertySchema["default"]; got != want {
				t.Errorf("base64-converter property %q default = %v, want %v", property, got, want)
			}
		}
	}
	for _, tool := range tools.Tools {
		if tool.Name != "generate-password" {
			continue
		}
		schema, ok := tool.InputSchema.(map[string]any)
		if !ok {
			t.Fatalf("generate-password input schema has type %T, want map[string]any", tool.InputSchema)
		}
		if _, ok := schema["anyOf"]; ok {
			t.Fatal("generate-password input schema exposes top-level anyOf variants")
		}
		properties, ok := schema["properties"].(map[string]any)
		if !ok {
			t.Fatalf("generate-password properties have type %T, want map[string]any", schema["properties"])
		}
		for property, want := range map[string]any{
			"length":    float64(16),
			"uppercase": true,
			"lowercase": true,
			"numbers":   true,
			"symbols":   true,
		} {
			propertySchema, ok := properties[property].(map[string]any)
			if !ok {
				t.Fatalf("generate-password property %q has type %T, want map[string]any", property, properties[property])
			}
			if got := propertySchema["default"]; got != want {
				t.Errorf("generate-password property %q default = %v, want %v", property, got, want)
			}
		}
	}

	result, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "base64-converter",
		Arguments: map[string]any{
			"input":     "abc",
			"operation": "encode",
			"format":    "standard",
		},
	})
	if err != nil {
		t.Fatalf("call base64-converter: %v", err)
	}
	if result.IsError {
		t.Fatalf("tool call returned an error: %#v", result.Content)
	}
	receivedMu.Lock()
	receivedMCPRequestID := receivedMCPRequestIDs["/api/base64converter"]
	receivedAPIRequestID := receivedAPIRequestIDs["/api/base64converter"]
	body := receivedBody
	receivedMu.Unlock()
	if receivedMCPRequestID == "" {
		t.Fatal("API request did not contain MCP-Request-ID")
	}
	if receivedAPIRequestID != receivedMCPRequestID {
		t.Fatalf("API request ID = %q, want same MCP correlation ID %q", receivedAPIRequestID, receivedMCPRequestID)
	}
	if body["input"] != "abc" || body["operation"] != "encode" || body["format"] != "standard" {
		t.Fatalf("API request body = %#v", body)
	}
	if got := result.Meta["request-id"]; got != receivedMCPRequestID {
		t.Fatalf("MCP response request ID = %v, want forwarded ID %q", got, receivedMCPRequestID)
	}
	if !bytes.Contains([]byte(result.Content[0].(*mcp.TextContent).Text), []byte("YWJj")) {
		t.Fatalf("MCP result content does not include API output: %#v", result.Content)
	}

	errorResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "format-converter",
		Arguments: map[string]any{
			"input": "invalid",
		},
	})
	if err != nil {
		t.Fatalf("call format-converter: %v", err)
	}
	if !errorResult.IsError {
		t.Fatalf("API error was not represented as an MCP tool error: %#v", errorResult)
	}
	receivedMu.Lock()
	errorMCPRequestID := receivedMCPRequestIDs["/api/formatconverter"]
	errorAPIRequestID := receivedAPIRequestIDs["/api/formatconverter"]
	receivedMu.Unlock()
	if errorMCPRequestID == "" || errorResult.Meta["request-id"] != errorMCPRequestID {
		t.Fatalf("MCP error result has incorrect MCP request ID: result=%v forwarded=%q", errorResult.Meta, errorMCPRequestID)
	}
	if errorAPIRequestID != errorMCPRequestID {
		t.Fatalf("API error request ID = %q, want same MCP correlation ID %q", errorAPIRequestID, errorMCPRequestID)
	}
	logs := logOutput.String()
	if !strings.Contains(logs, `msg=mcp-tool-failed`) ||
		!strings.Contains(logs, `request-id=`+errorMCPRequestID) ||
		!strings.Contains(logs, `tool-name=format-converter`) ||
		!strings.Contains(logs, `api-status=400`) {
		t.Fatalf("MCP tool failure log is missing identifiers or API status: %s", logs)
	}
	if strings.Contains(logs, "invalid input") {
		t.Fatalf("MCP tool failure log contains API error body: %s", logs)
	}

	invalidResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "base64-converter",
		Arguments: map[string]any{
			"input":     "abc",
			"operation": "invalid",
			"format":    "standard",
		},
	})
	if err != nil {
		t.Fatalf("call base64-converter with invalid operation: %v", err)
	}
	if !invalidResult.IsError {
		t.Fatalf("invalid enum value was not rejected by the MCP input schema: %#v", invalidResult)
	}
	receivedMu.Lock()
	apiCallsAfterInvalidInput := receivedAPICalls
	receivedMu.Unlock()
	if apiCallsAfterInvalidInput != 2 {
		t.Fatalf("invalid MCP input reached the API; received %d API calls, want 2", apiCallsAfterInvalidInput)
	}

	missingECDSABitsResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "generate-ssh-key",
		Arguments: map[string]any{
			"algorithm": "ecdsa",
			"ecdsaBits": "999",
		},
	})
	if err != nil {
		t.Fatalf("call generate-ssh-key with invalid ECDSA bits: %v", err)
	}
	if !missingECDSABitsResult.IsError {
		t.Fatalf("invalid ECDSA key size was not rejected by the MCP schema: %#v", missingECDSABitsResult)
	}
	receivedMu.Lock()
	apiCallsAfterMissingECDSABits := receivedAPICalls
	receivedMu.Unlock()
	if apiCallsAfterMissingECDSABits != 2 {
		t.Fatalf("invalid SSH input reached the API; received %d API calls, want 2", apiCallsAfterMissingECDSABits)
	}

	invalidPasswordResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "generate-password",
		Arguments: map[string]any{
			"length":    0,
			"uppercase": true,
			"lowercase": false,
			"numbers":   false,
			"symbols":   false,
		},
	})
	if err != nil {
		t.Fatalf("call generate-password with invalid length: %v", err)
	}
	if !invalidPasswordResult.IsError {
		t.Fatalf("invalid password length was not rejected by the MCP schema: %#v", invalidPasswordResult)
	}
	receivedMu.Lock()
	apiCallsAfterInvalidPassword := receivedAPICalls
	receivedMu.Unlock()
	if apiCallsAfterInvalidPassword != 2 {
		t.Fatalf("invalid password input reached the API; received %d API calls, want 2", apiCallsAfterInvalidPassword)
	}

	tooLongPasswordResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "generate-password",
		Arguments: map[string]any{
			"length":    1025,
			"uppercase": true,
			"lowercase": false,
			"numbers":   false,
			"symbols":   false,
		},
	})
	if err != nil {
		t.Fatalf("call generate-password with excessive length: %v", err)
	}
	if !tooLongPasswordResult.IsError {
		t.Fatalf("excessive password length was not rejected by the MCP schema: %#v", tooLongPasswordResult)
	}
	receivedMu.Lock()
	apiCallsAfterLongPassword := receivedAPICalls
	receivedMu.Unlock()
	if apiCallsAfterLongPassword != 2 {
		t.Fatalf("excessive password input reached the API; received %d API calls, want 2", apiCallsAfterLongPassword)
	}

	noCharacterClassesResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "generate-password",
		Arguments: map[string]any{
			"length":    12,
			"uppercase": false,
			"lowercase": false,
			"numbers":   false,
			"symbols":   false,
		},
	})
	if err != nil {
		t.Fatalf("call generate-password with no character classes: %v", err)
	}
	if !noCharacterClassesResult.IsError {
		t.Fatalf("password with no selected character classes was not rejected: %#v", noCharacterClassesResult)
	}
	receivedMu.Lock()
	apiCallsAfterNoCharacterClasses := receivedAPICalls
	receivedMu.Unlock()
	if apiCallsAfterNoCharacterClasses != 2 {
		t.Fatalf("invalid password input reached the API; received %d API calls, want 2", apiCallsAfterNoCharacterClasses)
	}
}

func TestTimeConverterSchemaRangesAreConditional(t *testing.T) {
	apiCalls := 0
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCalls++
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"timezone":"UTC"}`))
	}))
	defer api.Close()

	handler, err := NewHandler(api.URL)
	if err != nil {
		t.Fatalf("NewHandler() error = %v", err)
	}
	mcpHTTP := httptest.NewServer(handler)
	defer mcpHTTP.Close()

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "1.0.0"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{
		Endpoint:             mcpHTTP.URL,
		DisableStandaloneSSE: true,
	}, nil)
	if err != nil {
		t.Fatalf("connect MCP client: %v", err)
	}
	defer session.Close()

	tools, err := session.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatalf("list MCP tools: %v", err)
	}
	timeSchemaFound := false
	for _, tool := range tools.Tools {
		if tool.Name != "time-converter" {
			continue
		}
		timeSchemaFound = true
		schema, ok := tool.InputSchema.(map[string]any)
		if !ok {
			t.Fatalf("time-converter input schema has type %T, want map[string]any", tool.InputSchema)
		}
		if _, ok := schema["anyOf"]; ok {
			t.Fatal("time-converter input schema exposes top-level anyOf variants")
		}
	}
	if !timeSchemaFound {
		t.Fatal("time-converter tool was not returned by tools/list")
	}
	for _, tool := range tools.Tools {
		if tool.Name != "time-converter" {
			continue
		}
		schema := tool.InputSchema.(map[string]any)
		properties := schema["properties"].(map[string]any)
		for property, want := range map[string]any{
			"epochToHumanMode": true,
			"humanToEpochMode": false,
			"epochTime":        float64(0),
			"precision":        "seconds",
			"year":             float64(2025),
			"month":            float64(1),
			"day":              float64(1),
			"hour":             float64(0),
			"minute":           float64(0),
			"second":           float64(0),
			"millisecond":      float64(0),
			"microsecond":      float64(0),
			"timezone":         "UTC",
		} {
			propertySchema, ok := properties[property].(map[string]any)
			if !ok {
				t.Fatalf("time-converter property %q has type %T, want map[string]any", property, properties[property])
			}
			if got := propertySchema["default"]; got != want {
				t.Errorf("time-converter property %q default = %v, want %v", property, got, want)
			}
		}
	}

	epochResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "time-converter",
		Arguments: map[string]any{
			"epochToHumanMode": true,
			"humanToEpochMode": false,
			"epochTime":        0,
			"year":             0,
			"month":            0,
			"day":              0,
			"hour":             0,
			"minute":           0,
			"second":           0,
			"timezone":         "UTC",
		},
	})
	if err != nil {
		t.Fatalf("call time-converter in epoch mode: %v", err)
	}
	if epochResult.IsError {
		t.Fatalf("epoch conversion rejected unused human-time fields: %#v", epochResult)
	}

	bothModesResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "time-converter",
		Arguments: map[string]any{
			"epochToHumanMode": true,
			"humanToEpochMode": true,
			"epochTime":        0,
			"year":             2025,
			"month":            0,
			"day":              0,
			"hour":             0,
			"minute":           0,
			"second":           0,
			"timezone":         "UTC",
		},
	})
	if err != nil {
		t.Fatalf("call time-converter with both modes selected: %v", err)
	}
	if bothModesResult.IsError {
		t.Fatalf("epoch mode should take precedence over unused human-time fields: %#v", bothModesResult)
	}

	invalidHumanResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "time-converter",
		Arguments: map[string]any{
			"epochToHumanMode": false,
			"humanToEpochMode": true,
			"epochTime":        0,
			"year":             2025,
			"month":            0,
			"day":              1,
			"hour":             0,
			"minute":           0,
			"second":           0,
			"timezone":         "UTC",
		},
	})
	if err != nil {
		t.Fatalf("call time-converter with invalid month: %v", err)
	}
	if !invalidHumanResult.IsError {
		t.Fatalf("invalid month was not rejected by the MCP schema: %#v", invalidHumanResult)
	}
	if apiCalls != 2 {
		t.Fatalf("invalid human-time input reached the API; received %d calls, want 2", apiCalls)
	}

	noModeResult, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name: "time-converter",
		Arguments: map[string]any{
			"epochToHumanMode": false,
			"humanToEpochMode": false,
			"epochTime":        0,
			"year":             2025,
			"month":            1,
			"day":              1,
			"hour":             0,
			"minute":           0,
			"second":           0,
			"timezone":         "UTC",
		},
	})
	if err != nil {
		t.Fatalf("call time-converter with no mode selected: %v", err)
	}
	if !noModeResult.IsError {
		t.Fatalf("missing time conversion mode was not rejected by the MCP schema: %#v", noModeResult)
	}
	if apiCalls != 2 {
		t.Fatalf("time conversion without a mode reached the API; received %d calls, want 2", apiCalls)
	}
}
