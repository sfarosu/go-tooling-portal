package service

import (
	"fmt"
	"time"

	"github.com/sfarosu/go-tooling-portal/internal/helper"
)

// TimeConvertRequest contains the input data for time conversion operations.
type TimeConvertRequest struct {
	// Epoch to Human conversion fields
	EpochToHumanMode bool
	EpochTime        int64
	Precision        string // "seconds", "millis", or "micros"

	// Human to Epoch conversion fields
	HumanToEpochMode bool
	Year             int
	Month            int
	Day              int
	Hour             int
	Minute           int
	Second           int

	// Timezone for both conversions
	Timezone string
}

// TimeConvertResult contains the output data for time conversion operations.
type TimeConvertResult struct {
	Timezone                         string
	EpochToHuman                     bool
	HumanToEpoch                     bool
	CurrentUTCEpochTimeSeconds       int64
	CurrentUTCEpochTimeMilli         int64
	CurrentUTCEpochTimeMicro         int64
	CurrentUTCHumanTime              string
	CurrentLOCALHumanTime            string
	ConvertedEpochToHumanUTC         time.Time
	ConvertedEpochToHumanLOCAL       time.Time
	ConvertedHumanToEpochTimeSeconds int64
	ConvertedHumanToEpochTimeMilli   int64
	ConvertedHumanToEpochTimeMicro   int64
	ConvertedHumanToEpochTimeNano    int64
}

// TimeConvert performs time conversion operations based on the request.
// It returns a TimeConvertResult containing conversion results and current time information.
func TimeConvert(req *TimeConvertRequest) (*TimeConvertResult, error) {
	nowUTC := time.Now().UTC()

	var epochToHuman bool
	var humanToEpoch bool
	var convertedEpochToHumanUTC time.Time
	var convertedEpochToHumanLOCAL time.Time
	var insertedHumanTimeParser time.Time
	var err error

	// Load the timezone
	loc, err := time.LoadLocation(req.Timezone)
	if err != nil {
		return nil, fmt.Errorf("error loading location '%s': %w", req.Timezone, err)
	}

	// Process Epoch to Human conversion if requested
	if req.EpochToHumanMode {
		epochToHuman = true
		numDigits := helper.GetNumberDigitsAmmount(req.EpochTime)
		switch numDigits {
		case 9, 10:
			// Seconds precision
			convertedEpochToHumanUTC = time.Unix(req.EpochTime, 0).UTC()
			convertedEpochToHumanLOCAL = time.Unix(req.EpochTime, 0).In(loc)
		case 12, 13:
			// Milliseconds precision
			convertedEpochToHumanUTC = time.UnixMilli(req.EpochTime).UTC()
			convertedEpochToHumanLOCAL = time.UnixMilli(req.EpochTime).In(loc)
		case 15, 16:
			// Microseconds precision
			convertedEpochToHumanUTC = time.UnixMicro(req.EpochTime).UTC()
			convertedEpochToHumanLOCAL = time.UnixMicro(req.EpochTime).In(loc)
		default:
			return nil, fmt.Errorf("invalid epoch time precision: %d digits", numDigits)
		}
	} else if req.HumanToEpochMode {
		humanToEpoch = true
		// Create time from provided components in the specified timezone
		insertedHumanTimeParser = time.Date(req.Year, time.Month(req.Month), req.Day, req.Hour, req.Minute, req.Second, 0, loc)
	} else {
		return nil, fmt.Errorf("neither epochToHuman nor humanToEpoch mode specified")
	}

	// Build the result
	result := &TimeConvertResult{
		Timezone:                         req.Timezone,
		EpochToHuman:                     epochToHuman,
		HumanToEpoch:                     humanToEpoch,
		CurrentUTCEpochTimeSeconds:       nowUTC.Unix(),
		CurrentUTCEpochTimeMilli:         nowUTC.UnixMilli(),
		CurrentUTCEpochTimeMicro:         nowUTC.UnixMicro(),
		CurrentUTCHumanTime:              nowUTC.Format("2006-01-02 15:04:05 -0700 MST"),
		CurrentLOCALHumanTime:            nowUTC.In(loc).Format("2006-01-02 15:04:05 -0700 MST"),
		ConvertedEpochToHumanUTC:         convertedEpochToHumanUTC,
		ConvertedEpochToHumanLOCAL:       convertedEpochToHumanLOCAL,
		ConvertedHumanToEpochTimeSeconds: insertedHumanTimeParser.Unix(),
		ConvertedHumanToEpochTimeMilli:   insertedHumanTimeParser.UnixMilli(),
		ConvertedHumanToEpochTimeMicro:   insertedHumanTimeParser.UnixMicro(),
		ConvertedHumanToEpochTimeNano:    insertedHumanTimeParser.UnixNano(),
	}

	return result, nil
}
