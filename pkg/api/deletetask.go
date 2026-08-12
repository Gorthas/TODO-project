package api

import (
	"net/http"

	"github.com/Gorthas/TODO-project/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		writeJSON(w, map[string]any{
			"error": "task ID is not specified",
		})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, map[string]any{
			"error": err.Error(),
		})
		return
	}
	writeJSON(w, map[string]any{})
}
