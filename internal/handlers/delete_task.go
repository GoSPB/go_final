package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/GoSPB/go_final/internal/database"
	"net/http"
)

func DeleteTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		id := r.URL.Query().Get("id")
		if id == "" {
			json.NewEncoder(w).Encode(map[string]string{"error": "Не указан идентификатор задачи"})
			return
		}

		err := database.DeleteTask(db, id)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка при удалении задачи"})
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
