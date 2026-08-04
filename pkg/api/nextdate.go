package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	nowOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return dateOnly.After(nowOnly)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}

	if repeat == "y" {
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				return date.Format(dateFormat), nil
			}
		}
	} else {
		parts := strings.Split(repeat, " ")
		if len(parts) != 2 {
			return "", errors.New("invalid interval format")
		}

		if parts[0] != "d" {
			return "", errors.New("unknown interval format")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", err
		}

		if days <= 0 || days > 400 {
			return "", errors.New("invalid days interval")
		}

		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				return date.Format(dateFormat), nil
			}
		}
	}
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	var now time.Time
	var err error

	nowString := r.FormValue("now")
	if nowString == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateFormat, nowString)
		if err != nil {
			fmt.Fprintf(w, "error: %v\n", err)
			return
		}
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		fmt.Fprintf(w, "error: %v\n", err)
		return
	}

	if _, err := fmt.Fprintf(w, "%s\n", nextDate); err != nil {
		log.Printf("error: %v\n", err)
		return
	}
}
