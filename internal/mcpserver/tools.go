package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/metrics"
	"github.com/sfarosu/go-tooling-portal/internal/requestid"
)

type base64Input struct {
	Input     string `json:"input" jsonschema:"String to encode or decode"`
	Operation string `json:"operation" jsonschema:"Operation to perform: encode or decode"`
	Format    string `json:"format" jsonschema:"Base64 format: standard or url-compatible"`
}

type formatConverterInput struct {
	Input string `json:"input" jsonschema:"JSON or YAML string to convert"`
}

type htpasswdInput struct {
	Username  string `json:"username" jsonschema:"Username for the htpasswd entry"`
	Password  string `json:"password" jsonschema:"Password for the htpasswd entry"`
	Algorithm string `json:"algorithm" jsonschema:"Hash algorithm: apr1, 1, 5, or 6"`
}

type jsonPrettyInput struct {
	Input     string `json:"input" jsonschema:"JSON string to format or minify"`
	Operation string `json:"operation" jsonschema:"Operation to perform: pretty or unpretty"`
	Spaces    int    `json:"spaces" jsonschema:"Number of spaces for indentation"`
}

type passwordInput struct {
	Length    int  `json:"length" jsonschema:"Generated password length"`
	Uppercase bool `json:"uppercase" jsonschema:"Include uppercase letters"`
	Lowercase bool `json:"lowercase" jsonschema:"Include lowercase letters"`
	Numbers   bool `json:"numbers" jsonschema:"Include digits"`
	Symbols   bool `json:"symbols" jsonschema:"Include symbols"`
}

type sshKeyInput struct {
	Algorithm string `json:"algorithm" jsonschema:"SSH key algorithm: ed25519, ecdsa, or rsa"`
	ECDSABits string `json:"ecdsaBits,omitempty" jsonschema:"ECDSA key size: 256, 384, or 521"`
	RSABits   string `json:"rsaBits,omitempty" jsonschema:"RSA key size: 1024, 2048, 3072, or 4096"`
	Email     string `json:"email,omitempty" jsonschema:"Email or comment for the public key"`
}

type timeConverterInput struct {
	EpochToHumanMode bool   `json:"epochToHumanMode" jsonschema:"Convert from epoch to human-readable time"`
	HumanToEpochMode bool   `json:"humanToEpochMode" jsonschema:"Convert from human-readable time to epoch"`
	EpochTime        int64  `json:"epochTime" jsonschema:"Epoch timestamp to convert"`
	Precision        string `json:"precision,omitempty" jsonschema:"Epoch unit: seconds, millis, or micros"`
	Year             int    `json:"year" jsonschema:"Year for human-readable time conversion"`
	Month            int    `json:"month" jsonschema:"Month for human-readable time conversion (1-12)"`
	Day              int    `json:"day" jsonschema:"Day for human-readable time conversion (1-31)"`
	Hour             int    `json:"hour" jsonschema:"Hour for human-readable time conversion (0-23)"`
	Minute           int    `json:"minute" jsonschema:"Minute for human-readable time conversion (0-59)"`
	Second           int    `json:"second" jsonschema:"Second for human-readable time conversion (0-59)"`
	Millisecond      int    `json:"millisecond,omitempty" jsonschema:"Millisecond portion (0-999)"`
	Microsecond      int    `json:"microsecond,omitempty" jsonschema:"Additional microseconds (0-999)"`
	Timezone         string `json:"timezone" jsonschema:"IANA timezone identifier such as UTC or Europe/Paris"`
}

type urlDecoderInput struct {
	URL string `json:"url" jsonschema:"URL to parse and decode"`
}

type emptyInput struct{}

func registerTools(server *mcp.Server, client *resty.Client, apiBaseURL string) error {
	registrations := []func() error{
		func() error {
			return registerTool[base64Input](server, client, apiBaseURL, "base64-converter", "Encode or decode Base64 using standard or URL-safe encoding.", http.MethodPost, "/api/base64converter", func(schema *jsonschema.Schema) {
				setEnum(schema, "operation", "encode", "decode")
				setEnum(schema, "format", "standard", "url-compatible")
				setMinLength(schema, "input", 1)
				schema.Properties["input"].Default = json.RawMessage(`"replace-me"`)
				schema.Properties["operation"].Default = json.RawMessage(`"encode"`)
				schema.Properties["format"].Default = json.RawMessage(`"standard"`)
			})
		},
		func() error {
			return registerTool[formatConverterInput](server, client, apiBaseURL, "format-converter", "Convert a JSON string to YAML or a YAML string to JSON.", http.MethodPost, "/api/formatconverter", func(schema *jsonschema.Schema) {
				setMinLength(schema, "input", 1)
				schema.Properties["input"].Default = json.RawMessage(`"{\"replace\":\"me\"}"`)
			})
		},
		func() error {
			return registerTool[htpasswdInput](server, client, apiBaseURL, "generate-htpasswd", "Generate an htpasswd entry using the selected hash algorithm.", http.MethodPost, "/api/htpasswd", func(schema *jsonschema.Schema) {
				setEnum(schema, "algorithm", "apr1", "1", "5", "6")
				setMinLength(schema, "username", 1)
				setMinLength(schema, "password", 1)
				schema.Properties["username"].Default = json.RawMessage(`"replace-me"`)
				schema.Properties["password"].Default = json.RawMessage(`"replace-me"`)
				schema.Properties["algorithm"].Default = json.RawMessage(`"apr1"`)
			})
		},
		func() error {
			return registerTool[jsonPrettyInput](server, client, apiBaseURL, "json-pretty", "Format or minify JSON.", http.MethodPost, "/api/jsonpretty", func(schema *jsonschema.Schema) {
				setEnum(schema, "operation", "pretty", "unpretty")
				setMinLength(schema, "input", 1)
				setMinimum(schema, "spaces", 0)
				setMaximum(schema, "spaces", 16)
				schema.Properties["input"].Default = json.RawMessage(`"{\"key\":\"value\"}"`)
				schema.Properties["operation"].Default = json.RawMessage(`"pretty"`)
				schema.Properties["spaces"].Default = json.RawMessage("2")
			})
		},
		func() error {
			return registerTool[passwordInput](server, client, apiBaseURL, "generate-password", "Generate a random password with selected character sets.", http.MethodPost, "/api/passgen", func(schema *jsonschema.Schema) {
				setMinimum(schema, "length", 1)
				setMaximum(schema, "length", 1024)
				schema.Properties["length"].Default = json.RawMessage("16")
				schema.Properties["uppercase"].Default = json.RawMessage("true")
				schema.Properties["lowercase"].Default = json.RawMessage("true")
				schema.Properties["numbers"].Default = json.RawMessage("true")
				schema.Properties["symbols"].Default = json.RawMessage("true")
				schema.Not = &jsonschema.Schema{
					Properties: map[string]*jsonschema.Schema{
						"uppercase": {Const: jsonschema.Ptr[any](false)},
						"lowercase": {Const: jsonschema.Ptr[any](false)},
						"numbers":   {Const: jsonschema.Ptr[any](false)},
						"symbols":   {Const: jsonschema.Ptr[any](false)},
					},
					Required: []string{"uppercase", "lowercase", "numbers", "symbols"},
				}
			})
		},
		func() error {
			return registerTool[sshKeyInput](server, client, apiBaseURL, "generate-ssh-key", "Generate an Ed25519, ECDSA, or RSA SSH key pair.", http.MethodPost, "/api/sshgen", func(schema *jsonschema.Schema) {
				setEnum(schema, "algorithm", "ed25519", "ecdsa", "rsa")
				schema.Properties["algorithm"].Default = json.RawMessage(`"ed25519"`)
				schema.Properties["ecdsaBits"].Default = json.RawMessage(`"521"`)
				schema.Properties["rsaBits"].Default = json.RawMessage(`"4096"`)
				schema.Properties["email"].Default = json.RawMessage(`"replace-me@example.com"`)
				schema.AllOf = []*jsonschema.Schema{
					{
						If: &jsonschema.Schema{
							Properties: map[string]*jsonschema.Schema{
								"algorithm": {Const: jsonschema.Ptr[any]("ecdsa")},
							},
						},
						Then: &jsonschema.Schema{
							Required: []string{"ecdsaBits"},
							Properties: map[string]*jsonschema.Schema{
								"ecdsaBits": {Enum: []any{"256", "384", "521"}},
							},
						},
					},
					{
						If: &jsonschema.Schema{
							Properties: map[string]*jsonschema.Schema{
								"algorithm": {Const: jsonschema.Ptr[any]("rsa")},
							},
						},
						Then: &jsonschema.Schema{Required: []string{"rsaBits"}},
					},
				}
			})
		},
		func() error {
			return registerTool[timeConverterInput](server, client, apiBaseURL, "time-converter", "Convert between epoch timestamps and human-readable time in an IANA timezone.", http.MethodPost, "/api/timeconverter", func(schema *jsonschema.Schema) {
				setEnum(schema, "precision", "seconds", "millis", "micros")
				setMinLength(schema, "timezone", 2)
				setRange(schema, "millisecond", 0, 999)
				setRange(schema, "microsecond", 0, 999)
				schema.Properties["epochToHumanMode"].Default = json.RawMessage("true")
				schema.Properties["humanToEpochMode"].Default = json.RawMessage("false")
				schema.Properties["epochTime"].Default = json.RawMessage("0")
				schema.Properties["precision"].Default = json.RawMessage(`"seconds"`)
				schema.Properties["year"].Default = json.RawMessage("2025")
				schema.Properties["month"].Default = json.RawMessage("1")
				schema.Properties["day"].Default = json.RawMessage("1")
				schema.Properties["hour"].Default = json.RawMessage("0")
				schema.Properties["minute"].Default = json.RawMessage("0")
				schema.Properties["second"].Default = json.RawMessage("0")
				schema.Properties["millisecond"].Default = json.RawMessage("0")
				schema.Properties["microsecond"].Default = json.RawMessage("0")
				schema.Properties["timezone"].Default = json.RawMessage(`"UTC"`)
				schema.Not = &jsonschema.Schema{
					Properties: map[string]*jsonschema.Schema{
						"epochToHumanMode": {Const: jsonschema.Ptr[any](false)},
						"humanToEpochMode": {Const: jsonschema.Ptr[any](false)},
					},
					Required: []string{"epochToHumanMode", "humanToEpochMode"},
				}
				schema.AllOf = []*jsonschema.Schema{
					{
						If: &jsonschema.Schema{
							Properties: map[string]*jsonschema.Schema{
								"humanToEpochMode": {Const: jsonschema.Ptr[any](true)},
								"epochToHumanMode": {Const: jsonschema.Ptr[any](false)},
							},
							Required: []string{"humanToEpochMode", "epochToHumanMode"},
						},
						Then: &jsonschema.Schema{
							Properties: map[string]*jsonschema.Schema{
								"month":  {Minimum: jsonschema.Ptr(1.0), Maximum: jsonschema.Ptr(12.0)},
								"day":    {Minimum: jsonschema.Ptr(1.0), Maximum: jsonschema.Ptr(31.0)},
								"hour":   {Minimum: jsonschema.Ptr(0.0), Maximum: jsonschema.Ptr(23.0)},
								"minute": {Minimum: jsonschema.Ptr(0.0), Maximum: jsonschema.Ptr(59.0)},
								"second": {Minimum: jsonschema.Ptr(0.0), Maximum: jsonschema.Ptr(59.0)},
							},
						},
					},
				}
			})
		},
		func() error {
			return registerTool[urlDecoderInput](server, client, apiBaseURL, "url-decoder", "Parse a URL and decode its components and query parameters.", http.MethodPost, "/api/urldecoder", func(schema *jsonschema.Schema) {
				setMinLength(schema, "url", 1)
				schema.Properties["url"].Default = json.RawMessage(`"https://example.com/foo+bar%21?service=ab%2Bc&product=de%24f&service=second"`)
			})
		},
		func() error {
			return registerTool[emptyInput](server, client, apiBaseURL, "get-version", "Get the running application version.", http.MethodGet, "/api/version", nil)
		},
	}
	for _, register := range registrations {
		if err := register(); err != nil {
			return err
		}
	}
	return nil
}

func registerTool[In any](
	server *mcp.Server,
	client *resty.Client,
	apiBaseURL string,
	name string,
	description string,
	method string,
	path string,
	configureSchema func(*jsonschema.Schema),
) error {
	inputSchema, err := jsonschema.For[In](nil)
	if err != nil {
		return fmt.Errorf("infer input schema for tool %q: %w", name, err)
	}
	metrics.InitializeMCPTool(name)
	if configureSchema != nil {
		configureSchema(inputSchema)
	}
	mcp.AddTool(server, &mcp.Tool{
		Name:        name,
		Description: description,
		InputSchema: inputSchema,
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input In) (*mcp.CallToolResult, any, error) {
		var output map[string]any
		result, apiStatus, err := proxyAPI(ctx, client, apiBaseURL, method, path, input, &output)
		if err != nil {
			logger.Logger.Error(
				"mcp-tool-failed",
				"request-id", requestid.FromContext(ctx),
				"tool-name", name,
				"api-path", path,
				"api-status", apiStatus,
				"failure", apiFailureKind(apiStatus),
			)
			if result == nil {
				result = &mcp.CallToolResult{}
			}
			if result.Meta == nil {
				result.Meta = mcp.Meta{
					"request-id": requestid.FromContext(ctx),
				}
			}
			result.SetError(err)
			return result, nil, nil
		}
		return result, output, nil
	})
	return nil
}

func setEnum(schema *jsonschema.Schema, property string, values ...any) {
	schema.Properties[property].Enum = values
}

func setMinLength(schema *jsonschema.Schema, property string, minimum int) {
	schema.Properties[property].MinLength = jsonschema.Ptr(minimum)
}

func setMinimum(schema *jsonschema.Schema, property string, minimum float64) {
	schema.Properties[property].Minimum = jsonschema.Ptr(minimum)
}

func setMaximum(schema *jsonschema.Schema, property string, maximum float64) {
	schema.Properties[property].Maximum = jsonschema.Ptr(maximum)
}

func setRange(schema *jsonschema.Schema, property string, minimum, maximum float64) {
	setMinimum(schema, property, minimum)
	setMaximum(schema, property, maximum)
}

func proxyAPI(
	ctx context.Context,
	client *resty.Client,
	apiBaseURL string,
	method string,
	path string,
	input any,
	output any,
) (*mcp.CallToolResult, int, error) {
	mcpRequestID := requestid.MCPFromContext(ctx)
	if mcpRequestID == "" {
		return nil, 0, fmt.Errorf("MCP request ID is missing from context")
	}

	httpRequest := client.R().
		SetContext(ctx).
		SetHeader("MCP-Request-ID", mcpRequestID).
		SetHeader("X-Request-ID", mcpRequestID).
		SetResult(output)
	if method != http.MethodGet {
		httpRequest.SetBody(input)
	}
	response, err := httpRequest.Execute(method, apiBaseURL+path)
	if err != nil {
		return &mcp.CallToolResult{
			Meta: mcp.Meta{"request-id": mcpRequestID},
		}, 0, fmt.Errorf("API request failed: %w", err)
	}

	if response.IsError() {
		err := fmt.Errorf("API returned %s: %s", response.Status(), strings.TrimSpace(response.String()))
		return &mcp.CallToolResult{
			Meta: mcp.Meta{"request-id": mcpRequestID},
		}, response.StatusCode(), err
	}

	return &mcp.CallToolResult{
		Meta: mcp.Meta{"request-id": mcpRequestID},
	}, response.StatusCode(), nil
}

func apiFailureKind(status int) string {
	if status == 0 {
		return "transport"
	}
	return "api-response"
}
