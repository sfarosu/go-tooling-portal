package apis

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/service"
)

// URLDecoderInput is the request structure for the /api/urldecoder endpoint.
type URLDecoderInput struct {
	Body struct {
		URL string `json:"url" example:"https://example.com/foo+bar%21?service=ab%2Bc&product=de%24f&service=second" doc:"URL to parse and decode" minLength:"1"`
	}
}

// URLDecoderOutput is the response structure for the /api/urldecoder endpoint.
type URLDecoderOutput struct {
	Body struct {
		Result string `json:"result"`
	}
}

// RegisterURLDecoder registers the /api/urldecoder endpoint.
func RegisterURLDecoder(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "url-decoder",
		Summary:       "Parse and decode a URL (JSON)",
		Description:   "Parses a URL and returns its components and decoded query parameters in a YAML string inside a JSON response.",
		Method:        http.MethodPost,
		Path:          "/api/urldecoder",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"URL Decoder"},
		Responses: map[string]*huma.Response{
			"200": {Description: "JSON response containing decoded URL components as YAML."},
			"400": {Description: "Bad request - invalid URL."},
		},
	}, func(ctx context.Context, input *URLDecoderInput) (*URLDecoderOutput, error) {
		result, err := service.URLDecode(input.Body.URL)
		if err != nil {
			if logger.Logger != nil {
				logger.Logger.Error("failed to decode URL", "error", err)
			}
			return nil, huma.Error400BadRequest("URL decoding failed: " + err.Error())
		}

		resp := &URLDecoderOutput{}
		resp.Body.Result = result
		return resp, nil
	})
}
