package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/matifax1101/go_final_project/pkg/db"
)

func doneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, fmt.Errorf("не указан идентификатор"))
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err)
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, struct{}{})
		return
	}

	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := db.UpdateDate(next, id); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, struct{}{})
}
