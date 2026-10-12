package apis

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/version"
)

// VersionOutput is the response structure for the /api/version endpoint.
type VersionOutput struct {
	Body struct {
		Version string `json:"version"`
	}
}

// RegisterVersion registers the /api/version endpoint with the given Huma API
func RegisterVersion(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "get-version",
		Summary:       "Get the application version (JSON)",
		Description:   "Returns the running Go Tooling Portal application version as JSON.",
		Method:        http.MethodGet,
		Path:          "/api/version",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"API Version"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "JSON response containing the application version.",
			},
			"500": {
				Description: "The application version could not be retrieved.",
			},
		},
	}, func(ctx context.Context, input *struct{}) (*VersionOutput, error) {
		if version.Version == "" {
			logger.Logger.Error("internal server error, unable to retrieve api version information")
			return nil, huma.Error500InternalServerError("internal server error, unable to retrieve api version information")
		}

		resp := &VersionOutput{}
		resp.Body.Version = version.Version
		return resp, nil
	})
}
