package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/GoSPB/go_final/internal/database"
	"net/http"
)

func GetTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		id := r.URL.Query().Get("id")
		if id == "" {
			json.NewEncoder(w).Encode(map[string]string{"error": "Не указан идентификатор задачи"})
			return
		}

		task, err := database.GetTask(db, id)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Ошибка при получении задачи: %v", err)})
			return
		}

		if err := json.NewEncoder(w).Encode(task); err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка при формировании ответа"})
		}
	}
}
