package apis

import (
	"context"
	"html/template"
	"net/http"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sfarosu/go-tooling-portal/internal/service"
)

// FormatConverterHTMXInput is the request structure for the /api/formatconverter/htmx endpoint.
type FormatConverterHTMXInput struct {
	// HtmxHeader stores the value of the HX-Request header, which is sent by HTMX
	// when making AJAX requests. This allows the handler to detect if the request
	// originated from an HTMX-enabled frontend and respond with HTML instead of JSON
	HtmxHeader bool `header:"HX-Request" doc:"Set to true when calling this HTMX endpoint."`

	// Body contains the main input fields for format conversion.
	Body struct {
		Input string `json:"input" example:"{\"name\":\"Alice\",\"age\":30}" doc:"JSON or YAML string to convert to the other format" minLength:"1"`
	}
}

// Resolve implements the huma.RequestResolver interface.
// It is called by Huma before the handler to allow custom extraction or transformation
// of request data. Here, we extract the HX-Request header (sent by HTMX)
// and store it in the HtmxHeader field for later use in the handler
// See: https://huma.rocks/features/request-resolvers/
func (m *FormatConverterHTMXInput) Resolve(ctx huma.Context) []error {
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

// FormatConverterHTMXOutput is the response structure for the /api/formatconverter/htmx endpoint
type FormatConverterHTMXOutput struct {
	ContentType string `header:"Content-Type"`
	Body        []byte `json:"-"` // Body is not serialized to JSON, but used to return HTML content
}

// HTMX-compatible HTML fragment rendered dynamically by the API and injected into the page via HTMX
var formatconverterResultTmpl = template.Must(template.New("formatconverter-result").Parse(`
<div class="d-flex justify-content-center">
  <div class="input-group" style="max-width: 1024px; width: 100%;">
    <textarea class="form-control custom-output" id="formatconverter-result-input"
      readonly aria-label="Converted result" rows="4">{{.FormatConversion}}</textarea>
    <button class="btn btn-graphite" type="button"
      onclick="copyToClipboard('formatconverter-result-input')"
      aria-label="Copy result to clipboard" tabindex="-1">
      <i class="bi bi-clipboard"></i>
    </button>
  </div>
</div>
`))

// RegisterFormatConverterHtmx registers the /api/formatconverter/htmx endpoint with the given Huma API
func RegisterFormatConverterHtmx(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "format-converter-htmx",
		Summary:       "Convert JSON and YAML (HTMX)",
		Description:   "Converts JSON input to YAML or YAML input to JSON. Requires the HX-Request: true header and returns an HTML fragment for HTMX.",
		Method:        http.MethodPost,
		Path:          "/api/formatconverter/htmx",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Format Converter"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "HTML fragment containing the conversion result or a user-visible error.",
				Content: map[string]*huma.MediaType{
					"text/html": {},
				},
			},
			"400": {
				Description: "Invalid input or the required HX-Request: true header is missing.",
			},
		},
	}, func(ctx context.Context, input *FormatConverterHTMXInput) (*FormatConverterHTMXOutput, error) {
		// Process format conversion
		formatConversion, err := service.FormatConvert(input.Body.Input)
		if err != nil {
			logError(ctx,
				"failed to perform format conversion",
				"input", input.Body.Input,
				"error", err,
			)
			if input.HtmxHeader {
				return &FormatConverterHTMXOutput{
					ContentType: "text/html; charset=utf-8",
					Body:        []byte(generateHTMXError(err)),
				}, nil
			}
			return nil, huma.Error400BadRequest("format conversion failed: " + err.Error())
		}

		// Render the HTML template and return as a html block if this is an HTMX request
		if input.HtmxHeader {
			body, err := renderHTMXTemplate(formatconverterResultTmpl, map[string]string{
				"FormatConversion": formatConversion,
			})
			if err != nil {
				logError(ctx, "failed to render format conversion result", "error", err)
				return nil, huma.Error500InternalServerError("failed to render format conversion result")
			}
			resp := &FormatConverterHTMXOutput{
				ContentType: "text/html; charset=utf-8",
				Body:        body,
			}
			return resp, nil
		}

		return nil, huma.Error400BadRequest("This endpoint only supports HTMX requests with HX-Request header set to true; use /api/formatconverter for JSON responses")
	})
}
