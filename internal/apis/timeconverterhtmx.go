package apis

import (
	"bytes"
	"context"
	"net/http"
	"strconv"
	"text/template"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/service"
)

// TimeConverterHTMXInput is the request structure for the /api/timeconverter/htmx endpoint.
type TimeConverterHTMXInput struct {
	// HtmxHeader stores the value of the HX-Request header, which is sent by HTMX
	// when making AJAX requests. This allows the handler to detect if the request
	// originated from an HTMX-enabled frontend and respond with HTML instead of JSON
	HtmxHeader bool `header:"HX-Request"`

	// Body contains the main input fields for the time conversion.
	Body struct {
		EpochToHumanMode bool   `json:"epochToHumanMode" doc:"Convert from epoch to human time"`
		HumanToEpochMode bool   `json:"humanToEpochMode" doc:"Convert from human to epoch time"`
		EpochTime        int64  `json:"epochTime" doc:"Epoch timestamp (seconds, milliseconds, or microseconds)"`
		Year             int    `json:"year" doc:"Year for human time conversion"`
		Month            int    `json:"month" doc:"Month for human time conversion (1-12)"`
		Day              int    `json:"day" doc:"Day for human time conversion (1-31)"`
		Hour             int    `json:"hour" doc:"Hour for human time conversion (0-23)"`
		Minute           int    `json:"minute" doc:"Minute for human time conversion (0-59)"`
		Second           int    `json:"second" doc:"Second for human time conversion (0-59)"`
		Timezone         string `json:"timezone" example:"UTC" doc:"IANA timezone identifier (e.g., UTC, America/New_York, Europe/Paris)"`
	}
}

// Resolve implements the huma.RequestResolver interface.
// It is called by Huma before the handler to allow custom extraction or transformation
// of request data. Here, we extract the HX-Request header (sent by HTMX)
// and store it in the HtmxHeader field for later use in the handler
// See: https://huma.rocks/features/request-resolvers/
func (m *TimeConverterHTMXInput) Resolve(ctx huma.Context) []error {
	val := ctx.Header("HX-Request")
	parsed, err := strconv.ParseBool(val)
	if err != nil {
		// If parsing fails, assume it's not an HTMX request
		m.HtmxHeader = false
		return nil
	}
	m.HtmxHeader = parsed
	return nil
}

// TimeConverterHTMXOutput is the response structure for the /api/timeconverter/htmx endpoint
type TimeConverterHTMXOutput struct {
	ContentType string `header:"Content-Type"`
	Body        []byte `json:"-"` // Body is not serialized to JSON, but used to return HTML content
}

// HTMX-compatible HTML fragment rendered dynamically by the API and injected into the page via HTMX
var timeconverterResultTmpl = template.Must(template.New("timeconverter-result").Parse(`
<div class="alert alert-success mb-0">
  <div class="row g-3">
    <div class="col-12">
      <h6 class="mb-2">Conversion Results</h6>
    </div>
    {{if .EpochToHuman}}
    <div class="col-md-6">
      <label class="form-label fw-bold">UTC Time</label>
      <input type="text" class="form-control" value="{{.ConvertedEpochToHumanUTC}}" readonly>
    </div>
    <div class="col-md-6">
      <label class="form-label fw-bold">{{.Timezone}} Time</label>
      <input type="text" class="form-control" value="{{.ConvertedEpochToHumanLOCAL}}" readonly>
    </div>
    {{else if .HumanToEpoch}}
    <div class="col-md-6">
      <label class="form-label fw-bold">Epoch (Seconds)</label>
      <input type="text" class="form-control" value="{{.ConvertedHumanToEpochTimeSeconds}}" readonly>
    </div>
    <div class="col-md-6">
      <label class="form-label fw-bold">Epoch (Milliseconds)</label>
      <input type="text" class="form-control" value="{{.ConvertedHumanToEpochTimeMilli}}" readonly>
    </div>
    <div class="col-md-6">
      <label class="form-label fw-bold">Epoch (Microseconds)</label>
      <input type="text" class="form-control" value="{{.ConvertedHumanToEpochTimeMicro}}" readonly>
    </div>
    {{end}}
  </div>
</div>
`))

// RegisterTimeConverterHtmx registers the /api/timeconverter/htmx endpoint with the given Huma API
func RegisterTimeConverterHtmx(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "time-converter-htmx",
		Summary:       "time conversion - htmx",
		Description:   "Returns a suitable for HTMX injection HTML object containing the time conversion results.",
		Method:        http.MethodPost,
		Path:          "/api/timeconverter/htmx",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"TimeConverter"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Successfully converted the time, responding with HTML content for HTMX requests.",
				Content: map[string]*huma.MediaType{
					"text/html": {},
				},
			},
			"400": {
				Description: "Bad Request, invalid input (including missing HX-Request header).",
			},
		},
	}, func(ctx context.Context, input *TimeConverterHTMXInput) (*TimeConverterHTMXOutput, error) {
		// Create service request from input
		req := &service.TimeConvertRequest{
			EpochToHumanMode: input.Body.EpochToHumanMode,
			EpochTime:        input.Body.EpochTime,
			HumanToEpochMode: input.Body.HumanToEpochMode,
			Year:             input.Body.Year,
			Month:            input.Body.Month,
			Day:              input.Body.Day,
			Hour:             input.Body.Hour,
			Minute:           input.Body.Minute,
			Second:           input.Body.Second,
			Timezone:         input.Body.Timezone,
		}

		// Perform time conversion
		result, err := service.TimeConvert(req)
		if err != nil {
			if logger.Logger != nil {
				logger.Logger.Error(
					"failed to process time conversion",
					"epochToHumanMode", input.Body.EpochToHumanMode,
					"humanToEpochMode", input.Body.HumanToEpochMode,
					"timezone", input.Body.Timezone,
					"error", err,
				)
			}
			if input.HtmxHeader {
				// If it's a htmx request, return an HTML fragment with the error content (return code will be 200 as we want to display it in the UI)
				return &TimeConverterHTMXOutput{
					ContentType: "text/html; charset=utf-8",
					Body:        []byte(generateHTMXError(err)),
				}, nil
			}
			return nil, huma.Error400BadRequest("time conversion failed: " + err.Error())
		}

		// Render the HTML template and return as a html block if this is an HTMX request
		if input.HtmxHeader {
			var buf bytes.Buffer
			timeconverterResultTmpl.Execute(&buf, map[string]interface{}{
				"Timezone":                    result.Timezone,
				"EpochToHuman":                result.EpochToHuman,
				"HumanToEpoch":                result.HumanToEpoch,
				"ConvertedEpochToHumanUTC":    result.ConvertedEpochToHumanUTC.String(),
				"ConvertedEpochToHumanLOCAL":  result.ConvertedEpochToHumanLOCAL.String(),
				"ConvertedHumanToEpochTimeSeconds": result.ConvertedHumanToEpochTimeSeconds,
				"ConvertedHumanToEpochTimeMilli":   result.ConvertedHumanToEpochTimeMilli,
				"ConvertedHumanToEpochTimeMicro":   result.ConvertedHumanToEpochTimeMicro,
			})
			resp := &TimeConverterHTMXOutput{
				ContentType: "text/html; charset=utf-8",
				Body:        buf.Bytes(),
			}
			return resp, nil
		}

		return nil, huma.Error400BadRequest("This endpoint only supports HTMX requests with HX-Request header set to true; use /api/timeconverter for JSON responses")
	})
}
