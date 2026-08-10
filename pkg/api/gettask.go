package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Gorthas/TODO-project/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		writeJSON(w, map[string]any{
			"error": "task ID is not specified",
		})
		return
	}

	task, err := db.GetTask(id)
	if errors.Is(err, sql.ErrNoRows) {
		writeJSON(w, map[string]any{
			"error": "task not found",
		})
		return
	}
	if err != nil {
		writeJSON(w, map[string]any{
			"error": err.Error(),
		})
		return
	}
	writeJSON(w, task)
}
