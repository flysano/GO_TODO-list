package api

import (
	"GO_TODO-list/pkg/constants"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	dateStart, err := time.Parse(constants.DATE_FORMAT, dstart)
	if err != nil {
		return "", fmt.Errorf("failed to parsed date: %w", err)
	}

	dateStart = time.Date(dateStart.Year(), dateStart.Month(), dateStart.Day(), 0, 0, 0, 0, dateStart.Location())
	var rules []string

	if len(repeat) == 0 {
		return "", fmt.Errorf("incorrect data format: %w", err)
	}

	switch repeat[0] {
	case 'd', 'y', 'w', 'm':
		rules = strings.Split(repeat, " ")
	default:
		return "", fmt.Errorf("incorrect data format: %w", err)
	}

	if len(rules) == 0 {
		return "", fmt.Errorf("empty repeat rule")
	}

	switch rules[0] {
	case "d":
		if len(rules) < 2 {
			return "", fmt.Errorf("incorrect data format: %w", err)
		}
		interval, err := strconv.Atoi(rules[1])
		if err != nil {
			return "", fmt.Errorf("data conversion error: %w", err)
		}
		if interval <= 0 || interval > 400 {
			return "", fmt.Errorf("interval out of range (1-400)")
		}

		for {
			dateStart = dateStart.AddDate(0, 0, interval)
			if dateStart.After(now) {
				break
			}
		}

	case "w":
		if len(rules) < 2 {
			return "", fmt.Errorf("incorrect data format: %w", err)
		}
		var weekDays [8]bool
		weekDayStr := strings.Split(rules[1], ",")
		for _, wd := range weekDayStr {
			d, err := strconv.Atoi(wd)
			if err != nil || d < 1 || d > 7 {
				return "", fmt.Errorf("invalid weekday value: %s", wd)
			}
			weekDays[d] = true
		}

		for {
			dateStart = dateStart.AddDate(0, 0, 1)
			if dateStart.After(now) {
				wd := int(dateStart.Weekday())
				//корректировка индекса воскресенья с 0 на 7
				if wd == 0 {
					wd = 7
				}

				if weekDays[wd] {
					return dateStart.Format(constants.DATE_FORMAT), nil
				}
			}
		}

	case "m":
		if len(rules) < 2 {
			return "", fmt.Errorf("incorrect data format: %w", err)
		}
		var dayMonth [32]bool
		var monthYears [13]bool

		daysMonthStr := strings.Split(rules[1], ",")
		wantLast := false
		wantPenultimate := false

		for _, dm := range daysMonthStr {
			d, err := strconv.Atoi(dm)
			if err != nil {
				return "", fmt.Errorf("invalid weekday value: %s", dm)
			}
			switch {
			case d >= 1 && d <= 31:
				dayMonth[d] = true
			case d == -1:
				wantLast = true
			case d == -2:
				wantPenultimate = true
			default:
				return "", fmt.Errorf("day-of-month out of range: %d", d)
			}
		}

		if len(rules) > 2 {
			monthStrs := strings.Split(rules[2], ",")
			for _, monthStr := range monthStrs {
				m, err := strconv.Atoi(monthStr)
				if err != nil || m < 1 || m > 12 {
					return "", fmt.Errorf("invalid month value: %s", monthStr)
				}
				monthYears[m] = true
			}
		} else {
			for i := 1; i <= 12; i++ {
				monthYears[i] = true
			}
		}

		for {
			dateStart = dateStart.AddDate(0, 0, 1)
			if !dateStart.After(now) {
				continue
			}
			day := dateStart.Day()
			month := dateStart.Month()
			year := dateStart.Year()
			lastDayOfMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
			penultDay := lastDayOfMonth - 1

			if (dayMonth[day]) ||
				(wantLast && day == lastDayOfMonth) ||
				(wantPenultimate && day == penultDay) {
				if monthYears[int(month)] {
					return dateStart.Format(constants.DATE_FORMAT), nil
				}
			}
		}

	case "y":
		for {
			dateStart = dateStart.AddDate(1, 0, 0)
			if dateStart.After(now) {
				break
			}
		}
	default:
		return "", fmt.Errorf("incorrect repetition rule: %s", rules[0])
	}

	return dateStart.Format(constants.DATE_FORMAT), nil
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	var now time.Time
	var err error
	if nowStr != "" {
		now, err = time.Parse(constants.DATE_FORMAT, nowStr)
		if err != nil {
			http.Error(w, "invalid date format", http.StatusBadRequest)
			return
		}
	} else {
		now = time.Now()
	}
	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")
	if dstart == "" || repeat == "" {
		http.Error(w, "invalid input data", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, nextDate)
}
