package apis

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/service"
)

type HtpasswdInput struct {
	Body struct {
		Username  string `json:"username" example:"alice" doc:"Username for the htpasswd entry"`
		Password  string `json:"password" example:"S3cureP@ssw0rd" doc:"Password for the htpasswd entry"`
		Algorithm string `json:"algorithm" example:"apr1" doc:"Hashing algorithm to use (apr1, 1, 5, or 6)"`
	}
}

type HtpasswdOutput struct {
	Body struct {
		Htpasswd string `json:"htpasswd"`
	}
}

// RegisterHtpasswd registers the /api/htpasswd endpoint with the given Huma API
func RegisterHtpasswd(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "generate-htpasswd",
		Summary:       "Generate Htpasswd",
		Method:        http.MethodPost,
		Path:          "/api/htpasswd",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Htpassword"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Generates an htpasswd entry for the given username and password using the specified algorithm",
			},
			"400": {
				Description: "Bad Request - Invalid input or unsupported algorithm",
			},
		},
	}, func(ctx context.Context, input *HtpasswdInput) (*HtpasswdOutput, error) {
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

		resp := &HtpasswdOutput{}
		resp.Body.Htpasswd = generatedHtPassword
		return resp, nil
	})
}
