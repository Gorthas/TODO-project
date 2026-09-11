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

// Сравниваем календарные даты без учета времени суток.
// Возвращаем true, только если первая дата строго позже второй
func isDateAfter(first, second time.Time) bool {
	firstOnly := time.Date(first.Year(), first.Month(), first.Day(), 0, 0, 0, 0, time.UTC)
	secondOnly := time.Date(second.Year(), second.Month(), second.Day(), 0, 0, 0, 0, time.UTC)
	return firstOnly.After(secondOnly)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "y": // ежегодный повтор
		if len(parts) != 1 {
			return "", errors.New("invalid interval format")
		}

		for {
			date = date.AddDate(1, 0, 0)
			if isDateAfter(date, now) {
				return date.Format(dateFormat), nil
			}
		}
	case "d": // повтор через заданное количество дней
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
			if isDateAfter(date, now) {
				return date.Format(dateFormat), nil
			}
		}
	case "w": // повтор по выбранным дням недели
		if len(parts) != 2 {
			return "", errors.New("invalid interval format")
		}

		// Массив отмечает выбранные дни недели правила повтора
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
			if allowed[weekday] && isDateAfter(date, now) {
				return date.Format(dateFormat), nil
			}
		}
	case "m": // повтор по выбранным дням месяца и, при необходимости, месяцам
		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("invalid interval format")
		}

		// Массивы отмечают выбранные дни и месяцы правила повтора.
		// Последний (-1) и предпоследний (-2) дни месяца
		// хранятся отдельно в соответствующих флагах
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

			// Обычные значения задают конкретный день месяца,
			// -1 — последний день, -2 — предпоследний
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

		// Если в правиле указан список месяцев, используем только их;
		// иначе повтор разрешён во всех месяцах
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

		var possibleDate bool

		// Проверяем, существует ли хотя бы одна дата для выбранных дней и месяцев,
		// чтобы не уйти в бесконечный поиск. Берём високосный 2024 год,
		// чтобы 29 февраля считалось возможной датой
		for i := 1; i <= 12; i++ {
			if month[i] {
				daysInMonth := time.Date(
					2024,
					time.Month(i)+1,
					0,
					0, 0, 0, 0,
					time.UTC,
				).Day()
				for j := 1; j <= 31; j++ {
					if day[j] {
						if j <= daysInMonth {
							possibleDate = true
						}
					}
				}
			}
		}
		if !possibleDate && !lastDay && !penultDay {
			return "", errors.New("invalid combination of day and month")
		}

		// Ищем ближайшую дату повтора, двигаясь по одному дню вперед.
		// Для выбранного месяца проверяем обычные дни, а также
		// последний (-1) и предпоследний (-2) дни месяца
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
					if isDateAfter(date, now) {
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
			message := fmt.Sprintf("error: %v", err)
			http.Error(w, message, http.StatusBadRequest)
			return
		}
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		message := fmt.Sprintf("error: %v", err)
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	if _, err := fmt.Fprintf(w, "%s\n", nextDate); err != nil {
		log.Printf("error: %v\n", err)
		return
	}
}
