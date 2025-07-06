package apis

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/version"
)

type VersionOutput struct {
	Body struct {
		Version string `json:"version"`
	}
}

// RegisterVersion registers the /api/version endpoint with the given Huma API
func RegisterVersion(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "get-version",
		Summary:       "get version - json",
		Description:   "Returns a JSON object containg the api version.",
		Method:        http.MethodGet,
		Path:          "/api/version",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"ApiVersion"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Successfully retrieved the api version.",
			},
			"500": {
				Description: "Internal Server Error, Unable to retrieve api version information.",
			},
		},
	}, func(ctx context.Context, input *struct{}) (*VersionOutput, error) {
		if version.Version == "" {
			logger.Logger.Error("Internal Server Error - Unable to retrieve api version information")
			return nil, huma.Error500InternalServerError("Internal Server Error - Unable to retrieve api version information")
		}

		resp := &VersionOutput{}
		resp.Body.Version = version.Version
		return resp, nil
	})
}
