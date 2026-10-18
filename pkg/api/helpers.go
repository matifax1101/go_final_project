package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/matifax1101/go_final_project/pkg/db"
)

type errorResponse struct {
	Error string `json:"error"`
}

// writeJSONStatus сериализует ответ до отправки HTTP-статуса.
func writeJSONStatus(w http.ResponseWriter, status int, data any) {
	body, err := json.Marshal(data)
	if err != nil {
		log.Printf("encode JSON response: %v", err)
		http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if _, err := w.Write(append(body, '\n')); err != nil {
		log.Printf("write JSON response: %v", err)
	}
}

func writeJSON(w http.ResponseWriter, data any) {
	writeJSONStatus(w, http.StatusOK, data)
}

func writeErrorStatus(w http.ResponseWriter, status int, err error) {
	if status == http.StatusInternalServerError {
		log.Printf("internal API error: %v", err)
		writeJSONStatus(w, status, errorResponse{Error: "внутренняя ошибка сервера"})
		return
	}
	writeJSONStatus(w, status, errorResponse{Error: err.Error()})
}

func writeTaskError(w http.ResponseWriter, err error) {
	if errors.Is(err, db.ErrTaskNotFound) {
		writeErrorStatus(w, http.StatusNotFound, err)
		return
	}
	writeErrorStatus(w, http.StatusInternalServerError, err)
}

func methodNotAllowed(w http.ResponseWriter, allowed string) {
	w.Header().Set("Allow", allowed)
	writeErrorStatus(w, http.StatusMethodNotAllowed, errors.New("метод не поддерживается"))
}
