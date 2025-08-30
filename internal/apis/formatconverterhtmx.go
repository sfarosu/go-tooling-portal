package apis

import (
	"bytes"
	"context"
	"net/http"
	"strconv"
	"text/template"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/service"
)

// FormatConverterHTMXInput is the request structure for the /api/formatconverter/htmx endpoint.
type FormatConverterHTMXInput struct {
	// HtmxHeader stores the value of the HX-Request header, which is sent by HTMX
	// when making AJAX requests. This allows the handler to detect if the request
	// originated from an HTMX-enabled frontend and respond with HTML instead of JSON
	HtmxHeader bool `header:"HX-Request"`

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
  <div class="input-group" style="max-width: 680px; width: 100%;">
    <textarea class="form-control custom-output" id="formatconverter-result-input"
      readonly aria-label="Converted result" rows="4">{{.FormatConversionResult}}</textarea>
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
		Summary:       "convert format between json and yaml - htmx",
		Description:   "Convert data between JSON and YAML formats.",
		Method:        http.MethodPost,
		Path:          "/api/formatconverter/htmx",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"FormatConverter"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Successfull convertion, responding with HTML content for HTMX requests.",
				Content: map[string]*huma.MediaType{
					"text/html": {},
				},
			},
			"400": {
				Description: "Bad Request, invalid input (including missing HX-Request header).",
			},
		},
	}, func(ctx context.Context, input *FormatConverterHTMXInput) (*FormatConverterHTMXOutput, error) {
		formatConversionResult, err := service.FormatConvert(input.Body.Input)
		if err != nil {
			logger.Logger.Error(
				"failed to process format conversion",
				"input", input.Body.Input,
				"error", err,
			)
			if input.HtmxHeader {
				// If it's a htmx request, return an HTML fragment with the error content (return code will be 200 as we want to display it in the UI)
				alert := `
<div class="d-flex justify-content-center">
  <div class="input-group" style="max-width: 680px; width: 100%;">
    <textarea class="form-control custom-output is-invalid"
      readonly aria-label="Converted result">` + err.Error() + `</textarea>
  </div>
</div>
`
				return &FormatConverterHTMXOutput{
					ContentType: "text/html; charset=utf-8",
					Body:        []byte(alert),
				}, nil
			}
			return nil, huma.Error400BadRequest("format conversion failed: " + err.Error())
		}

		// Render the HTML template and return as a html block if this is an HTMX request
		if input.HtmxHeader {
			var buf bytes.Buffer
			formatconverterResultTmpl.Execute(&buf, map[string]string{
				"FormatConversionResult": formatConversionResult,
			})
			resp := &FormatConverterHTMXOutput{
				ContentType: "text/html; charset=utf-8",
				Body:        buf.Bytes(),
			}
			return resp, nil
		}

		return nil, huma.Error400BadRequest("This endpoint only supports HTMX requests with HX-Request header set to true; use /api/formatconverter for JSON responses")
	})
}
