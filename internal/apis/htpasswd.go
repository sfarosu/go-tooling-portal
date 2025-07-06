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
		Username  string `json:"username" example:"alice" doc:"Username for the htpasswd entry" minLength:"1"`
		Password  string `json:"password" example:"S3cureP@ssw0rd" doc:"Password for the htpasswd entry" minLength:"1"`
		Algorithm string `json:"algorithm" example:"apr1" doc:"Hashing algorithm to use (apr1, 1, 5, or 6)" enum:"apr1,1,5,6"`
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
		Summary:       "generate htpasswd - json",
		Description:   "Returns a JSON object containg an htpasswd entry for the given username and password using the specified algorithm.",
		Method:        http.MethodPost,
		Path:          "/api/htpasswd",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Htpassword"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Successfully generated an htpasswd entry using the specified algorithm for the given username and password.",
			},
			"400": {
				Description: "Bad Request, missing fields or unsupported hashing algorithm (must be one of: bcrypt, md5, sha1, or crypt).",
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
			return nil, huma.Error400BadRequest("htpasswd generation failed:  " + err.Error())
		}

		resp := &HtpasswdOutput{}
		resp.Body.Htpasswd = generatedHtPassword
		return resp, nil
	})
}
