package api

import (
	"GO_TODO-list/pkg/db"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func checkDate(task *db.Task) error {
	now := time.Now()
	var next string

	if task.Date == "" {
		task.Date = now.Format(DATE_FORMAT)
	}

	t, err := time.Parse(DATE_FORMAT, task.Date)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	if len(task.Repeat) > 0 {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("invalid input data: %w", err)
		}
	}

	if t.After(now) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(DATE_FORMAT)
		} else {
			task.Date = next
		}
	} else {
		task.Date = now.Format(DATE_FORMAT)
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task *db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	if len(task.Title) == 0 {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "empty title"})
		return
	}

	err = checkDate(task)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid date"})
		return
	}

	var id int64
	id, err = db.AddTask(task)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "error adding an task"})
		return
	}

	idStr := strconv.FormatInt(id, 10)

	err = WriteJSON(w, http.StatusCreated, map[string]string{"id": idStr})
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "failed to write JSON response"})
		return
	}

}
