package apis

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/sfarosu/go-tooling-portal/internal/logger"
	"github.com/sfarosu/go-tooling-portal/internal/service"
)

// TimeConverterInput is the request structure for the /api/timeconverter endpoint.
type TimeConverterInput struct {
	// Body contains the main input fields for the time conversion.
	Body struct {
		EpochToHumanMode bool   `json:"epochToHumanMode" doc:"Convert from epoch to human time" enum:"true,false"`
		HumanToEpochMode bool   `json:"humanToEpochMode" doc:"Convert from human to epoch time" enum:"true,false"`
		EpochTime        int64  `json:"epochTime" doc:"Epoch timestamp (seconds, milliseconds, or microseconds)" minLength:"9"`
		Year             int    `json:"year" doc:"Year for human time conversion" minLength:"4"`
		Month            int    `json:"month" doc:"Month for human time conversion (01-12)" minLength:"2"`
		Day              int    `json:"day" doc:"Day for human time conversion (01-31)" minLength:"2"`
		Hour             int    `json:"hour" doc:"Hour for human time conversion (0-23)" minLength:"2"`
		Minute           int    `json:"minute" doc:"Minute for human time conversion (0-59)" minLength:"2"`
		Second           int    `json:"second" doc:"Second for human time conversion (0-59)" minLength:"2"`
		Timezone         string `json:"timezone" example:"UTC" doc:"IANA timezone identifier (e.g., UTC, America/New_York, Europe/Paris)" minLength:"2"`
	}
}

// TimeConverterOutput is the response structure for the /api/timeconverter endpoint.
type TimeConverterOutput struct {
	Body struct {
		Timezone                         string `json:"timezone"`
		EpochToHuman                     bool   `json:"epochToHuman"`
		HumanToEpoch                     bool   `json:"humanToEpoch"`
		CurrentUTCEpochTimeSeconds       int64  `json:"currentUTCEpochTimeSeconds"`
		CurrentUTCEpochTimeMilli         int64  `json:"currentUTCEpochTimeMilli"`
		CurrentUTCEpochTimeMicro         int64  `json:"currentUTCEpochTimeMicro"`
		CurrentUTCHumanTime              string `json:"currentUTCHumanTime"`
		CurrentLOCALHumanTime            string `json:"currentLOCALHumanTime"`
		ConvertedEpochToHumanUTC         string `json:"convertedEpochToHumanUTC"`
		ConvertedEpochToHumanLOCAL       string `json:"convertedEpochToHumanLOCAL"`
		ConvertedHumanToEpochTimeSeconds int64  `json:"convertedHumanToEpochTimeSeconds"`
		ConvertedHumanToEpochTimeMilli   int64  `json:"convertedHumanToEpochTimeMilli"`
		ConvertedHumanToEpochTimeMicro   int64  `json:"convertedHumanToEpochTimeMicro"`
		ConvertedHumanToEpochTimeNano    int64  `json:"convertedHumanToEpochTimeNano"`
	}
}

// RegisterTimeConverter registers the /api/timeconverter endpoint with the given Huma API
func RegisterTimeConverter(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "time-converter",
		Summary:       "time conversion - json",
		Description:   "Returns a JSON object containing time conversion results for epoch to human or human to epoch conversion.",
		Method:        http.MethodPost,
		Path:          "/api/timeconverter",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"TimeConverter"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Successful response with the time conversion results.",
			},
			"400": {
				Description: "Bad request - invalid input.",
			},
		},
	}, func(ctx context.Context, input *TimeConverterInput) (*TimeConverterOutput, error) {
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
			return nil, huma.Error400BadRequest("time conversion failed: " + err.Error())
		}

		// Map result to output
		resp := &TimeConverterOutput{}
		resp.Body.Timezone = result.Timezone
		resp.Body.EpochToHuman = result.EpochToHuman
		resp.Body.HumanToEpoch = result.HumanToEpoch
		resp.Body.CurrentUTCEpochTimeSeconds = result.CurrentUTCEpochTimeSeconds
		resp.Body.CurrentUTCEpochTimeMilli = result.CurrentUTCEpochTimeMilli
		resp.Body.CurrentUTCEpochTimeMicro = result.CurrentUTCEpochTimeMicro
		resp.Body.CurrentUTCHumanTime = result.CurrentUTCHumanTime
		resp.Body.CurrentLOCALHumanTime = result.CurrentLOCALHumanTime
		resp.Body.ConvertedEpochToHumanUTC = result.ConvertedEpochToHumanUTC.String()
		resp.Body.ConvertedEpochToHumanLOCAL = result.ConvertedEpochToHumanLOCAL.String()
		resp.Body.ConvertedHumanToEpochTimeSeconds = result.ConvertedHumanToEpochTimeSeconds
		resp.Body.ConvertedHumanToEpochTimeMilli = result.ConvertedHumanToEpochTimeMilli
		resp.Body.ConvertedHumanToEpochTimeMicro = result.ConvertedHumanToEpochTimeMicro
		resp.Body.ConvertedHumanToEpochTimeNano = result.ConvertedHumanToEpochTimeNano

		return resp, nil
	})
}
