package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/GoSPB/go_final/internal/database"
	"github.com/GoSPB/go_final/internal/models"
)

func GetTasks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		tasks, err := database.GetTasks(db)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка при получении задач"})
			return
		}

		response := models.TasksResponse{Tasks: tasks}
		if tasks == nil {
			response.Tasks = []models.Task{}
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка при формировании ответа"})
		}
	}
}
