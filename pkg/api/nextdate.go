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

	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid interval format")
		}

		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				return date.Format(dateFormat), nil
			}
		}
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid interval format")
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
	case "w":
		if len(parts) != 2 {
			return "", errors.New("invalid interval format")
		}

		var allowed [8]bool
		weekdayIndex := strings.Split(parts[1], ",")

		for _, i := range weekdayIndex {
			day, err := strconv.Atoi(i)
			if err != nil {
				return "", errors.New("invalid days format")
			}
			if day < 1 || day > 7 {
				return "", errors.New("invalid days range")
			}

			allowed[day] = true
		}

		for {
			date = date.AddDate(0, 0, 1)
			weekday := int(date.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			if allowed[weekday] && afterNow(date, now) {
				return date.Format(dateFormat), nil
			}
		}
	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("invalid interval format")
		}

		var day [32]bool
		var month [13]bool
		var lastDay bool
		var penultDay bool

		dayIndex := strings.Split(parts[1], ",")

		for _, i := range dayIndex {
			dayConv, err := strconv.Atoi(i)
			if err != nil {
				return "", errors.New("invalid days format")
			}

			switch {
			case dayConv >= 1 && dayConv <= 31:
				day[dayConv] = true
			case dayConv == -1:
				lastDay = true
			case dayConv == -2:
				penultDay = true
			default:
				return "", errors.New("invalid days range")
			}
		}
		if len(parts) == 3 {
			monthIndex := strings.Split(parts[2], ",")

			for _, i := range monthIndex {
				monthConv, err := strconv.Atoi(i)
				if err != nil {
					return "", errors.New("invalid month format")
				}
				if monthConv >= 1 && monthConv <= 12 {
					month[monthConv] = true
				} else {
					return "", errors.New("invalid months range")
				}
			}
		} else {
			for i := 1; i <= 12; i++ {
				month[i] = true
			}
		}

		for {
			date = date.AddDate(0, 0, 1)

			checkMonth := int(date.Month())
			if month[checkMonth] {
				checkDay := date.Day()
				daysInMonth := time.Date(
					date.Year(),
					date.Month()+1,
					0,
					0, 0, 0, 0,
					time.UTC,
				).Day()

				isLastDay := lastDay && checkDay == daysInMonth
				isPenultDay := penultDay && checkDay == (daysInMonth-1)

				if day[checkDay] || isLastDay || isPenultDay {
					if afterNow(date, now) {
						return date.Format(dateFormat), nil
					}
				}
			}
		}
	default:
		return "", errors.New("unknown repeat format")
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
