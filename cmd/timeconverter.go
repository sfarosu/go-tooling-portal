package cmd

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sfarosu/go-tooling-portal/cmd/helper"
	"github.com/sfarosu/go-tooling-portal/cmd/tmpl"
)

var (
	timeConversions = promauto.NewCounter(prometheus.CounterOpts{
		Name: "time_conversions_total",
		Help: "The total number of time conversions",
	})
)

func timeconvert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/timeconvert", http.StatusSeeOther)
	}

	nowUTC := time.Now().UTC()

	data := struct {
		CurrentUTCEpochTimeSeconds int64
		CurrentUTCEpochTimeMilli   int64
		CurrentUTCEpochTimeMicro   int64
		CurrentUTCYear             int
		CurrentUTCMonth            string
		CurrentUTCDay              string
		CurrentUTCHour             string
		CurrentUTCMinute           string
		CurrentUTCSecond           string
	}{
		CurrentUTCEpochTimeSeconds: nowUTC.Unix(),
		CurrentUTCEpochTimeMilli:   nowUTC.UnixMilli(),
		CurrentUTCEpochTimeMicro:   nowUTC.UnixMicro(),
		CurrentUTCYear:             nowUTC.Year(),
		CurrentUTCMonth:            helper.AddSecondDigit(int(nowUTC.Month())),
		CurrentUTCDay:              helper.AddSecondDigit(nowUTC.Day()),
		CurrentUTCHour:             helper.AddSecondDigit(nowUTC.Hour()),
		CurrentUTCMinute:           helper.AddSecondDigit(nowUTC.Minute()),
		CurrentUTCSecond:           helper.AddSecondDigit(nowUTC.Second()),
	}

	log.Println(r.URL.String(), r.Method, r.RemoteAddr, r.Proto, r.Header.Get("User-Agent"))
	errExec := tmpl.Tpl.ExecuteTemplate(w, "timeconvert.html", data)
	if errExec != nil {
		log.Println("error executing template: ", errExec)
	}

}

func timeconvertProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/timeconvert", http.StatusSeeOther)
	}

	nowUTC := time.Now().UTC()

	var epochToHuman bool
	var humanToEpoch bool
	var timezone string
	var insertedEpochTime int64
	var convertedEpochToHumanUTC time.Time
	var convertedEpochToHumanLOCAL time.Time
	var loc *time.Location
	var insertedHumanTime string
	var insertedHumanTimeParser time.Time
	var err error

	if r.FormValue("epochToHuman") == "true" {
		epochToHuman = true
		timezone = r.FormValue("browserTimeZoneFromEpochForm")
		loc, err = time.LoadLocation(timezone)
		if err != nil {
			log.Printf("error loading location '%v': %v", timezone, err)
		}
		insertedEpochTime, _ = strconv.ParseInt(r.FormValue("epochtime"), 10, 64)
		switch {
		case helper.GetNumberDigitsAmmount(insertedEpochTime) == 9 || helper.GetNumberDigitsAmmount(insertedEpochTime) == 10:
			convertedEpochToHumanUTC = time.Unix(insertedEpochTime, 0).UTC()
			convertedEpochToHumanLOCAL = time.Unix(insertedEpochTime, 0).In(loc)
		case helper.GetNumberDigitsAmmount(insertedEpochTime) == 12 || helper.GetNumberDigitsAmmount(insertedEpochTime) == 13:
			convertedEpochToHumanUTC = time.UnixMilli(insertedEpochTime).UTC()
			convertedEpochToHumanLOCAL = time.UnixMilli(insertedEpochTime).In(loc)
		case helper.GetNumberDigitsAmmount(insertedEpochTime) == 15 || helper.GetNumberDigitsAmmount(insertedEpochTime) == 16:
			convertedEpochToHumanUTC = time.UnixMicro(insertedEpochTime).UTC()
			convertedEpochToHumanLOCAL = time.UnixMicro(insertedEpochTime).In(loc)
		}
	} else if r.FormValue("humanToEpoch") == "true" {
		humanToEpoch = true
		timezone = r.FormValue("browserTimeZoneFromHumanForm")
		loc, err = time.LoadLocation(timezone)
		if err != nil {
			log.Printf("error loading location '%v': %v", timezone, err)
		}
		insertedHumanTime = r.FormValue("year") + "-" + r.FormValue("month") + "-" + r.FormValue("day") + "T" + r.FormValue("hour") + ":" + r.FormValue("minute") + ":" + r.FormValue("second") + " " + "UTC"
		layout := "2006-01-02T15:04:05 MST"
		insertedHumanTimeParser, err = time.Parse(layout, insertedHumanTime)
		if err != nil {
			log.Printf("error parsing inserted human time '%v' using layout '%v': %v", insertedHumanTime, layout, err)
		}
	}

	data := struct {
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
	}{
		Timezone:                         timezone,
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

	log.Println(r.URL.String(), r.Method, r.RemoteAddr, r.Proto, r.Header.Get("User-Agent"))
	errExec := tmpl.Tpl.ExecuteTemplate(w, "timeconvert-process.html", data)
	if errExec != nil {
		log.Println("Error executing template: ", errExec)
	}

	timeConversions.Inc()
}
