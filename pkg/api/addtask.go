package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/matifax1101/go_final_project/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, err)
		return
	}

	if err := normalizeAndValidateTask(&task); err != nil {
		writeError(w, err)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, map[string]any{"id": id})
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, err)
		return
	}
	if task.ID == "" {
		writeError(w, fmt.Errorf("не указан идентификатор"))
		return
	}

	if err := normalizeAndValidateTask(&task); err != nil {
		writeError(w, err)
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, struct{}{})
}

func normalizeAndValidateTask(task *db.Task) error {
	if task.Title == "" {
		return fmt.Errorf("не указан заголовок задачи")
	}
	if len(task.Repeat) > 128 {
		return fmt.Errorf("правило повторения слишком длинное")
	}

	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	date, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return err
	}

	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		if !afterNow(date, now) {
			task.Date = next
		}
	}

	if !afterNow(date, now) && task.Repeat == "" {
		task.Date = now.Format(dateFormat)
	}

	return nil
}
