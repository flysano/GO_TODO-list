package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DATE_FORMAT = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	//парсинг даты в time.Time
	dateStart, err := time.Parse(DATE_FORMAT, dstart)
	if err != nil {
		return "", fmt.Errorf("failed to parsed date: %w", err)
	}

	//сплит строки с правилами повторений
	var rules []string
	if len(repeat) > 0 && (rune(repeat[0]) == 'd' || rune(repeat[0]) == 'y') {
		rules = strings.Split(repeat, " ")
	} else {
		return "", fmt.Errorf("incorrect data format: %w", err)
	}

	// if len(repeat) > 0 && (rune(repeat[0]) == 'd' || rune(repeat[0]) == 'w' || rune(repeat[0]) == 'm' || rune(repeat[0]) == 'y') {
	// 	rules = strings.Split(repeat, " ")
	// } else {
	// 	return "", fmt.Errorf("incorrect date format: %w", err)
	// }

	//обработка правил по типу день, неделя, месяц, год

	switch rules[0] {
	case "d":
		if len(rules) < 2 {
			return "", fmt.Errorf("incorrect data format: %w", err)
		}
		interval, err := strconv.Atoi(rules[1])
		if err != nil {
			return "", fmt.Errorf("data conversion error: %w", err)
		}
		if interval > 400 || interval <= 0 {
			return "", fmt.Errorf("the incorrect length has been set for the delay of the event: d <= 400 && != 0. Set: %d", interval)
		}

		for {
			dateStart = dateStart.AddDate(0, 0, interval)
			if dateStart.After(now) {
				break
			}
		}

	// case "w":
	// 	fmt.Print("В разработке...")
	// case "m":
	// 	fmt.Print("В разработке...")

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

	return dateStart.Format(DATE_FORMAT), nil
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	var now time.Time
	var err error
	if nowStr != "" {
		now, err = time.Parse(DATE_FORMAT, nowStr)
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
