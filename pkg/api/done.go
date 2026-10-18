package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/matifax1101/go_final_project/pkg/db"
)

func doneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}

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

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeTaskError(w, err)
			return
		}
		writeJSON(w, struct{}{})
		return
	}

	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeErrorStatus(w, http.StatusInternalServerError, err)
		return
	}
	if err := db.UpdateDate(next, id); err != nil {
		writeTaskError(w, err)
		return
	}
	writeJSON(w, struct{}{})
}
