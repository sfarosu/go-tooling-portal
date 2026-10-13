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
	Millisecond      int
	Microsecond      int

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
		precision := req.Precision
		if precision == "" {
			switch helper.GetNumberDigitsAmmount(req.EpochTime) {
			case 9, 10:
				precision = "seconds"
			case 12, 13:
				precision = "millis"
			case 15, 16:
				precision = "micros"
			default:
				return nil, fmt.Errorf("invalid epoch time precision: %d digits", helper.GetNumberDigitsAmmount(req.EpochTime))
			}
		}
		switch precision {
		case "seconds":
			// Seconds precision
			convertedEpochToHumanUTC = time.Unix(req.EpochTime, 0).UTC()
			convertedEpochToHumanLOCAL = time.Unix(req.EpochTime, 0).In(loc)
		case "millis":
			// Milliseconds precision
			convertedEpochToHumanUTC = time.UnixMilli(req.EpochTime).UTC()
			convertedEpochToHumanLOCAL = time.UnixMilli(req.EpochTime).In(loc)
		case "micros":
			// Microseconds precision
			convertedEpochToHumanUTC = time.UnixMicro(req.EpochTime).UTC()
			convertedEpochToHumanLOCAL = time.UnixMicro(req.EpochTime).In(loc)
		default:
			return nil, fmt.Errorf("invalid epoch time precision: %q", precision)
		}
	} else if req.HumanToEpochMode {
		humanToEpoch = true
		if req.Millisecond < 0 || req.Millisecond > 999 {
			return nil, fmt.Errorf("millisecond must be between 0 and 999")
		}
		if req.Microsecond < 0 || req.Microsecond > 999 {
			return nil, fmt.Errorf("microsecond must be between 0 and 999")
		}
		if req.Month < 1 || req.Month > 12 {
			return nil, fmt.Errorf("month must be between 1 and 12")
		}
		if req.Day < 1 || req.Day > 31 {
			return nil, fmt.Errorf("day must be between 1 and 31")
		}
		if req.Hour < 0 || req.Hour > 23 {
			return nil, fmt.Errorf("hour must be between 0 and 23")
		}
		if req.Minute < 0 || req.Minute > 59 {
			return nil, fmt.Errorf("minute must be between 0 and 59")
		}
		if req.Second < 0 || req.Second > 59 {
			return nil, fmt.Errorf("second must be between 0 and 59")
		}
		// Create time from provided components in the specified timezone
		nanosecond := req.Millisecond*int(time.Millisecond) + req.Microsecond*int(time.Microsecond)
		insertedHumanTimeParser = time.Date(req.Year, time.Month(req.Month), req.Day, req.Hour, req.Minute, req.Second, nanosecond, loc)
		year, month, day := insertedHumanTimeParser.Date()
		if year != req.Year || int(month) != req.Month || day != req.Day {
			return nil, fmt.Errorf("invalid date: %04d-%02d-%02d", req.Year, req.Month, req.Day)
		}
		if insertedHumanTimeParser.Hour() != req.Hour || insertedHumanTimeParser.Minute() != req.Minute || insertedHumanTimeParser.Second() != req.Second {
			return nil, fmt.Errorf("invalid local time: %02d:%02d:%02d does not exist in %s", req.Hour, req.Minute, req.Second, req.Timezone)
		}
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
