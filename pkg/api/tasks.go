package api

import (
	"net/http"

	"github.com/matifax1101/go_final_project/pkg/db"
)

const tasksLimit = 50

type tasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}

	tasks, err := db.Tasks(tasksLimit)
	if err != nil {
		writeErrorStatus(w, http.StatusInternalServerError, err)
		return
	}
	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}
	writeJSON(w, tasksResponse{Tasks: tasks})
}
