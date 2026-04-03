package apis

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/service"
)

// PassgenHTMXInput is the request structure for the /api/passgen/htmx endpoint.
type PassgenHTMXInput struct {
	// HtmxHeader stores the value of the HX-Request header, which is sent by HTMX
	// when making AJAX requests. This allows the handler to detect if the request
	// originated from an HTMX-enabled frontend and respond with HTML instead of JSON
	HtmxHeader bool `header:"HX-Request"`

	// Body contains the main input fields for the password generation.
	// The fields are strings because HTML form submissions send string values.
	Body struct {
		Length    string `json:"length" example:"16" doc:"Password length" minLength:"1"`
		Uppercase string `json:"uppercase" example:"true" doc:"Include uppercase letters"`
		Lowercase string `json:"lowercase" example:"true" doc:"Include lowercase letters"`
		Numbers   string `json:"numbers" example:"true" doc:"Include digits"`
		Symbols   string `json:"symbols" example:"true" doc:"Include symbols"`
	}
}

// Resolve implements the huma.RequestResolver interface.
// It is called by Huma before the handler to allow custom extraction or transformation
// of request data. Here, we extract the HX-Request header (sent by HTMX)
// and store it in the HtmxHeader field for later use in the handler
// See: https://huma.rocks/features/request-resolvers/
func (m *PassgenHTMXInput) Resolve(ctx huma.Context) []error {
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

// PassgenHTMXOutput is the response structure for the /api/passgen/htmx endpoint.
type PassgenHTMXOutput struct {
	ContentType string `header:"Content-Type"`
	Body        []byte `json:"-"` // Body is not serialized to JSON, but used to return HTML content
}

// HTMX-compatible HTML fragment rendered dynamically by the API and injected into the page via HTMX
var passgenResultTmpl = template.Must(template.New("passgen-result").Parse(`
<div class="d-flex justify-content-center">
  <div class="input-group" style="max-width: 680px; width: 100%;">
    <textarea class="form-control custom-output" id="passgen-result-input"
      readonly aria-label="Generated password" rows="1">{{.Password}}</textarea>
    <button class="btn btn-graphite" type="button"
      onclick="copyToClipboard('passgen-result-input')"
      aria-label="Copy password to clipboard" tabindex="-1">
      <i class="bi bi-clipboard"></i>
    </button>
  </div>
</div>
`))

// RegisterPassgenHtmx registers the /api/passgen/htmx endpoint with the given Huma API
func RegisterPassgenHtmx(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "generate-password-htmx",
		Summary:       "generate password - htmx",
		Description:   "Returns a suitable for HTMX injection HTML object containg the generated password entry.",
		Method:        http.MethodPost,
		Path:          "/api/passgen/htmx",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Passgen"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Successfully generated a password, responding with HTML content for HTMX requests.",
				Content: map[string]*huma.MediaType{
					"text/html": {},
				},
			},
			"400": {
				Description: "Bad Request, invalid input (including missing HX-Request header).",
			},
		},
	}, func(ctx context.Context, input *PassgenHTMXInput) (*PassgenHTMXOutput, error) {
		// Validate input
		length, err := strconv.Atoi(input.Body.Length)
		if err != nil || length <= 0 {
			switch input.HtmxHeader {
			case true:
				generateHTMXError(fmt.Errorf("length must be a valid positive integer"))
			case false:
				return nil, huma.Error400BadRequest("length must be a valid positive integer")
			}
		}

		// Checkboxes: "true" = checked, "false" = unchecked
		// We expect string values because HTMX form submissions send string values but we need to transform them to bool
		uppercase := input.Body.Uppercase == "true"
		lowercase := input.Body.Lowercase == "true"
		numbers := input.Body.Numbers == "true"
		symbols := input.Body.Symbols == "true"

		if !uppercase && !lowercase && !numbers && !symbols {
			switch input.HtmxHeader {
			case true:
				generateHTMXError(fmt.Errorf("at least one character class must be selected"))
			case false:
				return nil, huma.Error400BadRequest("at least one character class must be selected")
			}
		}

		// Generate password
		generatedPassword, err := service.GeneratePassword(length, uppercase, lowercase, numbers, symbols)
		if err != nil {
			logger.Logger.Error(
				"failed to generate password",
				"length", length,
				"uppercase", uppercase,
				"lowercase", lowercase,
				"numbers", numbers,
				"symbols", symbols,
				"error", err)
			if input.HtmxHeader {
				// If it's a htmx request, return an HTML fragment with the error content (return code will be 200 as we want to display it in the UI)
				return &PassgenHTMXOutput{
					ContentType: "text/html; charset=utf-8",
					Body:        []byte(generateHTMXError(err)),
				}, nil
			}
			return nil, huma.Error400BadRequest("failed to generate password: " + err.Error())
		}

		// Render the HTML template and return as a html block if this is an HTMX request
		if input.HtmxHeader {
			var buf bytes.Buffer
			passgenResultTmpl.Execute(&buf, map[string]string{
				"Password": generatedPassword,
			})
			resp := &PassgenHTMXOutput{
				ContentType: "text/html; charset=utf-8",
				Body:        buf.Bytes(),
			}
			return resp, nil
		}

		return nil, huma.Error400BadRequest("This endpoint only supports HTMX requests with HX-Request header set to true; use /api/passgen for JSON responses")
	})
}
