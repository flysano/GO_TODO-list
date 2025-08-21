package api

import (
	"GO_TODO-list/pkg/db"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// func getTasksHandler(w http.ResponseWriter, r *http.Request) {
// 	var tasks []*db.Task
// 	limit := 50

// 	tasks, err := db.Tasks(limit)
// 	if err != nil {
// 		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Task Retrieval Error"})
// 		return
// 	}

// 	WriteJSON(w, http.StatusOK, TasksResp{
// 		Tasks: tasks,
// 	})
// }

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []*db.Task
	var err error
	limit := 50
	search := r.URL.Query().Get("search")

	if strings.TrimSpace(search) == "" {
		tasks, err = db.Tasks(limit)
	} else {
		tasks, err = db.SearchTasks(search, limit)
	}
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Task Retrieval Error"})
		return
	}

	log.Printf("tasks received")
	WriteJSON(w, http.StatusOK, TasksResp{
		Tasks: tasks,
	})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "no id"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		return
	}
	log.Printf("task received")
	WriteJSON(w, http.StatusOK, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if task.ID == "" {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Не указан идентификатор задачи"})
		return
	}
	if len(task.Title) == 0 {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "empty title"})
		return
	}
	err = checkDate(&task)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "invalid date"})
		return
	}
	err = db.UpdateTask(&task)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "task not found"})
		return
	}
	log.Printf("task updated")
	WriteJSON(w, http.StatusOK, struct{}{})
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Printf("invalid ID")
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		log.Printf("faild getting task: %v", err)
		WriteJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			WriteJSON(w, http.StatusOK, struct{}{})
			return
		}
		WriteJSON(w, http.StatusOK, struct{}{})

	} else {
		now := time.Now()
		nextDate, err := NextDate(now, task.Date, task.Repeat)

		if err != nil {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid repeat"})
			return
		}

		if err = db.UpdateDate(nextDate, id); err != nil {
			WriteJSON(w, http.StatusOK, struct{}{})
			return
		}
		WriteJSON(w, http.StatusOK, struct{}{})
	}
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := db.DeleteTask(id)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			WriteJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		} else {
			WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
		return
	}
	log.Printf("task deleted")
	WriteJSON(w, http.StatusOK, struct{}{})
}
