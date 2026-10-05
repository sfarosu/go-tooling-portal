package apis

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/service"
)

// JsonPrettyInput is the request structure for the /api/jsonpretty endpoint.
type JsonPrettyInput struct {
	// Body contains the main input fields for the JSON prettifier.
	Body struct {
		Input     string `json:"Input" example:"{\"name\":\"Alice\",\"age\":30}" doc:"JSON input to prettify" minLength:"1"`
		Operation string `json:"operation" example:"pretty" doc:"Operation: pretty or unpretty" enum:"pretty,unpretty"`
		Spaces    int    `json:"spaces" example:"2" doc:"Number of spaces for indentation" minLength:"0"`
	}
}

// JsonPrettyOutput is the response structure for the /api/jsonpretty endpoint.
type JsonPrettyOutput struct {
	Body struct {
		Result string `json:"result"`
	}
}

// RegisterJsonPretty registers the /api/jsonpretty endpoint
func RegisterJsonPretty(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "json-pretty",
		Summary:       "prettify json",
		Description:   "Returns a JSON object containing the prettified JSON.",
		Method:        http.MethodPost,
		Path:          "/api/jsonpretty",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"JsonPretty"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Successfully prettified JSON, responding with JSON content",
			},
			"400": {
				Description: "Bad request, invalid input.",
			},
		},
	}, func(ctx context.Context, input *JsonPrettyInput) (*JsonPrettyOutput, error) {
		op := input.Body.Operation
		if op == "unpretty" {
			// Unpretty (minify) JSON
			result, err := service.UnprettyJSON(input.Body.Input)
			if err != nil {
				logger.Logger.Error(
					"failed to unpretty json",
					"input", input.Body.Input,
					"error", err,
				)
				return nil, huma.Error400BadRequest("failed to unpretty JSON: " + err.Error())
			}
			resp := &JsonPrettyOutput{}
			resp.Body.Result = result
			return resp, nil
		}

		// Default to pretty operation
		spaces := input.Body.Spaces
		if spaces <= 0 {
			return nil, huma.Error400BadRequest("spaces must be a valid non-negative integer")
		}

		prettyfiedJSON, err := service.PrettyJSONWithSpaces(input.Body.Input, input.Body.Spaces)
		if err != nil {
			logger.Logger.Error(
				"failed to prettify json",
				"input", input.Body.Input,
				"spaces", input.Body.Spaces,
				"error", err)
			return nil, huma.Error400BadRequest("failed to prettify JSON: " + err.Error())
		}
		resp := &JsonPrettyOutput{}
		resp.Body.Result = prettyfiedJSON
		return resp, nil
	})
}
