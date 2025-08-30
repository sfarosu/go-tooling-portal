package apis

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/service"
)

type FormatConverterInput struct {
	Body struct {
		Input string `json:"input" example:"{\"name\":\"Alice\",\"age\":30}" doc:"JSON or YAML string to convert to the other format" minLength:"1"`
	}
}

type FormatConverterOutput struct {
	Body struct {
		Result string `json:"result"`
	}
}

// RegisterFormatConverter registers the /api/formatconverter endpoint with the given Huma API
func RegisterFormatConverter(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "format-converter",
		Summary:       "convert format between json and yaml",
		Description:   "Convert data between JSON and YAML formats.",
		Method:        http.MethodPost,
		Path:          "/api/formatconverter",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"FormatConverter"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Successful response with the converted result.",
			},
			"400": {
				Description: "Bad Request, the input data is missing or malformed.",
			},
		},
	}, func(ctx context.Context, input *FormatConverterInput) (*FormatConverterOutput, error) {
		result, err := service.FormatConvert(input.Body.Input)
		if err != nil {
			logger.Logger.Error(
				"failed to process format conversion",
				"input", input.Body.Input,
				"error", err,
			)
			return nil, huma.Error400BadRequest("format conversion failed: " + err.Error())
		}

		resp := &FormatConverterOutput{}
		resp.Body.Result = result
		return resp, nil
	})
}
