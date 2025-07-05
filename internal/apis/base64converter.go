package apis

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/service"
)

type Base64ConverterInput struct {
	Body struct {
		Input     string `json:"input" example:"hello" doc:"String to encode or decode" minLength:"1"`
		Operation string `json:"operation" example:"encode" doc:"Operation to perform: encode or decode" enum:"encode,decode"`
		Format    string `json:"format" example:"standard" doc:"Base64 format: standard or url-compatible" enum:"standard,url-compatible"`
	}
}

type Base64ConverterOutput struct {
	Body struct {
		Result string `json:"result"`
	}
}

// RegisterBase64Converter registers the /api/base64converter endpoint with the given Huma API
func RegisterBase64Converter(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "base64-converter",
		Summary:       "base64 encode/decode - json",
		Description:   "Returns a JSON object containg a base64-encoded or decoded result depending on the selected operation.",
		Method:        http.MethodPost,
		Path:          "/api/base64converter",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Base64Converter"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Successful response with the base64 encoded or decoded result",
			},
			"400": {
				Description: "Bad Request, the input data is missing, malformed, or the specified operation is invalid (must be 'encode' or 'decode')",
			},
		},
	}, func(ctx context.Context, input *Base64ConverterInput) (*Base64ConverterOutput, error) {
		result, err := service.Base64Convert(input.Body.Input, input.Body.Operation, input.Body.Format)
		if err != nil {
			logger.Logger.Error(
				"failed to process base64 conversion",
				"input", input.Body.Input,
				"operation", input.Body.Operation,
				"format", input.Body.Format,
				"error", err,
			)
			return nil, huma.Error400BadRequest("base64 conversion failed: " + err.Error())
		}

		resp := &Base64ConverterOutput{}
		resp.Body.Result = result
		return resp, nil
	})
}
