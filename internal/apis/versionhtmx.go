package apis

import (
	"context"
	"html/template"
	"net/http"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/version"
)

// VersionHTMXInput is the request structure for the /api/version/htmx endpoint.
type VersionHTMXInput struct {
	// HtmxHeader stores the value of the HX-Request header, which is sent by HTMX
	// when making AJAX requests. This allows the handler to detect if the request
	// originated from an HTMX-enabled frontend and respond with HTML instead of JSON
	HtmxHeader bool `header:"HX-Request" doc:"Set to true when calling this HTMX endpoint."`
}

// Resolve implements the huma.RequestResolver interface.
// It is called by Huma before the handler to allow custom extraction or transformation
// of request data. Here, we extract the HX-Request header (sent by HTMX)
// and store it in the HtmxHeader field for later use in the handler
// See: https://huma.rocks/features/request-resolvers/
func (m *VersionHTMXInput) Resolve(ctx huma.Context) []error {
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

// VersionHTMXOutput is the response structure for the /api/version/htmx endpoint.
type VersionHTMXOutput struct {
	ContentType string `header:"Content-Type"`
	Body        []byte `json:"-"` // Body is not serialized to JSON, but used to return HTML content
}

// HTMX-compatible HTML fragment rendered dynamically by the API and injected into the page via HTMX
var versionResultTmpl = template.Must(template.New("version-result").Parse(`
<span class="badge">v{{.Version}}</span>
`))

// RegisterVersionHtmx registers the /api/version/htmx endpoint with the given Huma API
func RegisterVersionHtmx(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "get-version-htmx",
		Summary:       "Get the application version (HTMX)",
		Description:   "Returns the running Go Tooling Portal application version in an HTML fragment. Requires the HX-Request: true header.",
		Method:        http.MethodGet,
		Path:          "/api/version/htmx",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"API Version"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "HTML fragment containing the application version.",
				Content: map[string]*huma.MediaType{
					"text/html": {},
				},
			},
			"500": {
				Description: "The application version could not be retrieved.",
			},
		},
	}, func(ctx context.Context, input *VersionHTMXInput) (*VersionHTMXOutput, error) {
		// Validation: ensure we have version info
		if version.Version == "" {
			logger.Logger.Error("internal server error, unable to retrieve api version information")

			if input.HtmxHeader {
				// If it's a htmx request, return an HTML fragment with the error content (return code will be 200 as we want to display it in the UI)
				return &VersionHTMXOutput{
					ContentType: "text/html; charset=utf-8",
					Body:        []byte("internal server error, unable to retrieve api version information"),
				}, nil
			}
			return nil, huma.Error500InternalServerError("internal server error, unable to retrieve api version information")
		}

		// Render the HTML template and return as a html block if this is an HTMX request
		if input.HtmxHeader {
			body, err := renderHTMXTemplate(versionResultTmpl, map[string]string{
				"Version": version.Version,
			})
			if err != nil {
				if logger.Logger != nil {
					logger.Logger.Error("failed to render version result", "error", err)
				}
				return nil, huma.Error500InternalServerError("failed to render version result")
			}
			resp := &VersionHTMXOutput{
				ContentType: "text/html; charset=utf-8",
				Body:        body,
			}
			return resp, nil
		}

		return nil, huma.Error400BadRequest("This endpoint only supports HTMX requests with HX-Request header set to true; use /api/version for JSON responses")
	})
}
