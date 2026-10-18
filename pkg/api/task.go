package api

import (
	"errors"
	"net/http"

	"github.com/matifax1101/go_final_project/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		methodNotAllowed(w, "GET, POST, PUT, DELETE")
	}
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeErrorStatus(w, http.StatusBadRequest, errors.New("не указан идентификатор"))
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeTaskError(w, err)
		return
	}
	writeJSON(w, task)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeErrorStatus(w, http.StatusBadRequest, errors.New("не указан идентификатор"))
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeTaskError(w, err)
		return
	}
	writeJSON(w, struct{}{})
}
