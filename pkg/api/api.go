package api

import "net/http"

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("POST /api/task", addTaskHandler)
}
