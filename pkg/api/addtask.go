package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Gorthas/TODO-project/pkg/db"
)

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

func checkDate(task *db.Task) error {
	now := time.Now()
	var next string

	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return err
	}

	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = now.Format(dateFormat)
		} else {
			task.Date = next
		}
	}
	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		writeJSON(w, map[string]any{
			"error": err.Error(),
		})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]any{
			"error": "empty Title",
		})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]any{
			"error": err.Error(),
		})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]any{
			"error": err.Error(),
		})
		return
	}
	writeJSON(w, map[string]any{
		"id": strconv.FormatInt(id, 10),
	})
}
