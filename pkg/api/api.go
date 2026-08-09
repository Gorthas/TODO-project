package api

import "net/http"

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("POST /api/task", addTaskHandler)
	http.HandleFunc("GET /api/tasks", tasksHandler)
}
