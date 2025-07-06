package apis

import (
	"bytes"
	"context"
	"net/http"
	"strconv"
	"text/template"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/version"
)

// VersionHTMXInput is the request structure for the /api/version/htmx endpoint.
type VersionHTMXInput struct {
	// HtmxHeader stores the value of the HX-Request header, which is sent by HTMX
	// when making AJAX requests. This allows the handler to detect if the request
	// originated from an HTMX-enabled frontend and respond with HTML instead of JSON
	HtmxHeader bool `header:"HX-Request"`
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
		Summary:       "get version - htmx",
		Description:   "Returns a suitable for HTMX injection HTML object containg the api version.",
		Method:        http.MethodGet,
		Path:          "/api/version/htmx",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"ApiVersion"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Successfully retrieved the API version, responding with HTML content for HTMX requests.",
				Content: map[string]*huma.MediaType{
					"text/html": {},
				},
			},
			"500": {
				Description: "Internal Server Error, unable to retrieve api version information.",
			},
		},
	}, func(ctx context.Context, input *VersionHTMXInput) (*VersionHTMXOutput, error) {
		if version.Version == "" {
			logger.Logger.Error("Internal Server Error, unable to retrieve api version information")
			return nil, huma.Error500InternalServerError("Internal Server Error, unable to retrieve api version information")
		}

		// Render the HTML template and return as a html block if this is an HTMX request
		if input.HtmxHeader {
			var buf bytes.Buffer
			versionResultTmpl.Execute(&buf, map[string]string{
				"Version": version.Version,
			})
			resp := &VersionHTMXOutput{
				ContentType: "text/html; charset=utf-8",
				Body:        buf.Bytes(),
			}
			return resp, nil
		}

		return nil, huma.Error400BadRequest("This endpoint only supports HTMX requests with HX-Request header set to true; use /api/version for JSON responses")
	})
}
