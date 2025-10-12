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

// Base64ConverterHTMXInput is the request structure for the /api/base64converter/htmx endpoint.
type Base64ConverterHTMXInput struct {
	// HtmxHeader stores the value of the HX-Request header, which is sent by HTMX
	// when making AJAX requests. This allows the handler to detect if the request
	// originated from an HTMX-enabled frontend and respond with HTML instead of JSON
	HtmxHeader bool `header:"HX-Request"`

	Body struct {
		Input     string `json:"input" example:"hello" doc:"String to encode or decode" minLength:"1"`
		Operation string `json:"operation" example:"encode" doc:"Operation to perform: encode or decode" enum:"encode,decode"`
		Format    string `json:"format" example:"standard" doc:"Base64 format: standard or url-compatible" enum:"standard,url-compatible"`
	}
}

// Resolve implements the huma.RequestResolver interface.
// It is called by Huma before the handler to allow custom extraction or transformation
// of request data. Here, we extract the HX-Request header (sent by HTMX)
// and store it in the HtmxHeader field for later use in the handler
// See: https://huma.rocks/features/request-resolvers/
func (m *Base64ConverterHTMXInput) Resolve(ctx huma.Context) []error {
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

// Base64ConverterHTMXOutput is the response structure for the /api/base64converter/htmx endpoint
type Base64ConverterHTMXOutput struct {
	ContentType string `header:"Content-Type"`
	Body        []byte `json:"-"` // Body is not serialized to JSON, but used to return HTML content
}

// HTMX-compatible HTML fragment rendered dynamically by the API and injected into the page via HTMX
var base64converterResultTmpl = template.Must(template.New("base64converter-result").Parse(`
<div class="d-flex justify-content-center">
  <div class="input-group" style="max-width: 680px; width: 100%;">
    <textarea class="form-control custom-output" id="base64converter-result-input"
      readonly aria-label="Converted result" rows="4">{{.Base64ConversionResult}}</textarea>
    <button class="btn btn-graphite" type="button"
      onclick="copyToClipboard('base64converter-result-input')"
      aria-label="Copy result to clipboard" tabindex="-1">
      <i class="bi bi-clipboard"></i>
    </button>
  </div>
</div>
`))

// RegisterBase64ConverterHtmx registers the /api/base64converter/htmx endpoint with the given Huma API
func RegisterBase64ConverterHtmx(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "base64-converter-htmx",
		Summary:       "base64 encode/decode - htmx",
		Description:   "Returns a suitable for HTMX injection HTML object containg a base64-encoded or decoded result depending on the selected operation.",
		Method:        http.MethodPost,
		Path:          "/api/base64converter/htmx",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Base64Converter"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Successfully converted the input, responding with HTML content for HTMX requests.",
				Content: map[string]*huma.MediaType{
					"text/html": {},
				},
			},
			"400": {
				Description: "Bad Request, invalid input (including missing HX-Request header) or unsupported operation.",
			},
		},
	}, func(ctx context.Context, input *Base64ConverterHTMXInput) (*Base64ConverterHTMXOutput, error) {
		base64ConversionResult, err := service.Base64Convert(input.Body.Input, input.Body.Operation, input.Body.Format)
		if err != nil {
			logger.Logger.Error(
				"failed to process base64 conversion",
				"input", input.Body.Input,
				"operation", input.Body.Operation,
				"format", input.Body.Format,
				"error", err,
			)
			if input.HtmxHeader {
				// If it's a htmx request, return an HTML fragment with the error content (return code will be 200 as we want to display it in the UI)
				alert := `
<div class="d-flex justify-content-center">
  <div style="max-width: 680px; width: 100%;">
    <textarea class="form-control custom-output is-invalid"
      readonly aria-label="error">` + err.Error() + `</textarea>
  </div>
</div>
`
				return &Base64ConverterHTMXOutput{
					ContentType: "text/html; charset=utf-8",
					Body:        []byte(alert),
				}, nil
			}
			return nil, huma.Error400BadRequest("base64 conversion failed: " + err.Error())
		}

		// Render the HTML template and return as a html block if this is an HTMX request
		if input.HtmxHeader {
			var buf bytes.Buffer
			base64converterResultTmpl.Execute(&buf, map[string]string{
				"Base64ConversionResult": base64ConversionResult,
			})
			resp := &Base64ConverterHTMXOutput{
				ContentType: "text/html; charset=utf-8",
				Body:        buf.Bytes(),
			}
			return resp, nil
		}

		return nil, huma.Error400BadRequest("This endpoint only supports HTMX requests with HX-Request header set to true; use /api/base64converter for JSON responses")
	})
}
