package apis

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/service"
)

// URLDecoderHTMXInput is the request structure for the /api/urldecoder/htmx endpoint.
type URLDecoderHTMXInput struct {
	HtmxHeader bool `header:"HX-Request"`
	Body       struct {
		URL string `json:"url" example:"https://example.com/foo%21?query=a%2Bb" doc:"URL to parse and decode" minLength:"1"`
	}
}

// Resolve detects whether the request originated from HTMX.
func (input *URLDecoderHTMXInput) Resolve(ctx huma.Context) []error {
	value, err := strconv.ParseBool(ctx.Header("HX-Request"))
	if err != nil {
		input.HtmxHeader = false
		return nil
	}
	input.HtmxHeader = value
	return nil
}

// URLDecoderHTMXOutput returns an HTML fragment for HTMX requests.
type URLDecoderHTMXOutput struct {
	ContentType string `header:"Content-Type"`
	Body        []byte `json:"-"`
}

var urlDecoderResultTmpl = template.Must(template.New("urldecoder-result").Parse(`
<div class="row g-3">
  <div class="col-12">
    <label class="form-label fw-bold" for="urldecoder-result-input">Decoded URL components</label>
    <div class="input-group">
      <textarea class="form-control custom-output" id="urldecoder-result-input" readonly rows="10">{{.Result}}</textarea>
      <button class="btn btn-graphite" type="button" onclick="copyToClipboard('urldecoder-result-input')" title="Copy decoded URL" aria-label="Copy decoded URL">
        <i class="bi bi-clipboard"></i>
      </button>
    </div>
  </div>
</div>
`))

// RegisterURLDecoderHtmx registers the /api/urldecoder/htmx endpoint.
func RegisterURLDecoderHtmx(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "url-decoder-htmx",
		Summary:       "decode URL - htmx",
		Description:   "Returns decoded URL components as an HTML fragment for HTMX requests.",
		Method:        http.MethodPost,
		Path:          "/api/urldecoder/htmx",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"URLDecoder"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Successfully parsed and decoded the URL.",
				Content:     map[string]*huma.MediaType{"text/html": {}},
			},
			"400": {Description: "Bad request - invalid URL."},
		},
	}, func(ctx context.Context, input *URLDecoderHTMXInput) (*URLDecoderHTMXOutput, error) {
		result, err := service.URLDecode(input.Body.URL)
		if err != nil {
			if logger.Logger != nil {
				logger.Logger.Error("failed to decode URL", "error", err)
			}
			if input.HtmxHeader {
				return &URLDecoderHTMXOutput{
					ContentType: "text/html; charset=utf-8",
					Body:        []byte(generateHTMXError(err)),
				}, nil
			}
			return nil, huma.Error400BadRequest("URL decoding failed: " + err.Error())
		}

		if !input.HtmxHeader {
			return nil, huma.Error400BadRequest("This endpoint only supports HTMX requests with HX-Request header set to true; use /api/urldecoder for JSON responses")
		}

		var buffer bytes.Buffer
		if err := urlDecoderResultTmpl.Execute(&buffer, map[string]string{"Result": result}); err != nil {
			if logger.Logger != nil {
				logger.Logger.Error("failed to render URL decoder result", "error", err)
			}
			return nil, huma.Error500InternalServerError("failed to render URL decoder result")
		}
		return &URLDecoderHTMXOutput{
			ContentType: "text/html; charset=utf-8",
			Body:        buffer.Bytes(),
		}, nil
	})
}
