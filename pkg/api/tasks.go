package api

import (
	"GO_TODO-list/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {

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
