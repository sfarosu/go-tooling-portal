package apis

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/service"
)

// JsonPrettyHTMXInput is the request structure for the /api/jsonpretty/htmx endpoint.
type JsonPrettyHTMXInput struct {
	// HtmxHeader stores the value of the HX-Request header, which is sent by HTMX
	// when making AJAX requests. This allows the handler to detect if the request
	// originated from an HTMX-enabled frontend and respond with HTML instead of JSON
	HtmxHeader bool `header:"HX-Request" doc:"Set to true when calling this HTMX endpoint."`

	// Body contains the main input fields for the JSON prettifier.
	// The fields are strings because HTML form submissions send string values.
	Body struct {
		Input     string `json:"input" example:"{\"name\":\"Alice\",\"age\":30}" doc:"JSON input to prettify" minLength:"1"`
		Operation string `json:"operation" example:"pretty" doc:"Operation: pretty or unpretty" enum:"pretty,unpretty"`
		Spaces    string `json:"spaces" example:"2" doc:"Number of spaces for indentation (0-16)" minLength:"1"`
	}
}

// Resolve implements the huma.RequestResolver interface.
// It is called by Huma before the handler to allow custom extraction or transformation
// of request data. Here, we extract the HX-Request header (sent by HTMX)
// and store it in the HtmxHeader field for later use in the handler
// See: https://huma.rocks/features/request-resolvers/
func (m *JsonPrettyHTMXInput) Resolve(ctx huma.Context) []error {
	val := ctx.Header("HX-Request")
	parsed, err := strconv.ParseBool(val)
	if err != nil {
		// If parsing fails, assume it's not an HTMX request
		m.HtmxHeader = false
		return nil
	}
	m.HtmxHeader = parsed
	return nil
}

// JsonPrettyHTMXOutput is the response structure for the /api/jsonpretty/htmx endpoint.
type JsonPrettyHTMXOutput struct {
	ContentType string `header:"Content-Type"`
	Body        []byte `json:"-"` // Body is not serialized to JSON, but used to return HTML content
}

// HTMX-compatible HTML fragment rendered dynamically by the API and injected into the page via HTMX
var jsonPrettyResultTmpl = template.Must(template.New("jsonpretty-result").Parse(`
<div class="d-flex justify-content-center">
  <div class="input-group" style="width: 100%;">
    <textarea class="form-control custom-output" id="jsonpretty-result-input"
      readonly aria-label="Result JSON" rows="16">{{.ResultJSON}}</textarea>
    <button class="btn btn-graphite" type="button"
      onclick="copyToClipboard('jsonpretty-result-input')"
      aria-label="Copy result to clipboard" tabindex="-1">
      <i class="bi bi-clipboard"></i>
    </button>
  </div>
</div>
`))

// RegisterJsonPrettyHtmx registers the /api/jsonpretty/htmx endpoint with the given Huma API
func RegisterJsonPrettyHtmx(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "json-pretty-htmx",
		Summary:       "Format or minify JSON (HTMX)",
		Description:   "Formats JSON with the requested indentation or minifies it. Requires the HX-Request: true header and returns an HTML fragment for HTMX.",
		Method:        http.MethodPost,
		Path:          "/api/jsonpretty/htmx",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"JSON Prettifier"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "HTML fragment containing the formatted JSON or a user-visible error.",
				Content: map[string]*huma.MediaType{
					"text/html": {},
				},
			},
			"400": {
				Description: "Invalid input or the required HX-Request: true header is missing.",
			},
		},
	}, func(ctx context.Context, input *JsonPrettyHTMXInput) (*JsonPrettyHTMXOutput, error) {
		// Handle operation: pretty (default) or unpretty (minify)
		op := input.Body.Operation
		if op == "unpretty" {
			// Unpretty (minify) JSON
			result, err := service.UnprettyJSON(input.Body.Input)
			if err != nil {
				logError(ctx, "failed to unpretty json", "error", err)
				if input.HtmxHeader {
					return &JsonPrettyHTMXOutput{
						ContentType: "text/html; charset=utf-8",
						Body:        []byte(generateHTMXError(err)),
					}, nil
				}
				return nil, huma.Error400BadRequest("failed to unpretty JSON: " + err.Error())
			}
			if input.HtmxHeader {
				body, err := renderHTMXTemplate(jsonPrettyResultTmpl, map[string]string{"ResultJSON": result})
				if err != nil {
					if logger.Logger != nil {
						logError(ctx, "failed to render JSON prettifier result", "error", err)
					}
					return nil, huma.Error500InternalServerError("failed to render JSON prettifier result")
				}
				return &JsonPrettyHTMXOutput{ContentType: "text/html; charset=utf-8", Body: body}, nil
			}
			return nil, huma.Error400BadRequest("This endpoint only supports HTMX requests; use /api/jsonpretty for JSON responses")
		}

		// Default to pretty and validate spaces
		spaces, err := strconv.Atoi(input.Body.Spaces)
		if err != nil || spaces < 0 {
			message := fmt.Errorf("spaces must be a valid non-negative integer")
			if input.HtmxHeader {
				return &JsonPrettyHTMXOutput{
					ContentType: "text/html; charset=utf-8",
					Body:        []byte(generateHTMXError(message)),
				}, nil
			}
			return nil, huma.Error400BadRequest(message.Error())
		}

		// Prettify JSON
		result, err := service.PrettyJSONWithSpaces(input.Body.Input, spaces)
		if err != nil {
			logError(ctx,
				"failed to prettify json",
				"spaces", input.Body.Spaces,
				"error", err)
			if input.HtmxHeader {
				return &JsonPrettyHTMXOutput{
					ContentType: "text/html; charset=utf-8",
					Body:        []byte(generateHTMXError(err)),
				}, nil
			}
			return nil, huma.Error400BadRequest("failed to prettify JSON: " + err.Error())
		}

		// Render the HTML template and return as a html block if this is an HTMX request
		if input.HtmxHeader {
			body, err := renderHTMXTemplate(jsonPrettyResultTmpl, map[string]string{
				"ResultJSON": result,
			})
			if err != nil {
				if logger.Logger != nil {
					logError(ctx, "failed to render JSON prettifier result", "error", err)
				}
				return nil, huma.Error500InternalServerError("failed to render JSON prettifier result")
			}
			resp := &JsonPrettyHTMXOutput{
				ContentType: "text/html; charset=utf-8",
				Body:        body,
			}
			return resp, nil
		}

		return nil, huma.Error400BadRequest("This endpoint only supports HTMX requests; use /api/jsonpretty for JSON responses")
	})
}
