package api

import "net/http"

func Init() {
	http.HandleFunc("GET /api/nextdate", nextDayHandler)
	http.HandleFunc("POST /api/task", addTaskHandler)
	http.HandleFunc("GET /api/tasks", tasksHandler)
	http.HandleFunc("GET /api/task", getTaskHandler)
	http.HandleFunc("PUT /api/task", updateTaskHandler)
	http.HandleFunc("POST /api/task/done", doneTaskHandler)
	http.HandleFunc("DELETE /api/task", deleteTaskHandler)
}
