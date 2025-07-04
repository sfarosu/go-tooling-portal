package apis

import (
	"bytes"
	"context"
	"net/http"
	"text/template"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/service"
)

// HtpasswdInput is the request structure for the /api/htpasswd endpoint.
type HtpasswdHTMXInput struct {
	// HtmxHeader stores the value of the HX-Request header, which is sent by HTMX
	// when making AJAX requests. This allows the handler to detect if the request
	// originated from an HTMX-enabled frontend and respond with HTML instead of JSON.
	HtmxHeader string `header:"HX-Request"`

	// Body contains the main input fields for the htpasswd generation.
	Body struct {
		Username  string `json:"username" example:"alice" doc:"Username for the htpasswd entry" minLength:"1"`
		Password  string `json:"password" example:"S3cureP@ssw0rd" doc:"Password for the htpasswd entry" minLength:"1"`
		Algorithm string `json:"algorithm" example:"apr1" doc:"Hashing algorithm to use (apr1, 1, 5, or 6)" enum:"apr1,1,5,6"`
	}
}

// Resolve implements the huma.RequestResolver interface.
// It is called by Huma before the handler to allow custom extraction or transformation
// of request data. Here, we extract the HX-Request header (sent by HTMX)
// and store it in the HtmxHeader field for later use in the handler.
// See: https://huma.rocks/features/request-resolvers/
func (m *HtpasswdHTMXInput) Resolve(ctx huma.Context) []error {
	m.HtmxHeader = ctx.Header("HX-Request")

	return nil
}

// HtpasswdOutput is the response structure for the /api/htpasswd endpoint.
type HtpasswdHTMXOutput struct {
	ContentType string `header:"Content-Type"`
	Body        []byte `json:"-"` // Body is not serialized to JSON, but used to return HTML content
}

var htpasswdResultTmpl = template.Must(template.New("htpasswd-result").Parse(`
<div class="d-flex justify-content-center">
  <div class="input-group" style="max-width: 680px; width: 100%;">
    <input type="text" class="form-control custom-output" id="htpasswd-result-input"
      value="{{.Htpasswd}}" readonly aria-label="Generated htpasswd">
    <button class="btn btn-graphite" type="button"
      onclick="copyToClipboard('htpasswd-result-input')"
      aria-label="Copy htpasswd to clipboard" tabindex="-1">
      <i class="bi bi-clipboard"></i>
    </button>
  </div>
</div>
`))

// RegisterHtpasswdHtmx registers the /api/htpasswd/htmx endpoint with the given Huma API
func RegisterHtpasswdHtmx(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "generate-htpasswd-htmx",
		Summary:       "Generate htpasswd; returns HTML for HTMX",
		Method:        http.MethodPost,
		Path:          "/api/htpasswd/htmx",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Htpassword"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Generates an htpasswd entry for the given username and password using the specified algorithm and returns HTML for HTMX requests",
				Content: map[string]*huma.MediaType{
					"text/html": {},
				},
			},
			"400": {
				Description: "Bad Request - Invalid input (including missing HX-Request header) or unsupported algorithm",
			},
		},
	}, func(ctx context.Context, input *HtpasswdHTMXInput) (*HtpasswdHTMXOutput, error) {
		generatedHtPassword, err := service.GenerateHtpasswd(input.Body.Username, input.Body.Password, input.Body.Algorithm)
		if err != nil {
			logger.Logger.Error(
				"failed to generate htpasswd",
				"username", input.Body.Username,
				"algorithm", input.Body.Algorithm,
				"error", err,
			)
			return nil, huma.Error400BadRequest("failed to generate htpasswd: " + err.Error())
		}

		// Render the HTML template and return as a string if this is an HTMX request
		if input.HtmxHeader == "true" {
			var buf bytes.Buffer
			htpasswdResultTmpl.Execute(&buf, map[string]string{
				"Username":  input.Body.Username,
				"Password":  input.Body.Password,
				"Algorithm": input.Body.Algorithm,
				"Htpasswd":  generatedHtPassword,
			})
			resp := &HtpasswdHTMXOutput{
				ContentType: "text/html; charset=utf-8",
				Body:        buf.Bytes(),
			}
			return resp, nil
		}

		return nil, huma.Error400BadRequest("This endpoint only supports HTMX requests with HX-Request header set to true; use /api/htpasswd for JSON responses")
	})
}
