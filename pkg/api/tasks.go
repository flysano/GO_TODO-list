package api

import (
	"GO_TODO-list/pkg/db"
	"encoding/json"
	"net/http"
	"time"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// получение списка задач
func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []*db.Task

	tasks, err := db.Tasks(50)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Task Retrieval Error"})
		return
	}

	WriteJSON(w, http.StatusOK, TasksResp{
		Tasks: tasks,
	})
}

// получение одной задачи
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

	WriteJSON(w, http.StatusOK, task)
}

// обновление задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	if task.ID == "" {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор задачи"})
		return
	}

	if len(task.Title) == 0 {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "empty title"})
		return
	}

	err = checkDate(&task)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid date"})
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		return
	}

	WriteJSON(w, http.StatusOK, struct{}{})
}

// отметка задачи выполненной
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	//var task *db.Task
	var nextDate string

	id := r.URL.Query().Get("id")
	if id == "" {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
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
		nextDate, err = NextDate(now, task.Date, task.Repeat)
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

// удаление задачи
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := db.DeleteTask(id)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		return
	}

	WriteJSON(w, http.StatusOK, struct{}{})
}
