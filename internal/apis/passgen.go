package apis

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/service"
)

// PassgenInput is the request structure for the /api/passgen endpoint.
type PassgenInput struct {
	// Body contains the main input fields for the password generation.
	Body struct {
		Length    int  `json:"length" example:"16" doc:"Password length" minLength:"1"`
		Uppercase bool `json:"uppercase" example:"true" doc:"Include uppercase letters"`
		Lowercase bool `json:"lowercase" example:"true" doc:"Include lowercase letters"`
		Numbers   bool `json:"numbers" example:"true" doc:"Include digits"`
		Symbols   bool `json:"symbols" example:"true" doc:"Include symbols"`
	}
}

// PassgenOutput is the response structure for the /api/passgen endpoint.
type PassgenOutput struct {
	Body struct {
		Result string `json:"result"`
	}
}

// RegisterPassgen registers the /api/passgen endpoint
func RegisterPassgen(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "generate-password",
		Summary:       "generate random password - json",
		Description:   "Returns a JSON object containg the generated password entry.",
		Method:        http.MethodPost,
		Path:          "/api/passgen",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Passgen"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Successfully generated a password entry, responding with JSON content",
			},
			"400": {
				Description: "Bad request - invalid input.",
			},
		},
	}, func(ctx context.Context, input *PassgenInput) (*PassgenOutput, error) {
		// Validate input
		length := input.Body.Length
		if length <= 0 {
			return nil, huma.Error400BadRequest("length must be greater than zero")
		}

		if !input.Body.Uppercase && !input.Body.Lowercase && !input.Body.Numbers && !input.Body.Symbols {
			return nil, huma.Error400BadRequest("at least one character class must be selected")
		}

		// Generate password
		generatedPassword, err := service.GeneratePassword(length, input.Body.Uppercase, input.Body.Lowercase, input.Body.Numbers, input.Body.Symbols)
		if err != nil {
			logger.Logger.Error(
				"failed to generate password",
				"length", length,
				"uppercase", input.Body.Uppercase,
				"lowercase", input.Body.Lowercase,
				"numbers", input.Body.Numbers,
				"symbols", input.Body.Symbols,
				"error", err)
			return nil, huma.Error400BadRequest("failed to generate password: " + err.Error())
		}

		resp := &PassgenOutput{}
		resp.Body.Result = generatedPassword
		return resp, nil
	})
}
