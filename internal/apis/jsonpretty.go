package apis

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sfarosu/go-tooling-portal/internal/service"
)

// JsonPrettyInput is the request structure for the /api/jsonpretty endpoint.
type JsonPrettyInput struct {
	// Body contains the main input fields for the JSON prettifier.
	Body struct {
		Input     string `json:"input" example:"{\"name\":\"Alice\",\"age\":30}" doc:"JSON input to prettify" minLength:"1"`
		Operation string `json:"operation" example:"pretty" doc:"Operation: pretty or unpretty" enum:"pretty,unpretty"`
		Spaces    int    `json:"spaces" example:"2" doc:"Number of spaces for indentation (0-16)" minimum:"0"`
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
		Summary:       "Format or minify JSON (JSON)",
		Description:   "Formats JSON with the requested indentation or minifies it, returning the result as JSON.",
		Method:        http.MethodPost,
		Path:          "/api/jsonpretty",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"JSON Prettifier"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "JSON response containing the formatted or minified result.",
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
				logError(ctx,
					"failed to unpretty json",
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
		if spaces < 0 {
			return nil, huma.Error400BadRequest("spaces must be a valid non-negative integer")
		}

		prettyfiedJSON, err := service.PrettyJSONWithSpaces(input.Body.Input, input.Body.Spaces)
		if err != nil {
			logError(ctx,
				"failed to prettify json",
				"spaces", input.Body.Spaces,
				"error", err)
			return nil, huma.Error400BadRequest("failed to prettify JSON: " + err.Error())
		}
		resp := &JsonPrettyOutput{}
		resp.Body.Result = prettyfiedJSON
		return resp, nil
	})
}
