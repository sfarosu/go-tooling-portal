package apis

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/service"
)

// FormatConverterInput is the request structure for the /api/formatconverter endpoint.
type FormatConverterInput struct {
	// Body contains the main input fields for format conversion.
	Body struct {
		Input string `json:"input" example:"{\"name\":\"Alice\",\"age\":30}" doc:"JSON or YAML string to convert to the other format" minLength:"1"`
	}
}

// FormatConverterOutput is the response structure for the /api/formatconverter endpoint.
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
		Description:   "Returns a JSON object containg the conversion.",
		Method:        http.MethodPost,
		Path:          "/api/formatconverter",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"FormatConverter"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Successful conversion, responding with JSON content.",
			},
			"400": {
				Description: "Bad request - invalid input.",
			},
		},
	}, func(ctx context.Context, input *FormatConverterInput) (*FormatConverterOutput, error) {
		// Process format conversion
		formatConversion, err := service.FormatConvert(input.Body.Input)
		if err != nil {
			logger.Logger.Error(
				"failed to perform format conversion",
				"input", input.Body.Input,
				"error", err,
			)
			return nil, huma.Error400BadRequest("format conversion failed: " + err.Error())
		}

		resp := &FormatConverterOutput{}
		resp.Body.Result = formatConversion
		return resp, nil
	})
}
