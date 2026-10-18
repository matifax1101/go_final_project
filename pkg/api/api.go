package api

import "net/http"

// Init регистрирует обработчики API с конфигурацией, считанной при старте.
func Init(mux *http.ServeMux, password string, signingKey []byte) {
	mux.HandleFunc("/api/nextdate", nextDateHandler)
	mux.HandleFunc("/api/signin", signinHandler(password, signingKey))
	mux.HandleFunc("/api/task", auth(password, signingKey, taskHandler))
	mux.HandleFunc("/api/tasks", auth(password, signingKey, tasksHandler))
	mux.HandleFunc("/api/task/done", auth(password, signingKey, doneHandler))
}
