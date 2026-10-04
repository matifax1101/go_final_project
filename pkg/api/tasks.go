package api

import (
	"net/http"

	"github.com/matifax1101/go_final_project/pkg/db"
)

type tasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50)
	if err != nil {
		writeError(w, err)
		return
	}
	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}
	writeJSON(w, tasksResponse{Tasks: tasks})
}
