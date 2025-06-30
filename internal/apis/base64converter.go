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
		Input     string `json:"input" example:"hello" doc:"String to encode or decode"`
		Operation string `json:"operation" example:"encode" doc:"Operation to perform: encode or decode"`
		Format    string `json:"format" example:"standard" doc:"Base64 format: standard or url-compatible"`
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
		Summary:       "Base64 Encode/Decode",
		Method:        http.MethodPost,
		Path:          "/api/base64converter",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"Base64Converter"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Base64 encode or decode result",
			},
			"400": {
				Description: "Bad Request - Invalid input or operation",
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
			return nil, huma.Error400BadRequest("failed to process base64 conversion: " + err.Error())
		}

		resp := &Base64ConverterOutput{}
		resp.Body.Result = result
		return resp, nil
	})
}
